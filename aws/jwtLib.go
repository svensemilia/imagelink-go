package aws

import (
	"crypto/rsa"
	"errors"
	"fmt"

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
	var err error
	knownKeys, err = jwk.Fetch(iss + cognitoJwk)
	if err != nil {
		fmt.Println(err)
	}
}

func AuthenticateUser(idToken string) (string, error) {
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
		return "", err
	}

	err = ValidateToken(token)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	cast := token.Claims.(jwt.MapClaims)
	return cast["sub"].(string), nil
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
