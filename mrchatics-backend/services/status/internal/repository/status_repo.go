package repository

import (
	"database/sql"
	"fmt"
	"status-service/internal/model"
)

type StatusRepository struct {
	db *sql.DB
}

func NewStatusRepository(db *sql.DB) *StatusRepository {
	return &StatusRepository{db: db}
}

func (r *StatusRepository) GetStatus(userID int) (*model.StatusResponse, error) {
	var status model.StatusResponse
	query := `SELECT user_id, status, COALESCE(custom_status, ''), last_seen 
	          FROM user_statuses WHERE user_id = $1`
	
	err := r.db.QueryRow(query, userID).Scan(
		&status.UserID,
		&status.Status,
		&status.CustomStatus,
		&status.LastSeen,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	
	return &status, nil
}

func (r *StatusRepository) UpdateStatus(req *model.UpdateStatusRequest) error {
	query := `INSERT INTO user_statuses (user_id, status, custom_status, last_seen, updated_at)
	          VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	          ON CONFLICT (user_id) 
	          DO UPDATE SET status = $2, custom_status = $3, last_seen = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP`
	
	_, err := r.db.Exec(query, req.UserID, req.Status, req.CustomStatus)
	return err
}

func (r *StatusRepository) GetBatchStatuses(userIDs []int) (map[int]model.StatusResponse, error) {
	if len(userIDs) == 0 {
		return make(map[int]model.StatusResponse), nil
	}
	
	query := `SELECT user_id, status, COALESCE(custom_status, ''), last_seen 
	          FROM user_statuses WHERE user_id = ANY($1::int[])`
	
	rows, err := r.db.Query(query, userIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	statuses := make(map[int]model.StatusResponse)
	for rows.Next() {
		var status model.StatusResponse
		err := rows.Scan(&status.UserID, &status.Status, &status.CustomStatus, &status.LastSeen)
		if err != nil {
			return nil, err
		}
		statuses[status.UserID] = status
	}
	
	return statuses, nil
}

func (r *StatusRepository) GetOnlineUsers() ([]int, error) {
	query := `SELECT user_id FROM user_statuses 
	          WHERE status IN ('online', 'away', 'busy') 
	          AND last_seen > NOW() - INTERVAL '5 minutes'`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		users = append(users, userID)
	}
	
	return users, nil
}

func (r *StatusRepository) UpdateLastSeen(userID int) error {
	query := `UPDATE user_statuses SET last_seen = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
	          WHERE user_id = $1`
	
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *StatusRepository) CreateUserStatus(userID int) error {
	query := `INSERT INTO user_statuses (user_id, status, last_seen) 
	          VALUES ($1, 'offline', CURRENT_TIMESTAMP)
	          ON CONFLICT (user_id) DO NOTHING`
	
	_, err := r.db.Exec(query, userID)
	return err
}
