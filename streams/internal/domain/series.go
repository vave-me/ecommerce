package domain

import (
	"middleman/internal/ddd"
	"middleman/internal/es"
	"time"

	"github.com/stackus/errors"
)

const SeriesAggregate = "streams.Series"

var (
	ErrSeriesTitleIsBlank   = errors.Wrap(errors.ErrBadRequest, "the series title cannot be blank")
	ErrSeasonAlreadyExists  = errors.Wrap(errors.ErrBadRequest, "season already exists")
	ErrSeasonNotFound       = errors.Wrap(errors.ErrBadRequest, "season not found")
	ErrEpisodeAlreadyExists = errors.Wrap(errors.ErrBadRequest, "episode already exists in this season")
	ErrInvalidSeasonNumber  = errors.Wrap(errors.ErrBadRequest, "invalid season number")
	ErrInvalidEpisodeNumber = errors.Wrap(errors.ErrBadRequest, "invalid episode number")
)

// Series represents a TV series aggregate
type Series struct {
	es.Aggregate
	
	Title        string
	Description  string
	ThumbnailURL string
	Genre        []string
	Studio       string
	Seasons      map[int]*Season // SeasonNumber -> Season
	TotalSeasons int
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Season represents a season within a series
type Season struct {
	SeasonID      string
	SeasonNumber  int
	Title         string
	Description   string
	ThumbnailURL  string
	Episodes      map[int]*Episode // EpisodeNumber -> Episode
	TotalEpisodes int
	CreatedAt     time.Time
}

// Episode represents an episode within a season
type Episode struct {
	EpisodeID     string
	EpisodeNumber int
	StreamID      string // Reference to the actual stream
	Title         string
	Duration      int64
	AirDate       time.Time
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Series)(nil)

func NewSeries(id string) *Series {
	return &Series{
		Aggregate: es.NewAggregate(id, SeriesAggregate),
		Seasons:   make(map[int]*Season),
	}
}

// Key implements registry.Registerable
func (Series) Key() string { return SeriesAggregate }

// InitSeries initializes a new series
func (s *Series) InitSeries(
	title, description, thumbnailURL string,
	genre []string,
	studio string,
) (ddd.Event, error) {
	if title == "" {
		return nil, ErrSeriesTitleIsBlank
	}

	s.AddEvent(SeriesCreatedEvent, &SeriesCreated{
		SeriesID:     s.ID(),
		Title:        title,
		Description:  description,
		ThumbnailURL: thumbnailURL,
		Genre:        genre,
		Studio:       studio,
		CreatedAt:    time.Now(),
	})
	return ddd.NewEvent(SeriesCreatedEvent, s), nil
}

// AddSeason adds a new season to the series
func (s *Series) AddSeason(
	seasonNumber int,
	title, description, thumbnailURL string,
) (ddd.Event, error) {
	if seasonNumber <= 0 {
		return nil, ErrInvalidSeasonNumber
	}
	
	if _, exists := s.Seasons[seasonNumber]; exists {
		return nil, ErrSeasonAlreadyExists
	}

	seasonID := s.ID() + "-S" + string(rune(seasonNumber))
	
	s.AddEvent(SeasonAddedEvent, &SeasonAdded{
		SeriesID:     s.ID(),
		SeasonID:     seasonID,
		SeasonNumber: seasonNumber,
		Title:        title,
		Description:  description,
		ThumbnailURL: thumbnailURL,
		CreatedAt:    time.Now(),
	})
	return ddd.NewEvent(SeasonAddedEvent, s), nil
}

// AddEpisodeToSeason adds an episode to a specific season
func (s *Series) AddEpisodeToSeason(
	seasonNumber, episodeNumber int,
	streamID, title string,
	duration int64,
	airDate time.Time,
) (ddd.Event, error) {
	season, exists := s.Seasons[seasonNumber]
	if !exists {
		return nil, ErrSeasonNotFound
	}
	
	if episodeNumber <= 0 {
		return nil, ErrInvalidEpisodeNumber
	}
	
	if _, exists := season.Episodes[episodeNumber]; exists {
		return nil, ErrEpisodeAlreadyExists
	}

	episodeID := season.SeasonID + "-E" + string(rune(episodeNumber))
	
	s.AddEvent(EpisodeAddedToSeasonEvent, &EpisodeAddedToSeason{
		SeriesID:      s.ID(),
		SeasonID:      season.SeasonID,
		EpisodeID:     episodeID,
		EpisodeNumber: episodeNumber,
		StreamID:      streamID,
		AddedAt:       time.Now(),
	})
	return ddd.NewEvent(EpisodeAddedToSeasonEvent, s), nil
}

// ApplyEvent implements es.EventApplier
func (s *Series) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {
	case *SeriesCreated:
		s.Title = e.Title
		s.Description = e.Description
		s.ThumbnailURL = e.ThumbnailURL
		s.Genre = e.Genre
		s.Studio = e.Studio
		s.CreatedAt = e.CreatedAt
		s.UpdatedAt = e.CreatedAt
		s.Status = "active"

	case *SeasonAdded:
		s.Seasons[e.SeasonNumber] = &Season{
			SeasonID:     e.SeasonID,
			SeasonNumber: e.SeasonNumber,
			Title:        e.Title,
			Description:  e.Description,
			ThumbnailURL: e.ThumbnailURL,
			Episodes:     make(map[int]*Episode),
			CreatedAt:    e.CreatedAt,
		}
		s.TotalSeasons++
		s.UpdatedAt = e.CreatedAt

	case *EpisodeAddedToSeason:
		if season, exists := s.Seasons[s.getSeasonNumberFromID(e.SeasonID)]; exists {
			season.Episodes[e.EpisodeNumber] = &Episode{
				EpisodeID:     e.EpisodeID,
				EpisodeNumber: e.EpisodeNumber,
				StreamID:      e.StreamID,
			}
			season.TotalEpisodes++
		}
		s.UpdatedAt = e.AddedAt

	default:
		return errors.ErrInternal.Msgf(
			"%T received the event %s with unexpected payload %T",
			s, event.EventName(), e)
	}
	return nil
}

// ToSnapshot implements es.Snapshotter
func (s Series) ToSnapshot() es.Snapshot {
	seasons := make([]SeasonInfo, 0, len(s.Seasons))
	for _, season := range s.Seasons {
		episodes := make([]EpisodeInfo, 0, len(season.Episodes))
		for _, episode := range season.Episodes {
			episodes = append(episodes, EpisodeInfo{
				EpisodeID:     episode.EpisodeID,
				EpisodeNumber: episode.EpisodeNumber,
				StreamID:      episode.StreamID,
				Title:         episode.Title,
				Duration:      episode.Duration,
				AirDate:       episode.AirDate,
			})
		}
		
		seasons = append(seasons, SeasonInfo{
			SeasonID:      season.SeasonID,
			SeasonNumber:  season.SeasonNumber,
			Title:         season.Title,
			Description:   season.Description,
			ThumbnailURL:  season.ThumbnailURL,
			Episodes:      episodes,
			TotalEpisodes: season.TotalEpisodes,
			CreatedAt:     season.CreatedAt,
		})
	}

	return SeriesV1{
		SeriesID:     s.ID(),
		Title:        s.Title,
		Description:  s.Description,
		ThumbnailURL: s.ThumbnailURL,
		Genre:        s.Genre,
		Studio:       s.Studio,
		Seasons:      seasons,
		TotalSeasons: s.TotalSeasons,
		Status:       s.Status,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}

// ApplySnapshot implements es.Snapshotter
func (s *Series) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *SeriesV1:
		s.Title = ss.Title
		s.Description = ss.Description
		s.ThumbnailURL = ss.ThumbnailURL
		s.Genre = ss.Genre
		s.Studio = ss.Studio
		s.TotalSeasons = ss.TotalSeasons
		s.Status = ss.Status
		s.CreatedAt = ss.CreatedAt
		s.UpdatedAt = ss.UpdatedAt
		
		// Rebuild seasons map
		s.Seasons = make(map[int]*Season)
		for _, seasonInfo := range ss.Seasons {
			season := &Season{
				SeasonID:      seasonInfo.SeasonID,
				SeasonNumber:  seasonInfo.SeasonNumber,
				Title:         seasonInfo.Title,
				Description:   seasonInfo.Description,
				ThumbnailURL:  seasonInfo.ThumbnailURL,
				Episodes:      make(map[int]*Episode),
				TotalEpisodes: seasonInfo.TotalEpisodes,
				CreatedAt:     seasonInfo.CreatedAt,
			}
			
			for _, episodeInfo := range seasonInfo.Episodes {
				season.Episodes[episodeInfo.EpisodeNumber] = &Episode{
					EpisodeID:     episodeInfo.EpisodeID,
					EpisodeNumber: episodeInfo.EpisodeNumber,
					StreamID:      episodeInfo.StreamID,
					Title:         episodeInfo.Title,
					Duration:      episodeInfo.Duration,
					AirDate:       episodeInfo.AirDate,
				}
			}
			
			s.Seasons[seasonInfo.SeasonNumber] = season
		}

	default:
		return errors.ErrInternal.Msgf(
			"%T received the unexpected snapshot %T", s, snapshot)
	}
	return nil
}

// Helper method to extract season number from season ID
func (s *Series) getSeasonNumberFromID(seasonID string) int {
	for num, season := range s.Seasons {
		if season.SeasonID == seasonID {
			return num
		}
	}
	return -1
}