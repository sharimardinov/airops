package service

import (
	"airops/internal/store"
	"errors"
)

func mapStoreErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, store.ErrNotFound) {
		return store.ErrNotFound
	}
	return err
}
