package ntee

import (
	"fmt"
	"net/http"
	"reflect"
)

// Middleware includes the middleware function and the
// arguments to be passed.
type Middleware struct {
	// Func is the middleware function. The first argument of the
	// function should be a http.Handler and the return value should also be a
	// http.Handler ( middlewareFunc(Compress http.Handler, ...) http.Handler ).
	Function interface{}
	// Options are the arguments to pass to the middleware function.
	Options []interface{}
}

// Composer wraps the inner http.Handler and can be used to
// append middlewares to it.
type Composer struct {
	Handler http.Handler
}

// NewComposer returns a new Composer instance.
func NewComposer(h http.HandlerFunc) *Composer {
	return &Composer{http.Handler(h)}
}

// Use appends the middleware to the inner Compress. Accept the middleware
// function as the first argument and the first argument of the middleware
// function should be a http.Handler and the return value should also be
// http.Handler ( middlewareFunc(Compress http.Handler, ...) http.Handler ).
func (c *Composer) Use(function interface{}, opts ...interface{}) *Composer {
	t := reflect.TypeOf(function)

	if t.NumIn() == 0 || t.In(0).Name() != "Handler" {
		panic(fmt.Sprintf("compose: middleware function: %s"+
			" should take http.Compress as the first argument.", t))
	}

	if t.NumOut() != 1 || t.Out(0).Name() != "Handler" {
		panic(fmt.Sprintf("compose: middleware function: %s"+
			" should return http.Compress as the only return value.", t))
	}

	optionsValues := make([]reflect.Value, len(opts)+1)
	optionsValues[0] = reflect.ValueOf(c.Handler)

	for i, option := range opts {
		optionsValues[i+1] = reflect.ValueOf(option)
	}

	c.Handler = reflect.ValueOf(function).Call(optionsValues)[0].Interface().(http.Handler)

	return c
}
