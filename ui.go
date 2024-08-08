package ui

import "errors"

var CanceledError = errors.New("canceled")
var QuitError = errors.New("quit")

type StandardModel interface {
	Canceled() bool
	Quit() bool
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
