package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/traefik/traefik/v2/pkg/config/static"
	"github.com/traefik/traefik/v2/pkg/log"
)

type entryPointRepresentation struct {
	*static.EntryPoint
	Name string `json:"name,omitempty"`
}

func (h Handler) getEntryPoints(rw http.ResponseWriter, request *http.Request) {
	results := make([]entryPointRepresentation, 0, len(h.staticConfig.EntryPoints))

	for name, ep := range h.staticConfig.EntryPoints {
		results = append(results, entryPointRepresentation{
			EntryPoint: ep,
			Name:       name,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	rw.Header().Set("Content-Type", "application/json")

	pageInfo, err := pagination(request, len(results))
	if err != nil {
		writeError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set(nextPageHeader, strconv.Itoa(pageInfo.nextPage))

	err = json.NewEncoder(rw).Encode(results[pageInfo.startIndex:pageInfo.endIndex])
	if err != nil {
		log.FromContext(request.Context()).Error(err)
		writeError(rw, err.Error(), http.StatusInternalServerError)
	}
}

func (h Handler) getEntryPoint(rw http.ResponseWriter, request *http.Request) {
	scapedEntryPointID := mux.Vars(request)["entryPointID"]

	entryPointID, err := url.PathUnescape(scapedEntryPointID)
	if err != nil {
		writeError(rw, fmt.Sprintf("unable to decode entryPointID %q: %s", scapedEntryPointID, err), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	ep, ok := h.staticConfig.EntryPoints[entryPointID]
	if !ok {
		writeError(rw, fmt.Sprintf("entry point not found: %s", entryPointID), http.StatusNotFound)
		return
	}

	result := entryPointRepresentation{
		EntryPoint: ep,
		Name:       entryPointID,
	}

	err = json.NewEncoder(rw).Encode(result)
	if err != nil {
		log.FromContext(request.Context()).Error(err)
		writeError(rw, err.Error(), http.StatusInternalServerError)
	}
}
