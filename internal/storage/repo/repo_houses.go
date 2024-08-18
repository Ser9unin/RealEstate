package repository

import (
	"context"
)

const housesList = `-- name: HousesList :many
SELECT *
FROM houses
ORDER BY house_id
`

func (q *Queries) HousesList(ctx context.Context) ([]House, error) {
	rows, err := q.db.QueryContext(ctx, housesList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []House
	for rows.Next() {
		var i House
		if err := rows.Scan(
			&i.ID,
			&i.Address,
			&i.Year,
			&i.Developer,
			&i.CreateAt,
			&i.UpdateAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const house = `-- name: House :one
SELECT house_id, address, year, developer, created_at, updated_at
FROM houses
WHERE house_id = $1
`

func (q *Queries) House(ctx context.Context, id int) (House, error) {
	row := q.db.QueryRowContext(ctx, house, id)
	var i House
	err := row.Scan(
		&i.ID,
		&i.Address,
		&i.Year,
		&i.Developer,
		&i.CreateAt,
		&i.UpdateAt,
	)
	return i, err
}

// Создание дома:
// Только модератор имеет возможность создать дом используя endpoint /house/create. В случае успешного запроса возвращается полная информация о созданном доме
const newHouse = `-- name: NewHouse :one
INSERT INTO houses(address, year, developer) 
VALUES ($1, $2, $3)
RETURNING house_id, address, year, developer, created_at
`

func (q *Queries) NewHouse(ctx context.Context, arg House) (House, error) {
	row := q.db.QueryRowContext(ctx, newHouse,
		arg.Address,
		arg.Year,
		arg.Developer,
	)
	var i House
	err := row.Scan(
		&i.ID,
		&i.Address,
		&i.Year,
		&i.Developer,
		&i.CreateAt,
	)
	return i, err
}
