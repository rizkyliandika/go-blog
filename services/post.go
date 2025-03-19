package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Post struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string    `json:"title,omitempty" bson:"title,omitempty"`
	Content   string    `json:"content,omitempty" bson:"content,omitempty"`
	Category  string    `json:"category,omitempty" bson:"category,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Tags      []string  `json:"tags,omitempty" bson:"tags,omitempty"`
}

var client *mongo.Client

func New(mongo *mongo.Client) Post {
	client = mongo

	return Post{}
}

func returnCollectionPointer(collection string) *mongo.Collection {
	return client.Database("blog").Collection(collection)
}

func (p *Post) InsertPost(entry Post) error {
	collection := returnCollectionPointer("post")

	_, err := collection.InsertOne(context.TODO(), Post{
		Title:     entry.Title,
		Content:   entry.Content,
		Category:  entry.Category,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Tags:      entry.Tags,
	})

	if err != nil {
		log.Println("Error: ", err)
		return err
	}

	return nil
}

func (p *Post) GetAllPost() ([]Post, error) {
	collection := returnCollectionPointer("post")
	var posts []Post

	cursor, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var post Post
		cursor.Decode(&post)
		posts = append(posts, post)
	}

	return posts, nil
}

func (p *Post) GetPostById(id string) (Post, error) {
	collection := returnCollectionPointer("post")
	var post Post

	mongoId, err := primitive.ObjectIDFromHex(id)
	log.Println(mongoId)
	if err != nil {
		return Post{}, err
	}

	err = collection.FindOne(context.Background(), bson.M{"_id": mongoId}).Decode(&post)
	if err != nil {
		log.Println(err)
		return Post{}, err
	}

	return post, nil
}

func (p *Post) UpdatePost(entry Post) (*mongo.UpdateResult, error) {
	collection := returnCollectionPointer("post")
	mongoId, err := primitive.ObjectIDFromHex(entry.ID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: mongoId}}

	update := bson.D{{
		Key: "$set", Value: bson.M{
			"title":      entry.Title,
			"content":    entry.Content,
			"category":   entry.Category,
			"updated_at": time.Now(),
			"tags":       entry.Tags,
		},
	}}

	res, err := collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return res, nil
}

func (p *Post) DeletePost(id string) error {
	collection := returnCollectionPointer("post")
	mongoId, err := primitive.ObjectIDFromHex(id)

	filter := bson.D{{Key: "_id", Value: mongoId}}
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
