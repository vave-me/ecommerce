package domain

import (
	"middleman/internal/ddd"
	"middleman/internal/es"

	"github.com/stackus/errors"
)

const PostAggregate = "posts.Post"

// Domain-level errors (only keep what's relevant for Post)
var (
	ErrPostNameBlank = errors.Wrap(errors.ErrBadRequest, "the post name cannot be blank")
)

type Post struct {
	es.Aggregate
	UserID       string
	Name         string
	Description  string
	TypeOfPost   TypeOfPost
	Tags         []string
	Status       PostStatus
	Thumbnail    string
	UserType     UserType
	CategoryID   string
	CategorySlug string
	Lat          float64
	Lng          float64
}

// Compile-time check that Post implements the required interfaces.
var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Post)(nil)

// NewPost is a constructor for an empty Post aggregate with a given ID.
func NewPost(id string) *Post {
	return &Post{
		Aggregate: es.NewAggregate(id, PostAggregate),
	}
}

// Key implements registry.Registerable (used for referencing the aggregate type).
func (Post) Key() string {
	return PostAggregate
}

// InitPost raises an event to add/initialize a new post.
func (p *Post) InitPost(
	id string,
	name, description string, typeOfPost TypeOfPost, categoryID, categorySlug, userID string, userType UserType,
	tags []string,
	status PostStatus,
	thumbnail string,
	lat, lng float64,
) (ddd.Event, error) {

	p.AddEvent(PostAddedEvent, &PostAdded{
		ID:           id,
		UserID:       userID,
		UserType:     userType,
		Name:         name,
		Description:  description,
		TypeOfPost:   typeOfPost,
		CategoryID:   categoryID,
		CategorySlug: categorySlug,
		Tags:         tags,
		Status:       status,
		Thumbnail:    thumbnail,
		Lat:          lat,
		Lng:          lng,
	})
	return ddd.NewEvent(PostAddedEvent, p), nil
}

// UpdatePost raises an event to update an existing post’s fields.
func (p *Post) Update(
	name, description string, typeOfPost TypeOfPost, categoryID, categorySlug string,
	tags []string, status PostStatus,
	thumbnail string,
) (ddd.Event, error) {
	if name == "" {
		name = p.Name

	}
	if description == "" {
		description = p.Description
	}
	if typeOfPost == "" {
		typeOfPost = p.TypeOfPost
	}
	if categoryID == "" {
		categoryID = p.CategoryID
	}
	if categorySlug == "" {
		categorySlug = p.CategorySlug
	}
	if tags == nil {
		tags = p.Tags
	}
	if status == "" {
		status = p.Status
	}

	p.AddEvent(PostUpdatedEvent, &PostUpdated{
		Name:         name,
		Description:  description,
		TypeOfPost:   typeOfPost,
		CategoryID:   categoryID,
		CategorySlug: categorySlug,
		Tags:         tags,
		Status:       status,
		Thumbnail:    thumbnail,
	})
	return ddd.NewEvent(PostUpdatedEvent, p), nil
}

// ArchivePost raises an event indicating the post has been archived.
func (p *Post) Archive(userID string) (ddd.Event, error) {
	p.AddEvent(PostArchivedEvent, &PostArchived{
		PostID: p.ID(),
		UserID: userID,
	})
	return ddd.NewEvent(PostArchivedEvent, p), nil
}

func (p *Post) AddThumbnail(thumbnail string) (ddd.Event, error) {
	p.AddEvent(PostThumbnailAddedEvent, &PostThumbnailAdded{
		PostID:    p.ID(),
		Thumbnail: thumbnail,
	})
	return ddd.NewEvent(PostAddedEvent, p), nil
}
func (p *Post) UpdateThumbnail(thumbnail string) (ddd.Event, error) {
	p.AddEvent(PostThumbnailUpdatedEvent, &PostThumbnailUpdated{
		PostID:    p.ID(),
		Thumbnail: thumbnail,
	})
	return ddd.NewEvent(PostUpdatedEvent, p), nil
}

// RemovePost raises an event indicating the post has been removed.
func (p *Post) Remove(userID string) (ddd.Event, error) {
	p.AddEvent(PostRemovedEvent, &PostRemoved{
		PostID: p.ID(),
		UserID: userID,
	})
	return ddd.NewEvent(PostRemovedEvent, p), nil
}

// ApplyEvent applies domain events to mutate the in-memory aggregate state.
func (p *Post) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {

	case *PostAdded:
		p.UserID = e.UserID
		p.Name = e.Name
		p.Description = e.Description
		p.TypeOfPost = e.TypeOfPost
		p.CategoryID = e.CategoryID
		p.CategorySlug = e.CategorySlug
		p.Tags = e.Tags
		p.Status = e.Status
		p.Thumbnail = e.Thumbnail
		p.Lat = e.Lat
		p.Lng = e.Lng

	case *PostUpdated:
		p.Name = e.Name
		p.Description = e.Description
		p.TypeOfPost = e.TypeOfPost
		p.CategoryID = e.CategoryID
		p.CategorySlug = e.CategorySlug
		p.Tags = e.Tags
		p.Status = e.Status
		p.Thumbnail = e.Thumbnail

	case *PostArchived:
		// You might set p.Status to PostStatusArchived here
		p.Status = PostStatusArchived

	case *PostRemoved:
		// Possibly mark a `Removed` flag, or handle ephemeral state
		// For example: p.Status = "REMOVED" if that’s your domain rule

	default:
		return errors.ErrInternal.Msgf(
			"%T received an unexpected event %s with payload %T",
			p, event.EventName(), e,
		)
	}
	return nil
}

// ApplySnapshot restores the aggregate from a snapshot.
func (p *Post) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *PostV1:
		p.UserID = ss.UserID
		p.Name = ss.Name
		p.Description = ss.Description
		p.TypeOfPost = ss.TypeOfPost
		p.CategoryID = ss.CategoryID
		p.CategorySlug = ss.CategorySlug
		p.Tags = ss.Tags
		p.Status = ss.Status
		p.Thumbnail = ss.Thumbnail
		p.Lat = ss.Lat
		p.Lng = ss.Lng
	default:
		return errors.ErrInternal.Msgf(
			"%T received an unexpected snapshot %T",
			p, snapshot,
		)
	}
	return nil
}

// ToSnapshot captures the aggregate state in a snapshot struct for persistence.
func (p Post) ToSnapshot() es.Snapshot {
	return &PostV1{
		Name:         p.Name,
		Description:  p.Description,
		TypeOfPost:   p.TypeOfPost,
		CategoryID:   p.CategoryID,
		CategorySlug: p.CategorySlug,
		UserID:       p.UserID,
		UserType:     p.UserType,
		Tags:         p.Tags,
		Status:       p.Status,
		Thumbnail:    p.Thumbnail,
		Lat:          p.Lat,
		Lng:          p.Lng,
	}
}
