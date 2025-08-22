package domain

import "time"

type Profile struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	AvatarURL string    `json:"avatar_url"`
	UpdatedAt time.Time `json:"updated_at"`
}
