package application

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"middleman/baskets/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/es"
	"testing"
)

func TestApplication_AddItem(t *testing.T) {
	product := &domain.Product{
		ID:           "product-id",
		UserSellerID: "user-seller-id",
		Name:         "product-name",
		BasePrice:    10.00,
	}
	user := &domain.User{
		ID:        "user-id",
		FirstName: "user-first-name",
		LastName:  "user-last-name",
		Email:     "user-email",
	}

	type mocks struct {
		baskets   *domain.MockBasketRepository
		users     *domain.MockUserRepository
		catalog   *domain.MockCatalogRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}
	type args struct {
		ctx context.Context
		add AddItem
	}
	tests := map[string]struct {
		args    args
		on      func(f mocks)
		wantErr bool
	}{
		"Success": {
			args: args{
				ctx: context.Background(),
				add: AddItem{
					ID:        "basket-id",
					ProductID: "product-id",
					Quantity:  1,
				},
			},
			on: func(f mocks) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
					UserCustomerID: "user-customer-id",
					Items:          make(map[string]domain.Item),
					Status:         domain.BasketIsOpen,
				}, nil)
				f.products.On("Find", context.Background(), "product-id").Return(product, nil)
				f.users.On("Find", context.Background(), "user-seller-id").Return(user, nil)
				f.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(nil)
				// we do not specifically test publisher.On("Publish") unless you want to
				f.publisher.On("Publish", context.Background(), mock.AnythingOfType("ddd.event")).Return(nil)
			},
		},
		"NoBasket": {
			args: args{
				ctx: context.Background(),
				add: AddItem{
					ID:        "basket-id",
					ProductID: "product-id",
					Quantity:  1,
				},
			},
			on: func(f mocks) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(nil, fmt.Errorf("no basket"))
			},
			wantErr: true,
		},
		"NoProduct": {
			args: args{
				ctx: context.Background(),
				add: AddItem{
					ID:        "basket-id",
					ProductID: "product-id",
					Quantity:  1,
				},
			},
			on: func(f mocks) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
					UserCustomerID: "user-customer-id",
					Items:          make(map[string]domain.Item),
					Status:         domain.BasketIsOpen,
				}, nil)
				f.products.On("Find", context.Background(), "product-id").Return(nil, fmt.Errorf("no product"))
			},
			wantErr: true,
		},
		"NoUser": {
			args: args{
				ctx: context.Background(),
				add: AddItem{
					ID:        "basket-id",
					ProductID: "product-id",
					Quantity:  1,
				},
			},
			on: func(f mocks) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
					UserCustomerID: "user-customer-id",
					Items:          make(map[string]domain.Item),
					Status:         domain.BasketIsOpen,
				}, nil)
				f.products.On("Find", context.Background(), "product-id").Return(product, nil)
				f.users.On("Find", context.Background(), "user-seller-id").Return(nil, fmt.Errorf("no user seller"))
			},
			wantErr: true,
		},
		"SaveFailed": {
			args: args{
				ctx: context.Background(),
				add: AddItem{
					ID:        "basket-id",
					ProductID: "product-id",
					Quantity:  1,
				},
			},
			on: func(f mocks) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
					UserCustomerID: "user-customer-id",
					Items:          make(map[string]domain.Item),
					Status:         domain.BasketIsOpen,
				}, nil)
				f.products.On("Find", context.Background(), "product-id").Return(product, nil)
				f.users.On("Find", context.Background(), "user-seller-id").Return(user, nil)
				f.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(fmt.Errorf("save failed"))
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			m := mocks{
				baskets:   domain.NewMockBasketRepository(t),
				users:     domain.NewMockUserRepository(t),
				products:  domain.NewMockProductRepository(t),
				publisher: ddd.NewMockEventPublisher[ddd.Event](t),
			}
			a := New(m.baskets, m.users, m.products, domain.NewMockCatalogRepository(t), m.publisher)
			if tt.on != nil {
				tt.on(m)
			}

			if err := a.AddItem(tt.args.ctx, tt.args.add); (err != nil) != tt.wantErr {
				t.Errorf("AddItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_CancelBasket(t *testing.T) {
	type fields struct {
		baskets   *domain.MockBasketRepository
		users     *domain.MockUserRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}
	type args struct {
		ctx    context.Context
		cancel CancelBasket
	}
	tests := map[string]struct {
		args    args
		on      func(f fields)
		wantErr bool
	}{
		"Success": {
			args: args{
				ctx: context.Background(),
				cancel: CancelBasket{
					ID: "basket-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
					UserCustomerID: "user-customer-id",
					Items:          make(map[string]domain.Item),
					Status:         domain.BasketIsOpen,
				}, nil)
				f.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(nil)
				f.publisher.On("Publish", context.Background(), mock.AnythingOfType("ddd.event")).Return(nil)
			},
		},
		"NoBasket": {
			args: args{
				ctx: context.Background(),
				cancel: CancelBasket{
					ID: "basket-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(nil, fmt.Errorf("no basket"))
			},
			wantErr: true,
		},
		"SaveFailed": {
			args: args{
				ctx: context.Background(),
				cancel: CancelBasket{
					ID: "basket-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
					UserCustomerID: "user-customer-id",
					Items:          make(map[string]domain.Item),
					Status:         domain.BasketIsOpen,
				}, nil)
				f.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(fmt.Errorf("save failed"))
			},
			wantErr: true,
		},
		"PublishFailed": {
			args: args{
				ctx: context.Background(),
				cancel: CancelBasket{
					ID: "basket-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
					UserCustomerID: "customer-id",
					Items:          make(map[string]domain.Item),
					Status:         domain.BasketIsOpen,
				}, nil)
				f.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(nil)
				f.publisher.On("Publish", context.Background(), mock.AnythingOfType("ddd.event")).Return(fmt.Errorf("publish failed"))
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := fields{
				baskets:   domain.NewMockBasketRepository(t),
				users:     domain.NewMockUserRepository(t),
				products:  domain.NewMockProductRepository(t),
				publisher: ddd.NewMockEventPublisher[ddd.Event](t),
			}
			a := Application{
				baskets:   f.baskets,
				users:     f.users,
				products:  f.products,
				catalog:   domain.NewMockCatalogRepository(t),
				publisher: f.publisher,
			}
			if tt.on != nil {
				tt.on(f)
			}

			if err := a.CancelBasket(tt.args.ctx, tt.args.cancel); (err != nil) != tt.wantErr {
				t.Errorf("CancelBasket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_CheckoutBasket(t *testing.T) {
	store := &domain.User{
		ID:        "user-id",
		FirstName: "user-first-name",
		LastName:  "user-last-name",
		Email:     "user-email",
	}
	product := &domain.Product{
		ID:           "product-id",
		UserSellerID: "user-seller-id",
		Name:         "product-name",
		BasePrice:    10.00,
	}
	item := domain.Item{
		UserSellerID:   store.ID,
		ProductID:      product.ID,
		UserSellerName: store.FirstName,
		ProductName:    product.Name,
		ProductPrice:   product.BasePrice,
		Quantity:       10,
	}

	type fields struct {
		baskets   *domain.MockBasketRepository
		users     *domain.MockUserRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}
	type args struct {
		ctx      context.Context
		checkout CheckoutBasket
	}
	tests := map[string]struct {
		args    args
		on      func(f fields)
		wantErr bool
	}{
		"Success": {
			args: args{
				ctx: context.Background(),
				checkout: CheckoutBasket{
					ID:              "basket-id",
					PaymentMethodID: "payment-id", // Not used in domain now
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
					UserCustomerID: "user-customer-id",
					Items: map[string]domain.Item{
						product.ID: item,
					},
					Status: domain.BasketIsOpen,
				}, nil)
				f.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(nil)
				f.publisher.On("Publish", context.Background(), mock.AnythingOfType("ddd.event")).Return(nil)
			},
		},
		"NoBasket": {
			args: args{
				ctx: context.Background(),
				checkout: CheckoutBasket{
					ID:              "basket-id",
					PaymentMethodID: "payment-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(nil, fmt.Errorf("no basket"))
			},
			wantErr: true,
		},
		"EmptyBasket": {
			args: args{
				ctx: context.Background(),
				checkout: CheckoutBasket{
					ID:              "basket-id",
					PaymentMethodID: "payment-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
					UserCustomerID: "user-customer-id",
					Items:          make(map[string]domain.Item),
					Status:         domain.BasketIsOpen,
				}, nil)
			},
			wantErr: true, // domain requires at least one item
		},
		"SaveFailed": {
			args: args{
				ctx: context.Background(),
				checkout: CheckoutBasket{
					ID:              "basket-id",
					PaymentMethodID: "payment-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
					UserCustomerID: "user-customer-id",
					Items: map[string]domain.Item{
						product.ID: item,
					},
					Status: domain.BasketIsOpen,
				}, nil)
				f.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(fmt.Errorf("save failed"))
			},
			wantErr: true,
		},
		"PublishFailed": {
			args: args{
				ctx: context.Background(),
				checkout: CheckoutBasket{
					ID:              "basket-id",
					PaymentMethodID: "payment-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
					UserCustomerID: "customer-id",
					Items: map[string]domain.Item{
						product.ID: item,
					},
					Status: domain.BasketIsOpen,
				}, nil)
				f.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(nil)
				f.publisher.On("Publish", context.Background(), mock.AnythingOfType("ddd.event")).Return(fmt.Errorf("publish failed"))
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := fields{
				baskets:   domain.NewMockBasketRepository(t),
				users:     domain.NewMockUserRepository(t),
				products:  domain.NewMockProductRepository(t),
				publisher: ddd.NewMockEventPublisher[ddd.Event](t),
			}
			a := Application{
				baskets:   f.baskets,
				users:     f.users,
				products:  f.products,
				catalog:   domain.NewMockCatalogRepository(t),
				publisher: f.publisher,
			}
			if tt.on != nil {
				tt.on(f)
			}

			if err := a.CheckoutBasket(tt.args.ctx, tt.args.checkout); (err != nil) != tt.wantErr {
				t.Errorf("CheckoutBasket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_GetBasket(t *testing.T) {
	user := &domain.User{
		ID:        "user-id",
		FirstName: "user-first-name",
		LastName:  "user-last-name",
	}
	product := &domain.Product{
		ID:           "product-id",
		UserSellerID: "user-seller-id",
		Name:         "product-name",
		BasePrice:    10.00,
	}
	item := domain.Item{
		UserSellerID:   user.ID,
		ProductID:      product.ID,
		UserSellerName: user.FirstName,
		ProductName:    product.Name,
		ProductPrice:   product.BasePrice,
		Quantity:       10,
	}

	type fields struct {
		baskets   *domain.MockBasketRepository
		users     *domain.MockUserRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}
	type args struct {
		ctx context.Context
		get GetBasket
	}
	tests := map[string]struct {
		args    args
		on      func(f fields)
		want    *domain.Basket
		wantErr bool
	}{
		"GetBasket": {
			args: args{
				ctx: context.Background(),
				get: GetBasket{
					ID: "basket-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
					UserCustomerID: "user-customer-id",
					Items: map[string]domain.Item{
						product.ID: item,
					},
					Status: domain.BasketIsOpen,
				}, nil)
			},
			want: &domain.Basket{
				Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
				UserCustomerID: "user-customer-id",
				Items: map[string]domain.Item{
					product.ID: item,
				},
				Status: domain.BasketIsOpen,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := fields{
				baskets:   domain.NewMockBasketRepository(t),
				users:     domain.NewMockUserRepository(t),
				products:  domain.NewMockProductRepository(t),
				publisher: ddd.NewMockEventPublisher[ddd.Event](t),
			}
			a := Application{
				baskets:   f.baskets,
				users:     f.users,
				products:  f.products,
				catalog:   domain.NewMockCatalogRepository(t),
				publisher: f.publisher,
			}
			if tt.on != nil {
				tt.on(f)
			}

			got, err := a.GetBasket(tt.args.ctx, tt.args.get)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBasket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Compare the relevant fields
			assert.Equal(t, tt.want.UserCustomerID, got.UserCustomerID)
			assert.Equal(t, tt.want.Items, got.Items)
			assert.Equal(t, tt.want.Status, got.Status)
		})
	}
}

func TestApplication_RemoveItem(t *testing.T) {
	user := &domain.User{
		ID:        "user-id",
		FirstName: "user-first-name",
		LastName:  "user-last-name",
		Email:     "user-email",
	}
	product := &domain.Product{
		ID:           "product-id",
		UserSellerID: "user-seller-id",
		Name:         "product-name",
		BasePrice:    10.00,
	}
	item := domain.Item{
		UserSellerID:   user.ID,
		ProductID:      product.ID,
		UserSellerName: user.FirstName,
		ProductName:    product.Name,
		ProductPrice:   product.BasePrice,
		Quantity:       10,
	}

	type mocks struct {
		baskets   *domain.MockBasketRepository
		users     *domain.MockUserRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}
	type args struct {
		ctx    context.Context
		remove RemoveItem
	}
	tests := map[string]struct {
		args    args
		on      func(m mocks)
		wantErr bool
	}{
		"Success": {
			args: args{
				ctx: context.Background(),
				remove: RemoveItem{
					ID:        "basket-id",
					ProductID: product.ID,
					Quantity:  1,
				},
			},
			on: func(m mocks) {
				m.products.On("Find", context.Background(), product.ID).Return(product, nil)
				m.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
					UserCustomerID: "user-customer-id",
					Items: map[string]domain.Item{
						product.ID: item,
					},
					Status: domain.BasketIsOpen,
				}, nil)
				m.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(nil)
				m.publisher.On("Publish", context.Background(), mock.AnythingOfType("ddd.event")).Return(nil)
			},
		},
		"NoProduct": {
			args: args{
				ctx: context.Background(),
				remove: RemoveItem{
					ID:        "basket-id",
					ProductID: product.ID,
					Quantity:  1,
				},
			},
			on: func(m mocks) {
				m.products.On("Find", context.Background(), product.ID).Return(nil, fmt.Errorf("no product"))
			},
			wantErr: true,
		},
		"NoBasket": {
			args: args{
				ctx: context.Background(),
				remove: RemoveItem{
					ID:        "basket-id",
					ProductID: product.ID,
					Quantity:  1,
				},
			},
			on: func(m mocks) {
				m.products.On("Find", context.Background(), product.ID).Return(product, nil)
				m.baskets.On("Load", context.Background(), "basket-id").Return(nil, fmt.Errorf("no basket"))
			},
			wantErr: true,
		},
		"SaveFailed": {
			args: args{
				ctx: context.Background(),
				remove: RemoveItem{
					ID:        "basket-id",
					ProductID: product.ID,
					Quantity:  1,
				},
			},
			on: func(m mocks) {
				m.products.On("Find", context.Background(), product.ID).Return(product, nil)
				m.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
					UserCustomerID: "customer-id",
					Items: map[string]domain.Item{
						product.ID: item,
					},
					Status: domain.BasketIsOpen,
				}, nil)
				m.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(fmt.Errorf("save failed"))
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			m := mocks{
				baskets:   domain.NewMockBasketRepository(t),
				users:     domain.NewMockUserRepository(t),
				products:  domain.NewMockProductRepository(t),
				publisher: ddd.NewMockEventPublisher[ddd.Event](t),
			}
			a := Application{
				baskets:   m.baskets,
				users:     m.users,
				products:  m.products,
				catalog:   domain.NewMockCatalogRepository(t),
				publisher: m.publisher,
			}
			if tt.on != nil {
				tt.on(m)
			}

			if err := a.RemoveItem(tt.args.ctx, tt.args.remove); (err != nil) != tt.wantErr {
				t.Errorf("RemoveItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_StartBasket(t *testing.T) {
	type mocks struct {
		baskets   *domain.MockBasketRepository
		users     *domain.MockUserRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}
	type args struct {
		ctx   context.Context
		start StartBasket
	}
	tests := map[string]struct {
		args    args
		on      func(m mocks)
		wantErr bool
	}{
		"Success": {
			args: args{
				ctx: context.Background(),
				start: StartBasket{
					ID:             "basket-id",
					UserCustomerID: "user-customer-id",
				},
			},
			on: func(m mocks) {
				m.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
					UserCustomerID: "",
					Items:          make(map[string]domain.Item),
					Status:         "",
				}, nil)
				m.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(nil)
				m.publisher.On("Publish", context.Background(), mock.AnythingOfType("ddd.event")).Return(nil)
			},
		},
		"NoBasket": {
			args: args{
				ctx: context.Background(),
				start: StartBasket{
					ID:             "basket-id",
					UserCustomerID: "user-customer-id",
				},
			},
			on: func(m mocks) {
				m.baskets.On("Load", context.Background(), "basket-id").Return(nil, fmt.Errorf("no basket"))
			},
			wantErr: true,
		},
		"SaveFailed": {
			args: args{
				ctx: context.Background(),
				start: StartBasket{
					ID:             "basket-id",
					UserCustomerID: "customer-id",
				},
			},
			on: func(m mocks) {
				m.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
					UserCustomerID: "",
					Items:          make(map[string]domain.Item),
					Status:         "",
				}, nil)
				m.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(fmt.Errorf("save failed"))
			},
			wantErr: true,
		},
		"PublishFailed": {
			args: args{
				ctx: context.Background(),
				start: StartBasket{
					ID:             "basket-id",
					UserCustomerID: "user-customer-id",
				},
			},
			on: func(m mocks) {
				m.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
					UserCustomerID: "",
					Items:          make(map[string]domain.Item),
					Status:         "",
				}, nil)
				m.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(nil)
				m.publisher.On("Publish", context.Background(), mock.AnythingOfType("ddd.event")).Return(fmt.Errorf("publish failed"))
			},
			wantErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := mocks{
				baskets:   domain.NewMockBasketRepository(t),
				users:     domain.NewMockUserRepository(t),
				products:  domain.NewMockProductRepository(t),
				publisher: ddd.NewMockEventPublisher[ddd.Event](t),
			}
			a := Application{
				baskets:   m.baskets,
				users:     m.users,
				products:  m.products,
				catalog:   domain.NewMockCatalogRepository(t),
				publisher: m.publisher,
			}
			if tc.on != nil {
				tc.on(m)
			}

			if err := a.StartBasket(tc.args.ctx, tc.args.start); (err != nil) != tc.wantErr {
				t.Errorf("StartBasket() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
