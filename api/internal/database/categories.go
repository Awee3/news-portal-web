package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Category struct {
	KategoriID   int    `json:"kategori_id"`
	NamaKategori string `json:"nama_kategori"`
	CreatedAt    string `json:"created_at"`
}

type CategoryRequest struct {
	NamaKategori string `json:"nama_kategori"`
}

func CreateCategory(ctx context.Context, db *sql.DB, req *CategoryRequest) (*Category, error) {
	query := `
        INSERT INTO categories (nama_kategori, created_at)
        VALUES ($1, NOW())
        RETURNING kategori_id, nama_kategori, created_at
    `

	var category Category
	err := db.QueryRowContext(ctx, query, req.NamaKategori).Scan(
		&category.KategoriID, &category.NamaKategori, &category.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return &category, nil
}

func CreateCategoryTx(ctx context.Context, tx *sql.Tx, req *CategoryRequest) (*Category, error) {
	query := `
        INSERT INTO categories (nama_kategori, created_at)
        VALUES ($1, NOW())
        RETURNING kategori_id, nama_kategori, created_at
    `

	var category Category
	err := tx.QueryRowContext(ctx, query, req.NamaKategori).Scan(
		&category.KategoriID, &category.NamaKategori, &category.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return &category, nil
}

func GetCategoryByID(ctx context.Context, db *sql.DB, categoryID int) (*Category, error) {
	query := `
        SELECT kategori_id, nama_kategori, created_at
        FROM categories
        WHERE kategori_id = $1
    `

	var category Category
	err := db.QueryRowContext(ctx, query, categoryID).Scan(
		&category.KategoriID, &category.NamaKategori, &category.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return &category, nil
}

func GetCategoryByName(ctx context.Context, db *sql.DB, name string) (*Category, error) {
	query := `
        SELECT kategori_id, nama_kategori, created_at
        FROM categories
        WHERE nama_kategori = $1
    `

	var category Category
	err := db.QueryRowContext(ctx, query, name).Scan(
		&category.KategoriID, &category.NamaKategori, &category.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return &category, nil
}

func ListCategories(ctx context.Context, db *sql.DB) ([]Category, error) {
	query := `
        SELECT kategori_id, nama_kategori, created_at
        FROM categories
        ORDER BY nama_kategori ASC
    `

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.KategoriID, &category.NamaKategori, &category.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func ListCategoriesWithArticleCount(ctx context.Context, db *sql.DB) ([]map[string]interface{}, error) {
	query := `
        SELECT c.kategori_id, c.nama_kategori, c.created_at,
               COUNT(ak.artikel_id) as article_count
        FROM categories c
        LEFT JOIN artikel_kategori ak ON c.kategori_id = ak.kategori_id
        LEFT JOIN articles a ON ak.artikel_id = a.artikel_id AND a.status = 'published'
        GROUP BY c.kategori_id, c.nama_kategori, c.created_at
        ORDER BY c.nama_kategori ASC
    `

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []map[string]interface{}
	for rows.Next() {
		var kategoriID int
		var namaKategori, createdAt string
		var articleCount int

		err := rows.Scan(&kategoriID, &namaKategori, &createdAt, &articleCount)
		if err != nil {
			return nil, err
		}

		category := map[string]interface{}{
			"kategori_id":   kategoriID,
			"nama_kategori": namaKategori,
			"created_at":    createdAt,
			"article_count": articleCount,
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func UpdateCategory(ctx context.Context, db *sql.DB, categoryID int, req *CategoryRequest) (*Category, error) {
	query := `
        UPDATE categories 
        SET nama_kategori = $1
        WHERE kategori_id = $2
        RETURNING kategori_id, nama_kategori, created_at
    `

	var category Category
	err := db.QueryRowContext(ctx, query, req.NamaKategori, categoryID).Scan(
		&category.KategoriID, &category.NamaKategori, &category.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return &category, nil
}

func UpdateCategoryTx(ctx context.Context, tx *sql.Tx, categoryID int, req *CategoryRequest) (*Category, error) {
	query := `
        UPDATE categories 
        SET nama_kategori = $1
        WHERE kategori_id = $2
        RETURNING kategori_id, nama_kategori, created_at
    `

	var category Category
	err := tx.QueryRowContext(ctx, query, req.NamaKategori, categoryID).Scan(
		&category.KategoriID, &category.NamaKategori, &category.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return &category, nil
}

func DeleteCategory(ctx context.Context, db *sql.DB, categoryID int) error {
	// Check if category has articles
	var articleCount int
	countQuery := `
        SELECT COUNT(*)
        FROM artikel_kategori
        WHERE kategori_id = $1
    `
	err := db.QueryRowContext(ctx, countQuery, categoryID).Scan(&articleCount)
	if err != nil {
		return fmt.Errorf("failed to check category usage: %w", err)
	}

	if articleCount > 0 {
		return errors.New("cannot delete category that has articles")
	}

	res, err := db.ExecContext(ctx, "DELETE FROM categories WHERE kategori_id = $1", categoryID)
	if err != nil {
		return fmt.Errorf("error deleting category ID %d: %w", categoryID, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected for category ID %d delete: %w", categoryID, err)
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteCategoryTx(ctx context.Context, tx *sql.Tx, categoryID int) error {
	// Check if category has articles
	var articleCount int
	countQuery := `
        SELECT COUNT(*)
        FROM artikel_kategori
        WHERE kategori_id = $1
    `
	err := tx.QueryRowContext(ctx, countQuery, categoryID).Scan(&articleCount)
	if err != nil {
		return fmt.Errorf("failed to check category usage: %w", err)
	}

	if articleCount > 0 {
		return errors.New("cannot delete category that has articles")
	}

	res, err := tx.ExecContext(ctx, "DELETE FROM categories WHERE kategori_id = $1", categoryID)
	if err != nil {
		return fmt.Errorf("error executing delete for category ID %d in tx: %w", categoryID, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected for category ID %d delete in tx: %w", categoryID, err)
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func ForceDeleteCategory(ctx context.Context, db *sql.DB, categoryID int) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete category relations first
	_, err = tx.ExecContext(ctx, "DELETE FROM artikel_kategori WHERE kategori_id = $1", categoryID)
	if err != nil {
		return fmt.Errorf("failed to delete category relations: %w", err)
	}

	// Delete category
	res, err := tx.ExecContext(ctx, "DELETE FROM categories WHERE kategori_id = $1", categoryID)
	if err != nil {
		return fmt.Errorf("error deleting category ID %d: %w", categoryID, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected for category ID %d delete: %w", categoryID, err)
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}

func IsCategoryExists(ctx context.Context, db *sql.DB, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE nama_kategori = $1)`
	var exists bool
	err := db.QueryRowContext(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking category existence: %w", err)
	}
	return exists, nil
}
