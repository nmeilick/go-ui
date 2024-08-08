package pick

import (
	"errors"
	"fmt"
	"os"
	"strings"

	// Provides list model
	tea "github.com/charmbracelet/bubbletea" // Framework for building terminal applications
	"github.com/charmbracelet/lipgloss"      // Styles terminal UI components
	"github.com/nmeilick/go-ui"
)

// Model represents a selectable list component.
type Model struct {
	items             []string       // items is the list of items to select from.
	label             string         // label is the label for the list.
	cancelable        bool           // cancelable determines if selection can be canceled with escape key
	quitable          bool           // quitable determines if execution can be quit via ctrl+c
	selectedIdx       int            // selectedIdx is the index of the currently selected item.
	labelStyle        lipgloss.Style // labelStyle is the style for the label.
	selectedItemStyle lipgloss.Style // selectedItemStyle is the style for the selected item.
	normalItemStyle   lipgloss.Style // normalItemStyle is the style for the normal (unselected) items.
	selectedFormat    string         // selectedFormat is the format string for the selected item.
	normalFormat      string         // normalFormat is the format string for normal (unselected) items.
	horizontal        bool           // horizontal indicates if the items should be displayed horizontally.

	canceled bool // canceled indicates whether the selection was canceled
	quit     bool // quit indicates whether the selection was quit
}

// Canceled returns the canceled flag.
func (m *Model) Canceled() bool {
	return m.canceled
}

// Quit returns the quit flag.
func (m *Model) Quit() bool {
	return m.quit
}

// SelectedIdx returns the index of the selected item.
func (m *Model) SelectedIdx() int {
	return m.selectedIdx
}

// SelectedItem returns the selected item, or an empty string if no selection was performed.
func (m *Model) SelectedItem() string {
	if m.selectedIdx >= 0 && m.selectedIdx < len(m.items) {
		return m.items[m.selectedIdx]
	}
	return ""
}

// New creates and returns a new Model with the given items.
func New(items []string) *Model {
	return &Model{
		items:             items,
		label:             "",
		cancelable:        true,
		quitable:          true,
		selectedIdx:       0,
		labelStyle:        lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true), // Gold
		selectedItemStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")),            // Bright Green
		normalItemStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")),            // White
		selectedFormat:    "►%s◄",
		normalFormat:      " %s ",
		horizontal:        false,

		canceled: false,
		quit:     false,
	}
}

// WithLabel sets the label of the Model and returns a new Model with the updated label.
func (m *Model) WithLabel(label string) *Model {
	newModel := *m
	newModel.label = label
	return &newModel
}

// WithCancel sets the cancelable flag and returns a new Model with the updated flag.
func (m *Model) WithCancel(cancelable bool) *Model {
	newModel := *m
	newModel.cancelable = cancelable
	return &newModel
}

// WithQuit sets the quitable flag and returns a new Model with the updated flag.
func (m *Model) WithQuit(quitable bool) *Model {
	newModel := *m
	newModel.quitable = quitable
	return &newModel
}

// WithLabelStyle sets the style of the label and returns a new Model with the updated label style.
func (m *Model) WithLabelStyle(style lipgloss.Style) *Model {
	newModel := *m
	newModel.labelStyle = style
	return &newModel
}

// WithSelectedIndex sets the index of the initially selected item and returns a new Model with the updated selected index.
func (m *Model) WithSelectedIndex(i int) *Model {
	if i < 0 {
		i = 0
	} else if i > len(m.items)-1 {
		i = len(m.items) - 1
	}
	newModel := *m
	newModel.selectedIdx = i
	return &newModel
}

// WithSelectedItemStyle sets the style of the selected item and returns a new Model with the updated selected item style.
func (m *Model) WithSelectedItemStyle(style lipgloss.Style) *Model {
	newModel := *m
	newModel.selectedItemStyle = style
	return &newModel
}

// WithNormalItemStyle sets the style of the normal (unselected) items and returns a new Model with the updated normal item style.
func (m *Model) WithNormalItemStyle(style lipgloss.Style) *Model {
	newModel := *m
	newModel.normalItemStyle = style
	return &newModel
}

// WithLabelColor sets the color of the label and returns a new Model with the updated label color.
func (m *Model) WithLabelColor(color lipgloss.Color) *Model {
	newModel := *m
	newModel.labelStyle = lipgloss.NewStyle().Foreground(color)
	return &newModel
}

// WithSelectedItemColor sets the color of the selected item and returns a new Model with the updated selected item color.
func (m *Model) WithSelectedItemColor(color lipgloss.Color) *Model {
	newModel := *m
	newModel.selectedItemStyle = lipgloss.NewStyle().Foreground(color)
	return &newModel
}

// WithNormalItemColor sets the color of the normal (unselected) items and returns a new Model with the updated normal item color.
func (m *Model) WithNormalItemColor(color lipgloss.Color) *Model {
	newModel := *m
	newModel.normalItemStyle = lipgloss.NewStyle().Foreground(color)
	return &newModel
}

// WithSelectedFormat sets the format string for the selected item and returns a new Model with the updated selected format.
func (m *Model) WithSelectedFormat(format string) *Model {
	newModel := *m
	newModel.selectedFormat = format
	return &newModel
}

// WithNormalFormat sets the format string for normal (unselected) items and returns a new Model with the updated normal format.
func (m *Model) WithNormalFormat(format string) *Model {
	newModel := *m
	newModel.normalFormat = format
	return &newModel
}

// WithHorizontal sets whether the items should be displayed horizontally and returns a new Model with the updated horizontal setting.
func (m *Model) WithHorizontal(horizontal bool) *Model {
	newModel := *m
	newModel.horizontal = horizontal
	return &newModel
}

// Init initializes the Model and returns a nil command.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles user input and updates the list state by processing key messages and updating the selected index accordingly.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "j", "left":
			m.selectedIdx--
			if m.selectedIdx < 0 {
				m.selectedIdx = len(m.items) - 1
			}
		case "down", "k", "right":
			m.selectedIdx++
			if m.selectedIdx >= len(m.items) {
				m.selectedIdx = 0
			}
		case "enter":
			m.canceled = false
			m.quit = false
			return m, tea.Quit
		case "esc":
			if m.cancelable {
				m.selectedIdx = -1
				m.canceled = true
				m.quit = false
				return m, tea.Quit
			}
		case "ctrl+c":
			if m.quitable {
				m.selectedIdx = -1
				m.canceled = false
				m.quit = true
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

// View renders the list as a string, displaying the label and items with their respective styles.
func (m *Model) View() string {
	var b strings.Builder

	if m.label != "" {
		if m.horizontal {
			fmt.Fprintf(&b, "%s ", m.labelStyle.Render(m.label))
		} else {
			fmt.Fprintf(&b, "%s\n", m.labelStyle.Render(m.label))
		}
	}

	var items []string
	for i, item := range m.items {
		var line string
		var format string
		var style lipgloss.Style
		if i == m.selectedIdx {
			style = m.selectedItemStyle
			format = m.selectedFormat
		} else {
			style = m.normalItemStyle
			format = m.normalFormat
		}
		if !strings.Contains(format, "%s") {
			format += "%s"
		}
		line = fmt.Sprintf(format, style.Render(item))
		items = append(items, line)
	}

	if m.horizontal {
		fmt.Fprint(&b, strings.Join(items, "  "))
	} else {
		fmt.Fprint(&b, strings.Join(items, "\n"))
	}

	return b.String()
}

// Pick asks to pick an item and return its index or an error.
// Use errors.Is(ui.Canceled) or errors.Is(ui.Quit) to determine if the selection
// was canceled or aborting of the program was requested.
func Pick(label string, horizontal bool, idx int, items ...string) (int, error) {
	if len(items) == 0 {
		items = []string{"yes", "no"}
	}
	m := New(items).WithLabel(label).WithSelectedIndex(idx).WithHorizontal(horizontal)
	_, err := tea.NewProgram(m).Run()
	if err = ui.ErrorOrValidate(err, m); err != nil {
		return -1, err
	}
	return m.selectedIdx, nil
}

// Showcase demonstrates all features of the Model component by creating various list models and running interactive examples in the terminal.
func Showcase() {
	items := []string{"Apple", "Banana", "Cherry"}

	handle := func(m *Model) {
		_, err := tea.NewProgram(m).Run()
		err = ui.ErrorOrValidate(err, m)
		switch {
		case errors.Is(err, ui.CanceledError):
			fmt.Println("Canceled")
		case errors.Is(err, ui.QuitError):
			fmt.Println("Quit")
			os.Exit(0)
		case err != nil:
			fmt.Printf("Error running program: %v", err)
		default:
			fmt.Printf("Picked item: %s (Index: %d)\n", m.SelectedItem(), m.SelectedIdx())
		}
	}
	// Run interactive examples
	fmt.Println("=== Model Showcase ===")

	fmt.Println("\nDefault Style List (Use arrow keys to navigate, Enter to select):")
	// Create a vertical list with default style
	defaultStyleList := New(items).WithLabel("Default Style List")
	handle(defaultStyleList)

	fmt.Println("\nHorizontal List with Custom Colors (Use arrow keys to navigate, Enter to select):")
	// Create a horizontal list with custom colors
	horizontalList := New(items).
		WithLabel("Horizontal List").
		WithLabelColor(lipgloss.Color("#FF69B4")).        // Hot Pink
		WithSelectedItemColor(lipgloss.Color("#FF4500")). // OrangeRed
		WithNormalItemColor(lipgloss.Color("#98FB98")).   // PaleGreen
		WithHorizontal(true)
	handle(horizontalList)

	// Create a vertical list with custom formats
	customFormatList := New(items).
		WithLabel("Custom Format List").
		WithSelectedFormat("► %s ◄").
		WithNormalFormat("  %s  ")
	handle(customFormatList)
}
