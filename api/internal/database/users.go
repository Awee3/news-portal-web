package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Hidden from JSON response
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"` // Optional, default to "user"
}

type UserUpdateRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role,omitempty"`
}

type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserProfile struct {
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	// Additional profile fields can be added here
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(ctx context.Context, db *sql.DB, req *UserRequest) (*User, error) {
	// Set default role if not provided
	if req.Role == "" {
		req.Role = "user"
	}

	// Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
        INSERT INTO users (username, email, password, role, created_at)
        VALUES ($1, $2, $3, $4, NOW())
        RETURNING user_id, username, email, role, created_at
    `

	var user User
	err = db.QueryRowContext(ctx, query, req.Username, req.Email, hashedPassword, req.Role).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func CreateUserTx(ctx context.Context, tx *sql.Tx, req *UserRequest) (*User, error) {
	// Set default role if not provided
	if req.Role == "" {
		req.Role = "user"
	}

	// Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
        INSERT INTO users (username, email, password, role, created_at)
        VALUES ($1, $2, $3, $4, NOW())
        RETURNING user_id, username, email, role, created_at
    `

	var user User
	err = tx.QueryRowContext(ctx, query, req.Username, req.Email, hashedPassword, req.Role).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func GetUserByID(ctx context.Context, db *sql.DB, userID int) (*User, error) {
	query := `
        SELECT user_id, username, email, password, role, created_at
        FROM users
        WHERE user_id = $1
    `

	var user User
	err := db.QueryRowContext(ctx, query, userID).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(ctx context.Context, db *sql.DB, email string) (*User, error) {
	query := `
        SELECT user_id, username, email, password, role, created_at
        FROM users
        WHERE email = $1
    `

	var user User
	err := db.QueryRowContext(ctx, query, email).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(ctx context.Context, db *sql.DB, username string) (*User, error) {
	query := `
        SELECT user_id, username, email, password, role, created_at
        FROM users
        WHERE username = $1
    `

	var user User
	err := db.QueryRowContext(ctx, query, username).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func ListUsers(ctx context.Context, db *sql.DB, limit, offset int) ([]UserProfile, int, error) {
	// Count query
	countQuery := `SELECT COUNT(*) FROM users`

	// Select query
	selectQuery := `
        SELECT user_id, username, email, role, created_at
        FROM users
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

	// Get total count
	var totalCount int
	err := db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Execute main query
	rows, err := db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []UserProfile
	for rows.Next() {
		var user UserProfile
		err := rows.Scan(
			&user.UserID, &user.Username, &user.Email, &user.Role, &user.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, totalCount, nil
}

func ListUsersByRole(ctx context.Context, db *sql.DB, role string) ([]UserProfile, error) {
	query := `
        SELECT user_id, username, email, role, created_at
        FROM users
        WHERE role = $1
        ORDER BY username ASC
    `

	rows, err := db.QueryContext(ctx, query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserProfile
	for rows.Next() {
		var user UserProfile
		err := rows.Scan(
			&user.UserID, &user.Username, &user.Email, &user.Role, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func UpdateUser(ctx context.Context, db *sql.DB, userID int, req *UserUpdateRequest) (*User, error) {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Username != "" {
		setParts = append(setParts, fmt.Sprintf("username = $%d", argIndex))
		args = append(args, req.Username)
		argIndex++
	}
	if req.Email != "" {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, req.Email)
		argIndex++
	}
	if req.Role != "" {
		setParts = append(setParts, fmt.Sprintf("role = $%d", argIndex))
		args = append(args, req.Role)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil, errors.New("no fields to update")
	}

	query := fmt.Sprintf(`
        UPDATE users 
        SET %s
        WHERE user_id = $%d
        RETURNING user_id, username, email, role, created_at
    `, fmt.Sprintf("%s", setParts[0]), argIndex)

	// Add remaining setParts
	for i := 1; i < len(setParts); i++ {
		query = fmt.Sprintf(`
            UPDATE users 
            SET %s
            WHERE user_id = $%d
            RETURNING user_id, username, email, role, created_at
        `, fmt.Sprintf("%s, %s", setParts[0], setParts[i]), argIndex)
	}

	// Rebuild query properly
	setClause := ""
	for i, part := range setParts {
		if i > 0 {
			setClause += ", "
		}
		setClause += part
	}

	query = fmt.Sprintf(`
        UPDATE users 
        SET %s
        WHERE user_id = $%d
        RETURNING user_id, username, email, role, created_at
    `, setClause, argIndex)

	args = append(args, userID)

	var user User
	err := db.QueryRowContext(ctx, query, args...).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

func UpdateUserTx(ctx context.Context, tx *sql.Tx, userID int, req *UserUpdateRequest) (*User, error) {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Username != "" {
		setParts = append(setParts, fmt.Sprintf("username = $%d", argIndex))
		args = append(args, req.Username)
		argIndex++
	}
	if req.Email != "" {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, req.Email)
		argIndex++
	}
	if req.Role != "" {
		setParts = append(setParts, fmt.Sprintf("role = $%d", argIndex))
		args = append(args, req.Role)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil, errors.New("no fields to update")
	}

	setClause := ""
	for i, part := range setParts {
		if i > 0 {
			setClause += ", "
		}
		setClause += part
	}

	query := fmt.Sprintf(`
        UPDATE users 
        SET %s
        WHERE user_id = $%d
        RETURNING user_id, username, email, role, created_at
    `, setClause, argIndex)

	args = append(args, userID)

	var user User
	err := tx.QueryRowContext(ctx, query, args...).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

func ChangePassword(ctx context.Context, db *sql.DB, userID int, req *PasswordChangeRequest) error {
	// Get current password hash
	var currentHash string
	err := db.QueryRowContext(ctx, "SELECT password FROM users WHERE user_id = $1", userID).Scan(&currentHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	// Verify current password
	if !checkPasswordHash(req.CurrentPassword, currentHash) {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	newHash, err := hashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	_, err = db.ExecContext(ctx, "UPDATE users SET password = $1 WHERE user_id = $2", newHash, userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func DeleteUser(ctx context.Context, db *sql.DB, userID int) error {
	res, err := db.ExecContext(ctx, "DELETE FROM users WHERE user_id = $1", userID)
	if err != nil {
		return fmt.Errorf("error deleting user ID %d: %w", userID, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected for user ID %d delete: %w", userID, err)
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteUserTx(ctx context.Context, tx *sql.Tx, userID int) error {
	res, err := tx.ExecContext(ctx, "DELETE FROM users WHERE user_id = $1", userID)
	if err != nil {
		return fmt.Errorf("error executing delete for user ID %d in tx: %w", userID, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected for user ID %d delete in tx: %w", userID, err)
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func AuthenticateUser(ctx context.Context, db *sql.DB, req *LoginRequest) (*User, error) {
	user, err := GetUserByEmail(ctx, db, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !checkPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func IsEmailExists(ctx context.Context, db *sql.DB, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	err := db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking email existence: %w", err)
	}
	return exists, nil
}

func IsUsernameExists(ctx context.Context, db *sql.DB, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	var exists bool
	err := db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking username existence: %w", err)
	}
	return exists, nil
}

func GetUserStats(ctx context.Context, db *sql.DB) (map[string]int, error) {
	query := `
        SELECT 
            COUNT(CASE WHEN role = 'admin' THEN 1 END) as admin_count,
            COUNT(CASE WHEN role = 'editor' THEN 1 END) as editor_count,
            COUNT(CASE WHEN role = 'user' THEN 1 END) as user_count,
            COUNT(*) as total_count
        FROM users
    `

	var adminCount, editorCount, userCount, totalCount int
	err := db.QueryRowContext(ctx, query).Scan(&adminCount, &editorCount, &userCount, &totalCount)
	if err != nil {
		return nil, err
	}

	stats := map[string]int{
		"admin":  adminCount,
		"editor": editorCount,
		"user":   userCount,
		"total":  totalCount,
	}

	return stats, nil
}
