package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ellofae/deanery-gateway/core/controller"
	"github.com/ellofae/deanery-gateway/core/controller/middleware"
	"github.com/ellofae/deanery-gateway/core/domain"
	"github.com/ellofae/deanery-gateway/core/dto"
	"github.com/ellofae/deanery-gateway/core/models"
	"github.com/ellofae/deanery-gateway/core/session"
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
	mux.HandleFunc("/users/", middleware.AuthenticateMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		parsed_url := strings.TrimPrefix(r.URL.Path, "/users/")
		url_parts := strings.Split(parsed_url, "/")

		switch r.Method {
		case http.MethodGet:
			// endpoint /users/signup - GET
			if len(url_parts) == 1 {
				if url_parts[0] == "signup" {
					err = h.handleShowRegistrationPage(w, r)
					if err != nil {
						h.handleError(w, r, err)
					}

					return
				} else if url_parts[0] == "login" {
					err = h.handleShowLoginPage(w, r)
					if err != nil {
						h.handleError(w, r, err)
					}

					return
				} else if url_parts[0] == "profile" {
					err = h.handleShowProfilePage(w, r)
					if err != nil {
						h.handleError(w, r, err)
					}

					return
				}
			} else if len(url_parts) == 2 {
				if url_parts[0] == "signup" && url_parts[1] == "success" {
					err = h.handleSuccessfulRegistration(w, r)
					if err != nil {
						h.handleError(w, r, err)
					}

					return
				}
			}
		case http.MethodPost:
			// endpoint /users/signup - POST
			if len(url_parts) == 1 {
				if url_parts[0] == "signup" {
					err = h.handleRegisterUser(w, r)
					if err != nil {
						h.handleError(w, r, err)
					}

					return
				} else if url_parts[0] == "login" {
					err = h.handleLoginUser(w, r)
					if err != nil {
						h.handleError(w, r, err)
					}

					return
				}
			}
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})))
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

	http.Redirect(w, r, fmt.Sprintf("/users/signup/success?user_name=%s&email=%s&phone=%s", user.UserName, user.Email, user.Phone), http.StatusFound)

	return nil
}

func (h *ClientHandler) handleSuccessfulRegistration(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, _ := template.ParseFiles("templates/reg_success.html")

	user := &dto.UserRegistered{
		UserName: r.URL.Query().Get("user_name"),
		Email:    r.URL.Query().Get("email"),
		Phone:    r.URL.Query().Get("phone"),
	}

	w.WriteHeader(http.StatusCreated)
	if err := tmpl.Execute(w, user); err != nil {
		return err
	}

	return nil
}

func (h *ClientHandler) handleShowLoginPage(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, _ := template.ParseFiles("templates/login.html")

	w.WriteHeader(http.StatusOK)

	if err := tmpl.Execute(w, nil); err != nil {
		return err
	}

	return nil
}

func (h *ClientHandler) handleLoginUser(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "application/json")

	var err error

	userCode, err := strconv.Atoi(r.FormValue("login"))
	if err != nil {
		return err
	}

	user := &models.UserLogin{
		RecordCode: userCode,
		Password:   r.FormValue("password"),
	}

	json_data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	response, err := http.Post("http://localhost:8000/api/login", "application/json", bytes.NewReader(json_data))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to login user, status: %v", response.StatusCode)
	}

	tokens := &dto.Tokens{}
	if err = json.NewDecoder(response.Body).Decode(tokens); err != nil {
		return err
	}

	store := session.SessionStorage()
	session, err := store.Get(r, "session")
	if err != nil {
		return err
	}

	session.Values["access_token"] = fmt.Sprintf("%s %s", "Bearer", tokens.AccessToken)
	err = session.Save(r, w)
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/users/profile", http.StatusFound)
	return nil
}

func (h *ClientHandler) handleShowProfilePage(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl *template.Template

	storage := session.SessionStorage()
	session, err := storage.Get(r, "session")
	if err != nil {
		return err
	}

	code, err := strconv.Atoi(session.Values["record_code"].(string))
	if err != nil {
		return err
	}

	recordCode := &dto.RecordCode{
		Code: code,
	}

	json_data, err := json.Marshal(recordCode)
	if err != nil {
		return err
	}

	response, err := http.Post("http://localhost:8000/api/get_username", "application/json", bytes.NewReader(json_data))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("received status is not http.StatusOK, status: %v", response.StatusCode)
	}

	profileInfo := dto.ProfileInformation{}
	if err = json.NewDecoder(response.Body).Decode(&profileInfo); err != nil {
		return err
	}

	switch session.Values["role"].(string) {
	case "student":
		tmpl, _ = template.ParseFiles("templates/student.html")
	case "professor":
		tmpl, _ = template.ParseFiles("templates/professor.html")
	}

	w.WriteHeader(http.StatusOK)
	if err := tmpl.Execute(w, profileInfo); err != nil {
		return err
	}

	return nil
}
