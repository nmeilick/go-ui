package main

import (
	"github.com/nmeilick/go-ui/input"
	"github.com/nmeilick/go-ui/list"
	"github.com/nmeilick/go-ui/pick"
	"github.com/nmeilick/go-ui/textarea"
)

func main() {
	list.Showcase()
	textarea.Showcase()
	input.Showcase()
	pick.Showcase()
}
