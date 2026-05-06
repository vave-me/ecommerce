package domain

import (
	"context"
)

type ReservationRepository interface {
	Load(ctx context.Context, reservationID string) (*Reservation, error)
	Save(ctx context.Context, reservation *Reservation) error
}
