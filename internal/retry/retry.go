package retry

import (
	"fmt"
	"time"
)

func Retry(fn func() error, maxRetries int) error {
	for i := range maxRetries {
		if err := fn(); err != nil {
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		return nil
	}
	return fmt.Errorf("failed to execute function after %d attempts", maxRetries)
}
