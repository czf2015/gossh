package middlewares

import (
	config "gossh/config/v1"
	"gossh/libs/sessions"
	"gossh/libs/sessions/cookie"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Session() gin.HandlerFunc {
	var err error
	var sessionName = "session_id"
	var sessionPath = "/"
	var sessionHttpOnly = true
	var sessionSecure = false
	var sessionDomain = ""
	var sessionMaxAge = 3600 * 24
	var sessionSameSite = http.SameSiteLaxMode

	if len(config.Config["session"]["Name"]) > 0 {
		sessionName = config.Config["session"]["Name"]
	}

	if len(config.Config["session"]["Path"]) > 0 {
		sessionPath = config.Config["session"]["Path"]
	}

	if len(config.Config["session"]["Domain"]) > 0 {
		sessionDomain = config.Config["session"]["Domain"]
	}

	sessionHttpOnly, err = strconv.ParseBool(config.Config["session"]["HttpOnly"])
	if err != nil {
		sessionHttpOnly = true
	}

	sessionSecure, err = strconv.ParseBool(config.Config["session"]["Secure"])
	if err != nil {
		sessionSecure = false
	}

	sessionMaxAge, err = strconv.Atoi(config.Config["session"]["MaxAge"])
	if err != nil {
		sessionMaxAge = 3600
	}

	sessionSameSiteVal, err := strconv.Atoi(config.Config["session"]["SameSite"])
	if err != nil {
		sessionSameSite = http.SameSiteLaxMode
	} else {
		sessionSameSite = http.SameSite(sessionSameSiteVal)
	}

	// 加密cookie方式
	store := cookie.NewStore([]byte(config.Config["session"]["Secret"]))

	store.Options(sessions.Options{
		Path:     sessionPath,
		Domain:   sessionDomain,
		MaxAge:   sessionMaxAge,
		Secure:   sessionSecure,
		HttpOnly: sessionHttpOnly,
		SameSite: sessionSameSite,
	})
	return sessions.Sessions(sessionName, store)
}
