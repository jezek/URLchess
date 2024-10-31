//go:build js && wasm

package js

import "syscall/js"

const WRAPS = "syscall/js"

type Object = js.Value
type Func = js.Func

var Global func() Object = js.Global

var Undefined func() Object = js.Undefined

var FuncOf func(fn func(this Object, args []Object) any) Func = js.FuncOf

func IsUndefined(o Object) bool { return o.IsUndefined() }
