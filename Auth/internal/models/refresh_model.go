package model

import "time"

// RefreshTokens model
type RefreshToken struct {
    ID          int       `json:"id"`
    UserID      int       `json:"user_id"`
    TokenText   string    `json:"-"` // Never export the hash to JSON for security
    ClientType  string    `json:"client_type"` // 'web' or 'mobile'
    ExpiresAt   time.Time `json:"expires_at"`
    Revoked     bool      `json:"revoked"`
    CreatedAt   time.Time `json:"created_at"`
}