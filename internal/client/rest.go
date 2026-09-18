package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type creatable[P any] interface {
	CreatePayload() P
	fmt.Stringer
}

func getResourceList[T any, E any](c *ApiClient, endpoint string, unwrap func(E) T) (*[]T, error) {
	var envelopes []E

	// fetch list of envelopes for resource
	if err := c.do(context.Background(), http.MethodGet, endpoint, nil, &envelopes); err != nil {
		return nil, fmt.Errorf("getting %s: %w", endpoint, err)
	}

	// unwrap resources from envelopes
	items := make([]T, len(envelopes))
	for i, e := range envelopes {
		items[i] = unwrap(e)
	}

	return &items, nil
}

func findFirst[T any](items []T, match func(T) bool) (*T, bool) {
	// find first matching item from list
	for _, item := range items {
		if match(item) {
			return &item, true
		}
	}

	return nil, false
}

func create[T creatable[P], E any, P any](c *ApiClient, endpoint string, resource T, unwrap func(E) T) (*T, error) {
	var envelope E

	// try creating the resource
	if err := c.do(context.Background(), http.MethodPost, endpoint, resource.CreatePayload(), &envelope); err != nil {
		return nil, fmt.Errorf("creating %s: %w", resource, err)
	}

	// unwrap enveloped response
	result := unwrap(envelope)

	return &result, nil
}

func ensure[T creatable[P], P any](resource T, get func() (*T, error), create func() (*T, error)) (*T, error) {
	var apiErr *ApiError

	result, err := get()
	if err != nil {
		if !errors.As(err, &apiErr) {
			return nil, fmt.Errorf("getting %s: %w", resource, err)
		}

		switch apiErr.StatusCode {
		case http.StatusNotFound:
			return create()
		default:
			return nil, fmt.Errorf("getting %s: %w", resource, err)
		}
	}

	return result, nil
}
