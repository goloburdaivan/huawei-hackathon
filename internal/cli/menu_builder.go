package cli

type MenuBuilder struct {
	rootMenu *SubMenu
	current  *SubMenu
}

func NewMenuBuilder(description string) *MenuBuilder {
	rootMenu := &SubMenu{description: description, items: []MenuItem{}}
	return &MenuBuilder{
		rootMenu: rootMenu,
		current:  rootMenu,
	}
}

func (b *MenuBuilder) AddAction(description string, action func()) *MenuBuilder {
	b.current.items = append(b.current.items, &ActionMenuItem{
		description: description,
		action:      action,
	})
	return b
}

func (b *MenuBuilder) AddSubMenu(description string) *MenuBuilder {
	subMenu := &SubMenu{description: description, items: []MenuItem{}}
	b.current.items = append(b.current.items, subMenu)
	b.current = subMenu
	return b
}

func (b *MenuBuilder) EndSubMenu() *MenuBuilder {
	for _, item := range b.rootMenu.items {
		if submenu, ok := item.(*SubMenu); ok && submenu == b.current {
			b.current = b.rootMenu
			break
		}
	}
	return b
}

func (b *MenuBuilder) Build() *SubMenu {
	return b.rootMenu
}
