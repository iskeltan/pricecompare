package views

import (
	"encoding/json"
	"log"
	"net/http"

	"../helpers"
	"../models"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

var (
	key   = []byte("coksecretharbiden")
	store = sessions.NewCookieStore(key)
)

func AuthenticateMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "profiles")
		session.Values["authenticated"] = false
		session.Values["user"] = nil
		f(w, r)
	}
}

// LoginUser a view
func LoginUser(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Message string       `json:"message"`
		Success bool         `json:"success"`
		Token   models.Token `json:"token"`
	}

	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if r.Method != http.MethodPost {
		resp := Response{Message: "this method only use post", Success: false}
		ResponseJSON, _ := json.Marshal(resp)
		w.Write(ResponseJSON)
		return
	}

	decoder := json.NewDecoder(r.Body)
	request := Request{}
	decoder.Decode(&request)

	email := request.Email
	password := request.Password

	User := models.User{Email: email, Password: password}
	User = User.Get()
	resp := Response{Success: false}
	if User.ID == 0 {
		resp = Response{Message: "user not found", Success: false}
	} else {
		resp = Response{Message: "login is successful", Success: true}
		Token := models.Token{UserID: User.ID}
		Token = Token.Get()
		log.Print("---------------")
		log.Print(Token.UserID)
		log.Print("---------------")
		if Token.ID == 0 {
			resp = Response{Message: "password doesn't match", Success: false}
			NewToken := Token.CreateNew()
			if NewToken.ID != 0 {
				resp = Response{Message: "login is successful", Success: true, Token: NewToken}
			}
		} else {
			resp = Response{Message: "login is successful", Success: true, Token: Token}
		}

	}

	ResponseJSON, _ := json.Marshal(resp)
	w.Write(ResponseJSON)

	return
}

// RegisterUser a view
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Message string      `json:"message"`
		Success bool        `json:"success"`
		NewUser models.User `json:"new_user"`
	}

	type Request struct {
		Email         string `json:"email"`
		Password      string `json:"password"`
		PasswordAgain string `json:"password_again"`
	}

	if r.Method != http.MethodPost {
		resp := Response{Message: "this endpoint use only post", Success: false}
		ResponseJSON, _ := json.Marshal(resp)
		w.Write(ResponseJSON)
		return
	}

	decoder := json.NewDecoder(r.Body)
	request := Request{}
	decoder.Decode(&request)

	email := request.Email
	password := request.Password
	passwordAgain := request.PasswordAgain

	if helpers.IsEmpty(email) || helpers.IsEmpty(password) || helpers.IsEmpty(passwordAgain) {
		resp := Response{Message: "all fields required.", Success: false}
		ResponseJSON, _ := json.Marshal(resp)
		w.Write(ResponseJSON)
		return
	}

	if password != passwordAgain {
		resp := Response{Message: "passwords are not same, please try again", Success: false}
		ResponseJSON, _ := json.Marshal(resp)
		w.Write(ResponseJSON)
		return
	}

	User := models.User{Email: email}
	User = User.Get()
	if User.ID != 0 {
		resp := Response{Message: "this email already taken", Success: false}
		ResponseJSON, _ := json.Marshal(resp)
		w.Write(ResponseJSON)
		return
	}

	User.Password = password
	User.Email = email
	User = User.CreateNew()

	resp := Response{Message: "registration has been completed", Success: true, NewUser: User}
	ResponseJSON, _ := json.Marshal(resp)
	w.Write(ResponseJSON)
	return
}

// UpdateUser for update user record.
// This method now only update password
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Message string      `json:"message"`
		Success bool        `json:"success"`
		User    models.User `json:"new_user"`
	}

	type Request struct {
		Password string `json:"password"`
	}

	if r.Method != http.MethodPost {
		resp := Response{Message: "this endpoint use only post", Success: false}
		ResponseJSON, _ := json.Marshal(resp)
		w.Write(ResponseJSON)
		return
	}

	decoder := json.NewDecoder(r.Body)
	request := Request{}
	decoder.Decode(&request)

	password := request.Password

	if len(password) < 3 {
		resp := Response{Message: "password too short", Success: false}
		ResponseJSON, _ := json.Marshal(resp)
		w.Write(ResponseJSON)
		return
	}

	sessionUser := context.Get(r, "User").(models.User)
	sessionUser.Password = password
	sessionUser.Update()
	resp := Response{Message: "profile has been updated.", Success: true}

	ResponseJSON, _ := json.Marshal(resp)
	w.Write(ResponseJSON)
	return
}
