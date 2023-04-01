package website

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgconn"
)

type PostgreSQLClassicRepository struct {
	db *sql.DB
}

func NewPostgreSQLClassicRepository(db *sql.DB) *PostgreSQLClassicRepository {
	return &PostgreSQLClassicRepository{
		db: db,
	}
}

func (r *PostgreSQLClassicRepository) Migrate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS websites (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL,
			rank INT NOT NULL
		)
	`)
	return err
}

func (r *PostgreSQLClassicRepository) Create(ctx context.Context, website Website) (*Website, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO websites (name, url, rank) VALUES ($1, $2, $3) RETURNING id
	`, website.Name, website.URL, website.Rank).Scan(&id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" { //23505 is the code for unique_violation
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}
	website.ID = id

	return &website, nil
}

func (r *PostgreSQLClassicRepository) All(ctx context.Context) ([]Website, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM websites`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Website
	for rows.Next() {
		var website Website
		err = rows.Scan(&website.ID, &website.Name, &website.URL, &website.Rank)
		if err != nil {
			return nil, err
		}
		all = append(all, website)
	}

	return all, nil
}

func (r *PostgreSQLClassicRepository) GetByName(ctx context.Context, name string) (*Website, error) {
	var website Website
	err := r.db.QueryRowContext(ctx, `SELECT * FROM websites WHERE name = $1`, name).Scan(&website.ID, &website.Name, &website.URL, &website.Rank)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	return &website, nil
}

func (r *PostgreSQLClassicRepository) Update(ctx context.Context, id int64, update Website) (*Website, error) {
	res, err := r.db.ExecContext(ctx, `UPDATE websites SET name = $1, url = $2, rank = $3 WHERE id = $4`, update.Name, update.URL, update.Rank, id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" { //23505 is the code for unique_violation
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}
	rowsAffected, err := res.RowsAffected() //return the number of rows affected by the query
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &update, nil
}

func (r *PostgreSQLClassicRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM websites WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}
