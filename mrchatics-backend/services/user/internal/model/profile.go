package models

import (
    "time"
    "github.com/google/uuid"
)

type UserProfile struct {
    UserID     uuid.UUID  `json:"user_id"`
    FullName   *string    `json:"full_name,omitempty"`
    Bio        *string    `json:"bio,omitempty"`
    AvatarURL  *string    `json:"avatar_url,omitempty"`
    Phone      *string    `json:"phone,omitempty"`
    Location   *string    `json:"location,omitempty"`
    Website    *string    `json:"website,omitempty"`
    Company    *string    `json:"company,omitempty"`
    Position   *string    `json:"position,omitempty"`
    BirthDate  *time.Time `json:"birth_date,omitempty"`
    Gender     *string    `json:"gender,omitempty"`
    Language   string     `json:"language"`
    Timezone   string     `json:"timezone"`
    CreatedAt  time.Time  `json:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at"`
}

type UpdateProfileRequest struct {
    FullName  *string `json:"full_name"`
    Bio       *string `json:"bio"`
    Phone     *string `json:"phone"`
    Location  *string `json:"location"`
    Website   *string `json:"website"`
    Company   *string `json:"company"`
    Position  *string `json:"position"`
    BirthDate *string `json:"birth_date"` // ISO 8601 format
    Gender    *string `json:"gender"`
}
