package postgresql

import (
	"context"

	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/pkg/global"
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/pkg/logger"
	_ "github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const (
	getUserLayerSql = `
-- First result set.
select
	layer
from
	subscription_layer 
	join custom_subscription on subscription_layer.subscription_layer_id = custom_subscription.subscription_layer_id
	join subscription on subscription.custom_subscription_id = custom_subscription.custom_subscription_id
where
	subscription.user_id = ?
	and custom_subscription.author_id = ?
;
`

	// getUserRoleSQL - возвращает текстовое имя роли пользователя
	// Input: userId
	// Output: role_name ('Reader' or 'Author' or '')
	getUserRoleSQL = `
select 
	role_default_name
from 
	Role
	join People USING (role_id)
where People.user_id = $1
	`

	// getUserLayerOfAuthor - получение уровня подписки пользователя на определенного автора
	// Input: $1 - userId, $2 - authorId
	// Output: layer (int)
	getUserLayerOfAuthor = `
		select layer
		from Subscription
		join Custom_Subscription USING (custom_subscription_id)
		join Subscription_Layer ON Subscription_Layer.subscription_layer_id = Custom_Subscription.subscription_layer_id
		where Custom_Subscription.author_id = $2
		and Subscription.user_id = $1
		;
`
)

func (cr *ContentRepository) GetUserLayerForAuthor(ctx context.Context, userID string, authorID string) (int, error) {
	op := "internal.content.repository.user.GetUserLayerForAuthor"

	rows, err := cr.db.Query(ctx, getUserLayerSql, userID, authorID)
	if err != nil {
		return 0, errors.Wrap(err, op)
	}

	defer rows.Close()

	var (
		layer int
	)
	for rows.Next() {
		if err = rows.Scan(&layer); err != nil {
			return 0, errors.Wrap(err, op)
		}
		logger.StandardDebugF(ctx, op, "Got layer= %s user= %s author %s", layer, userID, authorID)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		return 0, errors.Wrap(global.ErrServer, op)
	}
	logger.StandardDebugF(ctx, op, "Got layer= %s user= %s author %s", layer, userID, authorID)

	return 0, nil
}

func (cr *ContentRepository) GetUserRole(ctx context.Context, userID string) (string, error) {
	op := "internal.content.repository.user.GetUserRole"

	logger.StandardDebugF(ctx, op, "Want to get user role userID=%v, db = %v", userID, cr.db)

	rows, err := cr.db.Query(ctx, getUserRoleSQL, userID)
	if err != nil {
		return "", errors.Wrap(err, op)
	}

	defer rows.Close()

	var (
		role string
	)
	for rows.Next() {
		if err = rows.Scan(&role); err != nil {
			return "", errors.Wrap(err, op)
		}
		logger.StandardDebugF(ctx, op, "Got layer= %s user= %s", role, userID)
		return role, nil
	}

	return "", nil
}

func (cr *ContentRepository) GetUserLayerOfAuthor(ctx context.Context, userID, authorID string) (int, error) {
	op := "internal.content.repository.post.GetUserLayerOfAuthor"

	logger.StandardDebugF(ctx, op, "Want to get user layer userID=%v, author = %v", userID, authorID)

	rows, err := cr.db.Query(ctx, getUserLayerOfAuthor, userID, authorID)
	if err != nil {
		return 0, errors.Wrap(err, op)
	}

	defer rows.Close()

	var (
		layer int
	)
	for rows.Next() {
		if err = rows.Scan(&layer); err != nil {
			return 0, errors.Wrap(err, op)
		}
		logger.StandardDebugF(ctx, op, "Got layer= %s user= %s", layer, userID)
		return layer, nil
	}

	return 0, nil
}
