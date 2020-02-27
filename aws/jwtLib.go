package aws

import (
    "fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

// func main(){
//	extractJWT()
// }

func ExtractSub(idToken string) string {
	token, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		fmt.Println("Parsing jwt token...")
		fmt.Println(token.Header["alg"])
		// to verify the token, you should provide here a key
		return nil, nil
	})

	if(err != nil){
		fmt.Println(err)
	}
	cast := token.Claims.(jwt.MapClaims)
	return cast["sub"].(string)
}

func extractJWT(){
	tokenString := "eyJraWQiOiJHbWNqRnB2WFFvbjBTVEdPcEdwRXdYTTBMXC9Nc0tlYmFTQ3REZ3ZZN2hwST0iLCJhbGciOiJSUzI1NiJ9.eyJhdF9oYXNoIjoiR014aW5pT3FGMF9rUnNSRGMzZUZGUSIsInN1YiI6IjVlNmVjYzI3LWI5ZmUtNGI2MC1hMWI4LWI3M2JhMzVlOGYxNSIsImF1ZCI6IjNxbmdhZXFodDYxaWU2amh1dG52NjkxOGU3IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImV2ZW50X2lkIjoiZDYyNmVmZjMtODgzOS0xMWU5LThlYjAtNmQxMGI5NTQ4YTY1IiwidG9rZW5fdXNlIjoiaWQiLCJhdXRoX3RpbWUiOjE1NTk4MTE3NjAsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC5ldS1jZW50cmFsLTEuYW1hem9uYXdzLmNvbVwvZXUtY2VudHJhbC0xX3Mwd1UybzN1biIsImNvZ25pdG86dXNlcm5hbWUiOiI1ZTZlY2MyNy1iOWZlLTRiNjAtYTFiOC1iNzNiYTM1ZThmMTUiLCJleHAiOjE1NTk4MTUzNjAsImlhdCI6MTU1OTgxMTc2MCwiZW1haWwiOiJzdmVuX2NhcmxpbkBhcmNvci5kZSJ9.ZvQPwi7QF4z9LTK8nkj_aRP9LarRcEFU7PrlG7TqYhdjML8Mc3FpqVG1YucBhqWj9lrqi0FGH1ysv4H3m0nAgFw6b3R1VNE623WdWwxTviqIogsZr-AZTA4cs18qAyp1ND1-bRSVxI_Im_Abm3Y07xCI0e_B9gqjngG9bis31rzaO5ICDkJjestlnmrrWpKAKgiQRjZU565W4VphXlwA763svv3hMl8yXw_5QMc36O5C8K0hVxGLKqKvGqm9aBmHrjJnMN4NqxcfGkU6pHxbj7XcA6eVrx1kEsN0OVXTF9O8ur0D9bPoESuhlt-HjgtMTDNRYC4g61RxOyUecNWWtA"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		fmt.Println("Parsing jwt token...")
		return nil, nil
	})

	if(err != nil){
		fmt.Println(err)
	}
	cast := token.Claims.(jwt.MapClaims)
	
	fmt.Println(cast["sub"])
	fmt.Println(cast["email"])
}

func getCode(r *http.Request, index int) (string) {
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

func GetObjectKey(r *http.Request, basepath string) (string) {
		var result string
		trimmedPath := strings.TrimRight(r.URL.Path, "/")
		trimmedPath = strings.TrimLeft(trimmedPath, "/")
		result = strings.Replace(trimmedPath, basepath, "", 1)
		return result
}