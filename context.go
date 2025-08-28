package mailstream

import "context"

type contextkey string

const ContextKey contextkey = "mailstream"

func WithContext(ctx context.Context, client *Client) context.Context {
	return context.WithValue(ctx, ContextKey, client)
}

func FromContext(ctx context.Context) *Client {
	client, ok := ctx.Value(ContextKey).(*Client)
	if !ok {
		return nil
	}
	return client
}
