package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ========================================
// MODELS
// ========================================

type User struct {
	UserID            int       `json:"user_id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Password          string    `json:"-"` // Never expose password
	Role              string    `json:"role"`
	TanggalDibuat     time.Time `json:"tanggal_dibuat"`
	TanggalDiperbarui time.Time `json:"tanggal_diperbarui"`
	CreatedAt         time.Time `json:"created_at"` // Alias for compatibility
}

type UserProfile struct {
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type UserResponse struct {
	UserID            int       `json:"user_id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Role              string    `json:"role"`
	TanggalDibuat     time.Time `json:"tanggal_dibuat"`
	TanggalDiperbarui time.Time `json:"tanggal_diperbarui"`
}

// ========================================
// REQUEST TYPES
// ========================================

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"`
}

type UserUpdateRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
}

type UserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
}

// ========================================
// HELPER METHODS
// ========================================

// ToPublic returns user data without sensitive information
func (u *User) ToPublic() UserResponse {
	return UserResponse{
		UserID:            u.UserID,
		Username:          u.Username,
		Email:             u.Email,
		Role:              u.Role,
		TanggalDibuat:     u.TanggalDibuat,
		TanggalDiperbarui: u.TanggalDiperbarui,
	}
}

// ToProfile returns user profile data
func (u *User) ToProfile() UserProfile {
	createdAt := u.TanggalDibuat
	if !u.CreatedAt.IsZero() {
		createdAt = u.CreatedAt
	}
	return UserProfile{
		UserID:    u.UserID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: createdAt,
	}
}

// ========================================
// PASSWORD FUNCTIONS
// ========================================

// HashPassword hashes a plain text password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerifyPassword checks if the provided password matches the hashed password
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// ValidatePassword is an alias for VerifyPassword
func ValidatePassword(hashedPassword, password string) bool {
	return VerifyPassword(hashedPassword, password)
}

// ========================================
// AUTHENTICATION
// ========================================

// AuthenticateUser validates credentials and returns user if valid
func AuthenticateUser(ctx context.Context, db *sql.DB, req *LoginRequest) (*User, error) {
	user, err := GetUserByEmail(ctx, db, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !VerifyPassword(user.Password, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// ========================================
// EXISTENCE CHECKS
// ========================================

// IsEmailExists checks if email already exists in database
func IsEmailExists(ctx context.Context, db *sql.DB, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	err := db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// IsUsernameExists checks if username already exists in database
func IsUsernameExists(ctx context.Context, db *sql.DB, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	var exists bool
	err := db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CheckUsernameExists checks if username exists excluding a specific user ID
func CheckUsernameExists(db *sql.DB, username string, excludeUserID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND user_id != $2)`
	var exists bool
	err := db.QueryRow(query, username, excludeUserID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CheckEmailExists checks if email exists excluding a specific user ID
func CheckEmailExists(db *sql.DB, email string, excludeUserID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND user_id != $2)`
	var exists bool
	err := db.QueryRow(query, email, excludeUserID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// ========================================
// USER CRUD OPERATIONS
// ========================================

// CreateUser creates a new user
func CreateUser(ctx context.Context, db *sql.DB, req *UserRequest) (*User, error) {
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	query := `
        INSERT INTO users (username, email, password, role)
        VALUES ($1, $2, $3, $4)
        RETURNING user_id, username, email, role, tanggal_dibuat, tanggal_diperbarui
    `

	var user User
	err = db.QueryRowContext(ctx, query, req.Username, req.Email, hashedPassword, role).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Role,
		&user.TanggalDibuat, &user.TanggalDiperbarui,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.CreatedAt = user.TanggalDibuat
	return &user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(ctx context.Context, db *sql.DB, id int) (*User, error) {
	query := `
        SELECT user_id, username, email, password, role, tanggal_dibuat, tanggal_diperbarui
        FROM users
        WHERE user_id = $1
    `

	var user User
	err := db.QueryRowContext(ctx, query, id).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Password, &user.Role,
		&user.TanggalDibuat, &user.TanggalDiperbarui,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.CreatedAt = user.TanggalDibuat
	return &user, nil
}

// GetUserByIDSimple retrieves a user by ID without context
func GetUserByIDSimple(db *sql.DB, id int) (*User, error) {
	query := `
        SELECT user_id, username, email, password, role, tanggal_dibuat, tanggal_diperbarui
        FROM users
        WHERE user_id = $1
    `

	var user User
	err := db.QueryRow(query, id).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Password, &user.Role,
		&user.TanggalDibuat, &user.TanggalDiperbarui,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.CreatedAt = user.TanggalDibuat
	return &user, nil
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(ctx context.Context, db *sql.DB, email string) (*User, error) {
	query := `
        SELECT user_id, username, email, password, role, tanggal_dibuat, tanggal_diperbarui
        FROM users
        WHERE email = $1
    `

	var user User
	err := db.QueryRowContext(ctx, query, email).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Password, &user.Role,
		&user.TanggalDibuat, &user.TanggalDiperbarui,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.CreatedAt = user.TanggalDibuat
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(ctx context.Context, db *sql.DB, username string) (*User, error) {
	query := `
        SELECT user_id, username, email, password, role, tanggal_dibuat, tanggal_diperbarui
        FROM users
        WHERE username = $1
    `

	var user User
	err := db.QueryRowContext(ctx, query, username).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Password, &user.Role,
		&user.TanggalDibuat, &user.TanggalDiperbarui,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.CreatedAt = user.TanggalDibuat
	return &user, nil
}

// GetAllUsers retrieves all users
func GetAllUsers(db *sql.DB) ([]UserResponse, error) {
	query := `
        SELECT user_id, username, email, role, tanggal_dibuat, tanggal_diperbarui
        FROM users
        ORDER BY tanggal_dibuat DESC
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserResponse
	for rows.Next() {
		var user UserResponse
		err := rows.Scan(
			&user.UserID, &user.Username, &user.Email, &user.Role,
			&user.TanggalDibuat, &user.TanggalDiperbarui,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if users == nil {
		users = []UserResponse{}
	}

	return users, nil
}

// UpdateUser updates a user's information
func UpdateUser(ctx context.Context, db *sql.DB, id int, req *UserUpdateRequest) (*User, error) {
	// Get existing user
	existing, err := GetUserByID(ctx, db, id)
	if err != nil {
		return nil, err
	}

	// Build update query dynamically
	username := existing.Username
	email := existing.Email
	role := existing.Role

	if req.Username != "" {
		username = req.Username
	}
	if req.Email != "" {
		email = req.Email
	}
	if req.Role != "" {
		role = req.Role
	}

	query := `
        UPDATE users
        SET username = $1, email = $2, role = $3
        WHERE user_id = $4
        RETURNING user_id, username, email, role, tanggal_dibuat, tanggal_diperbarui
    `

	var user User
	err = db.QueryRowContext(ctx, query, username, email, role, id).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Role,
		&user.TanggalDibuat, &user.TanggalDiperbarui,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Update password if provided
	if req.Password != "" {
		if err := UpdateUserPassword(db, id, req.Password); err != nil {
			return nil, err
		}
	}

	user.CreatedAt = user.TanggalDibuat
	return &user, nil
}

// UpdateUserBasic updates username and email only
func UpdateUserBasic(db *sql.DB, id int, username, email string) (*User, error) {
	query := `
        UPDATE users
        SET username = $1, email = $2
        WHERE user_id = $3
        RETURNING user_id, username, email, password, role, tanggal_dibuat, tanggal_diperbarui
    `

	var user User
	err := db.QueryRow(query, username, email, id).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Password, &user.Role,
		&user.TanggalDibuat, &user.TanggalDiperbarui,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	user.CreatedAt = user.TanggalDibuat
	return &user, nil
}

// UpdateUserPassword updates only the password
func UpdateUserPassword(db *sql.DB, userID int, newPassword string) error {
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `UPDATE users SET password = $1 WHERE user_id = $2`
	result, err := db.Exec(query, hashedPassword, userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// UpdateUserRole updates only the role
func UpdateUserRole(db *sql.DB, userID int, role string) error {
	query := `UPDATE users SET role = $1 WHERE user_id = $2`
	result, err := db.Exec(query, role, userID)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteUser deletes a user by ID
func DeleteUser(db *sql.DB, id int) error {
	query := `DELETE FROM users WHERE user_id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
