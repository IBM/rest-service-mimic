package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Code    int                    `json:"status_code"`
	Payload map[string]interface{} `json:"payload"`
}

type Route struct {
	Path        string            `json:"path"`
	Methods     []string          `json:"methods"`
	Headers     map[string]string `json:"headers"`
	Response    Response          `json:"response"`
	QueryParams map[string]string `json:"query_params"`
	CacheKey    string            `json:"cache_key"`
}

type Metaroute interface {
	Handle(w http.ResponseWriter, r *http.Request)
	GetPath() string
	GetMethods() []string
	GetHeaders() []string
	GetQueryParams() []string
}

type metaroute struct {
	config        Route
	instanceCache map[string]interface{}
}

func (meta metaroute) GetPath() string {
	return meta.config.Path
}

func (meta metaroute) GetMethods() []string {
	return meta.config.Methods
}

func (meta metaroute) GetHeaders() []string {
	headerPairs := []string{}

	for key, value := range meta.config.Headers {
		headerPairs = append(headerPairs, key, value)
	}

	return headerPairs
}

func (meta metaroute) GetQueryParams() []string {
	queryParamPairs := []string{}

	for key, value := range meta.config.QueryParams {
		queryParamPairs = append(queryParamPairs, key, value)
	}

	return queryParamPairs

}

func (meta metaroute) Handle(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received request [%s] %s -H %s -> %s \n", r.Method, r.RequestURI, r.Header, string(body))

	var requestMap map[string]interface{}
	json.Unmarshal(body, &requestMap)

	if meta.config.CacheKey != "" {
		meta.instanceCache[meta.config.CacheKey] = requestMap
	}

	payloadGenerator := ResponsePayloadGenerator{requestMap, meta.instanceCache}
	json, err := json.Marshal(payloadGenerator.Generate(meta.config.Response.Payload))
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(meta.config.Response.Code)
	w.Write(json)
}

type RouteConfigParser interface {
	Parse(pathToConfig string) ([]Metaroute, error)
}

type routeConfigParser struct{}

func CreateRouteConfigParser() RouteConfigParser {
	return routeConfigParser{}
}

func (configParser routeConfigParser) Parse(pathToConfig string) ([]Metaroute, error) {
	instanceCache := make(map[string]interface{})

	jsonFile, err := os.Open(pathToConfig)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var routes []Route
	err = json.Unmarshal(byteValue, &routes)
	if err != nil {
		return nil, err
	}

	var metaroutes []Metaroute
	for _, config := range routes {
		metaroutes = append(metaroutes, metaroute{config, instanceCache})
	}
	return metaroutes, nil
}
