package main

import (
	_ "encoding/base64"
	"fmt"
	"github.com/agpprastyo/todo-list-api/internal/database"
	"github.com/pascaldekloe/jwt"
	"net/http"
	"strconv"
	"time"
)

func (app *application) backgroundTask(r *http.Request, fn func() error) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		defer func() {
			err := recover()
			if err != nil {
				app.reportServerError(r, fmt.Errorf("%s", err))
			}
		}()

		err := fn()
		if err != nil {
			app.reportServerError(r, err)
		}
	}()
}

func (app *application) generateAccessToken(user *database.User) (string, time.Time, error) {
	var claims jwt.Claims
	claims.Subject = strconv.Itoa(user.ID)

	expiry := time.Now().Add(15 * time.Minute) // shorter lifetime for access tokens
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(expiry)

	claims.Issuer = app.config.baseURL
	claims.Audiences = []string{app.config.baseURL}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return string(jwtBytes), expiry, nil
}
