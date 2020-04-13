package main
import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/saskamegaprogrammist/dataBaseHomework/handlers"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)
func AccessLogMiddleware (mux *mux.Router,) http.HandlerFunc   {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
		rps.Inc()
	})
}

var (

	rps =prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "rps_total",
		})
)


func init() {

	prometheus.MustRegister(rps)

}


func main() {

	utils.CreateDataBaseConnection("alex", "alex", "95.163.251.171", "alex", 20)
	//utils.CreateDataBaseConnection("postgres", "1", "localhost", "project_techno_real", 20);
	//utils.InitDataBase();
	http.Handle("/metrics", promhttp.Handler())
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
	r.HandleFunc("/api/thread/{slug_or_id}/vote", handlers.Vote).Methods("POST")
	r.HandleFunc("/api/thread/{slug_or_id}/details", handlers.GetThread).Methods("GET")
	r.HandleFunc("/api/thread/{slug_or_id}/posts", handlers.GetPostsByThread).Methods("GET")
	r.HandleFunc("/api/thread/{slug_or_id}/details", handlers.UpdateThread).Methods("POST")

	r.HandleFunc("/api/post/{id:[0-9]+}/details", handlers.GetPost).Methods("GET")
	r.HandleFunc("/api/post/{id:[0-9]+}/details", handlers.UpdatePost).Methods("POST")

	r.HandleFunc("/api/service/status", handlers.GetStatus).Methods("GET")
	r.HandleFunc("/api/service/clear", handlers.Clear).Methods("POST")

	siteHandler := AccessLogMiddleware(r)

	http.Handle("/", siteHandler)

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
		return
	}

}
