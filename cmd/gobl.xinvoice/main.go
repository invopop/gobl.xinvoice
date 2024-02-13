// Package main provides a CLI interface for the library
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

// build data provided by goreleaser and mage setup
var (
	name    = "gobl.xinvoice"
	version = "dev"
	date    = ""
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := godotenv.Load(".env"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	return root().cmd().ExecuteContext(ctx)
}

func inputFilename(args []string) string {
	if len(args) > 0 && args[0] != "-" {
		return args[0]
	}
	return ""
}
