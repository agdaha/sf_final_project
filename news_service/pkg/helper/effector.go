package helper

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Effector func(context.Context) error

func Retry(effector Effector, count int, delay time.Duration) Effector {
	return func(ctx context.Context) error {
		for r := 0; ; r++ {
			err := effector(ctx)
			if err == nil || r >= count {
				return err
			}
			slog.Info(fmt.Sprintf("Неудачная попытка %d, повтор через %v", r+1, delay))

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
