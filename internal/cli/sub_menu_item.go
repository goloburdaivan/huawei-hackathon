package cli

import "fmt"

type SubMenu struct {
	description string
	items       []MenuItem
}

func (s *SubMenu) Display() {
	fmt.Println(s.description)
}

func (s *SubMenu) Execute() {
	for {
		fmt.Println("\n" + s.description)
		for i, item := range s.items {
			fmt.Printf("%d. ", i+1)
			item.Display()
		}
		fmt.Printf("%d. Назад/Выйти\n", len(s.items)+1)

		var choice int
		fmt.Scanln(&choice)

		if choice == len(s.items)+1 {
			return
		}

		if choice > 0 && choice <= len(s.items) {
			s.items[choice-1].Execute()
		} else {
			fmt.Println("Неверный выбор. Попробуйте снова.")
		}
	}
}
