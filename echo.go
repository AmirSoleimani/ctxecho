package ctxecho

import (
	"reflect"
	"unsafe"
)

const (
	keyFieldname = "key"
	valFieldname = "val"
)

// Context values are associated with a context and may have key, value
// fields associated with it. The value field can be ANY (interface{})
// type as well as the key. Returns a map of key, value pairs. The key
// is the name of the field in the context. The value is the value of
// the field.
func Inspect(ctx any) map[any]any {
	store := make(map[any]any)
	inspect(store, ctx)
	return store
}

func inspect(store map[any]any, ctx any) {
	contextValues := reflect.ValueOf(ctx)
	contextKeys := reflect.TypeOf(ctx)

	if reflect.TypeOf(ctx).Kind() == reflect.Ptr {
		contextValues = contextValues.Elem()
		contextKeys = contextKeys.Elem()
	}

	if contextKeys.Kind() != reflect.Struct {
		return
	}

	iterateCtx(store, contextValues, contextKeys)
}

func iterateCtx(store map[any]any, values reflect.Value, types reflect.Type) {
	var keyName any
	for i := 0; i < values.NumField(); i++ {
		fieldValue := values.Field(i)
		fieldName := types.Field(i).Name

		if fieldName == keyFieldname || fieldName == valFieldname {
			fieldValue = reflectNewAt(fieldValue)

			if fieldName == keyFieldname {
				keyName = fieldValue.Interface()
			} else {
				store[keyName] = fieldValue.Interface()
			}

			continue
		}

		switch fieldName {
		case "Context":
			inspect(store, fieldValue.Interface())
		case "cancelCtx", "timerCtx", "valueCtx":
			fieldValue = reflectNewAt(fieldValue)
			inspect(store, fieldValue.Interface())
		}
	}
}

func reflectNewAt(fieldValue reflect.Value) reflect.Value {
	return reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr())).Elem()
}
