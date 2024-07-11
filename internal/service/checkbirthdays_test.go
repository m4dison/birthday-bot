package service_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/m4dison/my-telegram-bot/internal/mocks"
	"github.com/m4dison/my-telegram-bot/internal/models"
	"github.com/m4dison/my-telegram-bot/internal/service"
)

func TestCheckBirthdays(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := service.NewUserService(mockRepo, &sync.Mutex{})

	// Test cases
	tests := []struct {
		name     string
		users    []models.User
		setup    func()
		expected []string
		wantErr  bool
	}{
		{
			name:  "no users",
			users: []models.User{},
			setup: func() {
				mockRepo.EXPECT().GetAllUsers(gomock.Any()).Return([]models.User{}, nil).Times(1)
			},
			expected: []string{},
			wantErr:  false,
		},
		{
			name: "no birthdays today",
			users: []models.User{
				{Name: "Alice", Birthday: time.Now().AddDate(0, 0, -1)}, // Yesterday
				{Name: "Bob", Birthday: time.Now().AddDate(0, 0, 1)},    // Tomorrow
			},
			setup: func() {
				mockRepo.EXPECT().GetAllUsers(gomock.Any()).Return([]models.User{
					{Name: "Alice", Birthday: time.Now().AddDate(0, 0, -1)},
					{Name: "Bob", Birthday: time.Now().AddDate(0, 0, 1)},
				}, nil).Times(1)
			},
			expected: []string{},
			wantErr:  false,
		},
		{
			name: "one birthday today",
			users: []models.User{
				{Name: "Alice", Birthday: time.Now()},
				{Name: "Bob", Birthday: time.Now().AddDate(0, 0, -1)}, // Yesterday
			},
			setup: func() {
				mockRepo.EXPECT().GetAllUsers(gomock.Any()).Return([]models.User{
					{Name: "Alice", Birthday: time.Now()},
					{Name: "Bob", Birthday: time.Now().AddDate(0, 0, -1)},
				}, nil).Times(1)
			},
			expected: []string{"Alice"},
			wantErr:  false,
		},
		{
			name: "multiple birthdays today",
			users: []models.User{
				{Name: "Alice", Birthday: time.Now()},
				{Name: "Bob", Birthday: time.Now()},
			},
			setup: func() {
				mockRepo.EXPECT().GetAllUsers(gomock.Any()).Return([]models.User{
					{Name: "Alice", Birthday: time.Now()},
					{Name: "Bob", Birthday: time.Now()},
				}, nil).Times(1)
			},
			expected: []string{"Alice", "Bob"},
			wantErr:  false,
		},
		{
			name:  "GetAllUsers returns error",
			users: []models.User{},
			setup: func() {
				mockRepo.EXPECT().GetAllUsers(gomock.Any()).Return(nil, errors.New("database error")).Times(1)
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			result, err := userService.CheckBirthdays(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				var resultNames []string
				for _, user := range result {
					resultNames = append(resultNames, user.Name)
				}
				assert.ElementsMatch(t, tt.expected, resultNames)
			}
		})
	}
}
