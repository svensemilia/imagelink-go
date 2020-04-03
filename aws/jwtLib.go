package aws

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat/go-jwx/jwk"
)

var (
	appClientId = "3qngaeqht61ie6jhutnv6918e7" // aud
	userPoolId  = "eu-central-1_s0wU2o3un"
	region      = "eu-central-1"
	iss         = "https://cognito-idp." + region + ".amazonaws.com/" + userPoolId
	cognitoJwk  = "/.well-known/jwks.json"
	knownKeys   *jwk.Set
)

func init() {
	fmt.Println("JWT init called")
	var err error
	knownKeys, err = jwk.Fetch(iss + cognitoJwk)
	if err != nil {
		fmt.Println(err)
	}
}

func ExtractSub(idToken string) string {
	token, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		key, err := getPublicKey(token.Header["kid"].(string))
		return key, err
	})
	//TODO handle error - authentication failure
	if err != nil {
		fmt.Println(err)
	}

	err = ValidateToken(token)
	if err != nil {
		fmt.Println(err)
	}
	cast := token.Claims.(jwt.MapClaims)
	return cast["sub"].(string)
}

func ValidateToken(token *jwt.Token) error {
	claims := token.Claims.(jwt.MapClaims)
	if claims["iss"] != iss {
		return fmt.Errorf("Token issuer do not match")
	}
	if claims["token_use"] != "id" {
		return fmt.Errorf("Not an ID token")
	}
	return nil
}

func getPublicKey(kid string) (*rsa.PublicKey, error) {
	keyRaw := knownKeys.LookupKeyID(kid)
	if len(keyRaw) == 0 {
		return nil, errors.New("could not find matching `kid` in well known tokens")
	}
	key, err := keyRaw[0].Materialize()
	if err != nil {
		return nil, err
	}
	rsaPublicKey := key.(*rsa.PublicKey)
	return rsaPublicKey, nil
}

func ExtractJWT() {
	tokenString := "eyJraWQiOiJHbWNqRnB2WFFvbjBTVEdPcEdwRXdYTTBMXC9Nc0tlYmFTQ3REZ3ZZN2hwST0iLCJhbGciOiJSUzI1NiJ9.eyJhdF9oYXNoIjoibENOcUxWTXM0eGE4WHNTWG5QcTlHUSIsInN1YiI6IjVlNmVjYzI3LWI5ZmUtNGI2MC1hMWI4LWI3M2JhMzVlOGYxNSIsImF1ZCI6IjNxbmdhZXFodDYxaWU2amh1dG52NjkxOGU3IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImV2ZW50X2lkIjoiYjAyYjk4NjktYjkxMC00MGJmLWI1NTMtZDllZmIxMzc4MDk0IiwidG9rZW5fdXNlIjoiaWQiLCJhdXRoX3RpbWUiOjE1ODU5MjAxMTIsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC5ldS1jZW50cmFsLTEuYW1hem9uYXdzLmNvbVwvZXUtY2VudHJhbC0xX3Mwd1UybzN1biIsImNvZ25pdG86dXNlcm5hbWUiOiI1ZTZlY2MyNy1iOWZlLTRiNjAtYTFiOC1iNzNiYTM1ZThmMTUiLCJleHAiOjE1ODU5MjM3MTIsImlhdCI6MTU4NTkyMDExMiwiZW1haWwiOiJzdmVuX2NhcmxpbkBhcmNvci5kZSJ9.ghJcXPKyI45fF8biIMaqeBOaqcbJQx8bsRYTwmRPh0lbMDFoCdWxjarj2T526k0dHYeTrf5LyuOwYhvlECGdzKDcxZI1wCi3kcUvQ5cP8mUOqttivmEebpgSGlbPtYa_mr_W0V_nC5OWNxNsQmleGtD33ygi3F3Iks-3lJaBkfX6jPkzmBvlhtHK9iOK0O6znOdOS0hkHKXzpdKHslQItWD2jrCQA8Cgxz_GN_24FN81VMknDtecQ9SNgmqJsj3_REqyozCV5AUJdmKDVwuU3-Jxcr3whDQu3dWHTF91LZFpjekpvwOKJIG6nEmJfEHgs2Bq96QTvER1xixpNza2Ag"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		keyRaw := knownKeys.LookupKeyID(token.Header["kid"].(string))
		if len(keyRaw) == 0 {
			return nil, errors.New("could not find matching `kid` in well known tokens")
		}
		key, err := keyRaw[0].Materialize()
		if err != nil {
			return nil, err
		}
		rsaPublicKey := key.(*rsa.PublicKey)
		return rsaPublicKey, nil
	})

	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println("Token:", token)
	cast := token.Claims.(jwt.MapClaims)
	//fmt.Println("Claims:", cast)
	fmt.Println(cast["sub"])
	fmt.Println(cast["email"])
	fmt.Println(token.Valid)
}

func getCode(r *http.Request, index int) string {
	var result string
	trimmedPath := strings.TrimRight(r.URL.Path, "/")
	trimmedPath = strings.TrimLeft(trimmedPath, "/")
	p := strings.Split(trimmedPath, "/")
	fmt.Println(len(p))
	if len(p) > index {
		result = p[index]
	}
	return result
}

func GetObjectKey(r *http.Request, basepath string) string {
	var result string
	trimmedPath := strings.TrimRight(r.URL.Path, "/")
	trimmedPath = strings.TrimLeft(trimmedPath, "/")
	result = strings.Replace(trimmedPath, basepath, "", 1)
	return result
}
