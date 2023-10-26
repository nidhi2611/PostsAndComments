package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	//api "gitlab.eng.vmware.com/nidhig1/goassignment-postandcomments/pkg/api"
	httpserver "gitlab.eng.vmware.com/nidhig1/goassignment-postandcomments/pkg/httpServer"
	errgrp "golang.org/x/sync/errgroup"
)

var httpAddr = ":3000"

const (
	httpTimeout = 3 * time.Second // timeouts used to protect the server
)

func Run(args []string, _ io.Writer) error {
	fmt.Println("Starting Application")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errGrp, errCtx := errgrp.WithContext(ctx)

	s := httpserver.NewServer(httpAddr, httpTimeout, httpTimeout)

	// http server
	errGrp.Go(func() error {
		return s.Start(errCtx)
	})

	//Graceful shutdown of Server
	errGrp.Go(func() error {
		return handleSignals(errCtx, cancel)
	})
	return errGrp.Wait()
}

func handleSignals(ctx context.Context, cancel context.CancelFunc) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sigCh:
		fmt.Printf("got signal, stopping")
		cancel()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
