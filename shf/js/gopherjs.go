//go:build js && ecmascript

package js

import "github.com/gopherjs/gopherjs/js"

const WRAPS = "github.com/gopherjs/gopherjs/js"

type Object = *js.Object
type Func func(e Object)

func (_ Func) Release() {}

var Global = func() Object { return js.Global }
var Undefined = func() Object { return js.Undefined }

var FuncOf = func(fn func(_ Object, _ []Object) any) Func {
	return func(e Object) {
		fn(e, []Object{e})
	}
}

func IsUndefined(o Object) bool { return o == js.Undefined }
