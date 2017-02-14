package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/database"
	"github.com/bonnyci/quartermaster/lib/keys"
	"github.com/bonnyci/quartermaster/web/engine"
)

type UserClaims struct {
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	jwt.StandardClaims
}

var (
	SignKey *keys.RSAKey
)

type AuthApiIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func init() {
	SignKey = keys.GetRSAKeySingleton()
}

func AuthTokenHandleFunc(w http.ResponseWriter, r *http.Request) {
	var in AuthApiIn
	if err := engine.Build(r.Body, &in); err != nil {
		jww.ERROR.Println(err)
		return
	}

	if in.Username == "" || in.Password == "" {
		efmt := "Invalid username or password."
		jww.ERROR.Println(efmt)
		http.Error(w, fmt.Errorf(efmt).Error(), http.StatusUnauthorized)
		return
	}

	if database.AuthUser(in.Username, in.Password, nil) {
		u, _ := database.GetUser(in.Username)
		g, _ := database.GetGroup("Admin")
		claims := UserClaims{
			in.Username,
			database.UserInGroup(g, u),
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 5).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		tokenString, err := token.SignedString(SignKey.Key())
		if err != nil {
			jww.ERROR.Println("Something went wrong with signing of the JWT token.")
			http.Error(w, fmt.Errorf("Authentication Failed!").Error(), http.StatusUnauthorized)
			return
		}

		jww.DEBUG.Printf("Authenticated User: %s", in.Username)
		jww.INFO.Printf("Token: %+v", tokenString)
		w.Write([]byte(tokenString))
	} else {
		http.Error(w, fmt.Errorf("Authentication Failed!").Error(), http.StatusUnauthorized)
	}
}

func ValidateToken(tokenStr string) (bool, *UserClaims) {

	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return &SignKey.Key().PublicKey, nil
	})

	if err != nil {
		return false, nil
	}

	claims := token.Claims.(*UserClaims)
	return token.Valid, claims
}

type AuthAPI struct {
	engine.API
}

func GetApi() *AuthAPI {
	return &AuthAPI{
		engine.APIBase{
			Handlers: engine.HandlersT{
				"/auth/": []engine.HandlersS{engine.MakeHandler("POST", AuthTokenHandleFunc)},
			},
		},
	}
}
