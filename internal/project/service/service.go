package service

import (
	"context"

	"github.com/devpies/saas-core/pkg/web"
)

func forEachT[T any](
	ctx context.Context,
	tmap web.TenantConnectionMap,
	handler func(ctx context.Context) ([]T, error),
) ([]T, error) {
	var results []T

	for _, v := range tmap {
		ctx = web.NewContext(ctx, &web.Values{TenantID: v.TenantID})
		result, err := handler(ctx)
		if err != nil {
			return nil, err
		}
		results = append(results, result...)
	}
	return results, nil
}
