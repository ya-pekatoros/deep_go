package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}
type PaymentService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	constructors map[string]reflect.Value
	singletons   map[string]reflect.Value
	instances    map[string]interface{}
}

func NewContainer() *Container {
	return &Container{
		constructors: make(map[string]reflect.Value),
		singletons:   make(map[string]reflect.Value),
		instances:    make(map[string]interface{}),
	}
}

func (c *Container) registerType(name string, constructor interface{}) {
	value := reflect.ValueOf(constructor)
	if !isValidConstructor(value) {
		return
	}

	c.constructors[name] = value
	delete(c.singletons, name)
	delete(c.instances, name)
}

func (c *Container) registerSingletonType(name string, constructor interface{}) {
	value := reflect.ValueOf(constructor)
	if !isValidConstructor(value) {
		return
	}

	c.singletons[name] = value
	delete(c.constructors, name)
	delete(c.instances, name)
}

func (c *Container) resolve(name string) (interface{}, error) {
	if instance, ok := c.instances[name]; ok {
		return instance, nil
	}

	if constructor, ok := c.singletons[name]; ok {
		instance := callConstructor(constructor)
		c.instances[name] = instance
		return instance, nil
	}

	if constructor, ok := c.constructors[name]; ok {
		return callConstructor(constructor), nil
	}

	return nil, fmt.Errorf("тип %q не зарегистрирован", name)
}

// Регистрация типа в контейнере, если ранее тип был зарегистрирован как синглон, то он будет перерегистрирован как обычный тип.
func Register[T any](c *Container, constructor func() T) {
	c.registerType(typeName[T](), constructor)
}

// Регистрация синглтона в контейнере, если ранее тип был зарегистрирован как обычный тип, то он будет перерегистрирован как синглон.
// Если вызывается повторно, то предыдущая регистрация будет перезаписана, а все ранее созданные экземпляры будут удалены.
func RegisterSingleton[T any](c *Container, constructor func() T) {
	c.registerSingletonType(typeName[T](), constructor)
}

func Resolve[T any](c *Container) (T, error) {
	var zero T

	name := typeName[T]()
	value, err := c.resolve(name)
	if err != nil {
		return zero, err
	}

	result, ok := value.(T)
	if !ok {
		return zero, fmt.Errorf("тип %q вернул значение неожиданного типа", name)
	}

	return result, nil
}

func typeName[T any]() string {
	return reflect.TypeFor[T]().String()
}

func isValidConstructor(value reflect.Value) bool {
	if value.Kind() != reflect.Func {
		return false
	}

	constructorType := value.Type()
	return constructorType.NumIn() == 0 && constructorType.NumOut() == 1
}

func callConstructor(constructor reflect.Value) interface{} {
	result := constructor.Call(nil)
	return result[0].Interface()
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	Register(container, func() *UserService {
		return &UserService{}
	})
	Register(container, func() *MessageService {
		return &MessageService{}
	})

	userService1, err := Resolve[*UserService](container)
	assert.NoError(t, err)
	userService2, err := Resolve[*UserService](container)
	assert.NoError(t, err)

	assert.False(t, userService1 == userService2)

	messageService, err := Resolve[*MessageService](container)
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := Resolve[*PaymentService](container)
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}

func TestDISingletonContainer(t *testing.T) {
	container := NewContainer()
	RegisterSingleton(container, func() *UserService {
		return &UserService{}
	})

	userService1, err := Resolve[*UserService](container)
	assert.NoError(t, err)
	userService2, err := Resolve[*UserService](container)
	assert.NoError(t, err)

	assert.True(t, userService1 == userService2)
}
