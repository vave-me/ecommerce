package models

type User struct {
	ID         string  `json:"id"`
	Email      string  `json:"email"`
	Username   string  `json:"user_name"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Enabled    bool    `json:"enabled"`
	GoogleID   string  `json:"google_id"`
	Location   string  `json:"location"`
	Lat        float32 `json:"lat"`
	Lng        float32 `json:"lng"`
	Thumbnail  string  `json:"thumbnail"`
	Background string  `json:"background"`
	Bio        string  `json:"bio"`
	Privacy    string  `json:"privacy"`
}

// BaseUser represents the simplified user information
type BaseUser struct {
	ID         string  `json:"id"`
	Username   string  `json:"user_name"`
	Thumbnail  string  `json:"thumbnail"`
	Lat        float32 `json:"lat"`
	Lng        float32 `json:"lng"`
	Location   string  `json:"location"`
	Bio        string  `json:"bio"`
	Privacy    string  `json:"privacy"`
	Background string  `json:"background"`
}

// UserUpdateInfo represents user update information
type UserUpdateInfo struct {
	Email     string `json:"email"`
	Username  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Location  string `json:"location"`
}

// LoginResponse represents login response
type LoginResponse struct {
	AccessToken string  `json:"access_token"`
	Token       string  `json:"token"`
	Username    string  `json:"user_name"`
	Lat         float32 `json:"lat"`
	Lng         float32 `json:"lng"`
}

// TokenResponse represents token refresh response
type TokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// MessageResponse represents generic message response
type MessageResponse struct {
	Message string `json:"message"`
}

// ClearTokensResponse represents clear tokens response
type ClearTokensResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
