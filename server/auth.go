package server

import (
	"net/http"
	"polaris/db"
	"polaris/log"
	"polaris/pkg/utils"
	"strings"
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

	auth := c.GetHeader("Authorization")
	if auth == "" {
		log.Infof("token is not present, abort")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	auth = strings.TrimPrefix(auth, "Bearer ")
	//log.Debugf("current token: %v", auth)
	token, err := jwt.ParseWithClaims(auth, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {	
		return []byte(s.jwtSerect), nil
	})
	if err != nil {
		log.Errorf("parse token error: %v", err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	if !token.Valid {
		log.Errorf("token is not valid: %v", auth)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	claim := token.Claims.(*jwt.RegisteredClaims)

	if time.Until(claim.ExpiresAt.Time) <= 0 {
		log.Infof("token is no longer valid: %s", auth)
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
	return gin.H{
		"token": sig,
	}, nil
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
