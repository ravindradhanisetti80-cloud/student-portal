package service

import (
	"context"
	// Needed for Login event timestamp
	"student-portal/internal/config"
	appErrors "student-portal/internal/errors"
	"student-portal/internal/logger" // Imported for structured logging
	"student-portal/internal/models"
	"student-portal/internal/repository"
	"student-portal/internal/utils"

	"go.uber.org/zap" // Imported for structured logging fields
)

// UserService defines the methods for user-related business operations.
type UserService interface {
	RegisterUser(ctx context.Context, req *models.RegisterRequest) (*models.UserResponse, error)
	LoginUser(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	GetUserByID(ctx context.Context, id int64) (*models.UserResponse, error)
	UpdateProfile(ctx context.Context, id int64, req *models.UpdateProfileRequest) (*models.UserResponse, error)
	UpdateUser(ctx context.Context, id int64, req *models.UpdateUserRequest) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, limit, offset int) ([]models.UserResponse, int64, error)
}

type userService struct {
	repo  repository.UserRepository
	cfg   *config.Config
	kafka *utils.KafkaProducer
}

// NewUserService creates a new UserService instance.
func NewUserService(repo repository.UserRepository, cfg *config.Config, kafka *utils.KafkaProducer) UserService {
	return &userService{repo: repo, cfg: cfg, kafka: kafka}
}

// publishAsync handles the non-blocking publication and logs any failure.
func (s *userService) publishAsync(eventFunc func(context.Context) error, eventName string, userID int64) {
	go func() {
		// Use a background context as the original request context may expire.
		if err := eventFunc(context.Background()); err != nil {
			logger.Logger.Error("Failed to publish Kafka event",
				zap.Error(err),
				zap.String("event", eventName),
				zap.Int64("user_id", userID),
			)
		} else {
			logger.Logger.Info("Successfully published Kafka event", zap.String("event", eventName), zap.Int64("user_id", userID))
		}
	}()
}

func (s *userService) RegisterUser(ctx context.Context, req *models.RegisterRequest) (*models.UserResponse, error) {
	// 1. Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	// 2. Create the User model
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	// 3. Save to repository
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err // Returns ErrEmailExists if unique constraint violated
	}

	// 4. Publish register event to Kafka asynchronously and with error logging
	s.publishAsync(
		func(ctx context.Context) error {
			return s.kafka.PublishRegisterEvent(ctx, user.ID, user.Email, user.Name, string(user.Role))
		},
		"user_registered",
		user.ID,
	)

	return user.ToResponsePtr(), nil
}

func (s *userService) LoginUser(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// 1. Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == appErrors.ErrNotFound {
			return nil, appErrors.ErrInvalidCredentials
		}
		return nil, err
	}

	// 2. Check password hash
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, appErrors.ErrInvalidCredentials
	}

	// 3. Generate JWT token
	token, err := utils.GenerateToken(s.cfg, user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	// 4. KAFKA: Publish Login Event
	s.publishAsync(
		func(ctx context.Context) error {
			return s.kafka.PublishLoginEvent(ctx, user.ID, user.Email, user.Name, user.Role)
		},
		"user_logged_in",
		user.ID,
	)

	return &models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, id int64, req *models.UpdateProfileRequest) (*models.UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}

	// The repository update method will handle the actual update and check for email conflicts.
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	// Publish update event to Kafka asynchronously
	s.publishAsync(
		func(ctx context.Context) error {
			return s.kafka.PublishUpdateEvent(ctx, user.ID, user.Email, user.Name, string(user.Role))
		},
		"user_updated",
		user.ID,
	)

	return user.ToResponsePtr(), nil
}

func (s *userService) UpdateUser(ctx context.Context, id int64, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Role != nil {
		user.Role = *req.Role
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	// Since this is a core UpdateUser (potentially by admin), we should also publish an event.
	s.publishAsync(
		func(ctx context.Context) error {
			return s.kafka.PublishUpdateEvent(ctx, user.ID, user.Email, user.Name, string(user.Role))
		},
		"user_updated_admin",
		user.ID,
	)

	return user.ToResponsePtr(), nil
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*models.UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user.ToResponsePtr(), nil
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	// You might want to publish a Delete event here as well!
	return s.repo.DeleteUser(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context, limit, offset int) ([]models.UserResponse, int64, error) {
	users, totalCount, err := s.repo.ListUsers(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	return userResponses, totalCount, nil
}
