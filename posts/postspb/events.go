package postspb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

// Channels
const (
	// The channel name or topic name if you use a message broker.
	PostAggregateChannel = "middleman.posts.events.Post"
)

// Post Event Names
const (
	PostAddedEvent            = "postsapi.PostAdded"
	PostUpdatedEvent          = "postsapi.PostUpdated"
	PostArchivedEvent         = "postsapi.PostArchived"
	PostRemovedEvent          = "postsapi.PostRemoved"
	PostThumbnailAddedEvent   = "postsapi.PostThumbnailAdded"
	PostThumbnailUpdatedEvent = "postsapi.PostThumbnailUpdated"
)

// In your real code, define the Go struct equivalents for these proto messages,
// e.g.:
// type PostAdded struct {
//     Id        string   `json:"id"`
//     UserId    string   `json:"user_id"`
//     Name     string   `json:"name"`
//     Description   string   `json:"description"`
//     Tags      []string `json:"tags"`
//     Status    string   `json:"status"`
//     Thumbnail string   `json:"thumbnail"`
//     Lat       float32  `json:"lat"`
//     Lng       float32  `json:"lng"`
// }
// etc.

// Registrations and Serde
func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Register new domain events
	if err := serde.Register(&PostAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&PostUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&PostRemoved{}); err != nil {
		return err
	}
	if err := serde.Register(&PostArchived{}); err != nil {
		return err
	}

	if err := serde.Register(&PostThumbnailAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&PostThumbnailUpdated{}); err != nil {
		return err
	}
	return nil
}

// The Go event structs should implement a Key() string method
// matching the constants above.

func (*PostAdded) Key() string            { return PostAddedEvent }
func (*PostUpdated) Key() string          { return PostUpdatedEvent }
func (*PostArchived) Key() string         { return PostArchivedEvent }
func (*PostRemoved) Key() string          { return PostRemovedEvent }
func (*PostThumbnailAdded) Key() string   { return PostThumbnailAddedEvent }
func (*PostThumbnailUpdated) Key() string { return PostThumbnailUpdatedEvent }

// Remove references to PostPriceIncreased, PostPriceDecreased,
// PostStockAdjusted, PostNegotiableToggled, PostSoldEvent, etc.
