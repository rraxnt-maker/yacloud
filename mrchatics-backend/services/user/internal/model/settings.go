package models

import (
    "time"
    "github.com/google/uuid"
)

type UserSettings struct {
    UserID               uuid.UUID `json:"user_id"`
    EmailNotifications   bool      `json:"email_notifications"`
    PushNotifications    bool      `json:"push_notifications"`
    SlackNotifications   bool      `json:"slack_notifications"`
    ProfileVisibility    string    `json:"profile_visibility"`
    ShowOnlineStatus     bool      `json:"show_online_status"`
    ShowLastSeen         bool      `json:"show_last_seen"`
    ShowEmail            bool      `json:"show_email"`
    ShowPhone            bool      `json:"show_phone"`
    Theme                string    `json:"theme"`
    CompactView          bool      `json:"compact_view"`
    TwoFactorEnabled     bool      `json:"two_factor_enabled"`
    SessionTimeoutMinutes int      `json:"session_timeout_minutes"`
    DefaultDashboard     string    `json:"default_dashboard"`
    ItemsPerPage         int       `json:"items_per_page"`
    CreatedAt            time.Time `json:"created_at"`
    UpdatedAt            time.Time `json:"updated_at"`
}

type UpdateSettingsRequest struct {
    EmailNotifications    *bool   `json:"email_notifications"`
    PushNotifications     *bool   `json:"push_notifications"`
    SlackNotifications    *bool   `json:"slack_notifications"`
    ProfileVisibility     *string `json:"profile_visibility"`
    ShowOnlineStatus      *bool   `json:"show_online_status"`
    ShowLastSeen          *bool   `json:"show_last_seen"`
    ShowEmail             *bool   `json:"show_email"`
    ShowPhone             *bool   `json:"show_phone"`
    Theme                 *string `json:"theme"`
    CompactView           *bool   `json:"compact_view"`
    TwoFactorEnabled      *bool   `json:"two_factor_enabled"`
    SessionTimeoutMinutes *int    `json:"session_timeout_minutes"`
    DefaultDashboard      *string `json:"default_dashboard"`
    ItemsPerPage          *int    `json:"items_per_page"`
}
