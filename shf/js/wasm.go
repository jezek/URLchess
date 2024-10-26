//go:build js && wasm

package js

import "syscall/js"

const WRAPS = "syscall/js"

type Object = js.Value
type Func = js.Func

//func (o *Object) Get(key string) *Object {
//	return &Object{o.Value.Get(key)}
//}

//func (o *Object) Set(key string, value interface{}) {
//	o.Value.Set(key, value)
//}

// Delete is the same.

//func (o *Object) Length() int {
//	return o.Value.Length()
//}
//
//func (o *Object) Index(i int) *Object {
//	return &Object{o.Value.Index(i)}
//}
//
//func (o *Object) SetIndex(i int, value interface{}) {
//	o.Value.SetIndex(i, value)
//}

//func (o *Object) Call(name string, args ...interface{}) *Object {
//	println("Call: ", o.Value.String(), name, args)
//	return &Object{o.Value.Call(name, args...)}
//}

//func (o *Object) Bool() bool {
//	return o.Value.Bool()
//}
//
//func (o *Object) String() string {
//	return o.Value.String()
//}
//
//func (o *Object) Int() int {
//	return o.Value.Int()
//}

//func (o *Object) Int64() int {
//	return o.Value.Int64()
//}
//
//func (o *Object) Uint64() int {
//	return o.Value.Uint64()
//}

//func (o *Object) Float() float64 {
//	return o.Value.Float()
//}

//func (o *Object) Interface() interface{} {
//	return o.Value.Interface()
//}

var Global func() Object = js.Global

var Undefined func() Object = js.Undefined

var FuncOf func(fn func(this Object, args []Object) any) Func = js.FuncOf
