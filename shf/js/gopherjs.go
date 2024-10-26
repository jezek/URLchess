//go:build js && ecmascript

package js

import "github.com/gopherjs/gopherjs/js"

const WRAPS = "github.com/gopherjs/gopherjs/js"

type Object = js.Object

var Global *Object = js.Global
var Undefined *Object = js.Undefined
