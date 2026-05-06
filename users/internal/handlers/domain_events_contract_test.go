//go:build contract

package handlers

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/pact-foundation/pact-go/v2/message"
	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/pact-foundation/pact-go/v2/provider"

	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
	"middleman/users/internal/application"
	"middleman/users/internal/application/commands"
	"middleman/users/internal/domain"
	"middleman/users/userspb"
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

func TestUsersProducer(t *testing.T) {
	var err error

	users := domain.NewFakeUserRepository()
	products := domain.NewFakeProductRepository()
	middleman := domain.NewFakeMiddlemanRepository()
	catalog := domain.NewFakeCatalogRepository()

	type rawEvent struct {
		Name    string
		Payload json.RawMessage
	}

	reg := registry.New()
	err = userspb.RegistrationsWithSerde(serdes.NewJsonSerde(reg))
	if err != nil {
		t.Fatal(err)
	}

	verifier := message.Verifier{}
	err = verifier.Verify(t, message.VerifyMessageRequest{
		VerifyRequest: provider.VerifyRequest{
			Provider:                   "users-pub",
			ProviderVersion:            "1.0.0",
			BrokerURL:                  pactBrokerURL,
			BrokerUsername:             pactUser,
			BrokerPassword:             pactPass,
			BrokerToken:                pactToken,
			PublishVerificationResults: true,
			AfterEach: func() error {
				users.Reset()
				products.Reset()
				return nil
			},
		},
		MessageHandlers: map[string]message.Handler{
			"a StoreCreated message": func(states []models.ProviderState) (message.Body, message.Metadata, error) {
				// Assign
				dispatcher := ddd.NewEventDispatcher[ddd.Event]()
				app := application.New(users, products, catalog, mall, dispatcher)
				publisher := am.NewFakeEventPublisher()
				handler := NewDomainEventHandlers(publisher)
				RegisterDomainEventHandlers(dispatcher, handler)

				// Act
				err := app.CreateStore(context.Background(), commands.CreateStore{
					ID:       "store-id",
					Name:     "NewStore",
					Location: "NewLocation",
				})
				if err != nil {
					return nil, nil, err
				}

				// Assert
				subject, event, err := publisher.Last()
				if err != nil {
					return nil, nil, err
				}

				return rawEvent{
						Name:    event.EventName(),
						Payload: reg.MustSerialize(event.EventName(), event.Payload()),
					}, map[string]any{
						"subject": subject,
					}, nil
			},
			"a UserRenamed message": func(states []models.ProviderState) (message.Body, message.Metadata, error) {
				dispatcher := ddd.NewEventDispatcher[ddd.Event]()
				app := application.New(users, products, catalog, mall, dispatcher)
				publisher := am.NewFakeEventPublisher()
				handler := NewDomainEventHandlers(publisher)
				RegisterDomainEventHandlers(dispatcher, handler)

				user := domain.NewStore("store-id")
				user.Name = "NewStore"
				user.Location = "NewLocation"
				users.Reset(store)

				err := app.RenamedUser(context.Background(), commands.RenameUser{
					ID:   "store-id",
					Name: "RebrandedStore",
				})
				if err != nil {
					return nil, nil, err
				}

				subject, event, err := publisher.Last()
				if err != nil {
					return nil, nil, err
				}

				return rawEvent{
						Name:    event.EventName(),
						Payload: reg.MustSerialize(event.EventName(), event.Payload()),
					}, map[string]any{
						"subject": subject,
					}, nil
			},
		},
	})

	if err != nil {
		t.Error(err)
	}
}
