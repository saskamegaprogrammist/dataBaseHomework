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
