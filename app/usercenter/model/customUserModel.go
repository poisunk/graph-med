package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	cacheGraphMedUsercenterUserEmailPrefix = "cache:graphMedUsercenter:user:email:"
)

func (m *customUserModel) FindOneByEmail(ctx context.Context, email string) (*User, error) {
	graphMedUsercenterUserIdKey := fmt.Sprintf("%s%v", cacheGraphMedUsercenterUserEmailPrefix, email)
	var resp User
	err := m.QueryRowCtx(ctx, &resp, graphMedUsercenterUserIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `email` = ? limit 1", userRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, email)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUserModel) FindOneByUserId(ctx context.Context, userId int64) (*User, error) {
	graphMedUsercenterUserIdKey := fmt.Sprintf("%s%v", cacheGraphMedUsercenterUserIdPrefix, userId)
	var resp User
	err := m.QueryRowCtx(ctx, &resp, graphMedUsercenterUserIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `user_id` = ? limit 1", userRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, userId)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUserModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}
