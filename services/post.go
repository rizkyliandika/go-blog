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

type UpdatedPost struct {
	TempData map[string]interface{}
}

// func (r *Post) Sanitize(model Post) UpdatedPost {
//     rValue := reflect.ValueOf(&model).Elem()
//     rType := rValue.Type()
// 	var newPost UpdatedPost

//     for i := 0; i < rType.NumField(); i++ {
// 		field := rType.Field(i)

// 		if field != "" {

// 		}
// 		jsonTag := field.Tag.Get("json")
// 		typeTag := field.Tag.Get("type")
// 		if jsonTag == item.Key && typeTag != "int" {
// 			rValue.FieldByName(field.Name).SetString(item.Value)
// 		}
// 		if jsonTag == item.Key && typeTag == "int" {
// 			if intValue, err := strconv.Atoi(item.Value); err == nil {
// 				rValue.FieldByName(field.Name).SetInt(int64(intValue))
// 			}
// 		}
// 	}
//     return newPost
// }

var client *mongo.Client

func New(mongo *mongo.Client) Post {
	client = mongo

	return Post{}
}

func returnCollectionPointer(collection string) *mongo.Collection {
	return client.Database("blog").Collection(collection)
}

func (p *Post) InsertPost(entry Post) (*mongo.InsertOneResult, error) {
	collection := returnCollectionPointer("post")

	result, err := collection.InsertOne(context.TODO(), Post{
		Title:     entry.Title,
		Content:   entry.Content,
		Category:  entry.Category,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Tags:      entry.Tags,
	})

	if err != nil {
		log.Println("Error: ", err)
		return nil, err
	}

	return result, nil
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

	mongoId, _ := primitive.ObjectIDFromHex(id)

	err := collection.FindOne(context.Background(), bson.M{"_id": mongoId}).Decode(&post)
	if err != nil {
		log.Fatal(err)
		return Post{}, err
	}

	return post, nil
}

func (p *Post) UpdatePost(entry Post) (*mongo.UpdateResult, error) {
	collection := returnCollectionPointer("post")
	mongoId, _ := primitive.ObjectIDFromHex(entry.ID)

	filter := bson.D{{Key: "_id", Value: mongoId}}
	bsonEntry, _ := bson.Marshal(bson.M{
		"title":      entry.Title,
		"content":    entry.Content,
		"category":   entry.Category,
		"updated_at": time.Now(),
		"tags":       entry.Tags,
	})
	var test Post
	err := bson.Unmarshal(bsonEntry, &test)

	update := bson.M{"$set": test}
	res, err := collection.UpdateOne(
		context.TODO(),
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
