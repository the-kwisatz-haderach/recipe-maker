package authservice

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func AuthDirective(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	tokenData := GetUser(ctx)
	if tokenData == nil {
		return nil, &gqlerror.Error{
			Message: "access denied",
		}
	}
	return next(ctx)
}
