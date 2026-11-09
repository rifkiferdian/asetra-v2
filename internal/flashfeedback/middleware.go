package flashfeedback

import (
	"net/http"
	"net/url"
	"regexp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	successKey = "flash_feedback_success"
	errorKey   = "flash_feedback_error"
)

var validCode = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

// Middleware converts legacy feedback query parameters into one-time,
// signed-session flash values. The browser is redirected to a clean URL before
// the request reaches a page handler. On the next request the values are made
// available only to the handler and immediately removed from the session.
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead {
			query := c.Request.URL.Query()
			hasFeedbackParam := query.Has("success") || query.Has("error")
			success, hasSuccess := validFeedback(query, "success")
			failure, hasError := validFeedback(query, "error")
			if hasFeedbackParam {
				session := sessions.Default(c)
				if hasSuccess {
					session.Set(successKey, success)
				}
				if hasError {
					session.Set(errorKey, failure)
				}
				if (hasSuccess || hasError) && session.Save() != nil {
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				query.Del("success")
				query.Del("error")
				cleanURL := *c.Request.URL
				cleanURL.RawQuery = query.Encode()
				c.Redirect(http.StatusSeeOther, cleanURL.RequestURI())
				c.Abort()
				return
			}
		}

		restore(c)
		c.Next()
	}
}

func validFeedback(query url.Values, key string) (string, bool) {
	values, exists := query[key]
	if !exists || len(values) != 1 || !validCode.MatchString(values[0]) {
		return "", false
	}
	return values[0], true
}

func restore(c *gin.Context) {
	session := sessions.Default(c)
	success, _ := session.Get(successKey).(string)
	failure, _ := session.Get(errorKey).(string)
	if success == "" && failure == "" {
		return
	}

	session.Delete(successKey)
	session.Delete(errorKey)
	if err := session.Save(); err != nil {
		return
	}

	query := c.Request.URL.Query()
	if validCode.MatchString(success) {
		query.Set("success", success)
	}
	if validCode.MatchString(failure) {
		query.Set("error", failure)
	}
	c.Request.URL.RawQuery = query.Encode()
}
