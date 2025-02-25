package api

import "fmt"

type OptFunc func(a *API)

func WithHTTPS(secure bool) OptFunc {
	return func(a *API) {
		a.secure = secure
	}
}

func WithPort(port int) OptFunc {
	return func(a *API) {
		a.server.Addr = fmt.Sprintf(":%d", port)
		a.port = port
	}
}
