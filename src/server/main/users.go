package main

var tokenMap = make(map[string]interface{}) //token to user

func checkToken(token string) bool {
	_, ok := tokenMap["foo"]
	return ok
}

func getToken(user string) string {

}