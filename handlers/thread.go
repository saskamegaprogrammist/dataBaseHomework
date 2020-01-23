package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/saskamegaprogrammist/dataBaseHomework/models"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
	"net/http"
	"strconv"
)

func CreateThread (writer http.ResponseWriter, req *http.Request) {
	var newThread models.Thread
	err := json.NewDecoder(req.Body).Decode(&newThread)
	if err != nil {
		//log.Println(err)
		utils.CreateAnswer(writer, 500, models.CreateError("cannot decode json"))
		return
	}
	forumSlug := mux.Vars(req)["slug"]
	newThread.Forum = forumSlug

	threadExists, err := newThread.CreateThread()
	if err != nil {
		if threadExists.Id != 0 {
			utils.CreateAnswer(writer, 409, threadExists)
			return
		} else {
			utils.CreateAnswer(writer, 404, models.CreateError(err.Error()))
		}
		return
	}
	utils.CreateAnswer(writer, 201, newThread)
}

func CreatePosts (writer http.ResponseWriter, req *http.Request) {
	var thread models.Thread
	threadInfo := mux.Vars(req)["slug_or_id"]
	threadId, err := strconv.Atoi(threadInfo)
	if err != nil {
		thread.Slug = threadInfo
	} else {
		thread.Id = int(threadId)
	}
	var newPosts [] models.Post
	err = json.NewDecoder(req.Body).Decode(&newPosts)
	if err != nil {
		//log.Println(err)
		utils.CreateAnswer(writer, 500, models.CreateError("cannot decode json"))
		return
	}
	posts, err, code := thread.CreatePosts(newPosts)
	if err != nil {
		//log.Println(err)
		switch code {
		case 1:
			utils.CreateAnswer(writer, 404, models.CreateError(err.Error()))
		case 2:
			utils.CreateAnswer(writer, 409, models.CreateError(err.Error()))
		case 3:
			utils.CreateAnswer(writer, 500, models.CreateError(err.Error()))
		}
		return
	}
	utils.CreateAnswer(writer, 201, posts)
}

func GetThread (writer http.ResponseWriter, req *http.Request) {
	var foundThread models.Thread
	threadInfo := mux.Vars(req)["slug_or_id"]
	threadId, err := strconv.Atoi(threadInfo)
	if err != nil {
		foundThread.Slug = threadInfo
	} else {
		foundThread.Id = int(threadId)
	}
	err = foundThread.GetThread()
	if err != nil {
		utils.CreateAnswer(writer, 404, models.CreateError(err.Error()))
		return
	}
	utils.CreateAnswer(writer, 200, foundThread)
}

func UpdateThread (writer http.ResponseWriter, req *http.Request) {
	var updatedThread models.Thread
	threadInfo := mux.Vars(req)["slug_or_id"]
	threadId, err := strconv.Atoi(threadInfo)
	if err != nil {
		updatedThread.Slug = threadInfo
	} else {
		updatedThread.Id = int(threadId)
	}
	err = json.NewDecoder(req.Body).Decode(&updatedThread)
	if err != nil {
		//log.Println(err)
		utils.CreateAnswer(writer, 500, models.CreateError("cannot decode json"))
		return
	}
	err = updatedThread.UpdateThread()
	if err != nil {
		//log.Println(err)
		utils.CreateAnswer(writer, 404, models.CreateError(err.Error()))
		return
	}
	utils.CreateAnswer(writer, 200, updatedThread)
}


func GetPostsByThread (writer http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	limit := query.Get("limit")
	since := query.Get("since")
	desc := query.Get("desc")
	sort := query.Get("sort")
	limitInt := -1
	descBool := false
	if limit != "" {
		limitInt, _ = strconv.Atoi(limit)
	}
	if desc == "true" {
		descBool = true
	}
	var thread models.Thread
	threadInfo := mux.Vars(req)["slug_or_id"]
	threadId, err := strconv.Atoi(threadInfo)
	if err != nil {
		thread.Slug = threadInfo
	} else {
		thread.Id = int(threadId)
	}

	posts := make([]models.Post, 0)
	posts, err = models.GetPostsByThread(limitInt, since, descBool, sort, thread)
	if err != nil {
		utils.CreateAnswer(writer, 404, models.CreateError(err.Error()))
		return
	}
	utils.CreateAnswer(writer, 200, posts)
}


func Vote(writer http.ResponseWriter, req *http.Request) {
	var newVote models.Vote
	var updatedThread models.Thread
	threadInfo := mux.Vars(req)["slug_or_id"]
	threadId, err := strconv.Atoi(threadInfo)
	if err != nil {
		updatedThread.Slug = threadInfo
	} else {
		updatedThread.Id = int(threadId)
	}
	err = json.NewDecoder(req.Body).Decode(&newVote)
	if err != nil {
		//log.Println(err)
		utils.CreateAnswer(writer, 500, models.CreateError("cannot decode json"))
		return
	}
	err = updatedThread.Vote(&newVote)
	if err != nil {
		utils.CreateAnswer(writer, 404, models.CreateError(err.Error()))
		return
	}
	utils.CreateAnswer(writer, 200, updatedThread)
}
