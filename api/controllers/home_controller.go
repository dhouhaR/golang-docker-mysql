package controllers

import (
	"net/http"

	"golang-docker-todo/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To TODO IT API")

}