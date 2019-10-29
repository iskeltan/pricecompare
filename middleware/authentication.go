package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"../models"
	"github.com/gorilla/context"
)

func Authenticated(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Response struct {
			Message string
			Success bool
		}
		resp := Response{Success: false}
		authentication := r.Header.Get("Authorization")

		if authentication != "" {
			tokenList := strings.Split(authentication, "Token ")
			if len(tokenList) != 2 {
				resp.Message = "authenticate fail"
				ResponseJSON, _ := json.Marshal(resp)
				w.Write(ResponseJSON)
				return
			}
			token := tokenList[1]

			GetToken := models.Token{Token: token}
			GetToken = GetToken.Get()

			if GetToken.ID == 0 {
				resp.Message = "authenticate fail"
				ResponseJSON, _ := json.Marshal(resp)
				w.Write(ResponseJSON)
				return
			}

			GetUser := models.User{}
			GetUser.ID = GetToken.UserID
			GetUser = GetUser.Get()
			context.Set(r, "User", GetUser)
			f(w, r)
		} else {
			resp.Message = "authenticate fail"
			ResponseJSON, _ := json.Marshal(resp)
			w.Write(ResponseJSON)
			return
		}
	}
}
