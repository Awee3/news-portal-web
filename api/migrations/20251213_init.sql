-- +goose Up

-- ========================================
-- ARTICLES TABLE - Tambah kolom untuk routing & content management
-- ========================================
ALTER TABLE articles 
  ADD COLUMN IF NOT EXISTS slug VARCHAR(255),
  ADD COLUMN IF NOT EXISTS excerpt TEXT,
  ADD COLUMN IF NOT EXISTS gambar_utama VARCHAR(255),
  ADD COLUMN IF NOT EXISTS penulis VARCHAR(100),
  ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'draft',
  ADD COLUMN IF NOT EXISTS tanggal_dibuat TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  ADD COLUMN IF NOT EXISTS tanggal_diperbarui TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;
 
-- Tambah unique constraint untuk slug
DO $$ 
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint 
    WHERE conname = 'articles_slug_unique'
  ) THEN
    ALTER TABLE articles ADD CONSTRAINT articles_slug_unique UNIQUE (slug);
  END IF;
END $$;

-- Tambah NOT NULL constraint untuk status (setelah kolom ada)
DO $$ 
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns 
    WHERE table_name = 'articles' 
      AND column_name = 'status'
      AND is_nullable = 'YES'
  ) THEN
    UPDATE articles SET status = 'draft' WHERE status IS NULL;
    ALTER TABLE articles ALTER COLUMN status SET NOT NULL;
  END IF;
END $$;

-- ========================================
-- CATEGORIES TABLE - Tambah deskripsi & timestamp
-- ========================================
ALTER TABLE categories 
  ADD COLUMN IF NOT EXISTS deskripsi TEXT,
  ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;

-- ========================================
-- TAGS TABLE - Tambah timestamp
-- ========================================
ALTER TABLE tags 
  ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;

-- ========================================
-- MEDIA TABLE - Tambah tipe media & timestamp
-- ========================================
ALTER TABLE media 
  ADD COLUMN IF NOT EXISTS tipe_media VARCHAR(50),
  ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;

-- ========================================
-- COMMENTS TABLE - Tambah support anonymous comments
-- ========================================
ALTER TABLE comments 
  ADD COLUMN IF NOT EXISTS nama_pengguna VARCHAR(100),
  ADD COLUMN IF NOT EXISTS tanggal_dibuat TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  ADD COLUMN IF NOT EXISTS tanggal_diperbarui TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;

-- Ubah user_id jadi nullable (untuk anonymous comments)
DO $$ 
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns 
    WHERE table_name = 'comments' 
      AND column_name = 'user_id' 
      AND is_nullable = 'NO'
  ) THEN
    ALTER TABLE comments ALTER COLUMN user_id DROP NOT NULL;
  END IF;
END $$;

-- Tambah NOT NULL constraint untuk status (setelah kolom ada)
DO $$ 
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns 
    WHERE table_name = 'comments' 
      AND column_name = 'status'
      AND is_nullable = 'YES'
  ) THEN
    UPDATE comments SET status = 'pending' WHERE status IS NULL;
    ALTER TABLE comments ALTER COLUMN status SET NOT NULL;
  END IF;
END $$;

-- ========================================
-- USERS TABLE - Tambah timestamps
-- ========================================
ALTER TABLE users 
  ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;

-- ========================================
-- INDEXES - Performance optimization
-- ========================================
DO $$ 
BEGIN
  -- Articles indexes
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_articles_slug') THEN
    CREATE INDEX idx_articles_slug ON articles(slug);
  END IF;
  
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_articles_status') THEN
    CREATE INDEX idx_articles_status ON articles(status);
  END IF;
  
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_articles_tanggal_publikasi') THEN
    CREATE INDEX idx_articles_tanggal_publikasi ON articles(tanggal_publikasi DESC);
  END IF;
  
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_articles_user_id') THEN
    CREATE INDEX idx_articles_user_id ON articles(user_id);
  END IF;
  
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_articles_tanggal_dibuat') THEN
    CREATE INDEX idx_articles_tanggal_dibuat ON articles(tanggal_dibuat DESC);
  END IF;
  
  -- Comments indexes
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_comments_artikel_id') THEN
    CREATE INDEX idx_comments_artikel_id ON comments(artikel_id);
  END IF;
  
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_comments_status') THEN
    CREATE INDEX idx_comments_status ON comments(status);
  END IF;
  
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_comments_user_id') THEN
    CREATE INDEX idx_comments_user_id ON comments(user_id);
  END IF;
  
  -- Categories index
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_categories_nama') THEN
    CREATE INDEX idx_categories_nama ON categories(nama_kategori);
  END IF;
  
  -- Tags index
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_tags_nama') THEN
    CREATE INDEX idx_tags_nama ON tags(nama_tag);
  END IF;
  
  -- Media index
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_media_artikel_id') THEN
    CREATE INDEX idx_media_artikel_id ON media(artikel_id);
  END IF;
END $$;

-- ========================================
-- TRIGGERS - Auto-update timestamps
-- ========================================

-- Function untuk update tanggal_diperbarui
CREATE OR REPLACE FUNCTION update_tanggal_diperbarui()
RETURNS TRIGGER AS $$
BEGIN
  NEW.tanggal_diperbarui = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk articles
DROP TRIGGER IF EXISTS trg_articles_update ON articles;
CREATE TRIGGER trg_articles_update
  BEFORE UPDATE ON articles
  FOR EACH ROW
  EXECUTE FUNCTION update_tanggal_diperbarui();

-- Trigger untuk comments
DROP TRIGGER IF EXISTS trg_comments_update ON comments;
CREATE TRIGGER trg_comments_update
  BEFORE UPDATE ON comments
  FOR EACH ROW
  EXECUTE FUNCTION update_tanggal_diperbarui();

-- Function untuk update updated_at di users
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk users
DROP TRIGGER IF EXISTS trg_users_update ON users;
CREATE TRIGGER trg_users_update
  BEFORE UPDATE ON users
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();

-- ========================================
-- DATA VALIDATION CONSTRAINTS
-- ========================================

-- Status validation untuk articles
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint 
    WHERE conname = 'articles_status_check'
  ) THEN
    ALTER TABLE articles 
      ADD CONSTRAINT articles_status_check 
      CHECK (status IN ('draft', 'published', 'archived'));
  END IF;
END $$;

-- Status validation untuk comments
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint 
    WHERE conname = 'comments_status_check'
  ) THEN
    ALTER TABLE comments 
      ADD CONSTRAINT comments_status_check 
      CHECK (status IN ('pending', 'approved', 'rejected'));
  END IF;
END $$;

-- +goose Down

-- Drop triggers
DROP TRIGGER IF EXISTS trg_users_update ON users;
DROP TRIGGER IF EXISTS trg_comments_update ON comments;
DROP TRIGGER IF EXISTS trg_articles_update ON articles;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at();
DROP FUNCTION IF EXISTS update_tanggal_diperbarui();

-- Drop constraints
ALTER TABLE comments DROP CONSTRAINT IF EXISTS comments_status_check;
ALTER TABLE articles DROP CONSTRAINT IF EXISTS articles_status_check;

-- Drop indexes
DROP INDEX IF EXISTS idx_media_artikel_id;
DROP INDEX IF EXISTS idx_tags_nama;
DROP INDEX IF EXISTS idx_categories_nama;
DROP INDEX IF EXISTS idx_comments_user_id;
DROP INDEX IF EXISTS idx_comments_status;
DROP INDEX IF EXISTS idx_comments_artikel_id;
DROP INDEX IF EXISTS idx_articles_tanggal_dibuat;
DROP INDEX IF EXISTS idx_articles_tanggal_publikasi;
DROP INDEX IF EXISTS idx_articles_status;
DROP INDEX IF EXISTS idx_articles_slug;
-- Skip idx_articles_user_id (might be from init migration)

-- Revert user_id nullable di comments
DO $$ 
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns 
    WHERE table_name = 'comments' 
      AND column_name = 'user_id' 
      AND is_nullable = 'YES'
  ) THEN
    UPDATE comments SET user_id = 1 WHERE user_id IS NULL;
    ALTER TABLE comments ALTER COLUMN user_id SET NOT NULL;
  END IF;
END $$;

-- Drop columns from users
ALTER TABLE users 
  DROP COLUMN IF EXISTS updated_at,
  DROP COLUMN IF EXISTS created_at;

-- Drop columns from comments
ALTER TABLE comments 
  DROP COLUMN IF EXISTS tanggal_diperbarui,
  DROP COLUMN IF EXISTS tanggal_dibuat,
  DROP COLUMN IF EXISTS nama_pengguna;

-- Drop columns from media
ALTER TABLE media 
  DROP COLUMN IF EXISTS created_at,
  DROP COLUMN IF EXISTS tipe_media;

-- Drop columns from tags
ALTER TABLE tags 
  DROP COLUMN IF EXISTS created_at;

-- Drop columns from categories
ALTER TABLE categories 
  DROP COLUMN IF EXISTS created_at,
  DROP COLUMN IF EXISTS deskripsi;

-- Drop unique constraint & columns from articles
ALTER TABLE articles DROP CONSTRAINT IF EXISTS articles_slug_unique;

ALTER TABLE articles 
  DROP COLUMN IF EXISTS tanggal_diperbarui,
  DROP COLUMN IF EXISTS tanggal_dibuat,
  DROP COLUMN IF EXISTS status,
  DROP COLUMN IF EXISTS penulis,
  DROP COLUMN IF EXISTS gambar_utama,
  DROP COLUMN IF EXISTS excerpt,
  DROP COLUMN IF EXISTS slug;