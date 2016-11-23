package main

import (
	"github.com/jason0x43/go-toggl"
)

func getUserToken(user, pass string) (string, error) {
	session, err := toggl.NewSession(user, pass)
	if err != nil {
		return "", err
	}
	return session.APIToken, nil
}
