package helper

import (
	"context"
	"errors"
	"testing"
	"time"
)

var count int

func Effector3Count(ctx context.Context) error {
	count++
	if count <= 3 {
		return errors.New("error")
	} else {
		return nil
	}
}

func EffectorWithOnlyErrors(ctx context.Context) error {
	return errors.New("error")
}

func TestRetry(t *testing.T) {
	t1 := time.Now()
	r := Retry(EffectorWithOnlyErrors, 3, 3*time.Millisecond)
	err := r(context.Background())
	d := time.Since(t1)
	if err == nil {
		t.Error("not error for effector with errors")
	}
	if d < 9*time.Millisecond {
		t.Error("fails duration")
	}
}

func TestRetrySuccess(t *testing.T) {
	t1 := time.Now()
	r := Retry(Effector3Count, 5, 3*time.Millisecond)
	err := r(context.Background())
	d := time.Since(t1)

	if err != nil {
		t.Error("error for success effector")
	}
	if d > 10*time.Millisecond {
		t.Error("fails duration")
	}
}
