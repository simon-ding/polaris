package server

import (
	"net/http"
	"polaris/db"
	"polaris/log"
	"polaris/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

func (s *Server) isAuthEnabled() bool {
	authEnabled := s.db.GetSetting(db.SettingAuthEnabled)
	return authEnabled == "true"
}

func (s *Server) authModdleware(c *gin.Context) {
	if !s.isAuthEnabled() {
		c.Next()
		return
	}
	token, err := c.Cookie("polaris_token")
	if err != nil {
		log.Errorf("token error: %v", err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	//log.Debugf("current token: %v", auth)
	tokenParsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSerect), nil
	})
	if err != nil {
		log.Errorf("parse token error: %v", err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	if !tokenParsed.Valid {
		log.Errorf("token is not valid: %v", token)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	claim := tokenParsed.Claims.(*jwt.RegisteredClaims)

	if time.Until(claim.ExpiresAt.Time) <= 0 {
		log.Infof("token is no longer valid: %s", token)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Next()

}

type LoginIn struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func (s *Server) Login(c *gin.Context) (interface{}, error) {
	var in LoginIn

	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}

	if !s.isAuthEnabled() {
		return nil, nil
	}

	user := s.db.GetSetting(db.SettingUsername)
	if user != in.User {
		return nil, errors.New("login fail")
	}
	password := s.db.GetSetting(db.SettingPassword)
	if !utils.VerifyPassword(in.Password, password) {
		return nil, errors.New("login fail")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "system",
		Subject:   in.User,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	})
	sig, err := token.SignedString([]byte(s.jwtSerect))
	if err != nil {
		return nil, errors.Wrap(err, "sign")
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("polaris_token", sig, 0, "/", "", false, false)
	return "success", nil
}

func (s *Server) Logout(c *gin.Context) (interface{}, error) {
	if !s.isAuthEnabled() {
		return nil, errors.New( "auth is not enabled")
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("polaris_token", "", -1, "/", "", false, false)
	return nil, nil
}

type EnableAuthIn struct {
	Enable   bool   `json:"enable"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (s *Server) EnableAuth(c *gin.Context) (interface{}, error) {
	var in EnableAuthIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}

	if in.Enable && (in.User == "" || in.Password == "") {
		return nil, errors.New("user password should not empty")
	}
	if !in.Enable {
		log.Infof("disable auth")
		s.db.SetSetting(db.SettingAuthEnabled, "false")
	} else {
		log.Info("enable auth")
		s.db.SetSetting(db.SettingAuthEnabled, "true")
		s.db.SetSetting(db.SettingUsername, in.User)

		hash, err := utils.HashPassword(in.Password)
		if err != nil {
			return nil, errors.Wrap(err, "hash password")
		}
		s.db.SetSetting(db.SettingPassword, hash)
	}
	return "success", nil
}

func (s *Server) GetAuthSetting(c *gin.Context) (interface{}, error) {
	enabled := s.db.GetSetting(db.SettingAuthEnabled)
	user := s.db.GetSetting(db.SettingUsername)

	return EnableAuthIn{
		Enable: enabled == "true",
		User:   user,
	}, nil
}
