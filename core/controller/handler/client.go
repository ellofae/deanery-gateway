package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/ellofae/deanery-gateway/core/controller"
	"github.com/ellofae/deanery-gateway/core/domain"
	"github.com/ellofae/deanery-gateway/core/models"
	"github.com/ellofae/deanery-gateway/pkg/logger"
)

type ClientHandler struct {
	logger  *log.Logger
	usecase domain.IClientUsecase
}

func NewClientHandler(clientUsecase domain.IClientUsecase) controller.IHandler {
	return &ClientHandler{
		logger:  logger.GetLogger(),
		usecase: clientUsecase,
	}
}

func (h *ClientHandler) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		var err error

		parsed_url := strings.TrimPrefix(r.URL.Path, "/users/")
		url_parts := strings.Split(parsed_url, "/")

		switch r.Method {
		case http.MethodGet:
			// endpoint /users/signup - GET
			if len(url_parts) == 1 && url_parts[0] == "signup" {
				err = h.handleShowRegistrationPage(w, r)
				if err != nil {
					h.handleError(w, r, err)
				}

				return
			} else if len(url_parts) == 2 && url_parts[0] == "signup" && url_parts[1] == "success" {

				err = h.handleSuccessfulRegistration(w, r)
				if err != nil {
					h.handleError(w, r, err)
				}

				return
			}
		case http.MethodPost:
			// endpoint /users/signup - POST
			if len(url_parts) == 1 && url_parts[0] == "signup" {
				err = h.handleRegisterUser(w, r)
				if err != nil {
					h.handleError(w, r, err)
				}

				return
			}
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
}

func (h *ClientHandler) handleError(w http.ResponseWriter, r *http.Request, err_occured error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.WriteHeader(http.StatusInternalServerError)

	tmpl, _ := template.ParseFiles("templates/error.html")
	_ = tmpl.Execute(w, models.APIError{ErrorMsg: err_occured.Error()})
}

func (h *ClientHandler) handleShowRegistrationPage(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	existingRoles := []models.Roles{}
	response, err := http.Get("http://localhost:8000/api/roles")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&existingRoles); err != nil {
		return err
	}

	tmpl, _ := template.ParseFiles("templates/registration.html")

	w.WriteHeader(http.StatusOK)

	if err := tmpl.Execute(w, models.RolesSelection{Roles: existingRoles}); err != nil {
		return err
	}

	return nil
}

func (h *ClientHandler) handleRegisterUser(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "application/json")

	var err error

	phoneNumber := r.FormValue("phone")
	if !strings.Contains(phoneNumber, "+") {
		phoneNumber = "+" + phoneNumber
	}

	user := &models.UserRegistration{
		UserName:   r.FormValue("user_name"),
		Email:      r.FormValue("email"),
		Phone:      phoneNumber,
		UserStatus: r.FormValue("user_status"),
	}

	json_data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	response, err := http.Post("http://localhost:8000/api/signup", "application/json", bytes.NewReader(json_data))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to register user, status: %v", response.StatusCode)
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h *ClientHandler) handleSuccessfulRegistration(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, _ := template.ParseFiles("templates/registration_success.html")

	w.WriteHeader(http.StatusCreated)
	if err := tmpl.Execute(w, nil); err != nil {
		return err
	}

	return nil
}
