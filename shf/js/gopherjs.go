//go:build js && !wasm

package js

import "github.com/gopherjs/gopherjs/js"

type Object = js.Object

var Global *Object = js.Global
var Undefined *Object = js.Undefined
