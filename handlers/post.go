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

func GetPost(writer http.ResponseWriter, req *http.Request) {
	var foundPost models.Post
	postId, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		//log.Println(err)
	}
	foundPost.Id = postId
	query := req.URL.Query()
	related := query.Get("related")
	postRelated, err := foundPost.GetPostRelated(related)
	if err != nil {
		utils.CreateAnswer(writer, 404, models.CreateError(err.Error()))
		return
	}
	utils.CreateAnswer(writer, 200, postRelated)

}

func UpdatePost(writer http.ResponseWriter, req *http.Request) {
	var updatedPost models.Post
	postId, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		//log.Println(err)
	}
	updatedPost.Id = postId
	err = json.NewDecoder(req.Body).Decode(&updatedPost)
	if err != nil {
		//log.Println(err)
		utils.CreateAnswer(writer, 500, models.CreateError("cannot decode json"))
		return
	}
	err = updatedPost.UpdatePost()
	if err != nil {
		//log.Println(err)
		utils.CreateAnswer(writer, 404, models.CreateError(err.Error()))
		return
	}
	utils.CreateAnswer(writer, 200, updatedPost)

}

