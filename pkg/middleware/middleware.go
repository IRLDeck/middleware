package middleware

import "context"

type Middleware func(context.Context, string) (context.Context, error)
