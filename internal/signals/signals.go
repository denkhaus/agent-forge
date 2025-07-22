// Package signals provides utilities for handling OS signals and creating signal-aware contexts.
package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// WithSignalContext creates a context that is cancelled when a signal is received.
func WithSignalContext(parent context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	ctx, cancel := context.WithCancel(parent)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)

	go func() {
		defer signal.Stop(sigChan)
		select {
		case sig := <-sigChan:

			log.Info("Signal received, cancelling context", zap.String("signal", sig.String()))

			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}

// WithInterruptContextFunc is a package-level variable that can be overridden for testing.
var WithInterruptContextFunc = WithInterruptContext

// WithInterruptContext creates a context that is cancelled on SIGINT or SIGTERM.
func WithInterruptContext(parent context.Context) (context.Context, context.CancelFunc) {
	return WithSignalContext(parent, syscall.SIGINT, syscall.SIGTERM)
}

// InterruptContext creates a context that is cancelled on SIGINT or SIGTERM.
func InterruptContext() (context.Context, context.CancelFunc) {
	return WithInterruptContext(context.Background())
}

// NotifyContext creates a context that is cancelled when any of the given signals are received
// This is similar to signal.NotifyContext but with better error handling.
func NotifyContext(parent context.Context, signals ...os.Signal) (context.Context, func()) {
	ctx, cancel := WithSignalContext(parent, signals...)
	return ctx, cancel
}

// WaitForSignal blocks until one of the specified signals is received.
func WaitForSignal(signals ...os.Signal) os.Signal {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)
	defer signal.Stop(sigChan)

	return <-sigChan
}

// WaitForInterrupt blocks until SIGINT or SIGTERM is received.
func WaitForInterrupt() os.Signal {
	return WaitForSignal(syscall.SIGINT, syscall.SIGTERM)
}
