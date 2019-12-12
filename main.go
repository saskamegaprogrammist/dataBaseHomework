package main
import (
	"github.com/gorilla/mux"
	"github.com/saskamegaprogrammist/dataBaseHomework/handlers"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
	"net/http"
)



func main() {
	//utils.CreateDataBaseConnection("docker", "docker", "localhost", "docker", 20);
	utils.CreateDataBaseConnection("postgres", "1", "localhost", "project_techno_real", 20);
	utils.InitDataBase();

	r := mux.NewRouter()
	r.HandleFunc("/api/user/{nickname}/create", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/api/user/{nickname}/profile", handlers.GetUser).Methods("GET")
	r.HandleFunc("/api/user/{nickname}/profile", handlers.UpdateUser).Methods("POST")

	r.HandleFunc("/api/forum/create", handlers.CreateForum).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/create", handlers.CreateThread).Methods("POST")

	r.HandleFunc("/api/forum/{slug}/details", handlers.GetForum).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/threads", handlers.GetThreadsByForum).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/users", handlers.GetUsersByForum).Methods("GET")

	r.HandleFunc("/api/thread/{slug_or_id}/create", handlers.CreatePosts).Methods("POST")
	//r.HandleFunc("/api/thread/{slug_or_id}/vote", handlers.Vote).Methods("POST")
	r.HandleFunc("/api/thread/{slug_or_id}/details", handlers.GetThread).Methods("GET")
	//r.HandleFunc("/api/thread/{slug_or_id}/posts", handlers.GetPosts).Methods("GET")
	r.HandleFunc("/api/thread/{slug_or_id}/details", handlers.UpdateThread).Methods("POST")
	//
	//r.HandleFunc("/api/post/{id:[0-9]+}/details", handlers.GetPost).Methods("GET")
	//r.HandleFunc("/api/post/{id:[0-9]+}/details", handlers.UpdatePost).Methods("POST")
	//
	//r.HandleFunc("/api/service/status", handlers.GetStatus).Methods("GET")
	//r.HandleFunc("/api/service/clear", handlers.Clear).Methods("POST")

	err := http.ListenAndServe(":5000", r)
	if err != nil {
		log.Fatal(err)
		return
	}

}
