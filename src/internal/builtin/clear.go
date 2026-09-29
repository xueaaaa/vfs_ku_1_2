package builtin

const (
	clearName    = "clear"
	msgConsoleOK = "console cleared"
)

func newClear(out Console) func([]string) (any, error) {
	return func([]string) (any, error) {
		out.Clear()
		out.Print(msgConsoleOK)
		return nil, nil
	}
}
