package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/streams/internal/domain"
)

type SetStreamPricing struct {
	StreamID       string
	RentalPrice    int64
	RentalDuration int64 // in hours
	PurchasePrice  int64
	PPVPrice       int64
}

type SetStreamPricingHandler struct {
	streams ddd.AggregateStore[*domain.Stream]
}

func NewSetStreamPricingHandler(streams ddd.AggregateStore[*domain.Stream]) SetStreamPricingHandler {
	return SetStreamPricingHandler{
		streams: streams,
	}
}

func (h SetStreamPricingHandler) SetStreamPricing(ctx context.Context, cmd SetStreamPricing) error {
	stream, err := h.streams.Load(ctx, cmd.StreamID)
	if err != nil {
		return err
	}

	event, err := stream.SetPricing(
		cmd.RentalPrice,
		cmd.RentalDuration,
		cmd.PurchasePrice,
		cmd.PPVPrice,
	)
	if err != nil {
		return err
	}
	stream.AddEvent(event)

	return h.streams.Save(ctx, stream)
}