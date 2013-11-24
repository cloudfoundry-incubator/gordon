package commands

import (
	"fmt"
	"io"
)

type UI interface {
	Say(string)
	Error(string)
}

type BasicUI struct {
	Writer io.Writer
}

func (ui BasicUI) Say(message string) {
	fmt.Fprintf(ui.Writer, "%s\n", message)
}

func (ui BasicUI) Error(message string) {
	fmt.Fprintf(ui.Writer, "%s\n", message)
}
