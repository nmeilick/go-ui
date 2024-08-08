# Terminal UI Components

This repository contains a collection of simple terminal UI components built using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework.

## Components

### Input

The `input` package provides a text input model with features like autocomplete suggestions.

#### Example Usage

```go
package main

import (
	"fmt"
	"github.com/nmeilick/go-ui/input"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	autocomplete := []string{"Apple", "Aardvark", "Banana", "Cherry", "Date", "Elderberry", "Fig", "Grape"}
	inputModel := input.New().WithPrompt("Enter a fruit: ").WithSuggestion(autocomplete)

	p := tea.NewProgram(inputModel)
	if model, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
	} else {
		if inputModel, ok := model.(*input.Model); ok {
			fmt.Printf("Final input: %s\n", inputModel.TextInput.Value())
		}
	}
}
```

### List

The `list` package provides a list model for displaying and selecting items.

#### Example Usage

```go
package main

import (
	"fmt"
	"github.com/nmeilick/go-ui/list"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	items := list.Items{
		&list.Item{Title: "Apple", Desc: "A sweet red fruit"},
		&list.Item{Title: "Banana", Desc: "A long yellow fruit"},
		&list.Item{Title: "Cherry", Desc: "A small red fruit"},
	}

	listModel := list.New(items...).WithSelectedIndex(0)

	p := tea.NewProgram(listModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
	} else {
		fmt.Printf("Selected item index: %d\n", listModel.SelectedItemIdx())
	}
}
```

### Pick

The `pick` package provides a simple interface for selecting an item from a list.

#### Example Usage

```go
package main

import (
	"fmt"
	"github.com/nmeilick/go-ui/pick"
)

func main() {
	items := []string{"Apple", "Banana", "Cherry"}
	selectedIdx, err := pick.Pick("Select a fruit:", false, 0, items...)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else if selectedIdx >= 0 {
		fmt.Printf("You picked: %s\n", items[selectedIdx])
	} else {
		fmt.Println("Selection was aborted.")
	}
}
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
