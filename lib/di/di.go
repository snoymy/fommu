package di

import (
	"reflect"
)

type dep struct {
    constructor interface{}
    object interface{}
}

type AnyFunc interface{}

type Container struct {
    deps map[string]*dep
}

func New() *Container {
    return &Container{
        deps: make(map[string]*dep),
    }
}

func (ctn *Container) Register(constructor AnyFunc) {
    typeOf := reflect.TypeOf(constructor)

    if typeOf.Kind() != reflect.Func {
        panic("Cannot register non function type constructor.")
    }

    if typeOf.NumOut() == 0 {
        panic("Constructor must return non-void value")
    }

    ctn.deps[typeOf.Out(0).String()] = &dep{constructor: constructor}
}

func (ctn *Container) Resolve(constructor AnyFunc) interface{} {
    typeOf := reflect.TypeOf(constructor)

    if typeOf.Kind() != reflect.Func {
        panic("Cannot register non function type constructor.")
    }
    
    var args []reflect.Value
    for i := 0; i < typeOf.NumIn(); i++ {
        dep := ctn.deps[typeOf.In(i).String()]
        if dep == nil {
            panic("Cannot resolve dependency: " + typeOf.In(i).String())
        }
        obj := dep.object

        if obj == nil {
            resolvedObj := ctn.Resolve(dep.constructor) 
            dep.object = resolvedObj
            obj = resolvedObj
        }

        args = append(args, reflect.ValueOf(obj))
    }

    v := reflect.ValueOf(constructor)
    r := v.Call(args)

    var resolvedObj interface{} = nil
    if len(r) > 0 {
        resolvedObj = r[0].Interface()
    }

    return resolvedObj
}

