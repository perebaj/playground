package website

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type gormWebsite struct {
	ID   int64  `gorm:"primary_key"`
	Name string `gorm:"uniqueIndex, not null"`
	URL  string `gorm:"not null"`
	Rank int64  `gorm:"not null"`
}

func (gormWebsite) TableName() string {
	return "websites"
}

type PostgresSQLGORMRepository struct {
	db *gorm.DB
}

func NewPostgresSQLGORMREpository(db *gorm.DB) *PostgresSQLGORMRepository {
	return &PostgresSQLGORMRepository{
		db: db,
	}
}

func (r *PostgresSQLGORMRepository) Migrate() error {
	m := &gormWebsite{}
	return r.db.AutoMigrate(&m)
}

func (r *PostgresSQLGORMRepository) Create(ctx context.Context, website Website) (*Website, error) {
	m := gormWebsite{
		Name: website.Name,
		URL:  website.URL,
		Rank: website.Rank,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" { //23505 is the code for unique_violation
				return nil, ErrDuplicate
			}
		}

		return nil, err
	}
	result := Website(m)
	return &result, nil
}

func (r *PostgresSQLGORMRepository) All(ctx context.Context) ([]Website, error) {
	var gormWebsite []gormWebsite
	if err := r.db.WithContext(ctx).Find(&gormWebsite).Error; err != nil {
		return nil, err
	}

	var result []Website
	for _, gw := range gormWebsite {
		result = append(result, Website(gw))
	}
	return result, nil
}
