package rest

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/stretchr/testify/assert"
	grpcstd "google.golang.org/grpc"
	"middleman/baskets/internal/application"
	"middleman/baskets/internal/domain"
	"middleman/baskets/internal/grpc"
	"middleman/internal/ddd"
	"middleman/internal/registry"
	"middleman/internal/rpc"
	"middleman/internal/web"
	"net"
	"net/http"
	"os"
	"testing"
)

var pactBrokerURL string
var pactUser string
var pactPass string
var pactToken string

func init() {
	getEnv := func(key, fallback string) string {
		if value, ok := os.LookupEnv(key); ok {
			return value
		}
		return fallback
	}

	pactBrokerURL = getEnv("PACT_URL", "http://127.0.0.1:9292")
	pactUser = getEnv("PACT_USER", "pactuser")
	pactPass = getEnv("PACT_PASS", "pactpass")
	pactToken = getEnv("PACT_TOKEN", "")
}

func TestProvider(t *testing.T) {
	var err error

	// init registry
	reg := registry.New()
	err = domain.Registrations(reg)
	if err != nil {
		t.Fatal(err)
	}
	// init repos
	baskets := domain.NewFakeBasketRepository()
	users := domain.NewFakeUserCacheRepository()
	products := domain.NewFakeProductCacheRepository()
	dispatcher := ddd.NewEventDispatcher[ddd.Event]()

	// init app
	app := application.New(baskets, users, products, dispatcher)

	// start grpc
	rpcConfig := rpc.RpcConfig{
		Host: "0.0.0.0",
		Port: ":9095",
	}
	grpcServer := grpcstd.NewServer()
	// start rest
	webConfig := web.WebConfig{
		Host: "0.0.0.0",
		Port: ":9090",
	}
	mux := chi.NewMux()

	err = grpc.RegisterServer(app, grpcServer)
	if err != nil {
		t.Fatal(err)
	}
	err = RegisterGateway(context.Background(), mux, rpcConfig.Address())
	if err != nil {
		t.Fatal(err)
	}

	// start up the GRPC API; we proxy the REST api through the GRPC API
	{
		listener, err := net.Listen("tcp", rpcConfig.Address())
		if err != nil {
			t.Fatal(err)
		}
		go func() {
			if err = grpcServer.Serve(listener); err != nil && err != grpcstd.ErrServerStopped {
				t.Error(err)
				return
			}
		}()
		defer func() {
			grpcServer.GracefulStop()
		}()
	}

	// start up the REST API
	{
		webServer := &http.Server{
			Addr:    webConfig.Address(),
			Handler: mux,
		}
		go func() {
			if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				t.Error(err)
				return
			}
		}()
		defer func() {
			if err := webServer.Shutdown(context.Background()); err != nil {
				t.Error(err)
				return
			}
		}()
	}

	verifier := provider.HTTPVerifier{}
	assert.NoError(t, verifier.VerifyProvider(t, provider.VerifyRequest{
		Provider:                   "baskets-api",
		ProviderBaseURL:            fmt.Sprintf("http://%s", webConfig.Address()),
		ProviderVersion:            "1.0.0",
		BrokerURL:                  pactBrokerURL,
		BrokerToken:                pactToken,
		BrokerUsername:             pactUser,
		BrokerPassword:             pactPass,
		PublishVerificationResults: true,
		AfterEach: func() error {
			baskets.Reset()
			products.Reset()
			users.Reset()
			return nil
		},
		StateHandlers: map[string]models.StateHandler{
			"a basket exists": func(setup bool, state models.ProviderState) (models.ProviderStateResponse, error) {
				basket := domain.NewBasket("basket-id")
				if v, exists := state.Parameters["id"]; exists {
					basket = domain.NewBasket(v.(string))
				}
				basket.Items = map[string]domain.Item{}
				basket.UserCustomerID = "user-customer-id"
				if v, exists := state.Parameters["userCustomerId"]; exists {
					basket.UserCustomerID = v.(string)
				}
				basket.Status = domain.BasketIsOpen
				if v, exists := state.Parameters["status"]; exists && domain.BasketStatus(v.(string)).String() != "" {
					basket.Status = domain.BasketStatus(v.(string))
				}
				baskets.Reset(basket)

				return nil, nil
			},
			"a product exists": func(setup bool, state models.ProviderState) (models.ProviderStateResponse, error) {
				product := &domain.Product{
					ID:           "product-id",
					UserSellerID: "user-seller-id",
					Name:         "TheProduct",
					Price:        10.00,
				}
				if v, exists := state.Parameters["id"]; exists {
					product.ID = v.(string)
				}
				if v, exists := state.Parameters["userSellerId"]; exists {
					product.UserSellerID = v.(string)
				}
				if v, exists := state.Parameters["name"]; exists {
					product.Name = v.(string)
				}
				if v, exists := state.Parameters["price"]; exists {
					product.Price = v.(float64)
				}
				products.Reset(product)
				return nil, nil
			},
			"a user seller exists": func(setup bool, state models.ProviderState) (models.ProviderStateResponse, error) {
				userSeller := &domain.User{
					ID:        "user-seller-id",
					FirstName: "user-first-name",
				}
				if v, exists := state.Parameters["id"]; exists {
					userSeller.ID = v.(string)
				}
				if v, exists := state.Parameters["name"]; exists {
					userSeller.FirstName = v.(string)
				}
				users.Reset(userSeller)
				return nil, nil
			},
		},
	}))
}
