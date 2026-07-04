package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/kingdev/newschain/internal/models"
)

type Store struct {
	db *sql.DB
}

func New(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS articles (
		id           TEXT PRIMARY KEY,
		title        TEXT NOT NULL,
		slug         TEXT NOT NULL UNIQUE,
		body         TEXT NOT NULL,
		author       TEXT NOT NULL,
		content_hash TEXT NOT NULL,
		ipfs_cid     TEXT,
		tx_hash      TEXT,
		published_at DATETIME NOT NULL,
		edit_count   INTEGER DEFAULT 0,
		retracted    BOOLEAN DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_articles_author ON articles(author);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("migrate schema: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) InsertArticle(a *models.Article) error {
	_, err := s.db.Exec(
		`INSERT INTO articles
			(id, title, slug, body, author, content_hash, ipfs_cid, tx_hash, published_at, edit_count, retracted)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.ID, a.Title, a.Slug, a.Body, a.Author, a.ContentHash, a.IPFSCid, a.TxHash,
		a.PublishedAt, a.EditCount, a.Retracted,
	)
	return err
}

func (s *Store) GetArticle(id string) (*models.Article, error) {
	row := s.db.QueryRow(
		`SELECT id, title, slug, body, author, content_hash, ipfs_cid, tx_hash, published_at, edit_count, retracted
		 FROM articles WHERE id = ?`, id,
	)
	return scanArticle(row)
}

func (s *Store) GetArticleBySlug(slug string) (*models.Article, error) {
	row := s.db.QueryRow(
		`SELECT id, title, slug, body, author, content_hash, ipfs_cid, tx_hash, published_at, edit_count, retracted
		 FROM articles WHERE slug = ?`, slug,
	)
	return scanArticle(row)
}

func (s *Store) ListArticles(limit, offset int) ([]*models.Article, error) {
	rows, err := s.db.Query(
		`SELECT id, title, slug, body, author, content_hash, ipfs_cid, tx_hash, published_at, edit_count, retracted
		 FROM articles ORDER BY published_at DESC LIMIT ? OFFSET ?`, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Article
	for rows.Next() {
		a := &models.Article{}
		if err := rows.Scan(&a.ID, &a.Title, &a.Slug, &a.Body, &a.Author, &a.ContentHash,
			&a.IPFSCid, &a.TxHash, &a.PublishedAt, &a.EditCount, &a.Retracted); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

func (s *Store) UpdateArticleContent(id, body, contentHash, ipfsCID, txHash string) error {
	_, err := s.db.Exec(
		`UPDATE articles SET body = ?, content_hash = ?, ipfs_cid = ?, tx_hash = ?, edit_count = edit_count + 1
		 WHERE id = ?`,
		body, contentHash, ipfsCID, txHash, id,
	)
	return err
}

func (s *Store) RetractArticle(id string) error {
	_, err := s.db.Exec(`UPDATE articles SET retracted = 1 WHERE id = ?`, id)
	return err
}

func scanArticle(row *sql.Row) (*models.Article, error) {
	a := &models.Article{}
	var publishedAt time.Time
	err := row.Scan(&a.ID, &a.Title, &a.Slug, &a.Body, &a.Author, &a.ContentHash,
		&a.IPFSCid, &a.TxHash, &publishedAt, &a.EditCount, &a.Retracted)
	if err != nil {
		return nil, err
	}
	a.PublishedAt = publishedAt
	return a, nil
}
