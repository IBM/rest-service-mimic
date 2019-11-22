package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"rest-service-mimic/routes"
)

type MockedHandler struct {
	Config        routes.Route
	InstanceCache map[string]interface{}
}

func (handler MockedHandler) Handle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received request [%s] %s -H %s -> %s \n\n", r.Method, r.RequestURI, r.Header, string(body))

	var requestMap map[string]interface{}
	json.Unmarshal(body, &requestMap)

	if handler.Config.CacheKey != "" {
		handler.InstanceCache[handler.Config.CacheKey] = requestMap
	}

	payloadGenerator := ResponsePayloadGenerator{requestMap, handler.InstanceCache}
	json, err := json.Marshal(payloadGenerator.Generate(handler.Config.Response.Payload))
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(handler.Config.Response.Code)
	w.Write(json)

}
