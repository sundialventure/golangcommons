package authentication

import (
	"bufio"

	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	//"github.com/grafana/grafana/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"sundialventure.com/common/redis"
	settings "sundialventure.com/common/settings"
)

//JWTAuthenticationBackend ...
type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

const (
	tokenDuration = 72
	expireOffset  = 3600
)

var (
	setting_dir = ""
)

var authBackendInstance *JWTAuthenticationBackend = nil

// InitJWTAuthenticationBackend ...
func InitJWTAuthenticationBackend(appConfig map[string]interface{}) *JWTAuthenticationBackend {

	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(appConfig),
			PublicKey:  getPublicKey(appConfig),
		}
	}

	return authBackendInstance
}

// GenerateToken ...
func (backend *JWTAuthenticationBackend) GenerateToken(userUUID string, userName string) (string, error) {
	expireToken := time.Now().Add(time.Hour * 24).Unix()
	//expireCookie := time.Now().Add(time.Hour * 24)

	claims := MyCustomClaims{
		userName,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			Subject:   userUUID,
			ExpiresAt: expireToken,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	/*token.Claims = jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix(),
		"iat": time.Now().Unix(),
		"sub": userUUID,
	}*/

	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		panic(err)
		return "", err
	}
	return tokenString, nil
}

//Authenticate ...
func (backend *JWTAuthenticationBackend) Authenticate(user_password, password string, auth_type string) (bool, error) {
	//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	/*testUser := util.User{
		UUID:     uuid.New(),
		Username: user.Username,
		Password: string(hashedPassword),
	}*/
	//user.UUID = uuid.New()
	//fmt.Println(bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(hashedPassword)))
	var err error
	if auth_type == "local" {
		err = bcrypt.CompareHashAndPassword([]byte(user_password), []byte(password))
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}
	return expireOffset
}

func (backend *JWTAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
	redisConn := redis.Connect()
	claims := token.Claims.(jwt.MapClaims)
	return redisConn.SetValue(tokenString, tokenString, backend.getTokenRemainingValidity(claims["exp"]))
}

func (backend *JWTAuthenticationBackend) IsInBlacklist(token string) bool {
	redisConn := redis.Connect()
	redisToken, _ := redisConn.GetValue(token)
	fmt.Println(redisConn)
	if redisToken == nil {
		return false
	}

	return true
}

func getPrivateKey(appConfig map[string]interface{}) *rsa.PrivateKey {
	fmt.Println(settings.Get().PrivateKeyPath)
	privateKeyFile, err := os.Open(appConfig["appdir"].(string) + "security/settings/keys/private_key") ///Users/amosunsunday/Documents/Official/xtremepay/Sources/services/src/xtremepay.com/backoffice/security/settings/
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

func getPublicKey(appConfig map[string]interface{}) *rsa.PublicKey {
	//settings.Get().PublicKeyPath
	///Users/amosunsunday/Documents/Official/xtremepay/Sources/services/src/xtremepay.com/backoffice/security/settings/keys/public_key.pub
	publicKeyFile, err := os.Open(appConfig["appdir"].(string) + "security/settings/keys/public_key.pub")
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	var size = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	publicKeyFile.Close()

	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		panic(err)
	}

	return rsaPub
}
