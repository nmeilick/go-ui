package input

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"      // Provides help view for key bindings
	"github.com/charmbracelet/bubbles/key"       // Manages key bindings
	"github.com/charmbracelet/bubbles/textinput" // Provides text input model
	tea "github.com/charmbracelet/bubbletea"     // Framework for building terminal applications
	"github.com/charmbracelet/lipgloss"          // Styles terminal UI components
)

// Model is the model handling user input.
type Model struct {
	textInput textinput.Model // textInput is the text input model.
	help      help.Model      // help is the help model for displaying key bindings.
	keymap    keymap          // keymap is for managing key bindings.
	abort     bool            // abort indicates if the input operation was aborted.
}

type keymap struct{}

// ShortHelp returns a list of key bindings for short help.
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "complete")),
		key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("ctrl+n", "next")),
		key.NewBinding(key.WithKeys("ctrl+p"), key.WithHelp("ctrl+p", "prev")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit")),
	}
}

// FullHelp returns a list of key bindings for full help.
func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

// New creates and returns a new Model with default settings.
func New() *Model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Prompt = ""
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 20
	ti.ShowSuggestions = true
	h := help.New()
	km := keymap{}

	return &Model{
		textInput: ti,
		help:      h,
		keymap:    km,
	}
}

// WithPrompt sets the prompt for the text input model and returns a new Model with the updated prompt.
func (m *Model) WithPrompt(s string) *Model {
	newModel := *m
	newModel.textInput.Prompt = s
	return &newModel
}

// WithPlaceholder sets the placeholder for the text input model and returns a new Model with the updated placeholder.
func (m *Model) WithPlaceholder(s string) *Model {
	newModel := *m
	newModel.textInput.Placeholder = s
	return &newModel
}

// WithPromptStyle sets the style of the prompt for the text input model and returns a new Model with the updated prompt style.
func (m *Model) WithPromptStyle(style lipgloss.Style) *Model {
	newModel := *m
	newModel.textInput.PromptStyle = style
	return &newModel
}

// WithCursorStyle sets the style of the cursor for the text input model and returns a new Model with the updated cursor style.
func (m *Model) WithCursorStyle(style lipgloss.Style) *Model {
	newModel := *m
	newModel.textInput.Cursor.Style = style
	return &newModel
}

// WithCharLimit sets the maximum allowed number of input characters and returns a new Model with the updated
// character limit.
func (m *Model) WithCharLimit(n int) *Model {
	newModel := *m
	newModel.textInput.CharLimit = n
	return &newModel
}

// WithWidth sets the width of the text input model and returns a new Model with the updated width.
func (m *Model) WithWidth(n int) *Model {
	newModel := *m
	newModel.textInput.Width = n
	return &newModel
}

// WithSuggestion sets the autocomplete suggestions for the text input model and returns a new Model with the
// updated suggestions.
func (m *Model) WithSuggestion(suggestions []string) *Model {
	newModel := *m
	newModel.textInput.SetSuggestions(suggestions)
	return &newModel
}

// Init initializes the Model, resets the abort flag, and returns a nil command.
func (m *Model) Init() tea.Cmd {
	m.abort = false
	return nil
}

// Update handles user input and updates the input state by processing key messages and updating the text input model
// accordingly.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, tea.Quit
		case "esc", "ctrl+c":
			m.abort = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the input widget as a string, displaying the prompt, text input, and help view for key bindings.
func (m *Model) View() string {
	return fmt.Sprintf(
		"Pick a Charm™ repo:\n\n  %s\n\n%s\n\n",
		m.textInput.View(),
		m.help.View(m.keymap),
	)
}

// Showcase demonstrates all features of the Model component by creating an input model with autocomplete
// suggestions and running an interactive example in the terminal.
func Showcase() {
	autocomplete := []string{"Apple", "Aardvark", "Banana", "Cherry", "Date", "Elderberry", "Fig", "Grape"}

	defaultStyleInput := New().WithPrompt("Default Style Input: ").WithSuggestion(autocomplete)
	// Run interactive examples
	fmt.Println("=== Model Showcase ===")

	fmt.Println("\nDefault Style Input (Type to see suggestions, Enter to select):")
	p := tea.NewProgram(defaultStyleInput)
	if model, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
	} else {
		if inputModel, ok := model.(*Model); ok {
			fmt.Printf("Final input: %s\n", inputModel.textInput.Value())
		}
	}
}
