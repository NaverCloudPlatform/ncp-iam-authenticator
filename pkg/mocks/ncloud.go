package mocks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
)

var (
	clustersUuidGet, _           = regexp.Compile("/clusters/([^/]+)")
	clustersUuidKubeconfigGet, _ = regexp.Compile("/clusters/([^/]+)/kubeconfig")
)

func newMockServer(clustersUuidGetData string, clustersUuidKubeconfigGetData string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			switch {
			case clustersUuidGet.MatchString(request.URL.Path):
				_, err := fmt.Fprintln(writer, clustersUuidGetData)
				if err != nil {
					panic(err)
				}
			case clustersUuidKubeconfigGet.MatchString(request.URL.Path):
				_, err := fmt.Fprintln(writer, clustersUuidKubeconfigGetData)
				if err != nil {
					panic(err)
				}
			default:
				panic("no match api path")
			}
		}
	}))
}
