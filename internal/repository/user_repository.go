// internal/repository/user_repository.go
package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	appErrors "student-portal/internal/errors"
	"student-portal/internal/models"
)

// UserRepository defines the methods for interacting with the users data store.
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, limit, offset int) ([]models.User, int64, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (name, email, password, role) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query, user.Name, user.Email, user.Password, user.Role).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 is unique violation
			return appErrors.ErrEmailExists
		}
		return appErrors.ErrInternalServerError
	}
	return nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, password, role, created_at, updated_at 
		FROM users 
		WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, appErrors.ErrNotFound
	}
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}
	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, password, role, created_at, updated_at 
		FROM users 
		WHERE email = $1
	`
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, appErrors.ErrNotFound
	}
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}
	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET name = $2, email = $3, role = $4, updated_at = NOW() 
		WHERE id = $1 
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query, user.ID, user.Name, user.Email, user.Role).Scan(&user.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return appErrors.ErrNotFound
	}
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 is unique violation
			return appErrors.ErrEmailExists
		}
		return appErrors.ErrInternalServerError
	}
	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int64) error {
	cmdTag, err := r.db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return appErrors.ErrInternalServerError
	}
	if cmdTag.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}

func (r *userRepository) ListUsers(ctx context.Context, limit, offset int) ([]models.User, int64, error) {
	// Query to count total users
	var totalCount int64
	countQuery := "SELECT COUNT(*) FROM users"
	err := r.db.QueryRow(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, appErrors.ErrInternalServerError
	}

	// Query to get paginated users
	usersQuery := `
		SELECT id, name, email, role, created_at, updated_at 
		FROM users 
		ORDER BY id 
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, usersQuery, limit, offset)
	if err != nil {
		return nil, 0, appErrors.ErrInternalServerError
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		user := models.User{}
		// Note: We don't select 'password' here as it's not needed for listing
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, appErrors.ErrInternalServerError
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, 0, appErrors.ErrInternalServerError
	}

	return users, totalCount, nil
}
