package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/products/internal/domain"
)

type AddVariant struct {
	ID           string
	ProductID    string
	SKU          string
	Barcode      string
	Condition    domain.ProductCondition
	VariantPrice int64
	CurrencyCode string
	Stock        int64
	Weight       int64
	Height       int64
	Width        int64
	Depth        int64
	Attributes   []domain.Attribute
	IsAvailable  bool
	HasOptions   bool
	Options      []domain.Option
}

type AddVariantHandler struct {
	variants  domain.VariantRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddVariantHandler(variants domain.VariantRepository, publisher ddd.EventPublisher[ddd.Event]) AddVariantHandler {

	return AddVariantHandler{
		variants:  variants,
		publisher: publisher,
	}
}

func (h AddVariantHandler) AddVariant(ctx context.Context, cmd AddVariant) error {
	variant, err := h.variants.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error adding product")
	}

	event, err := variant.InitVariant(cmd.ID, cmd.ProductID, cmd.SKU, cmd.Attributes, cmd.Condition, cmd.Barcode, cmd.VariantPrice, cmd.CurrencyCode, cmd.Stock, cmd.Weight, cmd.Height, cmd.Width, cmd.Depth, cmd.IsAvailable, cmd.HasOptions, cmd.Options)
	if err != nil {
		return errors.Wrap(err, "initializing product")
	}

	err = h.variants.Save(ctx, variant)
	if err != nil {
		return errors.Wrap(err, "error adding product")
	}

	return errors.Wrap(h.publisher.Publish(ctx, event), "publishing domain event")
}
