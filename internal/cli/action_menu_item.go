package cli

import "fmt"

type ActionMenuItem struct {
	description string
	action      func()
}

func (a *ActionMenuItem) Display() {
	fmt.Println(a.description)
}

func (a *ActionMenuItem) Execute() {
	a.action()
}
