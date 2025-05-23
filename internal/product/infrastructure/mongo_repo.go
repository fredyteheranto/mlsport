package infrastructure

import (
	"context"
	"log"
	"mlsport/config"
	"mlsport/internal/product/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoProductRepo struct {
	CollectionName string
}

func NewMongoProductRepo() *MongoProductRepo {
	return &MongoProductRepo{CollectionName: "products"}
}

func (r *MongoProductRepo) Create(p *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.GetDB().Collection(r.CollectionName)
	res, err := collection.InsertOne(ctx, p)
	if err != nil {
		return err
	}

	p.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *MongoProductRepo) FindAll() ([]domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result []domain.Product
	collection := config.GetDB().Collection(r.CollectionName)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("Error closing cursor: %v", err)
		}
	}()

	for cursor.Next(ctx) {
		var p domain.Product
		if err := cursor.Decode(&p); err != nil {
			continue
		}
		p.ID = p.ObjectID.Hex()
		result = append(result, p)
	}

	return result, nil
}

func (r *MongoProductRepo) FindByID(id string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var p domain.Product
	collection := config.GetDB().Collection(r.CollectionName)
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&p)
	if err != nil {
		return nil, err
	}

	p.ID = p.ObjectID.Hex()
	return &p, nil
}

func (r *MongoProductRepo) FindByCategory(cat string) ([]domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result []domain.Product
	collection := config.GetDB().Collection(r.CollectionName)
	cursor, err := collection.Find(ctx, bson.M{"category": cat})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("Error closing cursor: %v", err)
		}
	}()

	for cursor.Next(ctx) {
		var p domain.Product
		if err := cursor.Decode(&p); err != nil {
			continue
		}
		p.ID = p.ObjectID.Hex()
		result = append(result, p)
	}

	return result, nil
}

func (r *MongoProductRepo) GetCategories() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := config.GetDB().Collection(r.CollectionName)

	pipeline := []bson.M{
		{"$group": bson.M{"_id": "$category"}},
		{"$sort": bson.M{"_id": 1}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("Error closing cursor: %v", err)
		}
	}()

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	var categories []string
	for _, item := range result {
		if cat, ok := item["_id"].(string); ok {
			categories = append(categories, cat)
		}
	}

	return categories, nil
}

func (r *MongoProductRepo) Update(p *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		return err
	}

	collection := config.GetDB().Collection(r.CollectionName)
	_, err = collection.ReplaceOne(ctx, bson.M{"_id": objID}, p)
	return err
}

func (r *MongoProductRepo) Patch(id string, fields map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := config.GetDB().Collection(r.CollectionName)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": fields})
	return err
}

func (r *MongoProductRepo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := config.GetDB().Collection(r.CollectionName)
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (r *MongoProductRepo) GetMetrics() (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := config.GetDB().Collection(r.CollectionName)

	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":           nil,
				"total":         bson.M{"$sum": 1},
				"stock":         bson.M{"$sum": "$stock"},
				"average_price": bson.M{"$avg": "$price"},
			},
		},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("Error closing cursor: %v", err)
		}
	}()

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil || len(result) == 0 {
		return nil, err
	}

	topCatPipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":   "$category",
				"count": bson.M{"$sum": 1},
			},
		},
		{"$sort": bson.M{"count": -1}},
		{"$limit": 2},
	}

	catCursor, err := coll.Aggregate(ctx, topCatPipeline)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := catCursor.Close(ctx); err != nil {
			log.Printf("Error closing catCursor: %v", err)
		}
	}()

	var categories []bson.M
	if err := catCursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	var topCategories []string
	for _, cat := range categories {
		if name, ok := cat["_id"].(string); ok {
			topCategories = append(topCategories, name)
		}
	}

	data := map[string]interface{}{
		"total_products": result[0]["total"],
		"total_stock":    result[0]["stock"],
		"average_price":  result[0]["average_price"],
		"top_categories": topCategories,
	}

	return data, nil
}
