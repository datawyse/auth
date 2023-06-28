package app

import (
	"errors"
	"sync"
)

var StopPropagation = errors.New("event hook propagation stopped")

// Handler defines a hook handler function.
type Handler[T any] func(e T) error

// Hook defines a concurrent safe structure for handling event hooks
// (aka. callbacks propagation).
type Hook[T any] struct {
	// sync.RWMutex is a struct with no exported fields.
	mux sync.RWMutex

	// handlers is a slice of Handler[T] functions.
	handlers []Handler[T]
}

// PreAdd registers a new handler to the hook by prepending it to the existing queue.
func (h *Hook[T]) PreAdd(fn Handler[T]) {
	// Lock locks h for writing.
	h.mux.Lock()

	// Unlock unlocks h for writing.
	defer h.mux.Unlock()

	// minimize allocations by shifting the slice
	h.handlers = append(h.handlers, nil)
	// copy(dst, src []T) int copies src to dst and returns the number of elements copied.
	copy(h.handlers[1:], h.handlers)

	// The first element of a slice can be discovered by slicing it from 0 to 1:
	h.handlers[0] = fn
}

// Add registers a new handler to the hook by appending it to the existing queue.
func (h *Hook[T]) Add(fn Handler[T]) {
	// Lock locks h for writing.
	h.mux.Lock()

	// Unlock unlocks h for writing.
	defer h.mux.Unlock()

	// append(s []T, x ...T) []T appends x to s and returns the resulting slice.
	h.handlers = append(h.handlers, fn)
}

// Reset removes all registered handlers.
func (h *Hook[T]) Reset() {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.handlers = nil
}

// Trigger executes all registered hook handlers one by one
// with the specified `data` as an argument.
// Optionally, this method allows also to register additional one off
// handlers that will be temporarily appended to the handlers queue.
// The execution stops when:
// hook.StopPropagation is returned in one of the handlers
// any non-nil error is returned in one of the handlers
func (h *Hook[T]) Trigger(data T, oneOffHandlers ...Handler[T]) error {
	// RLock locks h for reading.
	h.mux.RLock()

	// RUnlock unlocks h for reading. It is a run-time error if h is not locked for reading on entry to RUnlock.
	// It is also a run-time error if h is not locked for reading immediately after RUnlock returns.
	handlers := make([]Handler[T], 0, len(h.handlers)+len(oneOffHandlers))

	// append(s []T, x ...T) []T appends x to s and returns the resulting slice.
	handlers = append(handlers, h.handlers...)

	// append(s []T, x ...T) []T appends x to s and returns the resulting slice.
	handlers = append(handlers, oneOffHandlers...)

	// unlock is not deferred to avoid deadlocks when Trigger is called recursive by the handlers
	h.mux.RUnlock()

	// range iterates over elements of a collection. For arrays and slices, it returns the index and the element value.
	for _, fn := range handlers {
		err := fn(data)
		if err == nil {
			continue
		}
		if errors.Is(err, StopPropagation) {
			return nil
		}
		return err
	}

	return nil
}
