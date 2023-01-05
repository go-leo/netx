package httpx

import (
	"net/http"
)

type CompositeHandler struct {
	handlers []*handlerWrapper
}

func (h *CompositeHandler) AddHandler(handler http.Handler, match func(request *http.Request) bool) {
	h.handlers = append(h.handlers, &handlerWrapper{handler: handler, match: match})
}

func (h *CompositeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	for _, handler := range h.handlers {
		if handler.match(request) {
			handler.handler.ServeHTTP(writer, request)
		}
	}
}

type handlerWrapper struct {
	handler http.Handler
	match   func(request *http.Request) bool
}

func (h *handlerWrapper) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.handler.ServeHTTP(writer, request)
}
