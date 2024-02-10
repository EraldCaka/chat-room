package repository

import (
	"context"
	"database/sql"
	"github.com/EraldCaka/chat-room/internal/types"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) types.Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	var lastInsertId int
	query := "INSERT INTO users (username, password, email) VALUES ($1, $2, $3) returning id"
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&lastInsertId)
	if err != nil {
		return &types.User{}, err
	}

	user.ID = lastInsertId
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	u := types.User{}
	query := "SELECT id, email, username, password FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Username, &u.Password)
	if err != nil {
		return &types.User{}, nil
	}

	return &u, nil
}
