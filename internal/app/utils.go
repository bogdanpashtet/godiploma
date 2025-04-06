package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

func start(ctx context.Context, server *http.Server) error {
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)

		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			errChan <- err
		}
	}()

	const interval = 50 * time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errChan:
			return err
		default:
		}

		conn, err := net.Dial("tcp", server.Addr)
		if err != nil {
			time.Sleep(interval)
			continue
		}

		_ = conn.Close()
		break
	}

	return nil
}
