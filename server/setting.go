package server

import (
	"polaris/log"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)


type setSettingIn struct {
	Key string `json:"key"`
	Value string `json:"value"`
}
func (s *Server) SetSetting(c *gin.Context) (interface{}, error) {
	var in setSettingIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	err := s.db.SetSetting(in.Key, in.Value)
	return nil, err
}

func (s *Server) GetSetting(c *gin.Context) (interface{}, error) {
	q := c.Query("key")
	log.Infof("query key: %v", q)
	if q == "" {
		return nil, nil
	}
	v := s.db.GetSetting(q)
	log.Infof("get value for key %v: %v", q, v)
	return gin.H{q: v}, nil
}
