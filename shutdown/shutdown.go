package shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mephistolie/chefbook-backend-common/log"
)

type Operation func(ctx context.Context) error

func Graceful(ctx context.Context, timeout time.Duration, ops map[string]Operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		log.Log(ctx, log.Event{
			Event:     "shutdown.started",
			Message:   "graceful shutdown started",
			Component: "shutdown",
		})

		timeoutFunc := time.AfterFunc(timeout, func() {
			log.LogWarn(ctx, log.Event{
				Event:     "shutdown.timeout",
				Message:   "graceful shutdown timed out",
				Component: "shutdown",
				Duration:  timeout,
			})
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.LogDebug(ctx, log.Event{
					Event:     "shutdown.operation.started",
					Message:   "shutdown operation started",
					Component: "shutdown",
					Operation: innerKey,
				})
				if err := innerOp(ctx); err != nil {
					log.LogError(ctx, log.Event{
						Event:     "shutdown.operation.failed",
						Message:   "shutdown operation failed",
						Component: "shutdown",
						Operation: innerKey,
					}, err)
					return
				}

				log.Log(ctx, log.Event{
					Event:     "shutdown.operation.completed",
					Message:   "shutdown operation completed",
					Component: "shutdown",
					Operation: innerKey,
				})
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}
