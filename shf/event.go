package shf

import "URLchess/shf/js"

type Event interface {
	Get(key string) js.Object
	Set(key string, value interface{})
	Delete(key string)
	Call(name string, args ...any) js.Object
}

type HashChangeEvent interface {
	Event
	NewURL() string
	OldURL() string
}

type hashChangeEvent struct {
	Event
}

func (hce *hashChangeEvent) NewURL() string {
	return hce.Get("newURL").String()
}
func (hce *hashChangeEvent) OldURL() string {
	return hce.Get("oldURL").String()
}
