package bookManage

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"libback/db"
	"net/http"
	"path/filepath"
)

type Service interface {
	Upload(c *gin.Context)
	Delete(c *gin.Context)
	Modify(c *gin.Context)
}

type service struct {
	db *gorm.DB
}

func MustNewService(database *gorm.DB) Service {
	return &service{database}
}

func (t *service) Upload(c *gin.Context) {
	f, _ := c.FormFile("file")
	basePath := "./"
	filename := basePath + filepath.Base(f.Filename)
	c.SaveUploadedFile(f, filename)
	if err := db.Ipt(t.db, f.Filename); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}

func (t *service) Delete(c *gin.Context) {
	var bk db.Book
	c.ShouldBind(&bk)
	t.db.Unscoped().Delete(&bk)
	c.Status(http.StatusOK)
}

func (t *service) Modify(c *gin.Context) {
	var bk db.Book
	c.ShouldBind(&bk)
	t.db.Save(&bk)
	c.Status(http.StatusOK)
}
