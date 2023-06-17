package utils

import (
	"fmt"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Verifier func() error

func NewGetterVerifier[T comparable](getter func() T, expected T) Verifier {
	return func() error {
		if got := getter(); got != expected {
			return fmt.Errorf("expected %v got %v", expected, got)
		}
		return nil
	}
}

func NewSliceGetterVerifier[S ~[]E, E comparable](getter func() S, expected S) Verifier {
	return func() error {
		if got := getter(); !slices.Equal(got, expected) {
			return fmt.Errorf("expected %v got %v", expected, got)
		}
		return nil
	}
}

func NewMapGetterVerifier[M ~map[K]V, K comparable, V comparable](getter func() M, expected M) Verifier {
	return func() error {
		if got := getter(); !maps.Equal(got, expected) {
			return fmt.Errorf("expected %v got %v", expected, got)
		}
		return nil
	}
}

func NewSetterVerifier[T comparable](a *T, setter func(v T), expected T) Verifier {
	return func() error {
		setter(expected)
		if got := *a; got != expected {
			return fmt.Errorf("expected %v got %v", expected, got)
		}
		return nil
	}
}

func NewPointSetterVerifier[T comparable](a **T, setter func(v *T), expected T) Verifier {
	return func() error {
		setter(&expected)
		if got := **a; got != expected {
			return fmt.Errorf("expected %v, got %v", expected, got)
		}
		return nil
	}
}

func NewSliceSetterVerifier[S ~[]E, E comparable](s *S, setter func(s S), expected S) Verifier {
	return func() error {
		setter(expected)
		if got := *s; !slices.Equal(got, expected) {
			return fmt.Errorf("expected %v got %v", expected, got)
		}
		return nil
	}
}

func NewMapSetterVerifier[M ~map[K]V, K comparable, V comparable](m *M, setter func(m M), expected M) Verifier {
	return func() error {
		setter(expected)
		if got := *m; !maps.Equal(got, expected) {
			return fmt.Errorf("expected %v got %v", expected, got)
		}
		return nil
	}
}
