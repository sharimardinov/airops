// map_errors.go
package usecase

import (
	"airops/internal/domain"
	"context"
	"errors"
)

func MapStoreErr(err error) error {
	if err == nil {
		return nil
	}

	// Уже доменные — не трогаем
	if errors.Is(err, domain.ErrNotFound) || errors.Is(err, domain.ErrInvalidArgument) {
		return err
	}

	// Нормальная история для сервисного слоя: пробрасываем как есть,
	// но в будущем тут можно маппить таймауты/отмену на свои доменные ошибки.
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	return err
}
