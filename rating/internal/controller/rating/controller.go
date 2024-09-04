package rating

import (
	"context"
	"errors"

	"github.com/Sahas001/movieapp/rating/internal/repository"
	model "github.com/Sahas001/movieapp/rating/pkg"
)

var ErrNotFound = errors.New("rating not found for the record")

// type RatingInput struct {
// 	recordID   model.RecordID
// 	recordType model.RecordType
// 	rating     *model.Rating
// }

type ratingRepository interface {
	Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error)
	Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

// Controller defines the rating service controller
type Controller struct {
	repo ratingRepository
}

// New creates a rating service controller
func New(repo ratingRepository) *Controller {
	return &Controller{repo}
}

func (c *Controller) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	ratings, err := c.repo.Get(ctx, recordID, recordType)
	if err != nil && err == repository.ErrNotFound {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}
	sum := float64(0)
	for _, r := range ratings {
		sum += float64(r.Value)
	}
	return sum / float64(len(ratings)), nil
}

func (c *Controller) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return c.repo.Put(ctx, recordID, recordType, rating)
}
