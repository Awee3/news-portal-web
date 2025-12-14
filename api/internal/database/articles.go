package database

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Article struct {
	ArtikelID         int        `json:"artikel_id"`
	Judul             string     `json:"judul"`
	Slug              string     `json:"slug"`
	Konten            string     `json:"konten"`
	Excerpt           *string    `json:"excerpt,omitempty"`
	GambarUtama       *string    `json:"gambar_utama,omitempty"`
	Penulis           *string    `json:"penulis,omitempty"`
	Status            string     `json:"status"`
	UserID            int        `json:"user_id"`
	TanggalPublikasi  *time.Time `json:"tanggal_publikasi,omitempty"`
	TanggalDibuat     time.Time  `json:"tanggal_dibuat"`
	TanggalDiperbarui time.Time  `json:"tanggal_diperbarui"`
	// Related data (populated separately)
	Kategori []Category `json:"kategori,omitempty"`
	Tags     []Tag      `json:"tags,omitempty"`
}

type ArticleInput struct {
	Judul            string `json:"judul"`
	Slug             string `json:"slug,omitempty"`
	Konten           string `json:"konten"`
	Excerpt          string `json:"excerpt,omitempty"`
	GambarUtama      string `json:"gambar_utama,omitempty"`
	Penulis          string `json:"penulis,omitempty"`
	Status           string `json:"status,omitempty"`
	TanggalPublikasi string `json:"tanggal_publikasi,omitempty"`
	KategoriIDs      []int  `json:"kategori_ids,omitempty"`
	TagIDs           []int  `json:"tag_ids,omitempty"`
}

type ArticleFilter struct {
	Status       string
	KategoriID   int
	KategoriName string // NEW: filter by category name
	TagID        int
	UserID       int
	Search       string
	Limit        int
	Offset       int
}

// GenerateSlug creates URL-friendly slug from title
func GenerateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters (keep alphanumeric and hyphens)
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	// Remove consecutive hyphens
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// EnsureUniqueSlug checks if slug exists and appends number if needed
func EnsureUniqueSlug(db *sql.DB, slug string, excludeID int) (string, error) {
	baseSlug := slug
	counter := 1

	for {
		var exists bool
		query := `SELECT EXISTS(SELECT 1 FROM articles WHERE slug = $1 AND artikel_id != $2)`
		err := db.QueryRow(query, slug, excludeID).Scan(&exists)
		if err != nil {
			return "", err
		}

		if !exists {
			return slug, nil
		}

		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}
}

// GetAllArticles retrieves articles with optional filters
func GetAllArticles(db *sql.DB, filter ArticleFilter) ([]Article, error) {
	query := `
        SELECT DISTINCT a.artikel_id, a.judul, a.slug, a.konten, a.excerpt, 
               a.gambar_utama, a.penulis, a.status, a.user_id, 
               a.tanggal_publikasi, a.tanggal_dibuat, a.tanggal_diperbarui
        FROM articles a
        LEFT JOIN artikel_kategori ak ON a.artikel_id = ak.artikel_id
        LEFT JOIN artikel_tag at ON a.artikel_id = at.artikel_id
    `

	// NEW: Join categories table if filtering by name
	if filter.KategoriName != "" {
		query += `
        LEFT JOIN categories c ON ak.kategori_id = c.kategori_id
        `
	}

	query += `
        WHERE 1=1
    `

	args := []interface{}{}
	argCount := 0

	if filter.Status != "" {
		argCount++
		query += fmt.Sprintf(" AND a.status = $%d", argCount)
		args = append(args, filter.Status)
	}

	if filter.KategoriID > 0 {
		argCount++
		query += fmt.Sprintf(" AND ak.kategori_id = $%d", argCount)
		args = append(args, filter.KategoriID)
	}

	// NEW: Filter by category name
	if filter.KategoriName != "" {
		argCount++
		query += fmt.Sprintf(" AND c.nama_kategori = $%d", argCount)
		args = append(args, filter.KategoriName)
	}

	if filter.TagID > 0 {
		argCount++
		query += fmt.Sprintf(" AND at.tag_id = $%d", argCount)
		args = append(args, filter.TagID)
	}

	if filter.UserID > 0 {
		argCount++
		query += fmt.Sprintf(" AND a.user_id = $%d", argCount)
		args = append(args, filter.UserID)
	}

	if filter.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (a.judul ILIKE $%d OR a.konten ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filter.Search+"%")
	}

	query += " ORDER BY a.tanggal_dibuat DESC"

	if filter.Limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.Limit)
	}

	if filter.Offset > 0 {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filter.Offset)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		err := rows.Scan(
			&a.ArtikelID, &a.Judul, &a.Slug, &a.Konten, &a.Excerpt,
			&a.GambarUtama, &a.Penulis, &a.Status, &a.UserID,
			&a.TanggalPublikasi, &a.TanggalDibuat, &a.TanggalDiperbarui,
		)
		if err != nil {
			return nil, err
		}

		// Fetch categories and tags for each article
		a.Kategori, _ = GetArticleCategories(db, a.ArtikelID)
		a.Tags, _ = GetArticleTags(db, a.ArtikelID)

		articles = append(articles, a)
	}

	return articles, nil
}

// GetArticleByID retrieves a single article by ID
func GetArticleByID(db *sql.DB, id int) (*Article, error) {
	query := `
        SELECT artikel_id, judul, slug, konten, excerpt, gambar_utama, 
               penulis, status, user_id, tanggal_publikasi, 
               tanggal_dibuat, tanggal_diperbarui
        FROM articles
        WHERE artikel_id = $1
    `

	var a Article
	err := db.QueryRow(query, id).Scan(
		&a.ArtikelID, &a.Judul, &a.Slug, &a.Konten, &a.Excerpt,
		&a.GambarUtama, &a.Penulis, &a.Status, &a.UserID,
		&a.TanggalPublikasi, &a.TanggalDibuat, &a.TanggalDiperbarui,
	)
	if err != nil {
		return nil, err
	}

	// Fetch related categories and tags
	a.Kategori, _ = GetArticleCategories(db, a.ArtikelID)
	a.Tags, _ = GetArticleTags(db, a.ArtikelID)

	return &a, nil
}

// GetArticleBySlug retrieves a single article by slug
func GetArticleBySlug(db *sql.DB, slug string) (*Article, error) {
	query := `
        SELECT artikel_id, judul, slug, konten, excerpt, gambar_utama, 
               penulis, status, user_id, tanggal_publikasi, 
               tanggal_dibuat, tanggal_diperbarui
        FROM articles
        WHERE slug = $1
    `

	var a Article
	err := db.QueryRow(query, slug).Scan(
		&a.ArtikelID, &a.Judul, &a.Slug, &a.Konten, &a.Excerpt,
		&a.GambarUtama, &a.Penulis, &a.Status, &a.UserID,
		&a.TanggalPublikasi, &a.TanggalDibuat, &a.TanggalDiperbarui,
	)
	if err != nil {
		return nil, err
	}

	// Fetch related categories and tags
	a.Kategori, _ = GetArticleCategories(db, a.ArtikelID)
	a.Tags, _ = GetArticleTags(db, a.ArtikelID)

	return &a, nil
}

// GetPublishedArticleBySlug retrieves a published article by slug (for public access)
func GetPublishedArticleBySlug(db *sql.DB, slug string) (*Article, error) {
	query := `
        SELECT artikel_id, judul, slug, konten, excerpt, gambar_utama, 
               penulis, status, user_id, tanggal_publikasi, 
               tanggal_dibuat, tanggal_diperbarui
        FROM articles
        WHERE slug = $1 AND status = 'published'
    `

	var a Article
	err := db.QueryRow(query, slug).Scan(
		&a.ArtikelID, &a.Judul, &a.Slug, &a.Konten, &a.Excerpt,
		&a.GambarUtama, &a.Penulis, &a.Status, &a.UserID,
		&a.TanggalPublikasi, &a.TanggalDibuat, &a.TanggalDiperbarui,
	)
	if err != nil {
		return nil, err
	}

	// Fetch related categories and tags
	a.Kategori, _ = GetArticleCategories(db, a.ArtikelID)
	a.Tags, _ = GetArticleTags(db, a.ArtikelID)

	return &a, nil
}

// CreateArticle creates a new article
func CreateArticle(db *sql.DB, input ArticleInput, userID int) (*Article, error) {
	// Generate slug if not provided
	slug := input.Slug
	if slug == "" {
		slug = GenerateSlug(input.Judul)
	}

	// Ensure slug is unique
	slug, err := EnsureUniqueSlug(db, slug, 0)
	if err != nil {
		return nil, err
	}

	// Set default status
	status := input.Status
	if status == "" {
		status = "draft"
	}

	// Set penulis from input or leave NULL
	var penulis *string
	if input.Penulis != "" {
		penulis = &input.Penulis
	}

	// Parse tanggal_publikasi
	var tanggalPublikasi *time.Time
	if input.TanggalPublikasi != "" {
		t, err := time.Parse(time.RFC3339, input.TanggalPublikasi)
		if err == nil {
			tanggalPublikasi = &t
		}
	}

	// If publishing, set tanggal_publikasi to now if not provided
	if status == "published" && tanggalPublikasi == nil {
		now := time.Now()
		tanggalPublikasi = &now
	}

	query := `
        INSERT INTO articles (judul, slug, konten, excerpt, gambar_utama, penulis, status, user_id, tanggal_publikasi)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING artikel_id, judul, slug, konten, excerpt, gambar_utama, penulis, status, user_id, tanggal_publikasi, tanggal_dibuat, tanggal_diperbarui
    `

	var a Article
	var excerpt, gambarUtama *string
	if input.Excerpt != "" {
		excerpt = &input.Excerpt
	}
	if input.GambarUtama != "" {
		gambarUtama = &input.GambarUtama
	}

	err = db.QueryRow(
		query,
		input.Judul, slug, input.Konten, excerpt, gambarUtama,
		penulis, status, userID, tanggalPublikasi,
	).Scan(
		&a.ArtikelID, &a.Judul, &a.Slug, &a.Konten, &a.Excerpt,
		&a.GambarUtama, &a.Penulis, &a.Status, &a.UserID,
		&a.TanggalPublikasi, &a.TanggalDibuat, &a.TanggalDiperbarui,
	)
	if err != nil {
		return nil, err
	}

	// Add categories
	if len(input.KategoriIDs) > 0 {
		for _, katID := range input.KategoriIDs {
			_, err := db.Exec(
				"INSERT INTO artikel_kategori (artikel_id, kategori_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
				a.ArtikelID, katID,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	// Add tags
	if len(input.TagIDs) > 0 {
		for _, tagID := range input.TagIDs {
			_, err := db.Exec(
				"INSERT INTO artikel_tag (artikel_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
				a.ArtikelID, tagID,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	// Fetch related data
	a.Kategori, _ = GetArticleCategories(db, a.ArtikelID)
	a.Tags, _ = GetArticleTags(db, a.ArtikelID)

	return &a, nil
}

// UpdateArticle updates an existing article
func UpdateArticle(db *sql.DB, id int, input ArticleInput) (*Article, error) {
	// Generate slug if provided or changed
	slug := input.Slug
	if slug == "" {
		slug = GenerateSlug(input.Judul)
	}

	// Ensure slug is unique (excluding current article)
	slug, err := EnsureUniqueSlug(db, slug, id)
	if err != nil {
		return nil, err
	}

	// Parse tanggal_publikasi
	var tanggalPublikasi *time.Time
	if input.TanggalPublikasi != "" {
		t, err := time.Parse(time.RFC3339, input.TanggalPublikasi)
		if err == nil {
			tanggalPublikasi = &t
		}
	}

	// If publishing for first time, set tanggal_publikasi
	if input.Status == "published" && tanggalPublikasi == nil {
		// Check if already published
		var existingStatus string
		var existingPubDate *time.Time
		db.QueryRow("SELECT status, tanggal_publikasi FROM articles WHERE artikel_id = $1", id).Scan(&existingStatus, &existingPubDate)

		if existingPubDate == nil {
			now := time.Now()
			tanggalPublikasi = &now
		} else {
			tanggalPublikasi = existingPubDate
		}
	}

	query := `
        UPDATE articles 
        SET judul = $1, slug = $2, konten = $3, excerpt = $4, gambar_utama = $5, 
            penulis = $6, status = $7, tanggal_publikasi = $8
        WHERE artikel_id = $9
        RETURNING artikel_id, judul, slug, konten, excerpt, gambar_utama, penulis, status, user_id, tanggal_publikasi, tanggal_dibuat, tanggal_diperbarui
    `

	var a Article
	var excerpt, gambarUtama, penulis *string
	if input.Excerpt != "" {
		excerpt = &input.Excerpt
	}
	if input.GambarUtama != "" {
		gambarUtama = &input.GambarUtama
	}
	if input.Penulis != "" {
		penulis = &input.Penulis
	}

	err = db.QueryRow(
		query,
		input.Judul, slug, input.Konten, excerpt, gambarUtama,
		penulis, input.Status, tanggalPublikasi, id,
	).Scan(
		&a.ArtikelID, &a.Judul, &a.Slug, &a.Konten, &a.Excerpt,
		&a.GambarUtama, &a.Penulis, &a.Status, &a.UserID,
		&a.TanggalPublikasi, &a.TanggalDibuat, &a.TanggalDiperbarui,
	)
	if err != nil {
		return nil, err
	}

	// Update categories
	if input.KategoriIDs != nil {
		// Remove existing
		db.Exec("DELETE FROM artikel_kategori WHERE artikel_id = $1", id)
		// Add new
		for _, katID := range input.KategoriIDs {
			db.Exec(
				"INSERT INTO artikel_kategori (artikel_id, kategori_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
				id, katID,
			)
		}
	}

	// Update tags
	if input.TagIDs != nil {
		// Remove existing
		db.Exec("DELETE FROM artikel_tag WHERE artikel_id = $1", id)
		// Add new
		for _, tagID := range input.TagIDs {
			db.Exec(
				"INSERT INTO artikel_tag (artikel_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
				id, tagID,
			)
		}
	}

	// Fetch related data
	a.Kategori, _ = GetArticleCategories(db, a.ArtikelID)
	a.Tags, _ = GetArticleTags(db, a.ArtikelID)

	return &a, nil
}

// DeleteArticle deletes an article by ID
func DeleteArticle(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM articles WHERE artikel_id = $1", id)
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

// GetArticleCategories retrieves categories for an article
func GetArticleCategories(db *sql.DB, artikelID int) ([]Category, error) {
	query := `
        SELECT c.kategori_id, c.nama_kategori, c.deskripsi, c.created_at
        FROM categories c
        JOIN artikel_kategori ak ON c.kategori_id = ak.kategori_id
        WHERE ak.artikel_id = $1
    `

	rows, err := db.Query(query, artikelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		err := rows.Scan(&c.KategoriID, &c.NamaKategori, &c.Deskripsi, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

// GetArticleTags retrieves tags for an article
func GetArticleTags(db *sql.DB, artikelID int) ([]Tag, error) {
	query := `
        SELECT t.tag_id, t.nama_tag, t.created_at
        FROM tags t
        JOIN artikel_tag at ON t.tag_id = at.tag_id
        WHERE at.artikel_id = $1
    `

	rows, err := db.Query(query, artikelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var t Tag
		err := rows.Scan(&t.TagID, &t.NamaTag, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return tags, nil
}

// GetArticlesByCategory retrieves articles by category ID
func GetArticlesByCategory(db *sql.DB, kategoriID int, limit int, offset int) ([]Article, error) {
	filter := ArticleFilter{
		Status:     "published",
		KategoriID: kategoriID,
		Limit:      limit,
		Offset:     offset,
	}
	return GetAllArticles(db, filter)
}

// GetArticlesByTag retrieves articles by tag ID
func GetArticlesByTag(db *sql.DB, tagID int, limit int, offset int) ([]Article, error) {
	filter := ArticleFilter{
		Status: "published",
		TagID:  tagID,
		Limit:  limit,
		Offset: offset,
	}
	return GetAllArticles(db, filter)
}
