package interfaces

import (
	"context"
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/custom_subscription/pkg/models"
)

type CustomSubscriptionRepository interface {
	GetLayersForAuthor(ctx context.Context, authorID string) ([]*models.SubscriptionLayer, error)
	GetFreeLayersForAuthor(ctx context.Context, authorID string) ([]*models.SubscriptionLayer, error)

	GetUserRole(ctx context.Context, userID string) (string, error)
	GetCustomSubscription(ctx context.Context, authorID string) ([]*models.CustomSubscription, error)

	CreateCustomSub(ctx context.Context, userID, title, description string, layer, cost int) error

	SearchAuthor(ctx context.Context, searchTerm string, limit, offset int) ([]string, error)

	GetCustomSubscriptionByTitle(ctx context.Context, authorID string, title string) (*models.CustomSubscription, error)
}
