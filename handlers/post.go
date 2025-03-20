package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rizkyliandika/go-blog/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	res := Response{
		Msg:  "Health Check",
		Code: 200,
	}

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		log.Panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(jsonResponse)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	var post services.Post
	posts, err := post.GetAllPost()
	if err != nil {
		log.Println(err)
		return
	}

	var res Response

	if len(posts) == 0 {
		res = Response{
			Msg:  "Entity not exist",
			Code: http.StatusUnprocessableEntity,
			Data: services.Post{},
		}
	} else {
		res = Response{
			Msg:  "Successfully get all data",
			Code: http.StatusOK,
			Data: posts,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	json.NewEncoder(w).Encode(res)
}

func getPostById(w http.ResponseWriter, r *http.Request) {
	var post services.Post
	id := chi.URLParam(r, "id")

	post, err := post.GetPostById(id)

	if err != nil {
		log.Println(err)
		return
	}

	var res = Response{
		Msg:  "Successfully get data",
		Code: http.StatusOK,
		Data: post,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var post services.Post
	var insertResult *mongo.InsertOneResult
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		log.Fatal(err)
	}
	insertResult, err = post.InsertPost(post)

	if err != nil {
		errResponse := Response{
			Msg:  "Error",
			Code: http.StatusInternalServerError,
		}
		json.NewEncoder(w).Encode(errResponse)
		w.WriteHeader(errResponse.Code)
		return
	}

	res := Response{
		Msg:  "Successfully Created Post",
		Code: http.StatusCreated,
		Data: insertResult,
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(jsonStr)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	var post services.Post
	var updateResult *mongo.UpdateResult

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		log.Println(err)
		return
	}

	updateResult, err = post.UpdatePost(post)
	if err != nil {
		errorRes := Response{
			Msg:  "Error update post",
			Code: http.StatusInternalServerError,
		}
		json.NewEncoder(w).Encode(errorRes)
		w.WriteHeader(errorRes.Code)
		return
	}

	res := Response{
		Msg:  "Successfully updated",
		Code: http.StatusOK,
		Data: updateResult,
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(jsonStr)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	var post services.Post
	id := chi.URLParam(r, "id")

	err := post.DeletePost(id)
	if err != nil {
		errorRes := Response{
			Msg:  "Error delete post",
			Code: http.StatusInternalServerError,
		}
		json.NewEncoder(w).Encode(errorRes)
		w.WriteHeader(errorRes.Code)
		return
	}

	res := Response{
		Msg:  "Successfully deleted",
		Code: http.StatusOK,
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(jsonStr)
}
