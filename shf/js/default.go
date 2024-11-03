//go:build !js || (!ecmascript && !wasm)

package js

const WRAPS = ""

type Object struct{}

// gopherjs: Get returns the object's property with the given key.
// wasm: Get returns the JavaScript property key of object o. It panics if o is not a JavaScript object.
func (o Object) Get(key string) Object { return Object{} }

// gopherjs: Set assigns the value to the object's property with the given key.
// wasm: Set sets the JavaScript property key of value value to ValueOf(key). It panics if o is not a JavaScript object. (Note: Wasm uses generics: func (o Object) Set(key string, value any))
func (o Object) Set(key string, value interface{}) {}

// gopherjs: Delete removes the object's property with the given key.
// wasm: Delete deletes the JavaScript property key of object o. It panics if o is not a JavaScript object.
func (o Object) Delete(key string) {}

// gopherjs: Call calls the object's method with the given name.
// wasm: Call does a JavaScript call to the method m of value v with the given arguments. It panics if v has no method m. The arguments get mapped to JavaScript values according to the ValueOf function. (Note: Wasm uses generics func (o Object) Call(m string, args ...any) Object)
func (o Object) Call(name string, args ...interface{}) Object { return Object{} }

// gopherjs: Bool returns the object converted to bool according to JavaScript type conversions.
// wasm: Bool returns the object o as a bool. It panics if o is not a JavaScript boolean.
func (o Object) Bool() bool { return false }

// gopherjs: String returns the object converted to string according to JavaScript type conversions.
// wasm: String returns the object o as a string. String is a special case because of Go's String method convention. Unlike the other getters, it does not panic if o's Type is not TypeString. Instead, it returns a string of the form "<T>" or "<T: V>" where T is o's type and V is a string representation of o's value.
func (o Object) String() string { return "" }

// gopherjs: Int returns the object converted to int according to JavaScript type conversions (parseInt).
// wasm: Int returns the object o truncated to an int. It panics if o is not a JavaScript number.
func (o Object) Int() int { return 0 }

// gopherjs: Float returns the object converted to float64 according to JavaScript type conversions (parseFloat).
// wasm: Float returns the object o as a float64. It panics if o is not a JavaScript number.
func (o Object) Float() float64 { return 0 }

// gopherjs: No such thing exists in gopherjs.
// wasm: Func is a wrapped Go function to be called by JavaScript.
type Func func(e Object)

// gopherjs: No such thing exists in gopherjs.
// wasm: Release frees up resources allocated for the function. The function must not be invoked after calling Release. It is allowed to call Release while the function is still running
func (c Func) Release() {}

// gopherjs: No such thing exists in gopherjs.
// wasm: FuncOf returns a function to be used by JavaScript. ...
var FuncOf = func(fn func(_ Object, _ []Object) any) Func {
	return nil
}

var Global = func() Object { return Object{} }
var Undefined = func() Object { return Object{} }

func IsUndefined(o Object) bool { return true }
