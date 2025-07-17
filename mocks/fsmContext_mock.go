package mocks

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
)

// type Context interface {
// 	// State returns current state for sender.
// 	State(ctx context.Context) (State, error)

// 	// SetState state for sender.
// 	SetState(ctx context.Context, state State) error

// 	// Finish state for sender and deletes data if arg provided.
// 	Finish(ctx context.Context, deleteData bool) error

// 	// Update data in storage. When data argument is nil it must
// 	// delete this item.
// 	Update(ctx context.Context, key string, data any) error

// 	// Data gets from storage and save it into `to` argument.
// 	// Destination argument must be a valid pointer.
// 	Data(ctx context.Context, key string, to any) error
// }

type MockFsmContext struct {
	storage      *MockStorage
	state        fsm.State
	defaultState fsm.State
}

func NewMockFsmContext(s *MockStorage, defaultState fsm.State) *MockFsmContext {
	return &MockFsmContext{storage: s, defaultState: defaultState}
}

// Methods bellow the line are for testing
// --------------------------------------------------------------------

// Methods bellow the line are to satisfy the [fsm.Context] interface
// --------------------------------------------------------------------

func (f *MockFsmContext) State(ctx context.Context) (fsm.State, error) {
	return f.state, nil
}

func (f *MockFsmContext) SetState(ctx context.Context, state fsm.State) error {
	f.state = state
	return nil
}

func (f *MockFsmContext) Finish(ctx context.Context, deleteData bool) error {
	f.state = f.defaultState
	if deleteData {
		f.storage.ClearData()
	}

	return nil
}

func (f *MockFsmContext) Update(ctx context.Context, key string, data any) error {
	f.storage.Set(key, data)
	return nil
}

func (f *MockFsmContext) Data(ctx context.Context, key string, to any) error {
	v := f.storage.Get(key)

	destValue := reflect.ValueOf(to)
	if destValue.Kind() != reflect.Ptr {
		return ErrNotPointer
	}
	if destValue.IsNil() || !destValue.IsValid() {
		return ErrInvalidValue
	}

	destElem := destValue.Elem()
	if !destElem.IsValid() {
		return ErrNotPointer
	}

	destType := destElem.Type()

	vType := reflect.TypeOf(v)
	if !vType.AssignableTo(destType) {
		return &ErrWrongTypeAssign{
			Expect: vType,
			Got:    destType,
		}
	}
	destElem.Set(reflect.ValueOf(v))
	return nil
}

// from original [repo](https://github.com/vitaliy-ukiru/fsm-telebot/blob/v2.x/pkg/storage/memory/errors.go)
var ErrNotPointer = errors.New("fsm/storage/memory: dest argument must be pointer")
var ErrInvalidValue = errors.New("fsm/storage/memory: dest value is nil or invalid")

type ErrWrongTypeAssign struct {
	Expect reflect.Type
	Got    reflect.Type
}

func (e ErrWrongTypeAssign) Error() string {
	return fmt.Sprintf("fsm/storage: wrong types, can't assign %s to %s", e.Expect, e.Got)
}
