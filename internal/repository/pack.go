package repository

import (
	"context"
	"re-partners/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PackRepository struct {
	pool *pgxpool.Pool
}

func NewPackRepository(pool *pgxpool.Pool) *PackRepository {
	return &PackRepository{pool: pool}
}

func (r *PackRepository) GetPackSizes(ctx context.Context) ([]model.PackSize, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, size FROM pack_sizes ORDER BY size")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	packs := make([]model.PackSize, 0)
	for rows.Next() {
		var pack model.PackSize
		if err := rows.Scan(&pack.ID, &pack.Size); err != nil {
			return nil, err
		}

		packs = append(packs, pack)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return packs, nil
}

func (r *PackRepository) AddPackSize(ctx context.Context, size int) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO pack_sizes (size) VALUES ($1) ON CONFLICT (size) DO NOTHING", size)
	return err
}

func (r *PackRepository) DeletePackSize(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM pack_sizes WHERE id = $1", id)
	return err
}
