package ui

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea" // Framework for building terminal applications
)

// CanceledError is return when the user decided to cancel input.
var CanceledError = errors.New("canceled")

// QuitError is returned when the user decided to quit the program.
var QuitError = errors.New("quit")

type StandardModel interface {
	Canceled() bool
	Quit() bool
}

func Run(m tea.Model, opts ...tea.ProgramOption) error {
	_, err := tea.NewProgram(m, opts...).Run()
	if m, ok := m.(StandardModel); ok {
		err = ErrorOrValidate(err, m)
	}
	return err
}

func ErrorOrValidate(err error, m StandardModel) error {
	switch {
	case err != nil:
		return err
	case m.Quit():
		return QuitError
	case m.Canceled():
		return CanceledError
	}
	return nil
}
