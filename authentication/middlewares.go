package authentication

import (
	"fmt"

	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//MyCustomClaims ...
type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// RequireTokenAuthentication ...
func RequireTokenAuthentication(appConfig map[string]interface{}) gin.HandlerFunc {
	authBackend := InitJWTAuthenticationBackend(appConfig)
	/*token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return authBackend.PublicKey, nil
		}
	})
	fmt.Println(err)*/

	return func(c *gin.Context) {
		authToken := strings.Replace(c.Request.Header.Get("Authorization"), "Bearer ", "", -1)
		_, err := jwt.ParseWithClaims(authToken, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			} else {
				return authBackend.PublicKey, nil
			}
		})

		if err != nil {
			fmt.Println(c.Request.Header.Get("Authorization"))
			c.AbortWithError(401, err)
		}
	}

	//if err == nil && token.Valid && !authBackend.IsInBlacklist(req.Header.Get("Authorization")) {
	//	next(rw, req)
	//} else {
	//	rw.WriteHeader(http.StatusUnauthorized)
	//}
}
