package cli

type MenuItem interface {
	Display()
	Execute()
}
