package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Comment struct {
	KomentarID        int       `json:"komentar_id"`
	Konten            string    `json:"konten"`
	NamaPengguna      *string   `json:"nama_pengguna,omitempty"`
	Status            string    `json:"status"`
	UserID            *int      `json:"user_id,omitempty"`
	ArtikelID         int       `json:"artikel_id"`
	TanggalDibuat     time.Time `json:"tanggal_dibuat"`
	TanggalDiperbarui time.Time `json:"tanggal_diperbarui"`
}

type CommentWithAuthor struct {
	Comment
	AuthorUsername *string `json:"author_username"`
	AuthorEmail    *string `json:"author_email"`
	ArticleTitle   string  `json:"article_title"`
}

type CommentRequest struct {
	ArtikelID int    `json:"artikel_id"`
	Konten    string `json:"konten"`
}

type CommentInput struct {
	Konten       string `json:"konten"`
	NamaPengguna string `json:"nama_pengguna,omitempty"`
}

// ========================================
// SIMPLE FUNCTIONS (used by server handlers)
// ========================================

// CreateCommentSimple creates a new comment (simpler signature for handlers)
func CreateCommentSimple(db *sql.DB, comment *Comment) (*Comment, error) {
	query := `
        INSERT INTO comments (konten, nama_pengguna, status, user_id, artikel_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING komentar_id, tanggal_dibuat, tanggal_diperbarui
    `

	err := db.QueryRow(
		query,
		comment.Konten,
		comment.NamaPengguna,
		comment.Status,
		comment.UserID,
		comment.ArtikelID,
	).Scan(&comment.KomentarID, &comment.TanggalDibuat, &comment.TanggalDiperbarui)

	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}

// GetCommentByIDSimple retrieves a single comment by ID (simpler signature)
func GetCommentByIDSimple(db *sql.DB, commentID int) (*Comment, error) {
	query := `
        SELECT komentar_id, konten, nama_pengguna, status, user_id, artikel_id, 
               tanggal_dibuat, tanggal_diperbarui
        FROM comments
        WHERE komentar_id = $1
    `

	var c Comment
	err := db.QueryRow(query, commentID).Scan(
		&c.KomentarID, &c.Konten, &c.NamaPengguna, &c.Status,
		&c.UserID, &c.ArtikelID, &c.TanggalDibuat, &c.TanggalDiperbarui,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}

	return &c, nil
}

// GetCommentsByUserID retrieves all comments by a specific user
func GetCommentsByUserID(db *sql.DB, userID int) ([]Comment, error) {
	query := `
        SELECT komentar_id, konten, nama_pengguna, status, user_id, artikel_id, 
               tanggal_dibuat, tanggal_diperbarui
        FROM comments
        WHERE user_id = $1
        ORDER BY tanggal_dibuat DESC
    `

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		err := rows.Scan(
			&c.KomentarID, &c.Konten, &c.NamaPengguna, &c.Status,
			&c.UserID, &c.ArtikelID, &c.TanggalDibuat, &c.TanggalDiperbarui,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if comments == nil {
		comments = []Comment{}
	}

	return comments, nil
}

// UpdateCommentSimple updates a comment's content and status (simpler signature)
func UpdateCommentSimple(db *sql.DB, commentID int, konten string, status string) (*Comment, error) {
	query := `
        UPDATE comments
        SET konten = $1, status = $2
        WHERE komentar_id = $3
        RETURNING komentar_id, konten, nama_pengguna, status, user_id, artikel_id, 
                  tanggal_dibuat, tanggal_diperbarui
    `

	var c Comment
	err := db.QueryRow(query, konten, status, commentID).Scan(
		&c.KomentarID, &c.Konten, &c.NamaPengguna, &c.Status,
		&c.UserID, &c.ArtikelID, &c.TanggalDibuat, &c.TanggalDiperbarui,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return &c, nil
}

// DeleteCommentSimple deletes a comment by ID (simpler signature, no ownership check)
func DeleteCommentSimple(db *sql.DB, commentID int) error {
	query := `DELETE FROM comments WHERE komentar_id = $1`
	result, err := db.Exec(query, commentID)
	if err != nil {
		return fmt.Errorf("error deleting comment: %w", err)
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

// GetAllComments retrieves all comments with optional status filter
func GetAllComments(db *sql.DB, status string, limit int, offset int) ([]Comment, error) {
	query := `
        SELECT komentar_id, konten, nama_pengguna, status, user_id, artikel_id, 
               tanggal_dibuat, tanggal_diperbarui
        FROM comments
        WHERE 1=1
    `
	args := []interface{}{}
	argIndex := 1

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	query += " ORDER BY tanggal_dibuat DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
		argIndex++
	}

	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, offset)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		err := rows.Scan(
			&c.KomentarID, &c.Konten, &c.NamaPengguna, &c.Status,
			&c.UserID, &c.ArtikelID, &c.TanggalDibuat, &c.TanggalDiperbarui,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if comments == nil {
		comments = []Comment{}
	}

	return comments, nil
}

// ========================================
// CONTEXT-AWARE FUNCTIONS (for transactions)
// ========================================

// CreateComment creates a new comment with context
func CreateComment(ctx context.Context, db *sql.DB, userID *int, req *CommentRequest) (*Comment, error) {
	// Check if article exists and is published
	var articleExists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM articles WHERE artikel_id = $1 AND status = 'published')`
	err := db.QueryRowContext(ctx, checkQuery, req.ArtikelID).Scan(&articleExists)
	if err != nil {
		return nil, fmt.Errorf("failed to check article existence: %w", err)
	}
	if !articleExists {
		return nil, errors.New("article not found or not published")
	}

	query := `
        INSERT INTO comments (artikel_id, user_id, konten, status)
        VALUES ($1, $2, $3, 'pending')
        RETURNING komentar_id, artikel_id, user_id, konten, status, tanggal_dibuat, tanggal_diperbarui
    `

	var comment Comment
	err = db.QueryRowContext(ctx, query, req.ArtikelID, userID, req.Konten).Scan(
		&comment.KomentarID, &comment.ArtikelID, &comment.UserID,
		&comment.Konten, &comment.Status, &comment.TanggalDibuat, &comment.TanggalDiperbarui)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return &comment, nil
}

// CreateCommentTx creates a new comment within a transaction
func CreateCommentTx(ctx context.Context, tx *sql.Tx, userID *int, req *CommentRequest) (*Comment, error) {
	// Check if article exists and is published
	var articleExists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM articles WHERE artikel_id = $1 AND status = 'published')`
	err := tx.QueryRowContext(ctx, checkQuery, req.ArtikelID).Scan(&articleExists)
	if err != nil {
		return nil, fmt.Errorf("failed to check article existence: %w", err)
	}
	if !articleExists {
		return nil, errors.New("article not found or not published")
	}

	query := `
        INSERT INTO comments (artikel_id, user_id, konten, status)
        VALUES ($1, $2, $3, 'pending')
        RETURNING komentar_id, artikel_id, user_id, konten, status, tanggal_dibuat, tanggal_diperbarui
    `

	var comment Comment
	err = tx.QueryRowContext(ctx, query, req.ArtikelID, userID, req.Konten).Scan(
		&comment.KomentarID, &comment.ArtikelID, &comment.UserID,
		&comment.Konten, &comment.Status, &comment.TanggalDibuat, &comment.TanggalDiperbarui)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return &comment, nil
}

// GetCommentByID retrieves a single comment by ID with context
func GetCommentByID(ctx context.Context, db *sql.DB, commentID int) (*CommentWithAuthor, error) {
	query := `
        SELECT c.komentar_id, c.artikel_id, c.user_id, c.konten, c.status,
               c.tanggal_dibuat, c.tanggal_diperbarui,
               u.username, u.email, a.judul
        FROM comments c
        LEFT JOIN users u ON c.user_id = u.user_id
        JOIN articles a ON c.artikel_id = a.artikel_id
        WHERE c.komentar_id = $1
    `

	var comment CommentWithAuthor
	err := db.QueryRowContext(ctx, query, commentID).Scan(
		&comment.KomentarID, &comment.ArtikelID, &comment.UserID,
		&comment.Konten, &comment.Status, &comment.TanggalDibuat, &comment.TanggalDiperbarui,
		&comment.AuthorUsername, &comment.AuthorEmail, &comment.ArticleTitle)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}

	return &comment, nil
}

// GetCommentsByArticleID retrieves all comments for an article
func GetCommentsByArticleID(db *sql.DB, artikelID int, status string) ([]Comment, error) {
	query := `
        SELECT komentar_id, konten, nama_pengguna, status, user_id, artikel_id, 
               tanggal_dibuat, tanggal_diperbarui
        FROM comments
        WHERE artikel_id = $1
    `
	args := []interface{}{artikelID}

	if status != "" {
		query += " AND status = $2"
		args = append(args, status)
	}

	query += " ORDER BY tanggal_dibuat DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		err := rows.Scan(
			&c.KomentarID, &c.Konten, &c.NamaPengguna, &c.Status,
			&c.UserID, &c.ArtikelID, &c.TanggalDibuat, &c.TanggalDiperbarui,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if comments == nil {
		comments = []Comment{}
	}

	return comments, nil
}

// GetApprovedCommentsByArticleID retrieves only approved comments for public view
func GetApprovedCommentsByArticleID(db *sql.DB, artikelID int) ([]Comment, error) {
	return GetCommentsByArticleID(db, artikelID, "approved")
}

// ListCommentsByArticle retrieves comments for an article with pagination
func ListCommentsByArticle(ctx context.Context, db *sql.DB, articleID int, limit, offset int) ([]CommentWithAuthor, error) {
	query := `
        SELECT c.komentar_id, c.artikel_id, c.user_id, c.konten, c.status,
               c.tanggal_dibuat, c.tanggal_diperbarui,
               u.username, u.email, a.judul
        FROM comments c
        LEFT JOIN users u ON c.user_id = u.user_id
        JOIN articles a ON c.artikel_id = a.artikel_id
        WHERE c.artikel_id = $1
        ORDER BY c.tanggal_dibuat ASC
        LIMIT $2 OFFSET $3
    `

	rows, err := db.QueryContext(ctx, query, articleID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []CommentWithAuthor
	for rows.Next() {
		var comment CommentWithAuthor
		err := rows.Scan(
			&comment.KomentarID, &comment.ArtikelID, &comment.UserID,
			&comment.Konten, &comment.Status, &comment.TanggalDibuat, &comment.TanggalDiperbarui,
			&comment.AuthorUsername, &comment.AuthorEmail, &comment.ArticleTitle)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if comments == nil {
		comments = []CommentWithAuthor{}
	}

	return comments, nil
}

// ListAllComments retrieves all comments with optional status filter (for admin)
func ListAllComments(db *sql.DB, status string, limit int, offset int) ([]Comment, error) {
	return GetAllComments(db, status, limit, offset)
}

// UpdateComment updates the content of an existing comment with ownership check
func UpdateComment(ctx context.Context, db *sql.DB, commentID int, userID *int, konten string) (*Comment, error) {
	// Check if comment exists and user owns it
	var existingUserID *int
	err := db.QueryRowContext(ctx, "SELECT user_id FROM comments WHERE komentar_id = $1", commentID).Scan(&existingUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}

	// Check ownership (only comment owner can edit)
	if existingUserID == nil && userID != nil {
		return nil, errors.New("unauthorized: anonymous comment cannot be edited")
	}
	if existingUserID != nil && userID != nil && *existingUserID != *userID {
		return nil, errors.New("unauthorized: you can only edit your own comments")
	}
	if existingUserID != nil && userID == nil {
		return nil, errors.New("unauthorized: login required to edit comment")
	}

	query := `
        UPDATE comments
        SET konten = $1, status = 'pending'
        WHERE komentar_id = $2
        RETURNING komentar_id, konten, nama_pengguna, status, user_id, artikel_id, 
                  tanggal_dibuat, tanggal_diperbarui
    `

	var comment Comment
	err = db.QueryRowContext(ctx, query, konten, commentID).Scan(
		&comment.KomentarID, &comment.Konten, &comment.NamaPengguna, &comment.Status,
		&comment.UserID, &comment.ArtikelID, &comment.TanggalDibuat, &comment.TanggalDiperbarui)
	if err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return &comment, nil
}

// UpdateCommentStatus updates the status of a comment (for moderation)
func UpdateCommentStatus(db *sql.DB, id int, status string) (*Comment, error) {
	query := `
        UPDATE comments
        SET status = $1
        WHERE komentar_id = $2
        RETURNING komentar_id, konten, nama_pengguna, status, user_id, artikel_id, 
                  tanggal_dibuat, tanggal_diperbarui
    `

	var c Comment
	err := db.QueryRow(query, status, id).Scan(
		&c.KomentarID, &c.Konten, &c.NamaPengguna, &c.Status,
		&c.UserID, &c.ArtikelID, &c.TanggalDibuat, &c.TanggalDiperbarui,
	)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// DeleteComment deletes a comment by ID with ownership check
func DeleteComment(ctx context.Context, db *sql.DB, commentID int, userID *int) error {
	// Check if comment exists and user owns it
	var existingUserID *int
	err := db.QueryRowContext(ctx, "SELECT user_id FROM comments WHERE komentar_id = $1", commentID).Scan(&existingUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return err
	}

	// Check ownership (only comment owner can delete)
	if existingUserID == nil && userID != nil {
		return errors.New("unauthorized: anonymous comment cannot be deleted")
	}
	if existingUserID != nil && userID != nil && *existingUserID != *userID {
		return errors.New("unauthorized: you can only delete your own comments")
	}
	if existingUserID != nil && userID == nil {
		return errors.New("unauthorized: login required to delete comment")
	}

	res, err := db.ExecContext(ctx, "DELETE FROM comments WHERE komentar_id = $1", commentID)
	if err != nil {
		return fmt.Errorf("error deleting comment ID %d: %w", commentID, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected for comment ID %d delete: %w", commentID, err)
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteCommentTx deletes a comment by ID within a transaction
func DeleteCommentTx(ctx context.Context, tx *sql.Tx, commentID int, userID *int) error {
	// Check if comment exists and user owns it
	var existingUserID *int
	err := tx.QueryRowContext(ctx, "SELECT user_id FROM comments WHERE komentar_id = $1", commentID).Scan(&existingUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return err
	}

	// Check ownership
	if existingUserID == nil && userID != nil {
		return errors.New("unauthorized: anonymous comment cannot be deleted")
	}
	if existingUserID != nil && userID != nil && *existingUserID != *userID {
		return errors.New("unauthorized: you can only delete your own comments")
	}
	if existingUserID != nil && userID == nil {
		return errors.New("unauthorized: login required to delete comment")
	}

	res, err := tx.ExecContext(ctx, "DELETE FROM comments WHERE komentar_id = $1", commentID)
	if err != nil {
		return fmt.Errorf("error executing delete for comment ID %d in tx: %w", commentID, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected for comment ID %d delete in tx: %w", commentID, err)
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetCommentCount retrieves the number of comments for an article
func GetCommentCount(ctx context.Context, db *sql.DB, articleID int) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE artikel_id = $1`
	var count int
	err := db.QueryRowContext(ctx, query, articleID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTotalCommentCount retrieves the total number of comments
func GetTotalCommentCount(ctx context.Context, db *sql.DB) (int, error) {
	query := `SELECT COUNT(*) FROM comments`
	var count int
	err := db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetPendingComments retrieves all pending comments for moderation
func GetPendingComments(db *sql.DB, limit int, offset int) ([]Comment, error) {
	return GetAllComments(db, "pending", limit, offset)
}
