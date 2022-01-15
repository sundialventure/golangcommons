package authentication

import (
	"fmt"

	jwt "github.com/golang-jwt/jwt"
	"github.com/gin-gonic/gin"
	"strings"
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
		fmt.Println(authToken)
		_, err := jwt.ParseWithClaims(authToken, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			} else {
				return authBackend.PublicKey, nil
			}
		})

		if err != nil {
			c.AbortWithError(401, err)
		}
	}

	//if err == nil && token.Valid && !authBackend.IsInBlacklist(req.Header.Get("Authorization")) {
	//	next(rw, req)
	//} else {
	//	rw.WriteHeader(http.StatusUnauthorized)
	//}
}
