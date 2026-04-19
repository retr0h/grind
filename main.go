// Package main is the grind CLI entry point. The CLI tree lives in the
// `cmd` package; implementation lives under `internal/grind`.
package main

import "github.com/retr0h/grind/cmd"

func main() {
	cmd.Execute()
}
