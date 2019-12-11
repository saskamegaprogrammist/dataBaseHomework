package handlers

import (
	"encoding/json"
	"github.com/saskamegaprogrammist/dataBaseHomework/models"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
	"net/http"
)

func CreateForum(writer http.ResponseWriter, req *http.Request) {
	var newForum models.Forum
	err := json.NewDecoder(req.Body).Decode(&newForum)
	if err != nil {
		log.Println(err)
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
