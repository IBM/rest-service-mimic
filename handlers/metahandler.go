package handlers

import (
	"container/ring"
	"net/http"
)

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type Metahandler struct {
	handlers *ring.Ring
}

func CreateMetahandler(handlers []Handler) Metahandler {
	handlersRing := ring.New(len(handlers))

	for i := 0; i < handlersRing.Len(); i++ {
		handlersRing.Value = handlers[i]
		handlersRing = handlersRing.Next()
	}

	return Metahandler{handlersRing}
}

func (metahandler *Metahandler) getCurrentHandler() Handler {
	return metahandler.handlers.Value.(Handler)
}

func (metahandler *Metahandler) Handle(w http.ResponseWriter, r *http.Request) {
	metahandler.getCurrentHandler().Handle(w, r)
	metahandler.handlers = metahandler.handlers.Next()
}
