package grpc

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"middleman/baskets/basketspb"
	"middleman/baskets/internal/application"
	"middleman/baskets/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/es"
	"net"
	"testing"
)

type serverSuite struct {
	mocks struct {
		baskets   *domain.MockBasketRepository
		users     *domain.MockUserRepository
		products  *domain.MockProductRepository
		catalog   *domain.MockCatalogRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}
	server *grpc.Server
	client basketspb.BasketServiceClient
	suite.Suite
}

func TestServer(t *testing.T) {
	suite.Run(t, &serverSuite{})
}
func (s *serverSuite) SetupSuite()    {}
func (s *serverSuite) TearDownSuite() {}

func (s *serverSuite) SetupTest() {
	const grpcTestPort = ":10912"

	var err error
	// create server
	s.server = grpc.NewServer()
	var listener net.Listener
	listener, err = net.Listen("tcp", grpcTestPort)
	if err != nil {
		s.T().Fatal(err)
	}

	// create mocks
	s.mocks = struct {
		baskets   *domain.MockBasketRepository
		users     *domain.MockUserRepository
		products  *domain.MockProductRepository
		catalog   *domain.MockCatalogRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}{
		baskets:   domain.NewMockBasketRepository(s.T()),
		users:     domain.NewMockUserRepository(s.T()),
		products:  domain.NewMockProductRepository(s.T()),
		catalog:   domain.NewMockCatalogRepository(s.T()),
		publisher: ddd.NewMockEventPublisher[ddd.Event](s.T()),
	}

	// create app
	app := application.New(s.mocks.baskets, s.mocks.users, s.mocks.products, s.mocks.catalog, s.mocks.publisher)

	// register app with server
	if err = RegisterServer(app, s.server); err != nil {
		s.T().Fatal(err)
	}
	go func(listener net.Listener) {
		err := s.server.Serve(listener)
		if err != nil {
			s.T().Fatal(err)
		}
	}(listener)

	// create client
	var conn *grpc.ClientConn
	conn, err = grpc.NewClient(grpcTestPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.T().Fatal(err)
	}
	s.client = basketspb.NewBasketServiceClient(conn)
}
func (s *serverSuite) TearDownTest() {
	s.server.GracefulStop()
}
func (s *serverSuite) TestBasketService_StartBasket() {
	s.mocks.baskets.On("Load", mock.Anything, mock.AnythingOfType("string")).Return(&domain.Basket{
		Aggregate: es.NewAggregate("basket-id", domain.BasketAggregate),
	}, nil)
	s.mocks.baskets.On("Save", mock.Anything, mock.AnythingOfType("*domain.Basket")).Return(nil)
	s.mocks.publisher.On("Publish", mock.Anything, mock.AnythingOfType("ddd.event")).Return(nil)

	_, err := s.client.StartBasket(context.Background(), &basketspb.StartBasketRequest{UserId: "user-customer-id"})
	s.Assert().NoError(err)
}
func (s *serverSuite) TestBasketService_CancelBasket() {
	s.mocks.baskets.On("Load", mock.Anything, "basket-id").Return(&domain.Basket{
		Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
		UserCustomerID: "user-customer-id",
		Status:         domain.BasketIsOpen,
	}, nil)
	s.mocks.baskets.On("Save", mock.Anything, mock.AnythingOfType("*domain.Basket")).Return(nil)
	s.mocks.publisher.On("Publish", mock.Anything, mock.AnythingOfType("ddd.event")).Return(nil)

	_, err := s.client.CancelBasket(context.Background(), &basketspb.CancelBasketRequest{BasketId: "basket-id"})
	s.Assert().NoError(err)
}

func (s *serverSuite) TestBasketService_CheckoutBasket() {
	s.mocks.baskets.On("Load", mock.Anything, "basket-id").Return(&domain.Basket{
		Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
		UserCustomerID: "user-customer-id",
		Items: map[string]domain.Item{
			"product-id": {
				UserSellerID:   "user-seller-id",
				ProductID:      "product-id",
				UserSellerName: "user-seller-name",
				ProductName:    "product-name",
				ProductPrice:   1.00,
				Quantity:       1,
			},
		},
		Status: domain.BasketIsOpen,
	}, nil)
	s.mocks.baskets.On("Save", mock.Anything, mock.AnythingOfType("*domain.Basket")).Return(nil)
	s.mocks.publisher.On("Publish", mock.Anything, mock.AnythingOfType("ddd.event")).Return(nil)

	_, err := s.client.CheckoutBasket(context.Background(), &basketspb.CheckoutBasketRequest{
		BasketId: "basket-id",
	})
	s.Assert().NoError(err)
}
func (s *serverSuite) TestBasketService_AddItem() {
	product := &domain.Product{
		ID:           "product-id",
		UserSellerID: "user-seller-id",
		Name:         "product-name",
		BasePrice:    10.00,
	}
	user := &domain.User{
		ID:        "user-id",
		FirstName: "user-seller-name",
	}
	s.mocks.baskets.On("Load", mock.Anything, "basket-id").Return(&domain.Basket{
		Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
		UserCustomerID: "user-customer-id",
		Items: map[string]domain.Item{
			"product-id": {
				UserSellerID:   "user-seller-id",
				ProductID:      "product-id",
				UserSellerName: "user-seller-name",
				ProductName:    "product-name",
				ProductPrice:   1.00,
				Quantity:       1,
			},
		},
		Status: domain.BasketIsOpen,
	}, nil)
	s.mocks.baskets.On("Save", mock.Anything, mock.AnythingOfType("*domain.Basket")).Return(nil)
	s.mocks.products.On("Find", mock.Anything, "product-id").Return(product, nil)
	s.mocks.users.On("Find", mock.Anything, "user-seller-id").Return(user, nil)

	_, err := s.client.AddItem(context.Background(), &basketspb.AddItemRequest{
		BasketId:  "basket-id",
		ProductId: "product-id",
		Quantity:  1,
	})
	s.Assert().NoError(err)
}

func (s *serverSuite) TestBasketService_RemoveItem() {
	product := &domain.Product{
		ID:           "product-id",
		UserSellerID: "user-seller-id",
		Name:         "product-name",
		BasePrice:    10.00,
	}

	s.mocks.baskets.On("Load", mock.Anything, "basket-id").Return(&domain.Basket{
		Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
		UserCustomerID: "user-customer-id",
		Items: map[string]domain.Item{
			"product-id": {
				UserSellerID:   "user-seller-id",
				ProductID:      "product-id",
				UserSellerName: "user-seller-name",
				ProductName:    "product-name",
				ProductPrice:   1.00,
				Quantity:       1,
			},
		},
		Status: domain.BasketIsOpen,
	}, nil)
	s.mocks.baskets.On("Save", mock.Anything, mock.AnythingOfType("*domain.Basket")).Return(nil)
	s.mocks.products.On("Find", mock.Anything, "product-id").Return(product, nil)

	_, err := s.client.RemoveItem(context.Background(), &basketspb.RemoveItemRequest{
		BasketId: "basket-id",
		ItemId:   "product-id",
		Quantity: 1,
	})
	s.Assert().NoError(err)
}
func (s *serverSuite) TestBasketService_GetBasket() {
	basket := &domain.Basket{
		Aggregate:      es.NewAggregate("basket-id", domain.BasketAggregate),
		UserCustomerID: "user-customer-id",
		Items: map[string]domain.Item{
			"product-id": {
				UserSellerID:   "user-seller-id",
				ProductID:      "product-id",
				UserSellerName: "user-seller-name",
				ProductName:    "product-name",
				ProductPrice:   1.00,
				Quantity:       1,
			},
		},
		Status: domain.BasketIsOpen,
	}
	s.mocks.baskets.On("Load", mock.Anything, "basket-id").Return(basket, nil)

	resp, err := s.client.GetBasket(context.Background(), &basketspb.GetBasketRequest{BasketId: "basket-id"})
	if s.Assert().NoError(err) {
		s.Assert().Equal(basket.ID(), resp.Basket.GetId())
	}
}
