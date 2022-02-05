package boot

import "context"

func NewContext(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return ctx
}

func GetEnv() string {
	return ""
}

func InitAPI(ctx context.Context, config string) error {
	return nil
}
