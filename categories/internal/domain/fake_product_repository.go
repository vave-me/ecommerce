package domain

import (
	"context"
)

type FakeCategoryRepository struct {
	categories map[string]*Category
}

func NewFakeCategoryRepository() *FakeCategoryRepository {
	return &FakeCategoryRepository{categories: map[string]*Category{}}
}

var _ CategoryRepository = (*FakeCategoryRepository)(nil)

func (r *FakeCategoryRepository) Load(ctx context.Context, categoryID string) (*Category, error) {
	if category, exists := r.categories[categoryID]; exists {
		return category, nil
	}

	return NewCategory(categoryID), nil
}

func (r *FakeCategoryRepository) Save(ctx context.Context, category *Category) error {
	for _, event := range category.Events() {
		if err := category.ApplyEvent(event); err != nil {
			return err
		}
	}

	r.categories[category.ID()] = category

	return nil
}

func (r *FakeCategoryRepository) Reset(categories ...*Category) {
	r.categories = make(map[string]*Category)

	for _, category := range categories {
		r.categories[category.ID()] = category
	}
}
