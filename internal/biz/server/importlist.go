package server

import (
	"fmt"
	"polaris/ent"
	"polaris/ent/importlist"
	"polaris/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) getAllImportLists(c *gin.Context) (interface{}, error) {
	lists, err := s.db.GetAllImportLists()
	return lists, err
}

type addImportlistIn struct {
	Name      string `json:"name" binding:"required"`
	Url       string `json:"url"`
	Type      string `json:"type"`
	Qulity    string `json:"qulity"`
	StorageId int    `json:"storage_id"`
}

func (s *Server) addImportlist(c *gin.Context) (interface{}, error) {
	var in addImportlistIn

	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "json")
	}
	utils.TrimFields(&in)
	
	st := s.db.GetStorage(in.StorageId)
	if st == nil {
		return nil, fmt.Errorf("storage id not exist: %v", in.StorageId)
	}
	err := s.db.AddImportlist(&ent.ImportList{
		Name:      in.Name,
		URL:       in.Url,
		Type:      importlist.Type(in.Type),
		Qulity:    in.Qulity,
		StorageID: in.StorageId,
	})
	if err != nil {
		return nil, err
	}
	return "success", nil
}

type deleteImportlistIn struct {
	ID int `json:"id"`
}

func (s *Server) deleteImportList(c *gin.Context) (interface{}, error) {
	var in deleteImportlistIn

	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "json")
	}
	s.db.DeleteImportlist(in.ID)
	return "sucess", nil
}
