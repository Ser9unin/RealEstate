package repository

import (
	"context"
	"database/sql"
)

const userByEmail = `-- name: UserByEmailAndPassword :one
SELECT uuid, hash_pass, role
FROM users
WHERE email = $1
`

func (q *Queries) UserByEmail(ctx context.Context, email string) (User, error) {
	var row *sql.Row
	var err error
	row = q.db.QueryRowContext(ctx, userByEmail, email)

	var i User
	err = row.Scan(
		&i.UserID,
		&i.HashPass,
		&i.Role,
	)
	return i, err
}

const userByID = `-- name: UserByIDAndPassword :one
SELECT uuid, hash_pass, role
FROM users
WHERE uuid = $1
`

func (q *Queries) UserByID(ctx context.Context, userID string) (User, error) {
	var row *sql.Row
	var err error
	row = q.db.QueryRowContext(ctx, userByID, userID)

	var i User
	err = row.Scan(
		&i.UserID,
		&i.HashPass,
		&i.Role,
	)
	return i, err
}

const userByIDAndRole = `-- name: UserByIDAndRole :one
SELECT EXISTS (SELECT 1
FROM users
WHERE uuid = $1 AND role = $2 LIMIT 1)
`

func (q *Queries) UserByIDAndRole(ctx context.Context, userID, userRole string) (bool, error) {
	var exists bool
	row, err := q.db.QueryContext(ctx, userByIDAndRole, userID, userRole)
	if err != nil {
		return false, err
	}
	for row.Next() {
		if err = row.Scan(&exists); err != nil {
			return false, err
		}
	}
	return exists, nil
}

// Создание квартиры:
// Создать квартиру может любой пользователь, используя endpoint /flat/create.
// При успешном запросе возвращается полная информация о квартире.
// Если жильё успешно создано через endpoint /flat/create, то объявление получает статус модерации created.
// У дома, в котором создали новую квартиру, обновляется дата последнего добавления жилья.
const newUser = `-- name: NewUser :one
INSERT INTO users(uuid, email, hash_pass, role) 
VALUES ($1, $2, $3, $4)
RETURNING uuid, role
`

func (q *Queries) NewUser(ctx context.Context, arg User) (User, error) {
	row := q.db.QueryRowContext(ctx, newUser,
		arg.UserID,
		arg.Email,
		arg.HashPass,
		arg.Role,
	)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Role,
	)
	return i, err
}

// Модерация квартиры:
// Статус модерации квартиры может принимать одно из четырёх значений: created, approved, declined, on moderation.
// Только модератор может изменить статус модерации квартиры с помощью endpoint /flat/update.
// При успешном запросе возвращается полная информация об обновленной квартире.
const updateUserRole = `-- name: UpdateUserRole :one
UPDATE users SET role = $1
WHERE uuid = $2
RETURNING uuid, role
`

func (q *Queries) UpdateUserRole(ctx context.Context, userRole, uuid string) (User, error) {
	var row *sql.Row
	var err error
	row = q.db.QueryRowContext(ctx, updateUserRole,
		userRole,
		uuid,
	)
	var i User
	err = row.Scan(
		&i.UserID,
		&i.Role,
	)
	return i, err
}
