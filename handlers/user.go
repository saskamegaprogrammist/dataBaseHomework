package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/saskamegaprogrammist/dataBaseHomework/models"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
	"net/http"
)

func GetUser(writer http.ResponseWriter, req *http.Request) {
	var foundUser models.User
	userNickname := mux.Vars(req)["nickname"]
	err := foundUser.GetUser(userNickname)
	if err != nil {
		utils.CreateAnswer(writer, 404, models.CreateError(err.Error()))
		return
	}
	utils.CreateAnswer(writer, 200, foundUser)
}

func CreateUser(writer http.ResponseWriter, req *http.Request) {
	var newUser models.User
	err := json.NewDecoder(req.Body).Decode(&newUser)
	if err != nil {
		log.Println(err)
		utils.CreateAnswer(writer, 500, models.CreateError("cannot decode json"))
		return
	}
	usersExisting, err := newUser.CreateUser()
	if err != nil {
		log.Println(err)
		utils.CreateAnswer(writer, 500, models.CreateError("internal error"))
		return
	}
	if usersExisting != nil {
		utils.CreateAnswer(writer, 409, usersExisting)
		return
	}
	utils.CreateAnswer(writer, 201, newUser)
}

func UpdateUser (writer http.ResponseWriter, req *http.Request) {
	var newUser models.User
	userNickname := mux.Vars(req)["nickname"]
	newUser.Nickname = userNickname
	err := json.NewDecoder(req.Body).Decode(&newUser)
	if err != nil {
		log.Println(err)
		utils.CreateAnswer(writer, 500, models.CreateError("cannot decode json"))
		return
	}
	err, code := newUser.UpdateUser()
	if err != nil {
		log.Println(err)
		switch code {
		case 1:
			utils.CreateAnswer(writer, 404, models.CreateError(err.Error()))
		case 2:
			utils.CreateAnswer(writer, 409, models.CreateError(err.Error()))
		}
		return
	}
	utils.CreateAnswer(writer, 200, newUser)
}