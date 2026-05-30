package models_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ds-backend/internal/domain/models"
)

// TestAppointmentScan tests the Scan method for Appointment type
func TestAppointmentScan(t *testing.T) {
	notes := "Test notes"

	tests := []struct {
		name    string
		input   any
		want    *models.Appointment
		wantErr bool
		errMsg  string
	}{
		{
			name:    "success - nil value",
			input:   nil,
			want:    nil,
			wantErr: false,
		},
		{
			name: "success - valid JSON bytes",
			input: []byte(`{
				"id": 1,
				"title": "Meeting",
				"description": "Team meeting",
				"start_date": "2024-01-15T10:00:00Z",
				"end_date": "2024-01-15T11:00:00Z",
				"status": "scheduled",
				"client_id": 1,
				"owner_id": 2,
				"property_id": 3,
				"created_at": "2024-01-15T09:00:00Z"
			}`),
			want: &models.Appointment{
				ID:          1,
				Title:       "Meeting",
				Description: "Team meeting",
				Status:      "scheduled",
				ClientID:    1,
				OwnerID:     2,
				PropertyID:  3,
			},
			wantErr: false,
		},
		{
			name: "success - valid JSON string",
			input: `{
				"id": 2,
				"title": "Property Viewing",
				"description": "Client viewing",
				"start_date": "2024-01-16T14:00:00Z",
				"end_date": "2024-01-16T15:00:00Z",
				"status": "completed",
				"notes": "Test notes",
				"client_id": 2,
				"owner_id": 3,
				"property_id": 4,
				"created_at": "2024-01-16T13:00:00Z"
			}`,
			want: &models.Appointment{
				ID:          2,
				Title:       "Property Viewing",
				Description: "Client viewing",
				Status:      "completed",
				Notes:       &notes,
				ClientID:    2,
				OwnerID:     3,
				PropertyID:  4,
			},
			wantErr: false,
		},
		{
			name: "success - minimal appointment",
			input: []byte(`{
				"id": 3,
				"title": "Appointment",
				"status": "scheduled",
				"client_id": 1,
				"owner_id": 1,
				"property_id": 1,
				"created_at": "2024-01-15T09:00:00Z"
			}`),
			want: &models.Appointment{
				ID:         3,
				Title:      "Appointment",
				Status:     "scheduled",
				ClientID:   1,
				OwnerID:    1,
				PropertyID: 1,
			},
			wantErr: false,
		},
		{
			name:    "error - invalid type (integer)",
			input:   123,
			want:    nil,
			wantErr: true,
			errMsg:  "cannot scan Appointment",
		},
		{
			name:    "error - invalid type (float)",
			input:   123.45,
			want:    nil,
			wantErr: true,
			errMsg:  "cannot scan Appointment",
		},
		{
			name:    "error - invalid type (bool)",
			input:   true,
			want:    nil,
			wantErr: true,
			errMsg:  "cannot scan Appointment",
		},
		{
			name:    "error - invalid JSON bytes",
			input:   []byte(`not a json`),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error - invalid JSON string",
			input:   `{invalid json}`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error - empty JSON object",
			input:   []byte(`{}`),
			want:    &models.Appointment{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a models.Appointment
			err := a.Scan(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				if tt.want != nil {
					assert.Equal(t, tt.want.ID, a.ID)
					assert.Equal(t, tt.want.Title, a.Title)
					assert.Equal(t, tt.want.Description, a.Description)
					assert.Equal(t, tt.want.Status, a.Status)
					assert.Equal(t, tt.want.ClientID, a.ClientID)
					assert.Equal(t, tt.want.OwnerID, a.OwnerID)
					assert.Equal(t, tt.want.PropertyID, a.PropertyID)
				}
			}
		})
	}
}

// TestAppointmentValue tests the Value method for Appointment type
func TestAppointmentValue(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	notes := "Test notes"

	tests := []struct {
		name    string
		input   *models.Appointment
		wantErr bool
		verify  func(t *testing.T, val any)
	}{
		{
			name: "success - full appointment",
			input: &models.Appointment{
				ID:          1,
				Title:       "Meeting",
				Description: "Team meeting",
				StartDate:   now,
				EndDate:     now.Add(1 * time.Hour),
				Status:      models.AppointmentStatusScheduled,
				Notes:       &notes,
				ClientID:    1,
				OwnerID:     2,
				PropertyID:  3,
				CreatedAt:   now,
			},
			wantErr: false,
			verify: func(t *testing.T, val any) {
				jsonBytes, ok := val.([]byte)
				require.True(t, ok)
				var result models.Appointment
				err := json.Unmarshal(jsonBytes, &result)
				require.NoError(t, err)
				assert.Equal(t, uint(1), result.ID)
				assert.Equal(t, "Meeting", result.Title)
				assert.Equal(t, "Team meeting", result.Description)
				assert.Equal(t, models.AppointmentStatusScheduled, result.Status)
				assert.Equal(t, uint(1), result.ClientID)
				assert.Equal(t, uint(2), result.OwnerID)
				assert.Equal(t, uint(3), result.PropertyID)
			},
		},
		{
			name: "success - minimal appointment",
			input: &models.Appointment{
				ID:         2,
				Title:      "Viewing",
				Status:     models.AppointmentStatusCompleted,
				ClientID:   1,
				OwnerID:    1,
				PropertyID: 1,
				CreatedAt:  now,
			},
			wantErr: false,
			verify: func(t *testing.T, val any) {
				jsonBytes, ok := val.([]byte)
				require.True(t, ok)
				var result models.Appointment
				err := json.Unmarshal(jsonBytes, &result)
				require.NoError(t, err)
				assert.Equal(t, uint(2), result.ID)
				assert.Equal(t, "Viewing", result.Title)
				assert.Equal(t, models.AppointmentStatusCompleted, result.Status)
			},
		},
		{
			name: "success - appointment with agent IDs",
			input: &models.Appointment{
				ID:         3,
				Title:      "Meeting",
				Status:     models.AppointmentStatusScheduled,
				ClientID:   1,
				OwnerID:    1,
				PropertyID: 1,
				AgentIDs:   []uint{1, 2, 3},
				CreatedAt:  now,
			},
			wantErr: false,
			verify: func(t *testing.T, val any) {
				jsonBytes, ok := val.([]byte)
				require.True(t, ok)
				var result models.Appointment
				err := json.Unmarshal(jsonBytes, &result)
				require.NoError(t, err)
				assert.Equal(t, []uint{1, 2, 3}, result.AgentIDs)
			},
		},
		{
			name:    "error - nil appointment",
			input:   nil,
			wantErr: true,
		},
		{
			name: "success - appointment with all status types",
			input: &models.Appointment{
				ID:         4,
				Title:      "Test",
				Status:     models.AppointmentStatusCancelled,
				ClientID:   1,
				OwnerID:    1,
				PropertyID: 1,
				CreatedAt:  now,
			},
			wantErr: false,
			verify: func(t *testing.T, val any) {
				jsonBytes, ok := val.([]byte)
				require.True(t, ok)
				var result models.Appointment
				err := json.Unmarshal(jsonBytes, &result)
				require.NoError(t, err)
				assert.Equal(t, models.AppointmentStatusCancelled, result.Status)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.input.Value()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				tt.verify(t, val)
			}
		})
	}
}

// TestAppointmentToResponse tests the ToResponse method for Appointment model
func TestAppointmentToResponse(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	notes := "Test notes"

	tests := []struct {
		name    string
		input   *models.Appointment
		verify  func(t *testing.T, resp *models.AppointmentDetail)
		wantErr bool
	}{
		{
			name: "success - full appointment to response",
			input: &models.Appointment{
				ID:          1,
				Title:       "Meeting",
				Description: "Team meeting",
				StartDate:   now,
				EndDate:     now.Add(1 * time.Hour),
				Status:      models.AppointmentStatusScheduled,
				Notes:       &notes,
				ClientID:    1,
				OwnerID:     2,
				PropertyID:  3,
				AgentIDs:    []uint{1, 2},
				CreatedAt:   now,
			},
			verify: func(t *testing.T, resp *models.AppointmentDetail) {
				assert.Equal(t, uint(1), resp.ID)
				assert.Equal(t, "Meeting", resp.Title)
				assert.Equal(t, "Team meeting", resp.Description)
				assert.Equal(t, now, resp.StartDate)
				assert.Equal(t, now.Add(1*time.Hour), resp.EndDate)
				assert.Equal(t, models.AppointmentStatusScheduled, resp.Status)
				assert.Equal(t, &notes, resp.Notes)
				assert.Equal(t, uint(3), resp.PropertyID)
				assert.Equal(t, uint(1), resp.ClientID)
				assert.NotNil(t, resp.OwnerID)
				assert.Equal(t, uint(2), *resp.OwnerID)
				assert.NotNil(t, resp.Agents)
			},
			wantErr: false,
		},
		{
			name: "success - minimal appointment to response",
			input: &models.Appointment{
				ID:         2,
				Title:      "Viewing",
				Status:     models.AppointmentStatusCompleted,
				ClientID:   1,
				OwnerID:    1,
				PropertyID: 1,
				CreatedAt:  now,
			},
			verify: func(t *testing.T, resp *models.AppointmentDetail) {
				assert.Equal(t, uint(2), resp.ID)
				assert.Equal(t, "Viewing", resp.Title)
				assert.Equal(t, models.AppointmentStatusCompleted, resp.Status)
				assert.Equal(t, uint(1), resp.PropertyID)
				assert.Equal(t, uint(1), resp.ClientID)
				assert.NotNil(t, resp.OwnerID)
				assert.Equal(t, uint(1), *resp.OwnerID)
				assert.NotNil(t, resp.Agents)
			},
			wantErr: false,
		},
		{
			name: "success - appointment with nil notes",
			input: &models.Appointment{
				ID:         3,
				Title:      "Meeting",
				Status:     models.AppointmentStatusScheduled,
				Notes:      nil,
				ClientID:   1,
				OwnerID:    1,
				PropertyID: 1,
				CreatedAt:  now,
			},
			verify: func(t *testing.T, resp *models.AppointmentDetail) {
				assert.Equal(t, uint(3), resp.ID)
				assert.Equal(t, "Meeting", resp.Title)
				assert.Equal(t, models.AppointmentStatusScheduled, resp.Status)
				assert.Nil(t, resp.Notes)
				assert.Equal(t, uint(1), resp.PropertyID)
				assert.Equal(t, uint(1), resp.ClientID)
				assert.NotNil(t, resp.OwnerID)
				assert.Equal(t, uint(1), *resp.OwnerID)
			},
			wantErr: false,
		},
		{
			name: "success - appointment with different status",
			input: &models.Appointment{
				ID:         4,
				Title:      "Cancelled Meeting",
				Status:     models.AppointmentStatusCancelled,
				ClientID:   1,
				OwnerID:    1,
				PropertyID: 1,
				CreatedAt:  now,
			},
			verify: func(t *testing.T, resp *models.AppointmentDetail) {
				assert.Equal(t, uint(4), resp.ID)
				assert.Equal(t, "Cancelled Meeting", resp.Title)
				assert.Equal(t, models.AppointmentStatusCancelled, resp.Status)
				assert.Equal(t, uint(1), resp.PropertyID)
				assert.Equal(t, uint(1), resp.ClientID)
				assert.NotNil(t, resp.OwnerID)
				assert.Equal(t, uint(1), *resp.OwnerID)
			},
			wantErr: false,
		},
		{
			name:    "error - nil appointment",
			input:   nil,
			verify:  nil,
			wantErr: true,
		},
		{
			name: "success - appointment with archived status",
			input: &models.Appointment{
				ID:         5,
				Title:      "Archived Meeting",
				Status:     models.AppointmentStatusArchived,
				ClientID:   1,
				OwnerID:    1,
				PropertyID: 1,
				CreatedAt:  now,
			},
			verify: func(t *testing.T, resp *models.AppointmentDetail) {
				assert.Equal(t, uint(5), resp.ID)
				assert.Equal(t, "Archived Meeting", resp.Title)
				assert.Equal(t, models.AppointmentStatusArchived, resp.Status)
				assert.Equal(t, uint(1), resp.PropertyID)
				assert.Equal(t, uint(1), resp.ClientID)
				assert.NotNil(t, resp.OwnerID)
				assert.Equal(t, uint(1), *resp.OwnerID)
			},
			wantErr: false,
		},
		{
			name: "success - appointment with no-show status",
			input: &models.Appointment{
				ID:         6,
				Title:      "No Show Meeting",
				Status:     models.AppointmentStatusNoShow,
				ClientID:   1,
				OwnerID:    1,
				PropertyID: 1,
				CreatedAt:  now,
			},
			verify: func(t *testing.T, resp *models.AppointmentDetail) {
				assert.Equal(t, uint(6), resp.ID)
				assert.Equal(t, "No Show Meeting", resp.Title)
				assert.Equal(t, models.AppointmentStatusNoShow, resp.Status)
				assert.Equal(t, uint(1), resp.PropertyID)
				assert.Equal(t, uint(1), resp.ClientID)
				assert.NotNil(t, resp.OwnerID)
				assert.Equal(t, uint(1), *resp.OwnerID)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				resp := tt.input.ToResponse()
				assert.Nil(t, resp)
			} else {
				resp := tt.input.ToResponse()
				require.NotNil(t, resp)
				if tt.verify != nil {
					tt.verify(t, resp)
				}
			}
		})
	}
}
