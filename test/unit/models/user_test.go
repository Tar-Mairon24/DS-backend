package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"ds-backend/internal/domain/models"
)

func TestToUserResponse_NilUser(t *testing.T) {
	var user *models.User = nil
	resp := user.ToUserResponse()
	assert.Nil(t, resp)
}

func TestToUserResponse_ValidUser(t *testing.T) {
	now := time.Now()
	user := &models.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: now,
		UpdatedAt: &now,
	}
	resp := user.ToUserResponse()
	assert.NotNil(t, resp)
	assert.Equal(t, user.ID, resp.ID)
	assert.Equal(t, user.Username, resp.Username)
	assert.Equal(t, user.Email, resp.Email)
	assert.Equal(t, user.CreatedAt, resp.CreatedAt)
	assert.Equal(t, user.UpdatedAt, resp.UpdatedAt)
}

// TestToUserResponse_TableDriven tests the ToUserResponse method with various user roles and scenarios
func TestToUserResponse_TableDriven(t *testing.T) {
	now := time.Now()
	phone := "+1234567890"
	notes := "Test notes"

	tests := []struct {
		name     string
		user     *models.User
		wantNil  bool
		validate func(t *testing.T, resp *models.UserResponse)
	}{
		{
			name:    "nil user returns nil response",
			user:    nil,
			wantNil: true,
		},
		{
			name: "admin role user",
			user: &models.User{
				ID:           1,
				Username:     "admin_user",
				Email:        "admin@example.com",
				Password:     "hashed_password",
				Role:         "admin",
				Verified:     true,
				Phone:        &phone,
				MFAActivated: true,
				Notes:        &notes,
				CreatedAt:    now,
				UpdatedAt:    &now,
			},
			wantNil: false,
			validate: func(t *testing.T, resp *models.UserResponse) {
				assert.Equal(t, uint(1), resp.ID)
				assert.Equal(t, "admin_user", resp.Username)
				assert.Equal(t, "admin@example.com", resp.Email)
				assert.Equal(t, "admin", resp.Role)
				assert.True(t, resp.Verified)
				assert.Equal(t, &phone, resp.Phone)
				assert.True(t, resp.MFAActivated)
				assert.Equal(t, &notes, resp.Notes)
				assert.Equal(t, now, resp.CreatedAt)
				assert.Equal(t, &now, resp.UpdatedAt)
			},
		},
		{
			name: "agente role user",
			user: &models.User{
				ID:           2,
				Username:     "agent_user",
				Email:        "agent@example.com",
				Password:     "hashed_password",
				Role:         "agente",
				Verified:     false,
				Phone:        nil,
				MFAActivated: false,
				Notes:        nil,
				CreatedAt:    now,
				UpdatedAt:    nil,
			},
			wantNil: false,
			validate: func(t *testing.T, resp *models.UserResponse) {
				assert.Equal(t, uint(2), resp.ID)
				assert.Equal(t, "agent_user", resp.Username)
				assert.Equal(t, "agent@example.com", resp.Email)
				assert.Equal(t, "agente", resp.Role)
				assert.False(t, resp.Verified)
				assert.Nil(t, resp.Phone)
				assert.False(t, resp.MFAActivated)
				assert.Nil(t, resp.Notes)
				assert.Equal(t, now, resp.CreatedAt)
				assert.Nil(t, resp.UpdatedAt)
			},
		},
		{
			name: "owner role user",
			user: &models.User{
				ID:           3,
				Username:     "owner_user",
				Email:        "owner@example.com",
				Password:     "hashed_password",
				Role:         "owner",
				Verified:     true,
				Phone:        &phone,
				MFAActivated: false,
				Notes:        &notes,
				CreatedAt:    now,
				UpdatedAt:    &now,
			},
			wantNil: false,
			validate: func(t *testing.T, resp *models.UserResponse) {
				assert.Equal(t, uint(3), resp.ID)
				assert.Equal(t, "owner_user", resp.Username)
				assert.Equal(t, "owner@example.com", resp.Email)
				assert.Equal(t, "owner", resp.Role)
				assert.True(t, resp.Verified)
				assert.Equal(t, &phone, resp.Phone)
				assert.False(t, resp.MFAActivated)
				assert.Equal(t, &notes, resp.Notes)
			},
		},
		{
			name: "user with minimal fields",
			user: &models.User{
				ID:        4,
				Username:  "minimal_user",
				Email:     "minimal@example.com",
				Role:      "client",
				CreatedAt: now,
			},
			wantNil: false,
			validate: func(t *testing.T, resp *models.UserResponse) {
				assert.Equal(t, uint(4), resp.ID)
				assert.Equal(t, "minimal_user", resp.Username)
				assert.Equal(t, "minimal@example.com", resp.Email)
				assert.Equal(t, "client", resp.Role)
				assert.False(t, resp.Verified)
				assert.Nil(t, resp.Phone)
				assert.False(t, resp.MFAActivated)
				assert.Nil(t, resp.Notes)
				assert.Equal(t, now, resp.CreatedAt)
				assert.Nil(t, resp.UpdatedAt)
			},
		},
		{
			name: "user with all optional fields populated",
			user: &models.User{
				ID:           5,
				Username:     "full_user",
				Email:        "full@example.com",
				Password:     "hashed_password",
				Role:         "admin",
				Verified:     true,
				Phone:        &phone,
				MFAActivated: true,
				Notes:        &notes,
				CreatedAt:    now,
				UpdatedAt:    &now,
			},
			wantNil: false,
			validate: func(t *testing.T, resp *models.UserResponse) {
				assert.NotNil(t, resp)
				assert.Equal(t, uint(5), resp.ID)
				assert.Equal(t, "full_user", resp.Username)
				assert.Equal(t, "full@example.com", resp.Email)
				assert.Equal(t, "admin", resp.Role)
				assert.True(t, resp.Verified)
				assert.NotNil(t, resp.Phone)
				assert.Equal(t, phone, *resp.Phone)
				assert.True(t, resp.MFAActivated)
				assert.NotNil(t, resp.Notes)
				assert.Equal(t, notes, *resp.Notes)
				assert.Equal(t, now, resp.CreatedAt)
				assert.NotNil(t, resp.UpdatedAt)
				assert.Equal(t, now, *resp.UpdatedAt)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := tt.user.ToUserResponse()

			if tt.wantNil {
				assert.Nil(t, resp)
			} else {
				assert.NotNil(t, resp)
				if tt.validate != nil {
					tt.validate(t, resp)
				}
			}
		})
	}
}
