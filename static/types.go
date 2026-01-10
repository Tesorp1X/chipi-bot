package static

import "strings"

type CallbackAction string

func (a CallbackAction) String() string {
	return string(a)
}

func (a CallbackAction) DataMatches(data string) bool {
	cringePrefix := "\f" + a.String()
	return data == cringePrefix || strings.HasPrefix(data, cringePrefix+"|")
}

const (
	//CallbackAction<name> CallbackAction = "name"
	CallbackActionKeep CallbackAction = "keep"
)
