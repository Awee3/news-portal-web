package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Comment struct {
	KomentarID int       `json:"komentar_id"`
	ArtikelID  int       `json:"artikel_id"`
	UserID     *int      `json:"user_id"`
	Konten     string    `json:"konten"`
	CreatedAt  time.Time `json:"created_at"`
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
        INSERT INTO comments (artikel_id, user_id, konten, created_at)
        VALUES ($1, $2, $3, NOW())
        RETURNING komentar_id, artikel_id, user_id, konten, created_at
    `

	var comment Comment
	err = db.QueryRowContext(ctx, query, req.ArtikelID, userID, req.Konten).Scan(
		&comment.KomentarID, &comment.ArtikelID, &comment.UserID,
		&comment.Konten, &comment.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return &comment, nil
}

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
        INSERT INTO comments (artikel_id, user_id, konten, created_at)
        VALUES ($1, $2, $3, NOW())
        RETURNING komentar_id, artikel_id, user_id, konten, created_at
    `

	var comment Comment
	err = tx.QueryRowContext(ctx, query, req.ArtikelID, userID, req.Konten).Scan(
		&comment.KomentarID, &comment.ArtikelID, &comment.UserID,
		&comment.Konten, &comment.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return &comment, nil
}

func GetCommentByID(ctx context.Context, db *sql.DB, commentID int) (*CommentWithAuthor, error) {
	query := `
        SELECT c.komentar_id, c.artikel_id, c.user_id, c.konten, c.created_at,
               u.username, u.email, a.judul
        FROM comments c
        LEFT JOIN users u ON c.user_id = u.user_id
        JOIN articles a ON c.artikel_id = a.artikel_id
        WHERE c.komentar_id = $1
    `

	var comment CommentWithAuthor
	err := db.QueryRowContext(ctx, query, commentID).Scan(
		&comment.KomentarID, &comment.ArtikelID, &comment.UserID,
		&comment.Konten, &comment.CreatedAt,
		&comment.AuthorUsername, &comment.AuthorEmail, &comment.ArticleTitle)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}

	return &comment, nil
}

func ListCommentsByArticle(ctx context.Context, db *sql.DB, articleID int) ([]CommentWithAuthor, error) {
	query := `
        SELECT c.komentar_id, c.artikel_id, c.user_id, c.konten, c.created_at,
               u.username, u.email, a.judul
        FROM comments c
        LEFT JOIN users u ON c.user_id = u.user_id
        JOIN articles a ON c.artikel_id = a.artikel_id
        WHERE c.artikel_id = $1
        ORDER BY c.created_at ASC
    `

	rows, err := db.QueryContext(ctx, query, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []CommentWithAuthor
	for rows.Next() {
		var comment CommentWithAuthor
		err := rows.Scan(
			&comment.KomentarID, &comment.ArtikelID, &comment.UserID,
			&comment.Konten, &comment.CreatedAt,
			&comment.AuthorUsername, &comment.AuthorEmail, &comment.ArticleTitle)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func ListAllComments(ctx context.Context, db *sql.DB, limit, offset int) ([]CommentWithAuthor, int, error) {
	// Count query
	countQuery := `SELECT COUNT(*) FROM comments`

	// Select query
	selectQuery := `
        SELECT c.komentar_id, c.artikel_id, c.user_id, c.konten, c.created_at,
               u.username, u.email, a.judul
        FROM comments c
        LEFT JOIN users u ON c.user_id = u.user_id
        JOIN articles a ON c.artikel_id = a.artikel_id
        ORDER BY c.created_at DESC
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

	var comments []CommentWithAuthor
	for rows.Next() {
		var comment CommentWithAuthor
		err := rows.Scan(
			&comment.KomentarID, &comment.ArtikelID, &comment.UserID,
			&comment.Konten, &comment.CreatedAt,
			&comment.AuthorUsername, &comment.AuthorEmail, &comment.ArticleTitle)
		if err != nil {
			return nil, 0, err
		}
		comments = append(comments, comment)
	}

	return comments, totalCount, nil
}

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
        SET konten = $1
        WHERE komentar_id = $2
        RETURNING komentar_id, artikel_id, user_id, konten, created_at
    `

	var comment Comment
	err = db.QueryRowContext(ctx, query, konten, commentID).Scan(
		&comment.KomentarID, &comment.ArtikelID, &comment.UserID,
		&comment.Konten, &comment.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return &comment, nil
}

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

func GetCommentCount(ctx context.Context, db *sql.DB, articleID int) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE artikel_id = $1`
	var count int
	err := db.QueryRowContext(ctx, query, articleID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetTotalCommentCount(ctx context.Context, db *sql.DB) (int, error) {
	query := `SELECT COUNT(*) FROM comments`
	var count int
	err := db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
