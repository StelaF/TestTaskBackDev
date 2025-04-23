package http_transport

import (
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"net"
	"net/http"
	"strings"
)

func (s *Server) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

func (s *Server) refreshHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		s.log.Warn(err.Error())

		http.Error(w, "cannot determine client IP", http.StatusBadRequest)
	}

	req := reqRefresh{}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		s.log.Warn(err.Error())

		http.Error(w, "cannot decode request", http.StatusBadRequest)
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		s.log.Warn("authorization header is missing")

		http.Error(w, "authorization header required", http.StatusUnauthorized)
		return
	}

	decodToken, err := base64.StdEncoding.DecodeString(req.Token)
	if err != nil {
		s.log.Error(err.Error())
		http.Error(w, "invalid refresh token", http.StatusBadRequest)
		return
	}

	splited := strings.Split(string(decodToken), ":")
	if len(splited) < 3 {
		s.log.Error("invalid refresh token")
		http.Error(w, "invalid refresh token", http.StatusBadRequest)
	}

	tokens, err := s.service.Refresh(ctx, splited, authHeader, ip)
	if err != nil {
		s.log.Error(err.Error())

		http.Error(w, "Error accessing user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tokens)
	if err != nil {
		s.log.Error("Encode token pairs failed")

		http.Error(w, "Error accessing user", http.StatusInternalServerError)
		return
	}

	return
}

func (s *Server) accessHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		s.log.Warn(err.Error())

		http.Error(w, "cannot determine client IP", http.StatusBadRequest)
	}

	userGuidString := r.URL.Query().Get("user_guid")
	if userGuidString == "" {
		s.log.Warn("User GUID empty")

		http.Error(w, "User GUID empty", http.StatusBadRequest)
		return
	}

	userGuid, err := uuid.Parse(userGuidString)

	tokens, err := s.service.Access(ctx, userGuid, ip)
	if err != nil {
		s.log.Error(err.Error())

		http.Error(w, "Error accessing user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tokens)
	if err != nil {
		s.log.Error("Encode token pairs failed")

		http.Error(w, "Error accessing user", http.StatusInternalServerError)
		return
	}

	return
}
