package controller

import "net/http"

type IHandler interface {
	RegisterHandlers(*http.ServeMux)
}
