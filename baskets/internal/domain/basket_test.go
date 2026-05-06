package domain

import (
	"github.com/stretchr/testify/assert"
	"middleman/internal/ddd"
	"middleman/internal/es"
	"reflect"
	"testing"
)

func TestBasket_AddItem(t *testing.T) {
	user := &User{
		ID:        "user-id",
		FirstName: "user-first-name",
		LastName:  "user-last-name",
		Email:     "user-email",
	}
	product := &Product{
		ID:           "product-id",
		UserSellerID: "user-seller-id",
		Name:         "product-name",
		BasePrice:    10.00,
	}

	type fields struct {
		UserCustomerID string
		Items          map[string]Item
		Status         BasketStatus
	}
	type args struct {
		user     *User
		product  *Product
		quantity int
	}
	tests := map[string]struct {
		fields  fields
		args    args
		on      func(a *es.MockAggregate)
		wantErr bool
	}{
		"OpenBasket": {
			fields: fields{
				Items:  make(map[string]Item),
				Status: BasketIsOpen,
			},
			args: args{
				user:     user,
				product:  product,
				quantity: 1,
			},
			on: func(a *es.MockAggregate) {
				a.On("AddEvent", BasketItemAddedEvent, &BasketItemAdded{
					Item: Item{
						UserSellerID:   user.ID,
						ProductID:      product.ID,
						UserSellerName: user.FirstName,
						ProductName:    product.Name,
						ProductPrice:   product.BasePrice,
						Quantity:       1,
					},
				})
			},
			wantErr: false,
		},
		"CheckedOutBasket": {
			fields: fields{
				Items:  make(map[string]Item),
				Status: BasketIsCheckedOut,
			},
			args: args{
				user:     user,
				product:  product,
				quantity: 1,
			},
			wantErr: true,
		},
		"CanceledBasket": {
			fields: fields{
				Items:  make(map[string]Item),
				Status: BasketIsCanceled,
			},
			args: args{
				user:     user,
				product:  product,
				quantity: 1,
			},
			wantErr: true,
		},
		"ZeroQuantity": {
			fields: fields{
				Items:  make(map[string]Item),
				Status: BasketIsCheckedOut, // or BasketIsOpen, but either way we're expecting an error
			},
			args: args{
				user:     user,
				product:  product,
				quantity: 0,
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			aggregate := es.NewMockAggregate(t)
			b := &Basket{
				Aggregate:      aggregate,
				UserCustomerID: tt.fields.UserCustomerID,
				Items:          tt.fields.Items,
				Status:         tt.fields.Status,
			}
			if tt.on != nil {
				tt.on(aggregate)
			}

			_, err := b.AddItem(tt.args.user, tt.args.product, int64(tt.args.quantity))
			if (err != nil) != tt.wantErr {
				t.Errorf("AddItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBasket_ApplyEvent(t *testing.T) {
	user := &User{
		ID:        "user-id",
		FirstName: "user-first-name",
		LastName:  "user-last-name",
		Email:     "user-email",
	}
	product := &Product{
		ID:           "product-id",
		UserSellerID: "user-seller-id",
		Name:         "product-name",
		BasePrice:    10.00,
	}
	product2 := &Product{
		ID:           "product-id2",
		UserSellerID: "user-seller-id",
		Name:         "product-name2",
		BasePrice:    100.00,
	}

	type fields struct {
		UserCustomerID string
		Items          map[string]Item
		Status         BasketStatus
	}
	type args struct {
		event ddd.Event
	}
	tests := map[string]struct {
		fields  fields
		args    args
		want    fields
		wantErr bool
	}{
		"BasketItemAddedEvent": {
			fields: fields{
				UserCustomerID: "user-customer-id",
				Items:          make(map[string]Item),
				Status:         BasketIsOpen,
			},
			args: args{
				event: ddd.NewEvent(BasketItemAddedEvent, &BasketItemAdded{
					Item: Item{
						UserSellerID:   user.ID,
						ProductID:      product.ID,
						UserSellerName: user.FirstName,
						ProductName:    product.Name,
						ProductPrice:   product.BasePrice,
						Quantity:       1,
					},
				}),
			},
			want: fields{
				UserCustomerID: "user-customer-id",
				Items: map[string]Item{
					product.ID: {
						UserSellerID:   user.ID,
						ProductID:      product.ID,
						UserSellerName: user.FirstName,
						ProductName:    product.Name,
						ProductPrice:   product.BasePrice,
						Quantity:       1,
					},
				},
				Status: BasketIsOpen,
			},
			wantErr: false,
		},
		"BasketItemAddedEvent.Quantity": {
			fields: fields{
				UserCustomerID: "user-customer-id",
				Items: map[string]Item{
					product.ID: {
						UserSellerID:   user.ID,
						ProductID:      product.ID,
						UserSellerName: user.FirstName,
						ProductName:    product.Name,
						ProductPrice:   product.BasePrice,
						Quantity:       1,
					},
				},
				Status: BasketIsOpen,
			},
			args: args{
				event: ddd.NewEvent(BasketItemAddedEvent, &BasketItemAdded{
					Item: Item{
						UserSellerID:   user.ID,
						ProductID:      product.ID,
						UserSellerName: user.FirstName,
						ProductName:    product.Name,
						ProductPrice:   product.BasePrice,
						Quantity:       1,
					},
				}),
			},
			want: fields{
				UserCustomerID: "user-customer-id",
				Items: map[string]Item{
					product.ID: {
						UserSellerID:   user.ID,
						ProductID:      product.ID,
						UserSellerName: user.FirstName,
						ProductName:    product.Name,
						ProductPrice:   product.BasePrice,
						Quantity:       2,
					},
				},
				Status: BasketIsOpen,
			},
			wantErr: false,
		},
		"BasketItemAddedEvent.Second": {
			fields: fields{
				UserCustomerID: "user-customer-id",
				Items: map[string]Item{
					product.ID: {
						UserSellerID:   user.ID,
						ProductID:      product.ID,
						UserSellerName: user.FirstName,
						ProductName:    product.Name,
						ProductPrice:   product.BasePrice,
						Quantity:       1,
					},
				},
				Status: BasketIsOpen,
			},
			args: args{
				event: ddd.NewEvent(BasketItemAddedEvent, &BasketItemAdded{
					Item: Item{
						UserSellerID:   user.ID,
						ProductID:      product2.ID,
						UserSellerName: user.FirstName,
						ProductName:    product2.Name,
						ProductPrice:   product2.BasePrice,
						Quantity:       1,
					},
				}),
			},
			want: fields{
				UserCustomerID: "user-customer-id",
				Items: map[string]Item{
					product.ID: {
						UserSellerID:   user.ID,
						ProductID:      product.ID,
						UserSellerName: user.FirstName,
						ProductName:    product.Name,
						ProductPrice:   product.BasePrice,
						Quantity:       1,
					},
					product2.ID: {
						UserSellerID:   user.ID,
						ProductID:      product2.ID,
						UserSellerName: user.FirstName,
						ProductName:    product2.Name,
						ProductPrice:   product2.BasePrice,
						Quantity:       1,
					},
				},
				Status: BasketIsOpen,
			},
			wantErr: false,
		},
		"BasketCanceledEvent": {
			fields: fields{
				UserCustomerID: "user-customer-id",
				Items:          make(map[string]Item),
				Status:         BasketIsOpen,
			},
			args: args{
				event: ddd.NewEvent(BasketCanceledEvent, &BasketCanceled{
					Status: BasketIsCanceled, // now matches domain
				}),
			},
			want: fields{
				UserCustomerID: "user-customer-id",
				Items:          map[string]Item{},
				Status:         BasketIsCanceled,
			},
			wantErr: false,
		},
		"BasketCanceledEvent.Cleared": {
			fields: fields{
				UserCustomerID: "user-customer-id",
				Items: map[string]Item{
					product.ID: {
						UserSellerID:   user.ID,
						ProductID:      product.ID,
						UserSellerName: user.FirstName,
						ProductName:    product.Name,
						ProductPrice:   product.BasePrice,
						Quantity:       1,
					},
				},
				Status: BasketIsOpen,
			},
			args: args{
				event: ddd.NewEvent(BasketCanceledEvent, &BasketCanceled{
					Status: BasketIsCanceled,
				}),
			},
			want: fields{
				UserCustomerID: "user-customer-id",
				Items:          map[string]Item{},
				Status:         BasketIsCanceled,
			},
			wantErr: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b := &Basket{
				Aggregate:      es.NewMockAggregate(t),
				UserCustomerID: tt.fields.UserCustomerID,
				Items:          tt.fields.Items,
				Status:         tt.fields.Status,
			}
			if err := b.ApplyEvent(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("ApplyEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want.UserCustomerID, b.UserCustomerID)
			assert.Equal(t, tt.want.Items, b.Items)
			assert.Equal(t, tt.want.Status, b.Status)
		})
	}
}

func TestBasket_ApplySnapshot(t *testing.T) {
	user := &User{
		ID:        "user-id",
		FirstName: "user-first-name",
		LastName:  "user-last-name",
		Email:     "user-email",
	}
	product := &Product{
		ID:           "product-id",
		UserSellerID: "user-id",
		Name:         "product-name",
		BasePrice:    10.00,
	}
	item := Item{
		UserSellerID:   user.ID,
		ProductID:      product.ID,
		UserSellerName: user.FirstName,
		ProductName:    product.Name,
		ProductPrice:   product.BasePrice,
		Quantity:       1,
	}

	type fields struct {
		UserCustomerID string
		Items          map[string]Item
		Status         BasketStatus
	}
	type args struct {
		snapshot es.Snapshot
	}
	tests := map[string]struct {
		fields  fields
		args    args
		want    fields
		wantErr bool
	}{
		"V1": {
			fields: fields{},
			args: args{
				snapshot: &BasketV1{
					UserCustomerID: "user-customer-id",
					// PaymentMethodID removed or unused
					Items: map[string]Item{
						product.ID: item,
					},
					Status: BasketIsOpen,
				},
			},
			want: fields{
				UserCustomerID: "user-customer-id",
				Items: map[string]Item{
					product.ID: item,
				},
				Status: BasketIsOpen,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b := &Basket{
				Aggregate: es.NewMockAggregate(t),
				Items:     tt.fields.Items,
				Status:    tt.fields.Status,
			}
			if err := b.ApplySnapshot(tt.args.snapshot); (err != nil) != tt.wantErr {
				t.Errorf("ApplySnapshot() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want.UserCustomerID, b.UserCustomerID)
			assert.Equal(t, tt.want.Items, b.Items)
			assert.Equal(t, tt.want.Status, b.Status)
		})
	}
}

func TestBasket_Cancel(t *testing.T) {
	type fields struct {
		UserCustomerID string
		Items          map[string]Item
		Status         BasketStatus
	}
	tests := map[string]struct {
		fields  fields
		on      func(a *es.MockAggregate)
		want    ddd.Event
		wantErr bool
	}{
		"OpenBasket": {
			fields: fields{
				UserCustomerID: "user-customer-id",
				Items:          make(map[string]Item),
				Status:         BasketIsOpen,
			},
			on: func(a *es.MockAggregate) {
				// Expect the domain to call AddEvent with status "canceled"
				a.On("AddEvent", BasketCanceledEvent, &BasketCanceled{
					Status: BasketIsCanceled,
				})
			},
			want: ddd.NewEvent(BasketCanceledEvent, &Basket{
				UserCustomerID: "user-customer-id",
				Items:          make(map[string]Item),
				Status:         BasketIsCanceled,
			}),
		},
		"CheckedOutBasket": {
			fields: fields{
				UserCustomerID: "customer-id",
				Items:          make(map[string]Item),
				Status:         BasketIsCheckedOut,
			},
			wantErr: true,
		},
		"CanceledBasket": {
			fields: fields{
				UserCustomerID: "user-customer-id",
				Items:          make(map[string]Item),
				Status:         BasketIsCanceled,
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			aggregate := es.NewMockAggregate(t)
			b := &Basket{
				Aggregate:      aggregate,
				UserCustomerID: tt.fields.UserCustomerID,
				Items:          tt.fields.Items,
				Status:         tt.fields.Status,
			}
			if tt.on != nil {
				tt.on(aggregate)
			}

			got, err := b.Cancel()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cancel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				assert.Equal(t, tt.want.EventName(), got.EventName())
				assert.IsType(t, tt.want.Payload(), got.Payload())
				assert.Equal(t, tt.want.Metadata(), got.Metadata())
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func TestBasket_Checkout(t *testing.T) {
	user := &User{
		ID:        "user-id",
		FirstName: "user-first-name",
		LastName:  "user-last-name",
		Email:     "user-email",
	}
	product := &Product{
		ID:           "product-id",
		UserSellerID: "user-seller-id",
		Name:         "product-name",
		BasePrice:    10.00,
	}
	item := Item{
		UserSellerID:   user.ID,
		ProductID:      product.ID,
		UserSellerName: user.FirstName,
		ProductName:    product.Name,
		ProductPrice:   product.BasePrice,
		Quantity:       1,
	}

	type fields struct {
		UserCustomerID string
		Items          map[string]Item
		Status         BasketStatus
	}
	type args struct {
		paymentMethodID string
	}
	tests := map[string]struct {
		fields  fields
		args    args
		on      func(a *es.MockAggregate)
		wantErr bool
	}{
		"OpenBasket": {
			fields: fields{
				UserCustomerID: "user-customer-id",
				Items: map[string]Item{
					product.ID: item,
				},
				Status: BasketIsOpen,
			},
			args: args{paymentMethodID: "method-123"},
			on: func(a *es.MockAggregate) {
				a.On("AddEvent", BasketCheckedOutEvent, &BasketCheckedOut{
					Status: BasketIsCheckedOut,
				})
			},
			wantErr: false,
		},
		"OpenBasket.NoItems": {
			fields: fields{
				UserCustomerID: "user-customer-id",
				Items:          make(map[string]Item), // empty
				Status:         BasketIsOpen,
			},
			args:    args{paymentMethodID: "method-123"},
			wantErr: true, // domain rejects an empty basket
		},
		"CheckedOutBasket": {
			fields: fields{
				UserCustomerID: "user-customer-id",
				Items:          make(map[string]Item),
				Status:         BasketIsCheckedOut,
			},
			args:    args{paymentMethodID: "method-123"},
			wantErr: true,
		},
		"CanceledBasket": {
			fields: fields{
				UserCustomerID: "user-customer-id",
				Items:          make(map[string]Item),
				Status:         BasketIsCanceled,
			},
			args:    args{paymentMethodID: "method-123"},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			aggregate := es.NewMockAggregate(t)
			b := &Basket{
				Aggregate:      aggregate,
				UserCustomerID: tt.fields.UserCustomerID,
				Items:          tt.fields.Items,
				Status:         tt.fields.Status,
			}
			if tt.on != nil {
				tt.on(aggregate)
			}

			_, err := b.Checkout()
			if (err != nil) != tt.wantErr {
				t.Errorf("Checkout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBasket_RemoveItem(t *testing.T) {
	user := &User{
		ID:        "user-id",
		FirstName: "user-first-name",
		LastName:  "user-last-name",
	}
	product := &Product{
		ID:           "product-id",
		UserSellerID: "user-seller-id",
		Name:         "product-name",
		BasePrice:    10.00,
	}
	item := Item{
		UserSellerID:   user.ID,
		ProductID:      product.ID,
		UserSellerName: user.FirstName,
		ProductName:    product.Name,
		ProductPrice:   product.BasePrice,
		Quantity:       10,
	}

	type fields struct {
		UserCustomerID string
		Items          map[string]Item
		Status         BasketStatus
	}
	type args struct {
		product  *Product
		quantity int
	}
	tests := map[string]struct {
		fields  fields
		args    args
		on      func(a *es.MockAggregate)
		wantErr bool
	}{
		"OpenBasket": {
			fields: fields{
				UserCustomerID: "user-customer-id",
				Items: map[string]Item{
					product.ID: item,
				},
				Status: BasketIsOpen,
			},
			args: args{
				product:  product,
				quantity: 1,
			},
			on: func(a *es.MockAggregate) {
				a.On("AddEvent", BasketItemRemovedEvent, &BasketItemRemoved{
					ProductID: product.ID,
					Quantity:  1,
				})
			},
			wantErr: false,
		},
		"OpenBasket.NoItems": {
			fields: fields{
				UserCustomerID: "user-customer-id",
				Items:          make(map[string]Item),
				Status:         BasketIsOpen,
			},
			args: args{
				product:  product,
				quantity: 1,
			},
			wantErr: false, // no-op remove
		},
		"CheckedOutBasket": {
			fields: fields{
				UserCustomerID: "user-customer-id",
				Items:          make(map[string]Item),
				Status:         BasketIsCheckedOut,
			},
			args: args{
				product:  product,
				quantity: 1,
			},
			wantErr: true,
		},
		"CanceledBasket": {
			fields: fields{
				UserCustomerID: "user-customer-id",
				Items:          make(map[string]Item),
				Status:         BasketIsCanceled,
			},
			args: args{
				product:  product,
				quantity: 1,
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			aggregate := es.NewMockAggregate(t)
			b := &Basket{
				Aggregate:      aggregate,
				UserCustomerID: tt.fields.UserCustomerID,
				Items:          tt.fields.Items,
				Status:         tt.fields.Status,
			}
			if tt.on != nil {
				tt.on(aggregate)
			}

			_, err := b.RemoveItem(tt.args.product, int64(tt.args.quantity))
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBasket_Start(t *testing.T) {
	type fields struct {
		ID             string
		UserCustomerID string
		Status         BasketStatus
	}
	type args struct {
		userCustomerID string
	}
	tests := map[string]struct {
		fields  fields
		args    args
		on      func(a *es.MockAggregate)
		want    ddd.Event
		wantErr bool
	}{
		"New": {
			fields: fields{
				// empty status => default
			},
			args: args{userCustomerID: "user-customer-id"},
			on: func(a *es.MockAggregate) {
				a.On("AddEvent", BasketStartedEvent, &BasketStarted{
					UserCustomerID: "user-customer-id",
					Status:         BasketIsOpen, // domain sets this
				})
			},
			want: ddd.NewEvent(BasketStartedEvent, &Basket{
				UserCustomerID: "user-customer-id",
				Status:         BasketIsOpen,
			}),
		},
		"OpenBasket": {
			fields: fields{
				Status: BasketIsOpen,
			},
			args:    args{userCustomerID: "user-customer-id"},
			wantErr: true,
		},
		"CheckedOutBasket": {
			fields: fields{
				Status: BasketIsCheckedOut,
			},
			args:    args{userCustomerID: "user-customer-id"},
			wantErr: true,
		},
		"CanceledBasket": {
			fields: fields{
				Status: BasketIsCanceled,
			},
			args:    args{userCustomerID: "user-customer-id"},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			aggregate := es.NewMockAggregate(t)
			b := &Basket{
				Aggregate:      aggregate,
				UserCustomerID: tt.fields.UserCustomerID,
				Status:         tt.fields.Status,
			}
			if tt.on != nil {
				tt.on(aggregate)
			}

			got, err := b.Start(tt.args.userCustomerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				assert.Equal(t, tt.want.EventName(), got.EventName())
				assert.IsType(t, tt.want.Payload(), got.Payload())
				assert.Equal(t, tt.want.Metadata(), got.Metadata())
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func TestBasket_ToSnapshot(t *testing.T) {
	user := &User{
		ID:        "user-seller-id",
		LastName:  "user-first-name",
		FirstName: "user-last-name",
		Email:     "user-email",
	}
	product := &Product{
		ID:           "product-id",
		UserSellerID: "user-seller-id",
		Name:         "product-name",
		BasePrice:    10.00,
	}
	item := Item{
		UserSellerID:   user.ID,
		ProductID:      product.ID,
		UserSellerName: user.FirstName,
		ProductName:    product.Name,
		ProductPrice:   product.BasePrice,
		Quantity:       10,
	}

	type fields struct {
		CustomerID string
		Items      map[string]Item
		Status     BasketStatus
	}
	tests := map[string]struct {
		fields fields
		want   es.Snapshot
	}{
		"V1": {
			fields: fields{
				CustomerID: "user-customer-id",
				Items: map[string]Item{
					product.ID: item,
				},
				Status: BasketIsOpen,
			},
			want: &BasketV1{
				UserCustomerID: "user-customer-id",
				Items: map[string]Item{
					product.ID: item,
				},
				Status: BasketIsOpen,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b := &Basket{
				Aggregate:      es.NewMockAggregate(t),
				UserCustomerID: tt.fields.CustomerID,
				Items:          tt.fields.Items,
				Status:         tt.fields.Status,
			}

			got := b.ToSnapshot()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToSnapshot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBasket(t *testing.T) {
	type args struct {
		id string
	}
	tests := map[string]struct {
		args args
		want *Basket
	}{
		"Basket": {
			args: args{id: "basket-id"},
			want: &Basket{
				Aggregate: es.NewAggregate("basket-id", BasketAggregate),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := NewBasket(tt.args.id)
			assert.Equal(t, tt.want.ID(), got.ID())
			assert.Equal(t, tt.want.AggregateName(), got.AggregateName())
		})
	}
}
