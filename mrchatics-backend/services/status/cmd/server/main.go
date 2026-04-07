package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"
    "bytes"
    "io/ioutil"

    _ "github.com/lib/pq"
)

type StatusResponse struct {
    UserID       int       `json:"user_id"`
    Status       string    `json:"status"`
    CustomStatus string    `json:"custom_status,omitempty"`
    LastSeen     time.Time `json:"last_seen"`
}

var cacheServiceURL string

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getDB() *sql.DB {
    host := getEnv("DB_HOST", "postgres")
    port := getEnv("DB_PORT", "5432")
    user := getEnv("DB_USER", "postgres")
    password := getEnv("DB_PASSWORD", "postgres")
    dbname := getEnv("DB_NAME", "status_service")

    connStr := "host=" + host + " port=" + port + " user=" + user +
        " password=" + password + " dbname=" + dbname + " sslmode=disable"

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    return db
}

func callCacheService(key string) (string, error) {
    url := cacheServiceURL + "/api/v1/cache/get/" + key
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", nil
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return "", err
    }

    if value, ok := result["value"]; ok {
        return value.(string), nil
    }
    return "", nil
}

func saveToCache(key string, value interface{}, ttl int) error {
    url := cacheServiceURL + "/api/v1/cache"
    data := map[string]interface{}{
        "key":         key,
        "value":       value,
        "ttl_seconds": ttl,
    }
    jsonData, _ := json.Marshal(data)

    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return nil
}

func deleteFromCache(key string) error {
    url := cacheServiceURL + "/api/v1/cache/del/" + key
    req, err := http.NewRequest(http.MethodDelete, url, nil)
    if err != nil {
        return err
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return nil
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "service": "status-service"})
}

func getStatus(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    parts := strings.Split(path, "/")

    var userIDStr string
    for i := len(parts) - 1; i >= 0; i-- {
        if parts[i] != "" && parts[i] != "api" && parts[i] != "v1" && parts[i] != "status" {
            userIDStr = parts[i]
            break
        }
    }

    if userIDStr == "" {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    cacheKey := "status:user:" + strconv.Itoa(userID)
    cachedStatus, err := callCacheService(cacheKey)

    if err == nil && cachedStatus != "" {
        log.Printf("Cache HIT for user %d", userID)
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("X-Cache", "HIT")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "user_id": userID,
            "status":  cachedStatus,
            "source":  "cache",
        })
        return
    }

    log.Printf("Cache MISS for user %d, fetching from DB", userID)

    db := getDB()
    defer db.Close()

    var status StatusResponse
    query := `SELECT user_id, status, COALESCE(custom_status, ''), last_seen 
              FROM user_statuses WHERE user_id = $1`

    err = db.QueryRow(query, userID).Scan(&status.UserID, &status.Status, &status.CustomStatus, &status.LastSeen)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    go saveToCache(cacheKey, status.Status, 60)

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Cache", "MISS")
    json.NewEncoder(w).Encode(status)
}

func updateStatus(w http.ResponseWriter, r *http.Request) {
    var req struct {
        UserID       int    `json:"user_id"`
        Status       string `json:"status"`
        CustomStatus string `json:"custom_status"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    db := getDB()
    defer db.Close()

    query := `INSERT INTO user_statuses (user_id, status, custom_status, last_seen, updated_at)
              VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
              ON CONFLICT (user_id) 
              DO UPDATE SET status = $2, custom_status = $3, last_seen = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP`

    _, err := db.Exec(query, req.UserID, req.Status, req.CustomStatus)
    if err != nil {
        http.Error(w, "Failed to update status", http.StatusInternalServerError)
        return
    }

    cacheKey := "status:user:" + strconv.Itoa(req.UserID)
    go deleteFromCache(cacheKey)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Status updated successfully"})
}

func getOnlineUsers(w http.ResponseWriter, r *http.Request) {
    db := getDB()
    defer db.Close()

    query := `SELECT user_id FROM user_statuses 
              WHERE status IN ('online', 'away', 'busy') 
              AND last_seen > NOW() - INTERVAL '5 minutes'`

    rows, err := db.Query(query)
    if err != nil {
        http.Error(w, "Failed to get online users", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []int
    for rows.Next() {
        var userID int
        rows.Scan(&userID)
        users = append(users, userID)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string][]int{"online_users": users})
}

func testCacheConnection(w http.ResponseWriter, r *http.Request) {
    testKey := "test_connection_" + strconv.FormatInt(time.Now().Unix(), 10)

    err := saveToCache(testKey, "connection_test", 10)
    if err != nil {
        json.NewEncoder(w).Encode(map[string]interface{}{
            "cache_service": "NOT_CONNECTED",
            "error":         err.Error(),
            "cache_url":     cacheServiceURL,
        })
        return
    }

    value, err := callCacheService(testKey)
    if err != nil || value == "" {
        json.NewEncoder(w).Encode(map[string]interface{}{
            "cache_service": "CONNECTION_ERROR",
            "error":         err,
        })
        return
    }

    deleteFromCache(testKey)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "cache_service": "CONNECTED",
        "test_key":      testKey,
        "test_value":    value,
        "cache_url":     cacheServiceURL,
    })
}

func main() {
    cacheServiceURL = getEnv("CACHE_SERVICE_URL", "http://cache-service:8082")
    port := getEnv("SERVER_PORT", "8081")

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        path := r.URL.Path
        log.Printf("Request: %s %s", r.Method, path)

        switch {
        case path == "/health" && r.Method == "GET":
            healthCheck(w, r)
        case strings.HasPrefix(path, "/api/v1/status/") && r.Method == "GET":
            getStatus(w, r)
        case path == "/api/v1/status" && r.Method == "PUT":
            updateStatus(w, r)
        case path == "/api/v1/status/online" && r.Method == "GET":
            getOnlineUsers(w, r)
        case path == "/api/v1/test/cache" && r.Method == "GET":
            testCacheConnection(w, r)
        default:
            http.NotFound(w, r)
        }
    })

    log.Printf("Status Service starting on port %s", port)
    log.Printf("Cache Service URL: %s", cacheServiceURL)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
