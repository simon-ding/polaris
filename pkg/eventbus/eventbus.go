package eventbus

import (
	"fmt"
	"polaris/pkg/utils"
	"reflect"
)

type EventBus struct {
	handlers utils.Map[string, []EventHandler]
}

type EventHandler struct {
	callback     reflect.Value
	async bool
}

func New() *EventBus {
	return &EventBus{
		handlers: utils.Map[string, []EventHandler]{},
	}
}

func (e *EventBus) Subscribe(event string, fn any) error{
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("handler must be a function: %v", reflect.TypeOf(fn).Kind())
	}
	if handlers, ok := e.handlers.Load(event); ok {
		handlers = append(handlers, EventHandler{
			callback: reflect.ValueOf(fn),})
		e.handlers.Store(event, handlers)
	} else {
		e.handlers.Store(event, []EventHandler{
			{callback: reflect.ValueOf(fn)},
		})
	}
	return nil
}

func (e *EventBus) SubscribeAsync(event string, fn any) error{
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("handler must be a function: %v", reflect.TypeOf(fn).Kind())
	}

	if handlers, ok := e.handlers.Load(event); ok {
		handlers = append(handlers, EventHandler{
			callback: reflect.ValueOf(fn), async: true,
		})
		e.handlers.Store(event, handlers)
	} else {
		e.handlers.Store(event, []EventHandler{
			{callback: reflect.ValueOf(fn), async: true},
		})
	}
	return nil
}

func (e *EventBus) Publish(event string, args... any) {
	if handlers, ok := e.handlers.Load(event); ok {
		for _, handler := range handlers {
			args1 := reflectArgs(handler,args...)
			if handler.async {
				go handler.callback.Call(args1)
			} else {
				handler.callback.Call(args1)
			}
		}
	}
}

func reflectArgs(handler EventHandler,args... any) []reflect.Value {
	funcType := handler.callback.Type()
	passedArguments := make([]reflect.Value, len(args))
	for i, v := range args {
		if v == nil {
			passedArguments[i] = reflect.New(funcType.In(i)).Elem()
		} else {
			passedArguments[i] = reflect.ValueOf(v)
		}
	}

	return passedArguments
}