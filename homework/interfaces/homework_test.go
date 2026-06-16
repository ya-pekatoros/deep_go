package main

import (
	"fmt"
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

type registration struct {
	constructor func() any
	singleton   bool
}

type Container struct {
	registrations map[string]registration
	instances     map[string]any
}

func NewContainer() *Container {
	return &Container{
		registrations: make(map[string]registration),
		instances:     make(map[string]any),
	}
}

func (c *Container) registerType(name string, constructor func() any) {
	c.registrations[name] = registration{
		constructor: constructor,
	}
	delete(c.instances, name)
}

func (c *Container) registerSingletonType(name string, constructor func() any) {
	c.registrations[name] = registration{
		constructor: constructor,
		singleton:   true,
	}
	delete(c.instances, name)
}

func (c *Container) resolve(name string) (any, error) {
	if instance, ok := c.instances[name]; ok {
		return instance, nil
	}

	item, ok := c.registrations[name]
	if !ok {
		return nil, fmt.Errorf("type %q is not registered", name)
	}

	instance := item.constructor()
	if item.singleton {
		c.instances[name] = instance
	}

	return instance, nil
}

// Регистрация типа в контейнере, если ранее тип был зарегистрирован как синглон, то он будет перерегистрирован как обычный тип.
func Register[T any](c *Container, constructor func() T) {
	c.registerType(typeName[T](), func() any {
		return constructor()
	})
}

// Регистрация синглтона в контейнере, если ранее тип был зарегистрирован как обычный тип, то он будет перерегистрирован как синглон.
// Если вызывается повторно, то предыдущая регистрация будет перезаписана, а все ранее созданные экземпляры будут удалены.
func RegisterSingleton[T any](c *Container, constructor func() T) {
	c.registerSingletonType(typeName[T](), func() any {
		return constructor()
	})
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
		return zero, fmt.Errorf("type %q returned a value of an unexpected type", name)
	}

	return result, nil
}

func typeName[T any]() string {
	var zero T
	return fmt.Sprintf("%T", zero)
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
