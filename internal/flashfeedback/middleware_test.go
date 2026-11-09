package flashfeedback

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func TestMiddlewareShowsFeedbackOnceOnCleanURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(sessions.Sessions("test_session", cookie.NewStore([]byte("01234567890123456789012345678901"))))
	router.Use(Middleware())
	router.GET("/page", func(c *gin.Context) { c.String(http.StatusOK, c.Query("success")) })

	first := httptest.NewRecorder()
	router.ServeHTTP(first, httptest.NewRequest(http.MethodGet, "/page?success=deleted", nil))
	if first.Code != http.StatusSeeOther || first.Header().Get("Location") != "/page" {
		t.Fatalf("legacy URL response = %d %q, want 303 /page", first.Code, first.Header().Get("Location"))
	}
	cookies := first.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("flash redirect did not set a session cookie")
	}

	secondRequest := httptest.NewRequest(http.MethodGet, "/page", nil)
	secondRequest.AddCookie(cookies[0])
	second := httptest.NewRecorder()
	router.ServeHTTP(second, secondRequest)
	if second.Code != http.StatusOK || second.Body.String() != "deleted" {
		t.Fatalf("flash response = %d %q, want 200 deleted", second.Code, second.Body.String())
	}

	clearedCookies := second.Result().Cookies()
	thirdRequest := httptest.NewRequest(http.MethodGet, "/page", nil)
	if len(clearedCookies) > 0 {
		thirdRequest.AddCookie(clearedCookies[0])
	}
	third := httptest.NewRecorder()
	router.ServeHTTP(third, thirdRequest)
	if third.Code != http.StatusOK || third.Body.String() != "" {
		t.Fatalf("reloaded response = %d %q, want empty feedback", third.Code, third.Body.String())
	}
}

func TestMiddlewareRejectsInvalidFeedbackCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(sessions.Sessions("test_session", cookie.NewStore([]byte("01234567890123456789012345678901"))))
	router.Use(Middleware())
	router.GET("/page", func(c *gin.Context) { c.String(http.StatusOK, c.Query("success")) })

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/page?success=%3Cscript%3E", nil))
	if response.Code != http.StatusSeeOther || response.Header().Get("Location") != "/page" {
		t.Fatalf("invalid code response = %d %q, want clean redirect", response.Code, response.Header().Get("Location"))
	}
}
