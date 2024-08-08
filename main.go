package main

import (
	"os"

	"github.com/nmeilick/go-ui/input"
	"github.com/nmeilick/go-ui/list"
	"github.com/nmeilick/go-ui/pick"
)

func main() {
	list.Showcase()
	os.Exit(0)
	input.Showcase()
	pick.Showcase()
}
