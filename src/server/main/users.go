package main

import "github.com/nu7hatch/gouuid"

var tokenMap = make(map[string]string) //token to user

func checkToken(token string) (user string, ok bool) {
	user, ok = tokenMap[token]
	return
}

func getNewToken(user string) (strToken string) { //TODO token expiry date
	token, err := uuid.NewV4()
	if err != nil {
		logFatalStr("Can't obtain a new token for " + user + " " + err.Error())
	}
	strToken = token.String()
	if _, ok := checkToken(strToken); ok { //get new token if it's already taken
		strToken = getNewToken(user) //hopefully doesn't stack overflow
	}
	tokenMap[strToken] = user
	return
}
