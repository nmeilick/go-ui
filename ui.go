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

func Run(m tea.Model) error {
	_, err := tea.NewProgram(m).Run()
	if m, ok := m.(StandardModel); ok {
		err = ErrorOrValidate(err, m)
	}
	return err
}

func ErrorOrValidate(err error, m StandardModel) error {
	switch {
	case err != nil:
		return err
	case m.Canceled():
		return CanceledError
	case m.Quit():
		return QuitError
	}
	return nil
}
