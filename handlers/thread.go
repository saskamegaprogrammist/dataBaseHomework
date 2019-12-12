package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/saskamegaprogrammist/dataBaseHomework/models"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

func CreateThread (writer http.ResponseWriter, req *http.Request) {
	var newThread models.Thread
	err := json.NewDecoder(req.Body).Decode(&newThread)
	if err != nil {
		log.Println(err)
		utils.CreateAnswer(writer, 500, models.CreateError("cannot decode json"))
		return
	}
	forumSlug := mux.Vars(req)["slug"]
	newThread.Forum = forumSlug
	if newThread.Slug == "" {
		title := newThread.Title
		title = strings.ReplaceAll(title, " ", "-")
		newThread.Slug = title + strconv.Itoa(rand.Intn(1000))
	}
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
		thread.Id = int32(threadId)
	}
	var newPosts [] models.Post
	err = json.NewDecoder(req.Body).Decode(&newPosts)
	if err != nil {
		log.Println(err)
		utils.CreateAnswer(writer, 500, models.CreateError("cannot decode json"))
		return
	}
	posts, err, code := thread.CreatePosts(newPosts)
	if err != nil {
		log.Println(err)
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
		foundThread.Id = int32(threadId)
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
		updatedThread.Id = int32(threadId)
	}
	err = json.NewDecoder(req.Body).Decode(&updatedThread)
	if err != nil {
		log.Println(err)
		utils.CreateAnswer(writer, 500, models.CreateError("cannot decode json"))
		return
	}
	err = updatedThread.UpdateThread()
	if err != nil {
		log.Println(err)
		utils.CreateAnswer(writer, 404, models.CreateError(err.Error()))
		return
	}
	utils.CreateAnswer(writer, 200, updatedThread)
}