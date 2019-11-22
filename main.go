package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func FallthroughHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received **UNMAPPED** request [%s] %s -H %s \n", r.Method, r.RequestURI, r.Header)
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	jsonPathPtr := flag.String("config", "config.json", "Path to th json mimic config")
	portPtr := flag.Int("port", 8000, "Port to listen on")
	flag.Parse()

	fmt.Println("Starting the Four Legged Mimic")

	configParser := CreateRouteConfigParser()
	metaroutes, err := configParser.Parse(*jsonPathPtr)

	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	for _, metaroute := range metaroutes {
		applyMetaroute(r, metaroute)
	}

	//r.HandleFunc("/", FallthroughHandler)
	r.NotFoundHandler = http.HandlerFunc(FallthroughHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *portPtr), r))
}

func applyMetaroute(router *mux.Router, metaroute Metaroute) {
	fmt.Printf("Adding %s -> %s \n", metaroute.GetMethods(), metaroute.GetPath())
	router.HandleFunc(metaroute.GetPath(), metaroute.Handle).Methods(metaroute.GetMethods()...).HeadersRegexp(metaroute.GetHeaders()...).Queries(metaroute.GetQueryParams()...)
}
