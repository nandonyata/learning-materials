package url

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// URL model for GORM
type URL struct {
	ID          string `gorm:"primaryKey;column:id" json:"id"`
	OriginalURL string `gorm:"column:original_url;not null" json:"url"`
}

// TableName specifies the table name for GORM
func (URL) TableName() string {
	return "url"
}

type ShortenParams struct {
	URL string // the URL to shorten
}

// Shorten shortens a URL.
//
//encore:api public method=POST path=/url
func Shorten(ctx context.Context, p *ShortenParams) (*URL, error) {
	id, err := generateID()
	if err != nil {
		return nil, err
	}

	url := &URL{
		ID:          id,
		OriginalURL: p.URL,
	}

	if err := gormDB.WithContext(ctx).Create(url).Error; err != nil {
		return nil, err
	}

	return &URL{ID: id, OriginalURL: p.URL}, nil
}

// generateID generates a random short ID.
func generateID() (string, error) {
	var data [6]byte // 6 bytes of entropy
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data[:]), nil
}

// Get retrieves the original URL for the id.
//
//encore:api public method=GET path=/url/:id
func Get(ctx context.Context, id string) (*URL, error) {
	var url URL
	err := gormDB.WithContext(ctx).Where("id = ?", id).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

// Database setup - still need Encore's database for connection management
// but without migrations directory
var db = sqldb.NewDatabase("url_db", sqldb.DatabaseConfig{Migrations: "./migrations"})

var gormDB *gorm.DB

// initGORM initializes GORM with Encore's database and runs AutoMigrate
func init() {
	// Get the stdlib database connection from Encore's sqldb
	stdDB := db.Stdlib()

	var err error
	gormDB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: stdDB,
	}), &gorm.Config{})

	if err != nil {
		panic("failed to initialize GORM: " + err.Error())
	}

	// Run GORM AutoMigrate to create/update tables
	if err := gormDB.AutoMigrate(&URL{}); err != nil {
		panic("failed to auto-migrate database: " + err.Error())
	}
}
