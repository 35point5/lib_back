package user

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"libback/db"
	"net/http"
	"strconv"
	"time"
)

type Service interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Auth(pri int) gin.HandlerFunc
	Profile(c *gin.Context)
	CookieLogin(c *gin.Context)
	Logout(c *gin.Context)
}

type service struct {
	db *gorm.DB
}

func MustNewService(database *gorm.DB) Service {
	return &service{database}
}
func UpdateCookie(usr *db.User) {
	h := md5.New()
	h.Write([]byte(usr.Name + strconv.FormatInt(time.Now().Unix(), 10) + "mogician"))
	usr.Cookie = hex.EncodeToString(h.Sum(nil))
}
func (t *service) UpdateInfo(usr *db.User) error {
	usr.Limit = db.CardLimit
	usr.Role = db.Guest
	UpdateCookie(usr)
	res := t.db.Save(&usr)
	return res.Error
}
func (t *service) Register(c *gin.Context) {
	var usr db.User
	err := c.ShouldBind(&usr)
	if err != nil || usr.Name == "" || usr.Password == "" {
		fmt.Println(1)
		c.Status(http.StatusBadRequest)
		return
	}
	res := t.db.Create(&usr)
	if res.Error != nil {
		fmt.Println(2)
		c.Status(http.StatusBadRequest)
		return
	}
	err = t.UpdateInfo(&usr)
	if err != nil {
		fmt.Println(3)
		c.Status(http.StatusBadRequest)
		return
	}
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("usr", usr.Cookie, 3600, "/", "mogician.cc", false, true)
	c.JSON(http.StatusOK, gin.H{"role": usr.Role, "username": usr.Name})
}
func (t *service) Login(c *gin.Context) {
	var usr, data db.User
	err := c.ShouldBind(&usr)
	if err != nil || usr.Name == "" || usr.Password == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	res := t.db.First(&data, "name = ?", usr.Name)
	if res.Error != nil || data.Password != usr.Password {
		c.Status(http.StatusBadRequest)
		return
	}
	UpdateCookie(&data)
	t.db.Save(&data)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("usr", data.Cookie, 3600, "/", "mogician.cc", false, true)
	c.JSON(http.StatusOK, gin.H{"role": data.Role, "username": data.Name})
}
func (t *service) Logout(c *gin.Context) {
	usrData, _ := c.Get("usr")
	usr := usrData.(db.User)
	UpdateCookie(&usr)
	t.db.Save(&usr)
	c.Status(http.StatusOK)
}
func (t *service) Auth(pri int) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("usr")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "请先登录"})
			return
		}
		var usr = db.User{Cookie: cookie}
		res := t.db.Where("cookie = ?", cookie).First(&usr)
		if res.Error != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "请先登录"})
			return
		}
		if usr.Role < pri {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "权限不足"})
			return
		}
		c.Set("usr", usr)
		c.Next()
	}
}

type ProfileItem struct {
	Title  string
	Number int
	Time   string
	ISBN   string
	CardID uint
}

func (t *service) Profile(c *gin.Context) {
	intf, _ := c.Get("usr")
	usr := intf.(db.User)
	var br []db.Borrow
	var crds []db.Card
	t.db.Where("user_id = ?", usr.ID).Find(&crds)

	var res []ProfileItem
	for _, crd := range crds {
		t.db.Where("card_id = ?", crd.ID).Find(&br)
		var qbook db.Book
		for _, v := range br {
			qbook.ISBN = v.ISBN
			t.db.First(&qbook)
			res = append(res, ProfileItem{
				Title:  qbook.Title,
				Number: v.Number,
				Time:   v.UpdatedAt.Format("2006-01-02 15:04:05"),
				ISBN:   qbook.ISBN,
				CardID: v.CardID,
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{"books": res})
}

func (t *service) CookieLogin(c *gin.Context) {
	usr, _ := c.Get("usr")
	c.JSON(http.StatusOK, gin.H{"role": usr.(db.User).Role, "username": usr.(db.User).Name})
}
