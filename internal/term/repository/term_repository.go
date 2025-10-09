package repository

import (
	"context"
	"errors"
	"term-service/internal/term/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TermRepository interface {
	Create(ctx context.Context, term *model.Term) (*model.Term, error)
	GetByID(ctx context.Context, id string) (*model.Term, error)
	Update(ctx context.Context, id string, term *model.Term) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]*model.Term, error)
	GetCurrentTerm(ctx context.Context) (*model.Term, error)
	GetAllByOrgID(ctx context.Context, orgID string) ([]*model.Term, error)
	GetCurrentTermByOrg(ctx context.Context, organizationID string) (*model.Term, error)
	GetAllByOrgID4App(ctx context.Context, orgID string) ([]*model.Term, error)
	GetPreviousTerm(ctx context.Context, orgID string, termID string) (*model.Term, error)
}

type termRepository struct {
	collection *mongo.Collection
}

func NewTermRepository(collection *mongo.Collection) TermRepository {
	return &termRepository{collection}
}

// Create inserts a new term
func (r *termRepository) Create(ctx context.Context, term *model.Term) (*model.Term, error) {
	now := time.Now()
	term.CreatedAt = now
	term.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, term)
	if err != nil {
		return nil, err
	}
	return term, nil
}

// GetByID finds a term by its ID
func (r *termRepository) GetByID(ctx context.Context, id string) (*model.Term, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	var term model.Term
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&term)
	if err != nil {
		return nil, err
	}
	return &term, nil
}

// Update modifies an existing term
func (r *termRepository) Update(ctx context.Context, id string, updated *model.Term) error {
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
			"published_teacher": updated.PublishedTeacher,
			"published_parent":  updated.PublishedParent,
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

// Delete removes a term by ID
func (r *termRepository) Delete(ctx context.Context, id string) error {
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

// GetAll returns all terms
func (r *termRepository) GetAll(ctx context.Context) ([]*model.Term, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var terms []*model.Term
	for cursor.Next(ctx) {
		var term model.Term
		if err := cursor.Decode(&term); err != nil {
			return nil, err
		}
		terms = append(terms, &term)
	}
	return terms, nil
}

// GetCurrentTerm returns the current active term (where now is between start_date and end_date)
func (r *termRepository) GetCurrentTerm(ctx context.Context) (*model.Term, error) {
	now := time.Now()

	filter := bson.M{
		"start_date": bson.M{"$lte": now},
		"end_date":   bson.M{"$gte": now},
	}

	var term model.Term
	err := r.collection.FindOne(ctx, filter).Decode(&term)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // No current term found, not an error
		}
		return nil, err
	}

	return &term, nil
}

func (r *termRepository) GetCurrentTermByOrg(ctx context.Context, organizationID string) (*model.Term, error) {
	now := time.Now()

	filter := bson.M{
		"organization_id": organizationID,
		"start_date":      bson.M{"$lte": now},
		"end_date":        bson.M{"$gte": now},
	}

	var term model.Term
	err := r.collection.FindOne(ctx, filter).Decode(&term)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // No current term found, not an error
		}
		return nil, err
	}

	return &term, nil
}

func (r *termRepository) GetAllByOrgID(ctx context.Context, orgID string) ([]*model.Term, error) {
	filter := bson.M{
		"organization_id": orgID,
	}

	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var terms []*model.Term
	if err := cur.All(ctx, &terms); err != nil {
		return nil, err
	}

	return terms, nil
}

func (r *termRepository) GetAllByOrgID4App(ctx context.Context, orgID string) ([]*model.Term, error) {
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

	var terms []*model.Term
	if err := cur.All(ctx, &terms); err != nil {
		return nil, err
	}

	return terms, nil
}

func (r *termRepository) GetPreviousTerm(ctx context.Context, orgID string, termID string) (*model.Term, error) {
	// Lấy term hiện tại để biết start_date của nó
	objectID, err := primitive.ObjectIDFromHex(termID)
	if err != nil {
		return nil, err
	}

	var current model.Term
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID, "organization_id": orgID}).Decode(&current)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	// Tìm term có start_date < current.StartDate, sắp xếp giảm dần để lấy term liền trước
	filter := bson.M{
		"organization_id": orgID,
		"start_date":      bson.M{"$lt": current.StartDate},
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "start_date", Value: -1}})

	var prev model.Term
	err = r.collection.FindOne(ctx, filter, opts).Decode(&prev)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // không có term trước đó
		}
		return nil, err
	}

	return &prev, nil
}
