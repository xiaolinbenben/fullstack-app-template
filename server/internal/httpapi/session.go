package httpapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	adminUsername = "admin"
	adminPassword = "admin123"
	sessionCookie = "admin_session"
	sessionTTL    = 7 * 24 * time.Hour
)

func writeSession(c *gin.Context, key string) {
	payload := fmt.Sprintf("%s|%d", adminUsername, time.Now().Add(sessionTTL).Unix())
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(sessionCookie, signSession(key, payload), int(sessionTTL.Seconds()), "/", "", false, true)
}

func clearSession(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(sessionCookie, "", -1, "/", "", false, true)
}

func currentAdmin(c *gin.Context, key string) (string, bool) {
	raw, err := c.Cookie(sessionCookie)
	if err != nil {
		return "", false
	}
	payload, ok := openSession(key, raw)
	if !ok {
		return "", false
	}
	name, expiryText, found := strings.Cut(payload, "|")
	if !found || name != adminUsername {
		return "", false
	}
	expiry, err := strconv.ParseInt(expiryText, 10, 64)
	if err != nil || time.Now().Unix() > expiry {
		return "", false
	}
	return name, true
}

func signSession(key, payload string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + base64.RawURLEncoding.EncodeToString(sessionMAC(key, payload))
}

func openSession(key, token string) (string, bool) {
	encoded, signature, found := strings.Cut(token, ".")
	if !found {
		return "", false
	}
	payload, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", false
	}
	got, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil || !hmac.Equal(got, sessionMAC(key, string(payload))) {
		return "", false
	}
	return string(payload), true
}

func sessionMAC(key, payload string) []byte {
	mac := hmac.New(sha256.New, []byte(key))
	_, _ = mac.Write([]byte(payload))
	return mac.Sum(nil)
}
