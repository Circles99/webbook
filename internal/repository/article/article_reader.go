package article

import (
	"context"
	"webbook/internal/domain"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
}
