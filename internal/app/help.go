package app

import (
	"fmt"
	"io"
)

func printHelp(w io.Writer, version string) {
	fmt.Fprintf(w, `jr-stack — JR Stack OPSX Edition (%s)

USAGE
  jr-stack                     Launch interactive TUI
  jr-stack <command> [flags]

COMMANDS
  install      Configure AI coding agents on this machine
  sync         Sync agent configs and skills to current version
  update       Check for available updates
  upgrade      Apply updates to managed tools
  restore      Restore a config backup
  version      Print version

FLAGS
  --help, -h    Show this help

Run 'jr-stack help' for this message.
Documentation: https://github.com/JuanCruzRobledo/jr-stack
`, version)
}
