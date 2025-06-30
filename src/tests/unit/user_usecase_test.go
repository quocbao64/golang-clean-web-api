package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/naeemaei/golang-clean-web-api/config"
	"github.com/naeemaei/golang-clean-web-api/domain/model"
	"github.com/naeemaei/golang-clean-web-api/usecase"
)

// TestUserRepository is a test implementation of UserRepository
type TestUserRepository struct {
	shouldReturnUser  bool
	shouldReturnError bool
	user              model.User
	err               error
}

func (t *TestUserRepository) ExistsMobileNumber(ctx context.Context, mobileNumber string) (bool, error) {
	return false, nil
}

func (t *TestUserRepository) ExistsUsername(ctx context.Context, username string) (bool, error) {
	return false, nil
}

func (t *TestUserRepository) ExistsEmail(ctx context.Context, email string) (bool, error) {
	return false, nil
}

func (t *TestUserRepository) FetchUserInfo(ctx context.Context, username string, password string) (model.User, error) {
	if t.shouldReturnError {
		return model.User{}, t.err
	}
	if t.shouldReturnUser {
		return t.user, nil
	}
	return model.User{}, errors.New("user not found")
}

func (t *TestUserRepository) GetDefaultRole(ctx context.Context) (int, error) {
	return 1, nil
}

func (t *TestUserRepository) CreateUser(ctx context.Context, u model.User) (model.User, error) {
	return u, nil
}

// createTestConfig creates a minimal config for testing
func createTestConfig() *config.Config {
	return &config.Config{
		Logger: config.LoggerConfig{
			Logger: "zap", // Use zap logger for testing
			Level:  "info",
		},
		JWT: config.JWTConfig{
			AccessTokenExpireDuration:  15 * time.Minute,
			RefreshTokenExpireDuration: 24 * time.Hour,
			Secret:                     "test-secret-key",
			RefreshSecret:              "test-refresh-secret-key",
		},
		Otp: config.OtpConfig{
			ExpireTime: 120 * time.Second,
			Digits:     6,
			Limiter:    5 * time.Second,
		},
	}
}

// TestUserUsecase_LoginByUsername tests the LoginByUsername method
func TestUserUsecase_LoginByUsername(t *testing.T) {
	// Test case: Success - Valid credentials
	t.Run("Success - Valid credentials", func(t *testing.T) {
		// Create test user
		testUser := model.User{
			BaseModel:    model.BaseModel{Id: 1},
			Username:     "testuser",
			FirstName:    "John",
			LastName:     "Doe",
			Email:        "john.doe@example.com",
			MobileNumber: "1234567890",
			UserRoles: &[]model.UserRole{
				{Role: model.Role{Name: "user"}},
			},
		}

		// Create test repository that returns the user
		testRepo := &TestUserRepository{
			shouldReturnUser: true,
			user:             testUser,
		}

		// Create config with proper logger configuration
		cfg := createTestConfig()

		// Create UserUsecase
		userUsecase := usecase.NewUserUsecase(cfg, testRepo)

		// Execute test
		result, err := userUsecase.LoginByUsername(context.Background(), "testuser", "testpass")

		// Assertions
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result == nil {
			t.Error("Expected token detail, got nil")
		}
		if result.AccessToken == "" {
			t.Error("Expected access token, got empty string")
		}
		if result.RefreshToken == "" {
			t.Error("Expected refresh token, got empty string")
		}
	})

	// Test case: Failure - Invalid credentials
	t.Run("Failure - Invalid credentials", func(t *testing.T) {
		// Create test repository that returns error
		testRepo := &TestUserRepository{
			shouldReturnError: true,
			err:               errors.New("user not found"),
		}

		// Create config with proper logger configuration
		cfg := createTestConfig()

		// Create UserUsecase
		userUsecase := usecase.NewUserUsecase(cfg, testRepo)

		// Execute test
		result, err := userUsecase.LoginByUsername(context.Background(), "invaliduser", "invalidpass")

		// Assertions
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if result != nil {
			t.Error("Expected nil result, got token detail")
		}
		if err.Error() != "user not found" {
			t.Errorf("Expected 'user not found' error, got %v", err)
		}
	})
}
