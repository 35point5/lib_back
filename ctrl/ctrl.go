package ctrl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"libback/bookManage"
	"libback/borrow"
	"libback/card"
	"libback/user"
)

type Ctrl interface {
	SetRouter() *gin.Engine
}
type ctrl struct {
	User       user.Service
	Borrow     borrow.Service
	Card       card.Service
	BookManage bookManage.Service
}

func MustNewCtrl(db *gorm.DB) Ctrl {
	return &ctrl{User: user.MustNewService(db), Borrow: borrow.MustNewService(db), Card: card.MustNewService(db), BookManage: bookManage.MustNewService(db)}
}

func (t *ctrl) SetRouter() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	//r := gin.New()
	r := gin.Default()
	//mime.AddExtensionType(".js", "application/javascript")
	//r.Static("/lib", "./static")
	r.POST("/lib/api/register", t.User.Register)
	r.POST("/lib/api/login", t.User.Login)
	AuthGroup := r.Group("/lib/api", t.User.Auth(1))
	{
		AuthGroup.POST("/profile", t.User.Profile)
		AuthGroup.POST("/lookup", t.Borrow.Lookup)
		AuthGroup.POST("/borrow", t.Borrow.Borrow)
		AuthGroup.POST("/return", t.Borrow.ReturnBook)
		AuthGroup.POST("/cookie", t.User.CookieLogin)
		AuthGroup.POST("/card/apply", t.Card.Apply)
		AuthGroup.POST("/card/refresh", t.Card.Refresh)
		AuthGroup.POST("/card/return", t.Card.Cancel)
		AuthGroup.POST("/logout", t.User.Logout)
	}
	AdminGroup := r.Group("/lib/api", t.User.Auth(2))
	{
		AdminGroup.POST("/upload", t.BookManage.Upload)
		AdminGroup.POST("/delete", t.BookManage.Delete)
		AdminGroup.POST("/modify", t.BookManage.Modify)
	}
	return r
}
