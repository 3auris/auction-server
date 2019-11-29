package api

import (
	"net/http"
)

func (s Server) handleRegistration() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		Success bool `json:"success"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var request request
		s.decode(w, r, request)
		response := response{Success: false}
		s.encodeAndRespond(w, r, response, 200)
	}
}

func (s Server) handleAuthentication() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		Token  string `json:"token"`
		Expire string `json:"expire"`
		Email  string `json:"email"`
	}

	var unauthorized = struct{ Error string `json:"error"` }{Error: "wrong credentials"}

	return func(w http.ResponseWriter, r *http.Request) {
		var request request
		s.decode(w, r, &request)

		session, err := s.app.Auth(request.Email, request.Password)
		if err != nil {
			s.encodeAndRespond(w, r, unauthorized, http.StatusBadRequest)
			return
		}

		s.encodeAndRespond(w, r, response{
			Token:  string(session.GetToken()),
			Expire: session.GetExpire(),
			Email:  request.Email,
		}, http.StatusOK)
	}
}