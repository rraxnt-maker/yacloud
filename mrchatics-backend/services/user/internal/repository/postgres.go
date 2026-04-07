package repository

import (
    "database/sql"
    "fmt"
    "time"
    "github.com/google/uuid"
    "your-project/internal/models"
)

type ProfileRepository struct {
    db *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
    return &ProfileRepository{db: db}
}

// GetProfile retrieves user profile by user_id
func (r *ProfileRepository) GetProfile(userID uuid.UUID) (*models.UserProfile, error) {
    query := `
        SELECT user_id, full_name, bio, avatar_url, phone, location, 
               website, company, position, birth_date, gender, 
               language, timezone, created_at, updated_at
        FROM user_profiles WHERE user_id = $1
    `
    
    profile := &models.UserProfile{}
    var birthDate sql.NullTime
    
    err := r.db.QueryRow(query, userID).Scan(
        &profile.UserID, &profile.FullName, &profile.Bio, &profile.AvatarURL,
        &profile.Phone, &profile.Location, &profile.Website, &profile.Company,
        &profile.Position, &birthDate, &profile.Gender, &profile.Language,
        &profile.Timezone, &profile.CreatedAt, &profile.UpdatedAt,
    )
    
    if err != nil {
        return nil, err
    }
    
    if birthDate.Valid {
        profile.BirthDate = &birthDate.Time
    }
    
    return profile, nil
}

// UpdateProfile updates user profile
func (r *ProfileRepository) UpdateProfile(userID uuid.UUID, req *models.UpdateProfileRequest) error {
    query := `
        INSERT INTO user_profiles (user_id, full_name, bio, phone, location, website, company, position, birth_date, gender)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        ON CONFLICT (user_id) DO UPDATE SET
            full_name = EXCLUDED.full_name,
            bio = EXCLUDED.bio,
            phone = EXCLUDED.phone,
            location = EXCLUDED.location,
            website = EXCLUDED.website,
            company = EXCLUDED.company,
            position = EXCLUDED.position,
            birth_date = EXCLUDED.birth_date,
            gender = EXCLUDED.gender,
            updated_at = NOW()
    `
    
    var birthDate interface{}
    if req.BirthDate != nil {
        parsed, err := time.Parse(time.RFC3339, *req.BirthDate)
        if err != nil {
            return err
        }
        birthDate = parsed
    } else {
        birthDate = nil
    }
    
    _, err := r.db.Exec(query, userID, req.FullName, req.Bio, req.Phone,
        req.Location, req.Website, req.Company, req.Position, birthDate, req.Gender)
    
    return err
}

// GetSettings retrieves user settings
func (r *ProfileRepository) GetSettings(userID uuid.UUID) (*models.UserSettings, error) {
    query := `
        SELECT user_id, email_notifications, push_notifications, slack_notifications,
               profile_visibility, show_online_status, show_last_seen, show_email, show_phone,
               theme, compact_view, two_factor_enabled, session_timeout_minutes,
               default_dashboard, items_per_page, created_at, updated_at
        FROM user_settings WHERE user_id = $1
    `
    
    settings := &models.UserSettings{}
    err := r.db.QueryRow(query, userID).Scan(
        &settings.UserID, &settings.EmailNotifications, &settings.PushNotifications,
        &settings.SlackNotifications, &settings.ProfileVisibility, &settings.ShowOnlineStatus,
        &settings.ShowLastSeen, &settings.ShowEmail, &settings.ShowPhone, &settings.Theme,
        &settings.CompactView, &settings.TwoFactorEnabled, &settings.SessionTimeoutMinutes,
        &settings.DefaultDashboard, &settings.ItemsPerPage, &settings.CreatedAt, &settings.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        // Create default settings if not exists
        return r.createDefaultSettings(userID)
    }
    
    return settings, err
}

func (r *ProfileRepository) createDefaultSettings(userID uuid.UUID) (*models.UserSettings, error) {
    query := `
        INSERT INTO user_settings (user_id)
        VALUES ($1)
        RETURNING user_id, email_notifications, push_notifications, slack_notifications,
                  profile_visibility, show_online_status, show_last_seen, show_email, show_phone,
                  theme, compact_view, two_factor_enabled, session_timeout_minutes,
                  default_dashboard, items_per_page, created_at, updated_at
    `
    
    settings := &models.UserSettings{}
    err := r.db.QueryRow(query, userID).Scan(
        &settings.UserID, &settings.EmailNotifications, &settings.PushNotifications,
        &settings.SlackNotifications, &settings.ProfileVisibility, &settings.ShowOnlineStatus,
        &settings.ShowLastSeen, &settings.ShowEmail, &settings.ShowPhone, &settings.Theme,
        &settings.CompactView, &settings.TwoFactorEnabled, &settings.SessionTimeoutMinutes,
        &settings.DefaultDashboard, &settings.ItemsPerPage, &settings.CreatedAt, &settings.UpdatedAt,
    )
    
    return settings, err
}

// UpdateSettings updates user settings
func (r *ProfileRepository) UpdateSettings(userID uuid.UUID, req *models.UpdateSettingsRequest) error {
    query := `
        UPDATE user_settings SET
            email_notifications = COALESCE($2, email_notifications),
            push_notifications = COALESCE($3, push_notifications),
            slack_notifications = COALESCE($4, slack_notifications),
            profile_visibility = COALESCE($5, profile_visibility),
            show_online_status = COALESCE($6, show_online_status),
            show_last_seen = COALESCE($7, show_last_seen),
            show_email = COALESCE($8, show_email),
            show_phone = COALESCE($9, show_phone),
            theme = COALESCE($10, theme),
            compact_view = COALESCE($11, compact_view),
            two_factor_enabled = COALESCE($12, two_factor_enabled),
            session_timeout_minutes = COALESCE($13, session_timeout_minutes),
            default_dashboard = COALESCE($14, default_dashboard),
            items_per_page = COALESCE($15, items_per_page),
            updated_at = NOW()
        WHERE user_id = $1
    `
    
    _, err := r.db.Exec(query, userID,
        req.EmailNotifications, req.PushNotifications, req.SlackNotifications,
        req.ProfileVisibility, req.ShowOnlineStatus, req.ShowLastSeen,
        req.ShowEmail, req.ShowPhone, req.Theme, req.CompactView,
        req.TwoFactorEnabled, req.SessionTimeoutMinutes, req.DefaultDashboard,
        req.ItemsPerPage)
    
    return err
}
