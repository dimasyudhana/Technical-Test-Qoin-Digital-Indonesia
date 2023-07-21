package controller

import (
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/user"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/utils/time"
)

type registerResponse struct {
	UserID string `json:"user_id,omitempty"`
}

type loginResponse struct {
	UserID string `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}

type profileResponse struct {
	UserID         string         `json:"user_id"`
	Username       string         `json:"username"`
	Email          string         `json:"email"`
	Role           string         `json:"role"`
	Status         string         `json:"status"`
	ProfilePicture string         `json:"profile_picture"`
	CreatedAt      time.LocalTime `json:"created_at"`
	UpdatedAt      time.LocalTime `json:"updated_at"`
}

func register(u user.UserCore) registerResponse {
	return registerResponse{
		UserID: u.UserID,
	}
}

func profileUser(u user.UserCore) profileResponse {
	return profileResponse{
		UserID:         u.UserID,
		Username:       u.Username,
		Email:          u.Email,
		Role:           u.Role,
		Status:         u.Status,
		ProfilePicture: u.ProfilePicture,
		CreatedAt:      time.LocalTime(u.CreatedAt),
		UpdatedAt:      time.LocalTime(u.UpdatedAt),
	}
}
