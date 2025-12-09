package main

import (
	"os"

	"github.com/gookit/slog"
	"github.com/kstsm/wb-sales-tracker/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}
