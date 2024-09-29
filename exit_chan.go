package gotk

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func ExitChan(errCh chan error, shutdown func() error) (err error) {
	var (
		count  int
		logger *slog.Logger
		sigCh  chan os.Signal
	)

	logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

	sigCh = make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	syncErrs := func(err error, count int) error {
		for i := 0; i < count; i++ {
			err = errors.Join(err, <-errCh)
		}
		return err
	}

	select {
	case e := <-errCh:
		logger.Error("... received from channel errch", "error", e)
		err = errors.Join(err, e)
		count -= 1
	case sig := <-sigCh:
		fmt.Println()
		logger.Info("... received from channel quit", "signal", sig.String())
		err = errors.Join(err, shutdown())
	}

	err = syncErrs(err, count)

	return err
}
