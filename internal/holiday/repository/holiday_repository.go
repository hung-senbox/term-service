package repository

import (
	"context"
	"errors"
	"term-service/internal/holiday/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HolidayRepository interface {
	Create(ctx context.Context, holiday *model.Holiday) (*model.Holiday, error)
	GetByID(ctx context.Context, id string) (*model.Holiday, error)
	Update(ctx context.Context, id string, holiday *model.Holiday) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]*model.Holiday, error)
	GetAllByOrgID(ctx context.Context, orgID string) ([]*model.Holiday, error)
	GetAllByOrgID4App(ctx context.Context, orgID string) ([]*model.Holiday, error)
}

type holidayRepository struct {
	collection *mongo.Collection
}

func NewHolidayRepository(collection *mongo.Collection) HolidayRepository {
	return &holidayRepository{collection}
}

// Create inserts a new holiday
func (r *holidayRepository) Create(ctx context.Context, holiday *model.Holiday) (*model.Holiday, error) {
	now := time.Now()
	holiday.CreatedAt = now
	holiday.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, holiday)
	if err != nil {
		return nil, err
	}
	return holiday, nil
}

// GetByID finds a holiday by its ID
func (r *holidayRepository) GetByID(ctx context.Context, id string) (*model.Holiday, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	var holiday model.Holiday
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&holiday)
	if err != nil {
		return nil, err
	}
	return &holiday, nil
}

// Update modifies an existing holiday
func (r *holidayRepository) Update(ctx context.Context, id string, updated *model.Holiday) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	updated.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"title":             updated.Title,
			"start_date":        updated.StartDate,
			"color":             updated.Color,
			"published_mobile":  updated.PublishedMobile,
			"published_desktop": updated.PublishedDesktop,
			"end_date":          updated.EndDate,
			"updated_at":        updated.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateByID(ctx, objectID, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// Delete removes a holiday by ID
func (r *holidayRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// GetAll returns all holidays
func (r *holidayRepository) GetAll(ctx context.Context) ([]*model.Holiday, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var holidays []*model.Holiday
	for cursor.Next(ctx) {
		var holiday model.Holiday
		if err := cursor.Decode(&holiday); err != nil {
			return nil, err
		}
		holidays = append(holidays, &holiday)
	}
	return holidays, nil
}

func (r *holidayRepository) GetAllByOrgID(ctx context.Context, orgID string) ([]*model.Holiday, error) {
	filter := bson.M{
		"organization_id": orgID,
	}

	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var holidays []*model.Holiday
	if err := cur.All(ctx, &holidays); err != nil {
		return nil, err
	}

	return holidays, nil
}

func (r *holidayRepository) GetAllByOrgID4App(ctx context.Context, orgID string) ([]*model.Holiday, error) {
	filter := bson.M{
		"organization_id":  orgID,
		"published_mobile": true,
	}

	// sort theo created_at ASC
	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})

	cur, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var holidays []*model.Holiday
	if err := cur.All(ctx, &holidays); err != nil {
		return nil, err
	}

	return holidays, nil
}
