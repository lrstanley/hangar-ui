// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package types

import "sync/atomic"

// NewAtomic returns a new Atomic[T] with a default value of val being set.
func NewAtomic[T any](val T) *Atomic[T] {
	v := &Atomic[T]{}
	return v.Store(val)
}

// Atomic is a wrapper around an atomic.Value, using generics.
type Atomic[T any] struct {
	value atomic.Value
}

// Load returns the value of the Atomic.
func (a *Atomic[T]) Load() T {
	return a.value.Load().(T)
}

// Store sets the value of the Atomic, and returns itself.
func (a *Atomic[T]) Store(val T) *Atomic[T] {
	a.value.Store(val)
	return a
}

// NewAtomicComparable returns a new AtomicComparable with a default value of
// val being set.
func NewAtomicComparable[T comparable](val T) *AtomicComparable[T] {
	v := &AtomicComparable[T]{}
	return v.Store(val)
}

// AtomicComparable is a wrapper around an atomic.Value, using generics, including
// some useful helpers for comparing values.
type AtomicComparable[T comparable] struct {
	value atomic.Value
}

// Load returns the value of the Atomic.
func (a *AtomicComparable[T]) Load() T {
	return a.value.Load().(T)
}

// Store sets the value of the Atomic, and returns itself.
func (a *AtomicComparable[T]) Store(val T) *AtomicComparable[T] {
	a.value.Store(val)
	return a
}

// Is returns true if the Atomic's value is equal to one of the given values.
func (a *AtomicComparable[T]) Is(val ...T) bool {
	v := a.Load()
	for _, vv := range val {
		if v == vv {
			return true
		}
	}
	return false
}
