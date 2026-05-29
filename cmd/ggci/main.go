// Copyright (c) 2026 Adam Walkiewicz
// SPDX-License-Identifier: MIT

package main

import (
	"os"

	"github.com/ajwalkiewicz/ggci/internal/app"
)

func main() {
	exitCode := app.Run(os.Args[1:])
	os.Exit(exitCode)
}
