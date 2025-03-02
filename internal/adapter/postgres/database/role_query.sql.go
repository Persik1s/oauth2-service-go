// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: role_query.sql

package database

import (
	"context"
)

const getRole = `-- name: GetRole :one
SELECT id, name, level FROM "roles" WHERE name = $1
`

func (q *Queries) GetRole(ctx context.Context, name string) (Role, error) {
	row := q.db.QueryRow(ctx, getRole, name)
	var i Role
	err := row.Scan(&i.ID, &i.Name, &i.Level)
	return i, err
}
