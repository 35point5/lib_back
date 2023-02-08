package card

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"libback/db"
	"net/http"
)

type Service interface {
	Apply(c *gin.Context)
	Refresh(c *gin.Context)
	Cancel(c *gin.Context)
}

type service struct {
	db *gorm.DB
}

func MustNewService(database *gorm.DB) Service {
	return &service{database}
}

func (t *service) Apply(c *gin.Context) {
	usrData, _ := c.Get("usr")
	usr := usrData.(db.User)
	if usr.Limit == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "读者证数目达到上限"})
		return
	}
	usr.Limit--
	t.db.Save(&usr)
	crd := db.Card{
		Limit:  db.BorrowLimit,
		UserID: usr.ID,
	}
	t.db.Create(&crd)
}

func (t *service) Refresh(c *gin.Context) {
	intf, _ := c.Get("usr")
	usr := intf.(db.User)
	var crds []db.Card
	fmt.Println(usr.ID)
	t.db.Where("user_id = ?", usr.ID).Find(&crds)
	c.JSON(http.StatusOK, gin.H{"cards": crds})
}

func (t *service) Cancel(c *gin.Context) {
	usrData, _ := c.Get("usr")
	usr := usrData.(db.User)
	var crd db.Card
	c.ShouldBind(&crd)
	br := db.Borrow{CardID: crd.ID}
	res := t.db.First(&br)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		usr.Limit++
		t.db.Unscoped().Delete(&crd)
		t.db.Save(&usr)
		c.Status(http.StatusOK)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请先归还卡下所借图书"})
	}
}
