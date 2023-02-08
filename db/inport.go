package db

import (
	"encoding/csv"
	"fmt"
	"gorm.io/gorm"
	"io"
	"math/rand"
	"os"
	"time"
)

func Ipt(db *gorm.DB, fileName string) error {
	rand.Seed(time.Now().Unix())
	fs, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fs.Close()
	r := csv.NewReader(fs)
	data := Book{}
	//针对大文件，一行一行的读取文件
	cnt := 0
	for {
		cnt = cnt + 1
		row, err := r.Read()
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}
		if cnt == 1 {
			continue
		}
		rd := rand.Int63() //[0:12]
		data.Title = row[0]
		data.Author = row[1]
		data.Publisher = row[2]
		data.KeyWords = row[3]
		data.Digest = row[4]
		data.Category = row[5]
		data.PublishTime = row[6]
		data.Remain = 3
		data.ISBN = fmt.Sprintf("%013d", rd)[0:12]
		if data.Title != "" {
			res := db.Create(&data)
			if res.Error != nil {
				fmt.Println(data.ISBN, res.Error, res.RowsAffected, cnt)
			}
		}
	}
	return nil
}
