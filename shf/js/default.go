//go:build !js && !wasm

package js

type Object struct{}

// Get returns the object's property with the given key.
func (o *Object) Get(key string) *Object { return nil }

// Set assigns the value to the object's property with the given key.
func (o *Object) Set(key string, value interface{}) {}

// Delete removes the object's property with the given key.
func (o *Object) Delete(key string) {}

// Length returns the object's "length" property, converted to int.
func (o *Object) Length() int { return 0 }

// Index returns the i'th element of an array.
func (o *Object) Index(i int) *Object { return nil }

// SetIndex sets the i'th element of an array.
func (o *Object) SetIndex(i int, value interface{}) {}

// Call calls the object's method with the given name.
func (o *Object) Call(name string, args ...interface{}) *Object { return nil }

// Invoke calls the object itself. This will fail if it is not a function.
func (o *Object) Invoke(args ...interface{}) *Object { return nil }

// New creates a new instance of this type object. This will fail if it not a function (constructor).
func (o *Object) New(args ...interface{}) *Object { return nil }

// Bool returns the object converted to bool according to JavaScript type conversions.
func (o *Object) Bool() bool { return false }

// String returns the object converted to string according to JavaScript type conversions.
func (o *Object) String() string { return "" }

// Int returns the object converted to int according to JavaScript type conversions (parseInt).
func (o *Object) Int() int { return 0 }

// Int64 returns the object converted to int64 according to JavaScript type conversions (parseInt).
func (o *Object) Int64() int64 { return 0 }

// Uint64 returns the object converted to uint64 according to JavaScript type conversions (parseInt).
func (o *Object) Uint64() uint64 { return 0 }

// Float returns the object converted to float64 according to JavaScript type conversions (parseFloat).
func (o *Object) Float() float64 { return 0 }

// Interface returns the object converted to interface{}. See table in package comment for details.
func (o *Object) Interface() interface{} { return nil }

// Unsafe returns the object as an uintptr, which can be converted via unsafe.Pointer. Not intended for public use.
func (o *Object) Unsafe() uintptr { return 0 }

var Global *Object
var Undefined *Object
