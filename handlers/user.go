package handlers

import (
	"encoding/json"
	"github.com/saskamegaprogrammist/dataBaseHomework/models"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
	"net/http"
)

func GetUser(writer http.ResponseWriter, req *http.Request) {

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
	return
}
