package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsonrw"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/fx"
)

type RepoId string

type Base struct {
	PID RepoId `json:"_id,omitempty" bson:"_id,omitempty"`
	ID  RepoId `json:"id,omitempty" bson:"id,omitempty"`
}

func (r Base) GetId() RepoId {
	return r.ID
}

type RepoBaseI interface {
	GetId() RepoId
}

type Repo[T RepoBaseI] interface {
	InsertOne(ctx context.Context, doc *T) (bool, error)
	FindOne(ctx context.Context, filter interface{}) (*T, error)
	Find(ctx context.Context, filter interface{}) ([]*T, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*T, error)
	DeleteOne(ctx context.Context, filter interface{}) (interface{}, error)
}

type MongoRepo[T RepoBaseI] struct {
	collection *mongo.Collection
}

// DeleteOne implements Repo.
func (m *MongoRepo[T]) DeleteOne(ctx context.Context, filter interface{}) (interface{}, error) {
	return m.collection.DeleteOne(ctx, filter)
}

// Find implements Repo.
func (m *MongoRepo[T]) Find(ctx context.Context, filter interface{}) ([]*T, error) {

	c, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer c.Close(ctx)

	var docs []*T
	for c.Next(ctx) {
		var doc T

		dec, err := bson.NewDecoder(bsonrw.NewBSONDocumentReader(c.Current))
		if err != nil {
			return nil, err
		}

		if err := dec.Decode(&doc); err != nil {
			return nil, err
		}
		docs = append(docs, &doc)
	}

	if err := c.Err(); err != nil {
		return nil, err
	}

	return docs, nil
}

// FindOne implements Repo.
func (m *MongoRepo[T]) FindOne(ctx context.Context, filter interface{}) (*T, error) {
	c := m.collection.FindOne(ctx, filter)
	r, err := c.Raw()
	if err != nil {
		return nil, err
	}

	dec, err := bson.NewDecoder(bsonrw.NewBSONDocumentReader(r))
	if err != nil {
		return nil, err
	}

	var doc T
	if err := dec.Decode(&doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

// InsertOne implements Repo.
func (m *MongoRepo[T]) InsertOne(ctx context.Context, doc *T) (bool, error) {
	// u := uuid.New().String()
	// if doc == nil {
	// 	return false, fmt.Errorf("no data")
	// }

	_, err := m.collection.InsertOne(ctx, doc)
	if err != nil {
		return false, err
	}

	return true, nil
}

// UpdateOne implements Repo.
func (m *MongoRepo[T]) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*T, error) {
	var result *T

	c := m.collection.FindOneAndUpdate(ctx, filter, update)
	r, err := c.Raw()
	if err != nil {
		return nil, err
	}

	dec, err := bson.NewDecoder(bsonrw.NewBSONDocumentReader(r))
	if err != nil {
		return nil, err
	}

	if err := dec.Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func NewMongoRepoFx[T RepoBaseI](dbname string, collectionName string) fx.Option {
	return fx.Provide(func(client *mongo.Client) Repo[T] {
		collection := client.Database(dbname).Collection(collectionName)

		return &MongoRepo[T]{
			collection: collection,
		}
	})
}
