package service

import (
	"context"

	"github.com/devpies/saas-core/internal/project/model"
	"github.com/devpies/saas-core/pkg/web"
)

type PList func(ctx context.Context) ([]model.Project, error)

func forEachT(ctx context.Context, handler PList) ([]model.Project, error) {
	var results []model.Project

	values, ok := web.FromContext(ctx)
	if !ok {
		return nil, web.CtxErr()
	}

	for _, v := range values.TenantMap {
		ctx = web.NewContext(ctx, &web.Values{TenantID: v.TenantID})
		result, err := handler(ctx)
		if err != nil {
			return nil, err
		}
		results = append(results, result...)
	}
	return results, nil
}
