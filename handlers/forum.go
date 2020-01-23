package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/saskamegaprogrammist/dataBaseHomework/models"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	//"log"
	"net/http"
	"strconv"
)

func CreateForum(writer http.ResponseWriter, req *http.Request) {
	var newForum models.Forum
	err := json.NewDecoder(req.Body).Decode(&newForum)
	if err != nil {
		//log.Println(err)
		utils.CreateAnswer(writer, 500, models.CreateError("cannot decode json"))
		return
	}
	forumExists, err := newForum.CreateForum()
	if err != nil {
		if forumExists.Id != 0 {
			utils.CreateAnswer(writer, 409, forumExists)
			return
		} else {
			utils.CreateAnswer(writer, 404, models.CreateError(err.Error()))
		}
		return
	}
	utils.CreateAnswer(writer, 201, newForum)
}

func GetForum(writer http.ResponseWriter, req *http.Request) {
	var foundForum models.Forum
	forumSlug := mux.Vars(req)["slug"]
	err := foundForum.GetForum(forumSlug)
	if err != nil {
		utils.CreateAnswer(writer, 404, models.CreateError(err.Error()))
		return
	}
	utils.CreateAnswer(writer, 200, foundForum)
}

func GetThreadsByForum (writer http.ResponseWriter, req *http.Request) {
	forumSlug := mux.Vars(req)["slug"]
	query := req.URL.Query()
	limit := query.Get("limit")
	since := query.Get("since")
	desc := query.Get("desc")
	limitInt := -1
	descBool := false
	if limit != "" {
		limitInt, _ = strconv.Atoi(limit)
	}
	if desc == "true" {
		descBool = true
	}
	threads := make([]models.Thread, 0)
	threads, err := models.GetThreadsByForum(limitInt, since, descBool, forumSlug)
	if err != nil {
		utils.CreateAnswer(writer, 404, models.CreateError(err.Error()))
		return
	}
	utils.CreateAnswer(writer, 200, threads)

}

func GetUsersByForum (writer http.ResponseWriter, req *http.Request) {
	users := make([]models.User, 0)
	forumSlug := mux.Vars(req)["slug"]
	query := req.URL.Query()
	limit := query.Get("limit")
	since := query.Get("since")
	desc := query.Get("desc")
	limitInt := -1
	descBool := false
	if limit != "" {
		limitInt, _ = strconv.Atoi(limit)
	}
	if desc == "true" {
		descBool = true
	}
	users, err := models.GetUsersByForum(limitInt, since, descBool, forumSlug)
	if err != nil {
		utils.CreateAnswer(writer, 404, models.CreateError(err.Error()))
		return
	}
	utils.CreateAnswer(writer, 200, users)
}