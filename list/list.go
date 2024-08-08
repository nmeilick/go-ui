package list

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"  // Provides list model
	tea "github.com/charmbracelet/bubbletea" // Framework for building terminal applications
	"github.com/charmbracelet/lipgloss"      // Styles terminal UI components
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// Item represents an item in the list.
type Item struct {
	title string // title is the title of the list item.
	desc  string // desc is the description of the list item.
}

// Items represents an array of items.
type Items []*Item

// Title returns the title of the list item.
func (i *Item) Title() string { return i.title }

// Description returns the description of the list item.
func (i *Item) Description() string { return i.desc }

// FilterValue returns the value used for filtering the list item.
func (i *Item) FilterValue() string { return i.title }

// Model represents the list model.
type Model struct {
	List        list.Model // List is the list model.
	selectedIdx int        // Selected is the index of the currently selected list item.
}

// New creates and returns a new Model with default settings.
func New(items ...*Item) *Model {
	var listItems []list.Item
	for _, i := range items {
		listItems = append(listItems, i)
	}
	l := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	return &Model{
		List: l,
	}
}

// WithItems sets the list items and returns a new Model with the updated items.
func (m *Model) WithItems(items ...list.Item) *Model {
	newModel := *m
	newModel.List.SetItems(items)
	return &newModel
}

// WithSelectedIndex sets the index of the initially selected item and returns a new Model with the updated selected index.
func (m *Model) WithSelectedIndex(i int) *Model {
	if i < 0 {
		i = 0
	} else if i > len(m.List.Items())-1 {
		i = len(m.List.Items()) - 1
	}
	newModel := *m
	newModel.List.Select(i)
	return &newModel
}

// Init initializes the Model and returns a nil command.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles user input and updates the list state by processing key messages and updating the selected item accordingly.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.selectedIdx = m.List.Index()
			return m, tea.Quit
		case "ctrl+c", "esc":
			m.selectedIdx = -1
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

// SelectedItemIdx returns the index of the currently selected item.
func (m *Model) SelectedItemIdx() int {
	return m.selectedIdx
}

// View renders the list as a string, displaying the list items with their respective styles.
func (m Model) View() string {
	return docStyle.Render(m.List.View())
}

// Showcase demonstrates all features of the Model component by creating a list model with some items and running an interactive example in the terminal.
func Showcase() {
	items := Items{
		&Item{title: "Apple", desc: "A sweet red fruit"},
		&Item{title: "Banana", desc: "A long yellow fruit"},
		&Item{title: "Cherry", desc: "A small red fruit"},
	}

	defaultList := New(items...).WithSelectedIndex(0)
	// Run interactive examples
	fmt.Println("=== List Showcase ===")

	fmt.Println("\nDefault List (Use arrow keys to navigate, Enter to select):")
	p := tea.NewProgram(defaultList)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
	}
	switch idx := defaultList.SelectedItemIdx(); idx {
	case -1:
		fmt.Println("Aborted")
	default:
		item := items[idx]
		fmt.Printf("Selected item: %s (#%d)\n", item.Title(), idx)
	}
}
