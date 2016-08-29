package weblog

import (
	"log"
	"net/http"

	_ "net/http/pprof"

	"github.com/gorilla/mux"
)

func newRouter(rh *requestHandler) *mux.Router {

	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:8099", nil))
	}()

	r := mux.NewRouter()
	r.HandleFunc("/healthz", rh.getHealthz).Methods("GET")
	r.HandleFunc("/healthz/", rh.getHealthz).Methods("GET")
	r.HandleFunc("/logs/{app}", rh.getLogs).Methods("GET")
	r.HandleFunc("/logs/{app}/", rh.getLogs).Methods("GET")
	r.HandleFunc("/logs/{app}", rh.deleteLogs).Methods("DELETE")
	r.HandleFunc("/logs/{app}/", rh.deleteLogs).Methods("DELETE")
	return r
}
