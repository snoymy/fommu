package structdi

import (
	"reflect"
	"unsafe"
)

type dep struct {
    constructor interface{}
    object interface{}
}

type Config struct {
    ResourceName string
    Primary bool
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

func (ctn *Container) Register(constructor AnyFunc, args ...any) {
    config := Config{}
    if args != nil && len(args) > 0 {
        if c, ok := args[0].(Config); ok {
            config = c
        } else {
            panic("Unexpected argument, expected Config.")
        }
    }

    typeOf := reflect.TypeOf(constructor)

    if typeOf.Kind() != reflect.Func {
        panic("Cannot register non function type constructor.")
    }

    if typeOf.NumOut() == 0 {
        panic("Constructor must return non-void value")
    }

    typeName := typeOf.Out(0).String()  
    depName := typeName

    if config.ResourceName != "" {
        depName = config.ResourceName
    }

    if ctn.deps[depName] != nil {
        panic("Conflict dependency: " + depName)
    }

    ctn.deps[depName] = &dep{constructor: constructor}
    if config.Primary {
        if ctn.deps[typeName] == nil {
            ctn.deps[typeName] = &dep{constructor: constructor}
        }
    }
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

        if reflect.TypeOf(resolvedObj).Kind() == reflect.Ptr {
            if reflect.Indirect(reflect.ValueOf(resolvedObj)).Kind() == reflect.Struct {
                val := reflect.Indirect(reflect.ValueOf(resolvedObj))
                for i := 0; i < val.NumField(); i++ {
                    field := val.Type().Field(i)
                    fieldValue := val.Field(i)
                    tag, ok := field.Tag.Lookup("injectable")
                    if !ok {
                        continue 
                    }

                    depName := field.Type.String()
                    if tag != "" {
                        depName = tag
                    }

                    dep := ctn.deps[depName]
                    if dep == nil {
                        panic("Cannot resolve dependency: " + depName)
                    }
                    obj := dep.object

                    if obj == nil {
                        resolvedObj := ctn.Resolve(dep.constructor) 
                        dep.object = resolvedObj
                        obj = resolvedObj
                    }

                    // if reflect.TypeOf(obj).Implements(fieldValue.Type()) {
                    //     panic(fmt.Sprintf("Struct %s, Cannot inject type %s into %s.", val.Type().Name(), reflect.TypeOf(obj).String(), fieldValue.Type().String()))
                    // }

                    if fieldValue.CanSet() {
                        fieldValue.Set(reflect.ValueOf(obj))
                    } else {
                        fieldPtr := reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr()))
                        fieldPtr.Elem().Set(reflect.ValueOf(obj))
                    }
                }
            }
        }
    }

    return resolvedObj
}
