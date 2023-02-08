package borrow

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"libback/db"
	"net/http"
	"reflect"
)

type Service interface {
	Lookup(c *gin.Context)
	Borrow(c *gin.Context)
	ReturnBook(c *gin.Context)
}
type service struct {
	db   *gorm.DB
	info []string
}

func MustNewService(database *gorm.DB) Service {
	var temp []string
	var Pricerecord []db.Tabel_infor
	database.Raw("desc books").Scan(&Pricerecord)
	for _, v := range Pricerecord {
		temp = append(temp, v.Field)
	}
	fmt.Println(temp)
	return &service{database, temp}
}

type queryModel struct {
	BookInfo db.Book
	Keywords string
}

func (t *service) Lookup(c *gin.Context) {
	fmt.Println(c.Get("usr"))
	var query queryModel
	var qbook db.Book
	var data []db.Book
	err := c.ShouldBind(&query)
	qbook = query.BookInfo
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	tp := reflect.TypeOf(qbook)
	val := reflect.ValueOf(qbook)
	questr := ""
	cnt := 0
	for i := 0; i < tp.NumField(); i++ {
		if reflect.Kind(reflect.String) == tp.Field(i).Type.Kind() && val.Field(i).String() != "" {
			if cnt > 0 {
				questr = questr + " and "
			}
			cnt++
			questr = questr + t.db.NamingStrategy.ColumnName("books", tp.Field(i).Name) + " like '%" + val.Field(i).String() + "%'"
			//t.db.Where(t.info[i]+" like ?", "%"+val.Field(i).String()+"%").Find(&que)
			//data = append(data, que...)
		}
	}
	if query.Keywords != "" {
		if cnt > 0 {
			questr = questr + " and "
		}
		cnt++
		questr += "("
		cnt = 0
		for i, _ := range t.info {
			if reflect.Kind(reflect.String) == tp.Field(i).Type.Kind() {
				if cnt > 0 {
					questr = questr + " or "
				}
				cnt++
				questr += t.db.NamingStrategy.ColumnName("books", tp.Field(i).Name) + " like '%" + query.Keywords + "%'"
			}
		}
		questr += ")"
	}
	fmt.Println(questr)
	if questr != "" {
		t.db.Where(questr).Limit(100).Find(&data)
	} else {
		c.Status(http.StatusBadRequest)
		return
	}
	//t.db.Where(&qbook).Find(&data)
	//fmt.Println(data)
	usr, _ := c.Get("usr")
	c.JSON(http.StatusOK, gin.H{"Num": len(data), "Data": data, "role": usr.(db.User).Role})
}

func (t *service) Borrow(c *gin.Context) {
	var br db.Borrow

	err := c.ShouldBind(&br)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	var qbook db.Book
	res := t.db.Model(&qbook).Where(&br, "ISBN").First(&qbook)

	var crd db.Card
	res = t.db.Where("ID = ?", br.CardID).First(&crd)

	if res.Error != nil || qbook.Remain < br.Number || crd.Limit < br.Number {
		c.Status(http.StatusBadRequest)
		return
	}

	old := br
	qbook.Remain -= br.Number
	crd.Limit -= br.Number
	t.db.Unscoped().Where(&br, "ISBN", "card_id").Attrs("Number", 0).FirstOrCreate(&old)
	old.DeletedAt = gorm.DeletedAt{}
	//t.db.Unscoped().Where(&br, "ISBN", "CardID").First(&old)
	old.Number += br.Number
	fmt.Println(br)
	fmt.Println(qbook)
	err = t.db.Transaction(func(tx *gorm.DB) error {
		tx.Save(&crd)
		tx.Save(&qbook)
		tx.Unscoped().Save(&old)
		return tx.Error
	})
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
}

type ReturnModel struct {
	ISBN   string
	Number int
}

func (t *service) ReturnBook(c *gin.Context) {
	var rt db.Borrow
	err := c.ShouldBind(&rt)
	br := db.Borrow{
		CardID: rt.CardID,
		ISBN:   rt.ISBN,
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "还书失败"})
		return
	}

	t.db.First(&br)
	if br.Number < rt.Number {
		c.JSON(http.StatusBadRequest, gin.H{"message": "还书失败"})
		return
	}

	var crd db.Card
	t.db.Where("ID = ?", rt.CardID).First(&crd)

	qbook := db.Book{ISBN: rt.ISBN}
	t.db.First(&qbook)

	crd.Limit += rt.Number
	qbook.Remain += rt.Number
	br.Number -= rt.Number
	fmt.Println(br.Number)
	_ = t.db.Transaction(func(tx *gorm.DB) error {
		tx.Save(&crd)
		tx.Save(&qbook)
		tx.Save(&br)
		return nil
	})
}
