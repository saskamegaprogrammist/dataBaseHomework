package handlers

import (
	"github.com/saskamegaprogrammist/dataBaseHomework/models"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"net/http"
)

func GetStatus (writer http.ResponseWriter, req *http.Request) {
	var dataBase models.DB
	err := dataBase.GetStatus()
	if err != nil {
		utils.CreateAnswer(writer, 404, models.CreateError(err.Error()))
		return
	}
	utils.CreateAnswer(writer, 200, dataBase)
}

func Clear (writer http.ResponseWriter, req *http.Request) {
	var dataBase models.DB
	err := dataBase.Clear()
	if err != nil {
		utils.CreateAnswer(writer, 404, models.CreateError(err.Error()))
		return
	}
	utils.CreateAnswer(writer, 200, dataBase)
}