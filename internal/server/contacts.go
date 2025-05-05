package server

import (
	"Contacter/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (s *Server) getContactsHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page := 0
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}

	searchString := r.URL.Query().Get("search")

	matchAllStr := r.URL.Query().Get("matchAll")
	matchAll := false
	if matchAllStr != "" {
		matchAll, _ = strconv.ParseBool(matchAllStr)
	}

	contacts, count, err := s.db.GetContacts(searchString, matchAll, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Contacts []models.Contact `json:"contacts"`
		Count    int              `json:"count"`
	}{
		Contacts: contacts,
		Count:    count,
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) putContactsHandler(w http.ResponseWriter, r *http.Request) {
	var contact models.Contact
	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.db.EditContact(contact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) postContactsHandler(w http.ResponseWriter, r *http.Request) {
	var contact models.Contact
	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contact, err = s.db.CreateContact(contact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(contact)
}

func (s *Server) deleteContactsHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := s.db.DeleteContact(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
