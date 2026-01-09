package middleware

import "context"

type ctxKeyLogger struct{}

func WithLogger(ctx context.Context, lg *JSONLogger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger{}, lg)
}

func LoggerFrom(ctx context.Context) *JSONLogger {
	if v := ctx.Value(ctxKeyLogger{}); v != nil {
		if lg, ok := v.(*JSONLogger); ok && lg != nil {
			return lg
		}
	}
	return NewJSONLogger()
}
