package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"rest-service-mimic/handlers"
	"rest-service-mimic/routes"
)

type Metaroute interface {
	Handle(w http.ResponseWriter, r *http.Request)
	GetPath() string
	GetMethods() []string
	GetHeaders() []string
	GetQueryParams() []string
}

type metaroute struct {
	config  routes.Route
	handler handlers.Metahandler
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
	meta.handler.Handle(w, r)
}

type RouteConfigParser interface {
	Parse(pathToConfig string) ([]Metaroute, error)
}

type routeConfigParser struct{}

func CreateRouteConfigParser() RouteConfigParser {
	return routeConfigParser{}
}

func (configParser routeConfigParser) Parse(pathToConfig string) ([]Metaroute, error) {
	jsonFile, err := os.Open(pathToConfig)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var routes []routes.Route
	err = json.Unmarshal(byteValue, &routes)
	if err != nil {
		return nil, err
	}

	var metaroutes []Metaroute
	for _, config := range routes {
		metaroutes = append(metaroutes, metaroute{config, createMetahandler(config)})
	}
	return metaroutes, nil
}

func createMetahandler(config routes.Route) handlers.Metahandler {
	var allHandlers []handlers.Handler
	if config.Response.Proxy.Host != "" {
		allHandlers = append(allHandlers, handlers.ProxyHandler{config})
	} else {
		instanceCache := make(map[string]interface{})
		allHandlers = append(allHandlers, handlers.MockedHandler{config, instanceCache})
	}

	return handlers.CreateMetahandler(allHandlers)
}
