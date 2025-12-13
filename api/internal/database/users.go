package database

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID            int       `json:"user_id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Role              string    `json:"role"`
	Password          string    `json:"-"` // Never expose password in JSON
	TanggalDibuat     time.Time `json:"tanggal_dibuat"`
	TanggalDiperbarui time.Time `json:"tanggal_diperbarui"`
}

type UserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
}

type UserResponse struct {
	UserID            int       `json:"user_id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Role              string    `json:"role"`
	TanggalDibuat     time.Time `json:"tanggal_dibuat"`
	TanggalDiperbarui time.Time `json:"tanggal_diperbarui"`
}

// ToResponse converts User to UserResponse (without password)
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		UserID:            u.UserID,
		Username:          u.Username,
		Email:             u.Email,
		Role:              u.Role,
		TanggalDibuat:     u.TanggalDibuat,
		TanggalDiperbarui: u.TanggalDiperbarui,
	}
}

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

// GetAllUsers retrieves all users
func GetAllUsers(db *sql.DB) ([]User, error) {
	query := `
        SELECT user_id, username, email, role, password, tanggal_dibuat, tanggal_diperbarui
        FROM users
        ORDER BY tanggal_dibuat DESC
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.UserID, &u.Username, &u.Email, &u.Role, &u.Password, &u.TanggalDibuat, &u.TanggalDiperbarui)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// GetUserByID retrieves a single user by ID
func GetUserByID(db *sql.DB, id int) (*User, error) {
	query := `
        SELECT user_id, username, email, role, password, tanggal_dibuat, tanggal_diperbarui
        FROM users
        WHERE user_id = $1
    `

	var u User
	err := db.QueryRow(query, id).Scan(&u.UserID, &u.Username, &u.Email, &u.Role, &u.Password, &u.TanggalDibuat, &u.TanggalDiperbarui)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	query := `
        SELECT user_id, username, email, role, password, tanggal_dibuat, tanggal_diperbarui
        FROM users
        WHERE email = $1
    `

	var u User
	err := db.QueryRow(query, email).Scan(&u.UserID, &u.Username, &u.Email, &u.Role, &u.Password, &u.TanggalDibuat, &u.TanggalDiperbarui)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	query := `
        SELECT user_id, username, email, role, password, tanggal_dibuat, tanggal_diperbarui
        FROM users
        WHERE username = $1
    `

	var u User
	err := db.QueryRow(query, username).Scan(&u.UserID, &u.Username, &u.Email, &u.Role, &u.Password, &u.TanggalDibuat, &u.TanggalDiperbarui)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// CreateUser creates a new user with hashed password
func CreateUser(db *sql.DB, input UserInput) (*User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Set default role if not provided
	role := input.Role
	if role == "" {
		role = "user"
	}

	query := `
        INSERT INTO users (username, email, password, role)
        VALUES ($1, $2, $3, $4)
        RETURNING user_id, username, email, role, password, tanggal_dibuat, tanggal_diperbarui
    `

	var u User
	err = db.QueryRow(query, input.Username, input.Email, string(hashedPassword), role).Scan(
		&u.UserID, &u.Username, &u.Email, &u.Role, &u.Password, &u.TanggalDibuat, &u.TanggalDiperbarui,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// UpdateUser updates an existing user
func UpdateUser(db *sql.DB, id int, input UserInput) (*User, error) {
	// Start building the query
	query := `UPDATE users SET username = $1, email = $2`
	args := []interface{}{input.Username, input.Email}
	argCount := 2

	// If password is provided, hash and update it
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		argCount++
		query += ", password = $" + string(rune('0'+argCount))
		args = append(args, string(hashedPassword))
	}

	// If role is provided, update it
	if input.Role != "" {
		argCount++
		query += ", role = $" + string(rune('0'+argCount))
		args = append(args, input.Role)
	}

	argCount++
	query += " WHERE user_id = $" + string(rune('0'+argCount))
	args = append(args, id)

	query += " RETURNING user_id, username, email, role, password, tanggal_dibuat, tanggal_diperbarui"

	var u User
	err := db.QueryRow(query, args...).Scan(
		&u.UserID, &u.Username, &u.Email, &u.Role, &u.Password, &u.TanggalDibuat, &u.TanggalDiperbarui,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// UpdateUserBasic updates username and email only (used by handleUpdateCurrentUser)
func UpdateUserBasic(db *sql.DB, id int, username, email string) (*User, error) {
	query := `
        UPDATE users 
        SET username = $1, email = $2
        WHERE user_id = $3
        RETURNING user_id, username, email, role, password, tanggal_dibuat, tanggal_diperbarui
    `

	var u User
	err := db.QueryRow(query, username, email, id).Scan(
		&u.UserID, &u.Username, &u.Email, &u.Role, &u.Password, &u.TanggalDibuat, &u.TanggalDiperbarui,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// DeleteUser deletes a user by ID
func DeleteUser(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM users WHERE user_id = $1", id)
	if err != nil {
		return err
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

// VerifyPassword checks if the provided password matches the hashed password
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// UpdateUserPassword updates only the password for a user
func UpdateUserPassword(db *sql.DB, userID int, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `UPDATE users SET password = $1 WHERE user_id = $2`
	result, err := db.Exec(query, string(hashedPassword), userID)
	if err != nil {
		return err
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

// UpdateUserRole updates only the role for a user
func UpdateUserRole(db *sql.DB, userID int, role string) error {
	query := `UPDATE users SET role = $1 WHERE user_id = $2`
	result, err := db.Exec(query, role, userID)
	if err != nil {
		return err
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
