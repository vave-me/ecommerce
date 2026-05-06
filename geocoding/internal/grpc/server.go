package grpc

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"middleman/geocoding/geocodingpb"
	"middleman/geocoding/internal/application"
	"middleman/geocoding/internal/application/commands"
	"middleman/geocoding/internal/constants"
	"middleman/internal/di"
	"middleman/internal/geo"
)

type server struct {
	app application.App
	geocodingpb.UnimplementedGeocodingServiceServer
}

// RegisterServer registers this service implementation with the given gRPC server.
func RegisterServer(
	ctx context.Context,
	app application.Application,
	registrar grpc.ServiceRegistrar,
) error {
	geocodingpb.RegisterGeocodingServiceServer(registrar, server{app: app})
	return nil
}

func (s server) GeocodeAddress(ctx context.Context, req *geocodingpb.GeocodeAddressRequest) (*geocodingpb.GeocodeAddressResponse, error) {

	log.Printf("[GeocodingService] GeocodingAddress: %s",
		req.GetAddress())

	//client := di.Get(ctx, constants.GoogleGeocode).(*geo.GoogleGeocodingClient)
	nominati := di.Get(ctx, constants.NominatimGeocode).(*geo.NominatimGeocodingClient)

	addr, err := nominati.ForwardGeocode(ctx, req.GetAddress())
	if err != nil {
		return nil, err
	}
	cmd := commands.GeocodeAddress{
		Address:   addr.DisplayName,
		Latitude:  addr.Latitude,
		Longitude: addr.Longitude,
	}

	if err := s.app.GeocodeAddress(ctx, cmd); err != nil {
		return nil, err
	}

	return &geocodingpb.GeocodeAddressResponse{
		Address:   addr.DisplayName,
		Latitude:  float32(addr.Latitude),
		Longitude: float32(addr.Longitude),
	}, nil
}

func (s server) SuggestAddress(ctx context.Context, req *geocodingpb.SuggestAddressRequest) (*geocodingpb.SuggestAddressResponse, error) {
	log.Printf("[GeocodingService] SuggestAddress: %s", req.GetAddress())

	// Retrieve the Nominatim geocoding client from the dependency injector.
	nominati := di.Get(ctx, constants.NominatimGeocode).(*geo.NominatimGeocodingClient)

	// Call the SuggestCities method (which returns a slice of geo.Suggestion).
	suggestions, err := nominati.SuggestAddresses(ctx, req.GetAddress(), 10)
	if err != nil {
		return nil, err
	}

	// Convert the JSON suggestions to the proto format.
	protoSuggestions := addressToProto(suggestions)

	// Return the response with the converted suggestions.
	return &geocodingpb.SuggestAddressResponse{
		SuggestionAddress: protoSuggestions,
	}, nil
}

// SuggestCity handles the gRPC request for city suggestions.
func (s server) SuggestCity(ctx context.Context, req *geocodingpb.SuggestCityRequest) (*geocodingpb.SuggestCityResponse, error) {
	log.Printf("[GeocodingService] SuggestCity: %s", req.GetName())

	// Retrieve the Nominatim geocoding client from the dependency injector.
	nominati := di.Get(ctx, constants.NominatimGeocode).(*geo.NominatimGeocodingClient)

	// Call the SuggestCities method (which returns a slice of geo.Suggestion).
	suggestions, err := nominati.SuggestCities(ctx, req.GetName(), 10)
	if err != nil {
		return nil, err
	}

	// Convert the JSON suggestions to the proto format.
	protoSuggestions := suggestionsToProto(suggestions)

	// Return the response with the converted suggestions.
	return &geocodingpb.SuggestCityResponse{
		SuggestedCities: protoSuggestions,
	}, nil
}

// suggestionsToProto converts a slice of geo.Suggestion to a slice of geocodingpb.SuggestionCity.
func suggestionsToProto(suggestions []geo.Suggestion) []*geocodingpb.SuggestionCity {
	protoSuggestions := make([]*geocodingpb.SuggestionCity, 0, len(suggestions))
	for _, s := range suggestions {
		protoSuggestions = append(protoSuggestions, &geocodingpb.SuggestionCity{
			SuggestedCity: s.DisplayName, // mapping the JSON display name to the proto field
			Latitude:      float32(s.Latitude),
			Longitude:     float32(s.Longitude),
		})
	}
	return protoSuggestions
}

// suggestionsToProto converts a slice of geo.Suggestion to a slice of geocodingpb.SuggestionCity.
func addressToProto(suggestions []geo.Suggestion) []*geocodingpb.SuggestionAddress {
	protoSuggestions := make([]*geocodingpb.SuggestionAddress, 0, len(suggestions))
	for _, s := range suggestions {
		protoSuggestions = append(protoSuggestions, &geocodingpb.SuggestionAddress{
			SuggestedAddress: s.DisplayName, // mapping the JSON display name to the proto field
			Latitude:         float32(s.Latitude),
			Longitude:        float32(s.Longitude),
		})
	}
	return protoSuggestions
}
