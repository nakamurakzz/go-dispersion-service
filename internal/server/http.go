package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func NewHTTPServer(addr string) *http.Server {
	httpsrv := newHttpServer()
	r := mux.NewRouter()
	r.HandleFunc("/", httpsrv.handleCreateRecords).Methods("POST")
	r.HandleFunc("/", httpsrv.handleFindRecords).Methods("GET")
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

type httpServer struct {
	Log *Log
}

func newHttpServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

type CreateRecordsRequest struct {
	Record Record `json:"record"`
}

type CreateRecordsResponse struct {
	Offset int64 `json:"offset"`
}

type FindRecordsRequest struct {
	Offset int64 `json:"offset"`
}

type FindRecordsResponse struct {
	Record Record `json:"record"`
}

func (s *httpServer) handleCreateRecords(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateRecordsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	offset, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := CreateRecordsResponse{Offset: offset}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) handleFindRecords(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req FindRecordsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	record, err := s.Log.Read(req.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res := FindRecordsResponse{Record: record}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
