package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// Получение списка квартир по номеру дома:
// Используя endpoint /house/{id}, обычный пользователь и модератор могут получить список квартир по номеру дома.
// Только обычный пользователь увидит все квартиры со статусом модерации approved,
// а модератор — жильё с любым статусом модерации.
const flatsListForAll = `-- name: FlatsList :many
SELECT id, house_id, price, rooms, status
FROM flats
WHERE house_id = $1 AND status = 'approved'
ORDER BY id
`
const flatsListForModerator = `-- name: FlatsList :many
SELECT id, house_id, price, rooms, status
FROM flats
WHERE house_id = $1
ORDER BY id
`

func (q *Queries) FlatsList(ctx context.Context, houseID int, userRole string) ([]Flat, error) {
	var rows *sql.Rows
	var err error
	if userRole == moderator {
		rows, err = q.db.QueryContext(ctx, flatsListForModerator, houseID)
	} else {
		rows, err = q.db.QueryContext(ctx, flatsListForAll, houseID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Flat
	for rows.Next() {
		var i Flat
		if err := rows.Scan(
			&i.ID,
			&i.HouseID,
			&i.Price,
			&i.Rooms,
			&i.Status,
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

const flatForAll = `-- name: Flat :one
SELECT id, house_id, price, rooms, status
FROM flats
WHERE house_id = $1 AND id = $2 AND status = 'approved'
`

const flatForModerator = `-- name: Flat :one
SELECT id, house_id, price, rooms, status
FROM flats
WHERE house_id = $1 AND id = $2
`

func (q *Queries) Flat(ctx context.Context, userRole string, houseID, flatID int) (Flat, error) {
	var row *sql.Row
	var err error
	if userRole == moderator {
		row = q.db.QueryRowContext(ctx, flatForModerator, houseID, flatID)
	} else {
		row = q.db.QueryRowContext(ctx, flatForAll, houseID, flatID)
	}
	var i Flat
	err = row.Scan(
		&i.ID,
		&i.HouseID,
		&i.Price,
		&i.Rooms,
		&i.Status,
	)
	return i, err
}

// Создание квартиры:
// Создать квартиру может любой пользователь, используя endpoint /flat/create.
// При успешном запросе возвращается полная информация о квартире.
// Если жильё успешно создано через endpoint /flat/create, то объявление получает статус модерации created.
// У дома, в котором создали новую квартиру, обновляется дата последнего добавления жилья.
const newFlat = `-- name: NewFlat :one
INSERT INTO flats(house_id, price, rooms, status) 
VALUES ($1, $2, $3, $4)
RETURNING id, house_id, price, rooms, status
`

func (q *Queries) NewFlat(ctx context.Context, arg Flat) (Flat, error) {
	row := q.db.QueryRowContext(ctx, newFlat,
		arg.HouseID,
		arg.Price,
		arg.Rooms,
		"created",
	)
	var i Flat
	err := row.Scan(
		&i.ID,
		&i.HouseID,
		&i.Price,
		&i.Rooms,
		&i.Status,
	)
	return i, err
}

// Модерация квартиры:
// Статус модерации квартиры может принимать одно из четырёх значений: created, approved, declined, on moderation.
// Только модератор может изменить статус модерации квартиры с помощью endpoint /flat/update.
// При успешном запросе возвращается полная информация об обновленной квартире.

const flatStatus = `-- name: UpdateFlatStatus :one
SELECT status, COALESCE(moderator, 'no moderator')
FROM flats
WHERE house_id = $1 AND id = $2
`
const updateFlatStatus = `-- name: UpdateFlatStatus :one
UPDATE flats SET status = $1, moderator = $2
WHERE house_id = $3 AND id = $4
RETURNING id, house_id, price, rooms, status
`

func (q *Queries) UpdateFlatStatus(ctx context.Context, user User, status string, flat Flat) (Flat, error) {
	var row *sql.Row
	var err error
	if user.Role == moderator {
		row = q.db.QueryRowContext(ctx, flatStatus,
			flat.HouseID,
			flat.ID,
		)
	}
	var i Flat
	err = row.Scan(
		&i.Status,
		&i.ModeratorID,
	)
	if err != nil {
		return i, err
	}
	if i.Status != "on moderate" || i.ModeratorID == user.UserID {
		row = q.db.QueryRowContext(ctx, updateFlatStatus,
			status,
			user.UserID,
			flat.HouseID,
			flat.ID,
		)
	} else {
		return Flat{}, fmt.Errorf("you don't have permission for this action")
	}
	err = row.Scan(
		&i.ID,
		&i.HouseID,
		&i.Price,
		&i.Rooms,
		&i.Status,
	)
	return i, err
}
