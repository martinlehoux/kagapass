package status

import (
	"github.com/charmbracelet/lipgloss"
)

type Code int

const (
	None Code = iota
	SuccessCode
	ErrorCode
)

type Status struct {
	code    Code
	message string
}

func Success(message string) Status {
	return Status{code: SuccessCode, message: message}
}

func Error(message string) Status {
	return Status{code: ErrorCode, message: message}
}

func (status Status) Render() string {
	switch status.code {
	case SuccessCode:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#32D74B")).Render(status.message)
	case ErrorCode:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render(status.message)
	case None:
		return ""
	default:
		panic("unknown status code")
	}
}
