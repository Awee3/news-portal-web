package database

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"
    "unicode"
)

type Article struct {
    ArtikelID        int       `json:"artikel_id"`
    UserID           int       `json:"user_id"`
    Judul            string    `json:"judul"`
    Slug             string    `json:"slug"`
    Konten           string    `json:"konten"`
    FeaturedImageURL *string   `json:"featured_image_url"`
    Status           string    `json:"status"`
    ViewCount        int       `json:"view_count"`
    TanggalPublikasi *time.Time `json:"tanggal_publikasi"`
    MetaTitle        *string   `json:"meta_title"`
    MetaDescription  *string   `json:"meta_description"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}

type ArticleWithAuthor struct {
    Article
    AuthorUsername string `json:"author_username"`
    AuthorEmail    string `json:"author_email"`
}

type ArticleRequest struct {
    Judul            string   `json:"judul"`
    Konten           string   `json:"konten"`
    FeaturedImageURL *string  `json:"featured_image_url"`
    Status           string   `json:"status"`
    CategoryIDs      []int    `json:"category_ids"`
    TagIDs           []int    `json:"tag_ids"`
    MetaTitle        *string  `json:"meta_title"`
    MetaDescription  *string  `json:"meta_description"`
}

func CreateArticle(ctx context.Context, db *sql.DB, userID int, req *ArticleRequest) (*Article, error) {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    article, err := CreateArticleTx(ctx, tx, userID, req)
    if err != nil {
        return nil, err
    }

    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }

    return article, nil
}

func CreateArticleTx(ctx context.Context, tx *sql.Tx, userID int, req *ArticleRequest) (*Article, error) {
    // Generate slug
    slug := generateSlug(req.Judul)
    
    // Ensure slug uniqueness
    originalSlug := slug
    counter := 1
    for {
        var existingID int
        err := tx.QueryRowContext(ctx, "SELECT artikel_id FROM articles WHERE slug = $1", slug).Scan(&existingID)
        if err == sql.ErrNoRows {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
        }
        slug = fmt.Sprintf("%s-%d", originalSlug, counter)
        counter++
    }

    // Set publish date if status is published
    var publishDate *time.Time
    if req.Status == "published" {
        now := time.Now()
        publishDate = &now
    }

    query := `
        INSERT INTO articles (user_id, judul, slug, konten, featured_image_url, 
            status, meta_title, meta_description, view_count, tanggal_publikasi, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
        RETURNING artikel_id, created_at, updated_at
    `

    var article Article
    err := tx.QueryRowContext(ctx, query, userID, req.Judul, slug, req.Konten,
        req.FeaturedImageURL, req.Status, req.MetaTitle, req.MetaDescription, 0, publishDate).Scan(
        &article.ArtikelID, &article.CreatedAt, &article.UpdatedAt)
    if err != nil {
        return nil, fmt.Errorf("failed to insert article: %w", err)
    }

    // Set article fields
    article.UserID = userID
    article.Judul = req.Judul
    article.Slug = slug
    article.Konten = req.Konten
    article.FeaturedImageURL = req.FeaturedImageURL
    article.Status = req.Status
    article.ViewCount = 0
    article.TanggalPublikasi = publishDate
    article.MetaTitle = req.MetaTitle
    article.MetaDescription = req.MetaDescription

    // Insert categories
    if len(req.CategoryIDs) > 0 {
        if err := insertArticleCategoriesTx(ctx, tx, article.ArtikelID, req.CategoryIDs); err != nil {
            return nil, err
        }
    }

    // Insert tags
    if len(req.TagIDs) > 0 {
        if err := insertArticleTagsTx(ctx, tx, article.ArtikelID, req.TagIDs); err != nil {
            return nil, err
        }
    }

    return &article, nil
}

func GetArticleBySlug(ctx context.Context, db *sql.DB, slug string) (*ArticleWithAuthor, error) {
    query := `
        SELECT a.artikel_id, a.user_id, a.judul, a.slug, a.konten,
            a.featured_image_url, a.status, a.view_count, a.tanggal_publikasi,
            a.created_at, a.updated_at, a.meta_title, a.meta_description,
            u.username, u.email
        FROM articles a
        JOIN users u ON a.user_id = u.user_id
        WHERE a.slug = $1 AND a.status = 'published'
    `

    var article ArticleWithAuthor
    err := db.QueryRowContext(ctx, query, slug).Scan(
        &article.ArtikelID, &article.UserID, &article.Judul, &article.Slug,
        &article.Konten, &article.FeaturedImageURL, &article.Status,
        &article.ViewCount, &article.TanggalPublikasi, &article.CreatedAt,
        &article.UpdatedAt, &article.MetaTitle, &article.MetaDescription,
        &article.AuthorUsername, &article.AuthorEmail,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("article not found")
        }
        return nil, err
    }

    // Increment view count (async)
    go incrementViewCount(db, article.ArtikelID)

    return &article, nil
}

func GetArticleByID(ctx context.Context, db *sql.DB, articleID int) (*ArticleWithAuthor, error) {
    query := `
        SELECT a.artikel_id, a.user_id, a.judul, a.slug, a.konten,
            a.featured_image_url, a.status, a.view_count, a.tanggal_publikasi,
            a.created_at, a.updated_at, a.meta_title, a.meta_description,
            u.username, u.email
        FROM articles a
        JOIN users u ON a.user_id = u.user_id
        WHERE a.artikel_id = $1
    `

    var article ArticleWithAuthor
    err := db.QueryRowContext(ctx, query, articleID).Scan(
        &article.ArtikelID, &article.UserID, &article.Judul, &article.Slug,
        &article.Konten, &article.FeaturedImageURL, &article.Status,
        &article.ViewCount, &article.TanggalPublikasi, &article.CreatedAt,
        &article.UpdatedAt, &article.MetaTitle, &article.MetaDescription,
        &article.AuthorUsername, &article.AuthorEmail,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("article not found")
        }
        return nil, err
    }

    return &article, nil
}

func ListArticles(ctx context.Context, db *sql.DB, limit, offset int, categoryID *int, status string) ([]ArticleWithAuthor, int, error) {
    // Base query for counting
    countQuery := `
        SELECT COUNT(*)
        FROM articles a
        WHERE 1=1
    `
    
    // Base query for selecting
    selectQuery := `
        SELECT a.artikel_id, a.user_id, a.judul, a.slug, a.konten,
            a.featured_image_url, a.status, a.view_count, a.tanggal_publikasi,
            a.created_at, a.updated_at, a.meta_title, a.meta_description,
            u.username, u.email
        FROM articles a
        JOIN users u ON a.user_id = u.user_id
        WHERE 1=1
    `

    args := []interface{}{}
    argIndex := 1

    // Apply filters
    if status != "" {
        condition := fmt.Sprintf(" AND a.status = $%d", argIndex)
        countQuery += condition
        selectQuery += condition
        args = append(args, status)
        argIndex++
    }

    if categoryID != nil {
        condition := fmt.Sprintf(` AND a.artikel_id IN (
            SELECT artikel_id FROM artikel_kategori WHERE kategori_id = $%d
        )`, argIndex)
        countQuery += condition
        selectQuery += condition
        args = append(args, *categoryID)
        argIndex++
    }

    // Get total count
    var totalCount int
    err := db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
    if err != nil {
        return nil, 0, err
    }

    // Add ordering and pagination
    selectQuery += ` ORDER BY a.tanggal_publikasi DESC, a.created_at DESC LIMIT $%d OFFSET $%d`
    args = append(args, limit, offset)

    // Execute main query
    rows, err := db.QueryContext(ctx, selectQuery, args...)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var articles []ArticleWithAuthor
    for rows.Next() {
        var article ArticleWithAuthor
        err := rows.Scan(
            &article.ArtikelID, &article.UserID, &article.Judul, &article.Slug,
            &article.Konten, &article.FeaturedImageURL, &article.Status,
            &article.ViewCount, &article.TanggalPublikasi, &article.CreatedAt,
            &article.UpdatedAt, &article.MetaTitle, &article.MetaDescription,
            &article.AuthorUsername, &article.AuthorEmail,
        )
        if err != nil {
            return nil, 0, err
        }
        articles = append(articles, article)
    }

    return articles, totalCount, nil
}

func UpdateArticle(ctx context.Context, db *sql.DB, articleID int, userID int, req *ArticleRequest) (*Article, error) {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    article, err := UpdateArticleTx(ctx, tx, articleID, userID, req)
    if err != nil {
        return nil, err
    }

    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }

    return article, nil
}

func UpdateArticleTx(ctx context.Context, tx *sql.Tx, articleID int, userID int, req *ArticleRequest) (*Article, error) {
    // Check if article exists and user owns it
    var existingUserID int
    var currentSlug string
    err := tx.QueryRowContext(ctx, "SELECT user_id, slug FROM articles WHERE artikel_id = $1", articleID).Scan(&existingUserID, &currentSlug)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("article not found")
        }
        return nil, err
    }

    if existingUserID != userID {
        return nil, errors.New("unauthorized: you can only edit your own articles")
    }

    // Generate new slug if title changed
    newSlug := generateSlug(req.Judul)
    if newSlug != currentSlug {
        // Ensure slug uniqueness
        originalSlug := newSlug
        counter := 1
        for {
            var existingID int
            err := tx.QueryRowContext(ctx, "SELECT artikel_id FROM articles WHERE slug = $1 AND artikel_id != $2", newSlug, articleID).Scan(&existingID)
            if err == sql.ErrNoRows {
                break
            }
            if err != nil {
                return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
            }
            newSlug = fmt.Sprintf("%s-%d", originalSlug, counter)
            counter++
        }
    }

    // Set publish date if changing to published
    var publishDate *time.Time
    if req.Status == "published" {
        var currentStatus string
        var currentPublishDate *time.Time
        err := tx.QueryRowContext(ctx, "SELECT status, tanggal_publikasi FROM articles WHERE artikel_id = $1", articleID).Scan(&currentStatus, &currentPublishDate)
        if err != nil {
            return nil, err
        }

        if currentStatus != "published" || currentPublishDate == nil {
            now := time.Now()
            publishDate = &now
        } else {
            publishDate = currentPublishDate
        }
    }

    // Update article
    query := `
        UPDATE articles SET 
            judul = $1, slug = $2, konten = $3, featured_image_url = $4,
            status = $5, meta_title = $6, meta_description = $7, 
            tanggal_publikasi = $8, updated_at = NOW()
        WHERE artikel_id = $9
        RETURNING artikel_id, user_id, created_at, updated_at, view_count
    `

    var article Article
    err = tx.QueryRowContext(ctx, query, req.Judul, newSlug, req.Konten, req.FeaturedImageURL,
        req.Status, req.MetaTitle, req.MetaDescription, publishDate, articleID).Scan(
        &article.ArtikelID, &article.UserID, &article.CreatedAt, &article.UpdatedAt, &article.ViewCount)
    if err != nil {
        return nil, err
    }

    // Delete existing categories and tags
    if _, err := tx.ExecContext(ctx, "DELETE FROM artikel_kategori WHERE artikel_id = $1", articleID); err != nil {
        return nil, err
    }
    if _, err := tx.ExecContext(ctx, "DELETE FROM artikel_tag WHERE artikel_id = $1", articleID); err != nil {
        return nil, err
    }

    // Insert new categories and tags
    if len(req.CategoryIDs) > 0 {
        if err := insertArticleCategoriesTx(ctx, tx, articleID, req.CategoryIDs); err != nil {
            return nil, err
        }
    }
    if len(req.TagIDs) > 0 {
        if err := insertArticleTagsTx(ctx, tx, articleID, req.TagIDs); err != nil {
            return nil, err
        }
    }

    // Set article fields
    article.Judul = req.Judul
    article.Slug = newSlug
    article.Konten = req.Konten
    article.FeaturedImageURL = req.FeaturedImageURL
    article.Status = req.Status
    article.TanggalPublikasi = publishDate
    article.MetaTitle = req.MetaTitle
    article.MetaDescription = req.MetaDescription

    return &article, nil
}

func DeleteArticle(ctx context.Context, db *sql.DB, articleID int, userID int) error {
    // Check if article exists and user owns it
    var existingUserID int
    err := db.QueryRowContext(ctx, "SELECT user_id FROM articles WHERE artikel_id = $1", articleID).Scan(&existingUserID)
    if err != nil {
        if err == sql.ErrNoRows {
            return sql.ErrNoRows
        }
        return err
    }

    if existingUserID != userID {
        return errors.New("unauthorized: you can only delete your own articles")
    }

    // Delete article (CASCADE will handle related tables)
    res, err := db.ExecContext(ctx, "DELETE FROM articles WHERE artikel_id = $1", articleID)
    if err != nil {
        return fmt.Errorf("error deleting article ID %d: %w", articleID, err)
    }

    count, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected for article ID %d delete: %w", articleID, err)
    }
    if count == 0 {
        return sql.ErrNoRows
    }

    return nil
}

func DeleteArticleTx(ctx context.Context, tx *sql.Tx, articleID int, userID int) error {
    // Check if article exists and user owns it
    var existingUserID int
    err := tx.QueryRowContext(ctx, "SELECT user_id FROM articles WHERE artikel_id = $1", articleID).Scan(&existingUserID)
    if err != nil {
        if err == sql.ErrNoRows {
            return sql.ErrNoRows
        }
        return err
    }

    if existingUserID != userID {
        return errors.New("unauthorized: you can only delete your own articles")
    }

    res, err := tx.ExecContext(ctx, "DELETE FROM articles WHERE artikel_id = $1", articleID)
    if err != nil {
        return fmt.Errorf("error executing delete for article ID %d in tx: %w", articleID, err)
    }

    count, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected for article ID %d delete in tx: %w", articleID, err)
    }
    if count == 0 {
        return sql.ErrNoRows
    }

    return nil
}

// Helper functions
func insertArticleCategoriesTx(ctx context.Context, tx *sql.Tx, articleID int, categoryIDs []int) error {
    if len(categoryIDs) == 0 {
        return nil
    }

    query := "INSERT INTO artikel_kategori (artikel_id, kategori_id) VALUES "
    values := []interface{}{}
    placeholders := []string{}

    for i, categoryID := range categoryIDs {
        placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
        values = append(values, articleID, categoryID)
    }

    query += strings.Join(placeholders, ", ")
    _, err := tx.ExecContext(ctx, query, values...)
    return err
}

func insertArticleTagsTx(ctx context.Context, tx *sql.Tx, articleID int, tagIDs []int) error {
    if len(tagIDs) == 0 {
        return nil
    }

    query := "INSERT INTO artikel_tag (artikel_id, tag_id) VALUES "
    values := []interface{}{}
    placeholders := []string{}

    for i, tagID := range tagIDs {
        placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
        values = append(values, articleID, tagID)
    }

    query += strings.Join(placeholders, ", ")
    _, err := tx.ExecContext(ctx, query, values...)
    return err
}

func incrementViewCount(db *sql.DB, articleID int) {
    query := "UPDATE articles SET view_count = view_count + 1 WHERE artikel_id = $1"
    db.Exec(query, articleID)
}

func generateSlug(title string) string {
    // Convert to lowercase
    slug := strings.ToLower(title)
    
    // Remove special characters and replace spaces
    var result strings.Builder
    for _, r := range slug {
        if unicode.IsLetter(r) || unicode.IsDigit(r) {
            result.WriteRune(r)
        } else if unicode.IsSpace(r) || r == '-' {
            result.WriteRune('-')
        }
    }
    
    // Clean up multiple dashes and trim
    slug = result.String()
    slug = strings.ReplaceAll(slug, "--", "-")
    slug = strings.Trim(slug, "-")
    
    // Limit length
    if len(slug) > 100 {
        slug = slug[:100]
        slug = strings.Trim(slug, "-")
    }
    
    return slug
}