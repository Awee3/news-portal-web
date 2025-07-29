package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Tag struct {
	TagID   int    `json:"tag_id"`
	NamaTag string `json:"nama_tag"`
	// Hapus CreatedAt field karena tidak ada di database
}

type TagRequest struct {
	NamaTag string `json:"nama_tag"`
}

func CreateTag(ctx context.Context, db *sql.DB, req *TagRequest) (*Tag, error) {
	query := `
        INSERT INTO tags (nama_tag)
        VALUES ($1)
        RETURNING tag_id, nama_tag
    `

	var tag Tag
	err := db.QueryRowContext(ctx, query, req.NamaTag).Scan(
		&tag.TagID, &tag.NamaTag)
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return &tag, nil
}

func CreateTagTx(ctx context.Context, tx *sql.Tx, req *TagRequest) (*Tag, error) {
	query := `
        INSERT INTO tags (nama_tag)
        VALUES ($1)
        RETURNING tag_id, nama_tag
    `

	var tag Tag
	err := tx.QueryRowContext(ctx, query, req.NamaTag).Scan(
		&tag.TagID, &tag.NamaTag)
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return &tag, nil
}

func GetTagByID(ctx context.Context, db *sql.DB, tagID int) (*Tag, error) {
	query := `
        SELECT tag_id, nama_tag
        FROM tags
        WHERE tag_id = $1
    `

	var tag Tag
	err := db.QueryRowContext(ctx, query, tagID).Scan(
		&tag.TagID, &tag.NamaTag)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}

	return &tag, nil
}

func GetTagByName(ctx context.Context, db *sql.DB, name string) (*Tag, error) {
	query := `
        SELECT tag_id, nama_tag
        FROM tags
        WHERE nama_tag = $1
    `

	var tag Tag
	err := db.QueryRowContext(ctx, query, name).Scan(
		&tag.TagID, &tag.NamaTag)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}

	return &tag, nil
}

func ListTags(ctx context.Context, db *sql.DB) ([]Tag, error) {
	query := `
        SELECT tag_id, nama_tag
        FROM tags
        ORDER BY nama_tag ASC
    `

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.TagID, &tag.NamaTag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func ListTagsWithArticleCount(ctx context.Context, db *sql.DB) ([]map[string]interface{}, error) {
	query := `
        SELECT t.tag_id, t.nama_tag,
               COUNT(at.artikel_id) as article_count
        FROM tags t
        LEFT JOIN artikel_tag at ON t.tag_id = at.tag_id
        LEFT JOIN articles a ON at.artikel_id = a.artikel_id AND a.status = 'published'
        GROUP BY t.tag_id, t.nama_tag
        ORDER BY t.nama_tag ASC
    `

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []map[string]interface{}
	for rows.Next() {
		var tagID int
		var namaTag string
		var articleCount int

		err := rows.Scan(&tagID, &namaTag, &articleCount)
		if err != nil {
			return nil, err
		}

		tag := map[string]interface{}{
			"tag_id":        tagID,
			"nama_tag":      namaTag,
			"article_count": articleCount,
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func ListPopularTags(ctx context.Context, db *sql.DB, limit int) ([]map[string]interface{}, error) {
	query := `
        SELECT t.tag_id, t.nama_tag,
               COUNT(at.artikel_id) as article_count
        FROM tags t
        LEFT JOIN artikel_tag at ON t.tag_id = at.tag_id
        LEFT JOIN articles a ON at.artikel_id = a.artikel_id AND a.status = 'published'
        GROUP BY t.tag_id, t.nama_tag
        HAVING COUNT(at.artikel_id) > 0
        ORDER BY article_count DESC, t.nama_tag ASC
        LIMIT $1
    `

	rows, err := db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []map[string]interface{}
	for rows.Next() {
		var tagID int
		var namaTag string
		var articleCount int

		err := rows.Scan(&tagID, &namaTag, &articleCount)
		if err != nil {
			return nil, err
		}

		tag := map[string]interface{}{
			"tag_id":        tagID,
			"nama_tag":      namaTag,
			"article_count": articleCount,
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func SearchTags(ctx context.Context, db *sql.DB, keyword string) ([]Tag, error) {
	query := `
        SELECT tag_id, nama_tag
        FROM tags
        WHERE LOWER(nama_tag) LIKE LOWER($1)
        ORDER BY nama_tag ASC
    `

	searchTerm := "%" + keyword + "%"
	rows, err := db.QueryContext(ctx, query, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.TagID, &tag.NamaTag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func UpdateTag(ctx context.Context, db *sql.DB, tagID int, req *TagRequest) (*Tag, error) {
	query := `
        UPDATE tags 
        SET nama_tag = $1
        WHERE tag_id = $2
        RETURNING tag_id, nama_tag
    `

	var tag Tag
	err := db.QueryRowContext(ctx, query, req.NamaTag, tagID).Scan(
		&tag.TagID, &tag.NamaTag)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tag not found")
		}
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	return &tag, nil
}

func UpdateTagTx(ctx context.Context, tx *sql.Tx, tagID int, req *TagRequest) (*Tag, error) {
	query := `
        UPDATE tags 
        SET nama_tag = $1
        WHERE tag_id = $2
        RETURNING tag_id, nama_tag
    `

	var tag Tag
	err := tx.QueryRowContext(ctx, query, req.NamaTag, tagID).Scan(
		&tag.TagID, &tag.NamaTag)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tag not found")
		}
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	return &tag, nil
}

func DeleteTag(ctx context.Context, db *sql.DB, tagID int) error {
	// Check if tag has articles
	var articleCount int
	countQuery := `
        SELECT COUNT(*)
        FROM artikel_tag
        WHERE tag_id = $1
    `
	err := db.QueryRowContext(ctx, countQuery, tagID).Scan(&articleCount)
	if err != nil {
		return fmt.Errorf("failed to check tag usage: %w", err)
	}

	if articleCount > 0 {
		return errors.New("cannot delete tag that has articles")
	}

	res, err := db.ExecContext(ctx, "DELETE FROM tags WHERE tag_id = $1", tagID)
	if err != nil {
		return fmt.Errorf("error deleting tag ID %d: %w", tagID, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected for tag ID %d delete: %w", tagID, err)
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteTagTx(ctx context.Context, tx *sql.Tx, tagID int) error {
	// Check if tag has articles
	var articleCount int
	countQuery := `
        SELECT COUNT(*)
        FROM artikel_tag
        WHERE tag_id = $1
    `
	err := tx.QueryRowContext(ctx, countQuery, tagID).Scan(&articleCount)
	if err != nil {
		return fmt.Errorf("failed to check tag usage: %w", err)
	}

	if articleCount > 0 {
		return errors.New("cannot delete tag that has articles")
	}

	res, err := tx.ExecContext(ctx, "DELETE FROM tags WHERE tag_id = $1", tagID)
	if err != nil {
		return fmt.Errorf("error executing delete for tag ID %d in tx: %w", tagID, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected for tag ID %d delete in tx: %w", tagID, err)
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func ForceDeleteTag(ctx context.Context, db *sql.DB, tagID int) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete tag relations first
	_, err = tx.ExecContext(ctx, "DELETE FROM artikel_tag WHERE tag_id = $1", tagID)
	if err != nil {
		return fmt.Errorf("failed to delete tag relations: %w", err)
	}

	// Delete tag
	res, err := tx.ExecContext(ctx, "DELETE FROM tags WHERE tag_id = $1", tagID)
	if err != nil {
		return fmt.Errorf("error deleting tag ID %d: %w", tagID, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected for tag ID %d delete: %w", tagID, err)
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}

func IsTagExists(ctx context.Context, db *sql.DB, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM tags WHERE nama_tag = $1)`
	var exists bool
	err := db.QueryRowContext(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking tag existence: %w", err)
	}
	return exists, nil
}

func GetOrCreateTags(ctx context.Context, db *sql.DB, tagNames []string) ([]int, error) {
	if len(tagNames) == 0 {
		return []int{}, nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var tagIDs []int

	for _, tagName := range tagNames {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}

		// Try to get existing tag
		var tagID int
		err := tx.QueryRowContext(ctx, "SELECT tag_id FROM tags WHERE nama_tag = $1", tagName).Scan(&tagID)
		if err == sql.ErrNoRows {
			// Create new tag
			err = tx.QueryRowContext(ctx,
				"INSERT INTO tags (nama_tag) VALUES ($1) RETURNING tag_id",
				tagName).Scan(&tagID)
			if err != nil {
				return nil, fmt.Errorf("failed to create tag %s: %w", tagName, err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("failed to get tag %s: %w", tagName, err)
		}

		tagIDs = append(tagIDs, tagID)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tagIDs, nil
}

func GetTagsByArticleID(ctx context.Context, db *sql.DB, articleID int) ([]Tag, error) {
	query := `
        SELECT t.tag_id, t.nama_tag
        FROM tags t
        JOIN artikel_tag at ON t.tag_id = at.tag_id
        WHERE at.artikel_id = $1
        ORDER BY t.nama_tag ASC
    `

	rows, err := db.QueryContext(ctx, query, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.TagID, &tag.NamaTag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
