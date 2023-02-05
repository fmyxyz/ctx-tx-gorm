package gormx

import (
	"context"
	"log"
	"testing"

	"github.com/fmyxyz/ctx-tx/test"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestWithTx(t *testing.T) {

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	//db = db.Debug()
	Register(db)

	//resetAll(db)
	test.Update88 = update88
	test.Update99 = update99
	test.Update = update
	for _, tt := range test.Tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.Fc(context.Background())
			eq88 := getValId(db, 88) == 88
			eq99 := getValId(db, 99) == 99
			resetAll(db)
			if eq88 != tt.Eq88 || eq99 != tt.Eq99 || (err != nil) != tt.WantErr {
				t.Errorf("WithTx() {eq88 = %v, want %v} {eq99 = %v, want %v} {error = %v, want has error %v}", eq88, tt.Eq88, eq99, tt.Eq99, err, tt.WantErr)
			}
			if err != nil {
				t.Log(err)
			}
		})
	}

	//resetAll(db)

}

func openDB() (db *gorm.DB, err error) {

	db, err = gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true&loc=Local"), &gorm.Config{
		AllowGlobalUpdate: false,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	if err != nil {
		return db, err
	}

	return db, nil
}

func getValId(db *gorm.DB, id int) (val int) {
	dest := &test.Item{
		ID: id,
	}
	err := db.Model(test.Item{}).Find(dest).Error
	if err != nil {
		return 0
	}
	return dest.Qty
}

func reset(db *gorm.DB, id int) error {
	dest := &test.Item{
		ID:  id,
		Qty: id + 1,
	}
	err := db.Model(test.Item{}).Where("id=?", id).Updates(dest).Error
	if err != nil {
		return err
	}
	return nil
}

func resetAll(db *gorm.DB) error {
	log.Println("resetAll")
	reset(db, 88)
	reset(db, 99)
	return nil
}

func update(ctx context.Context, id, num int) error {
	gormDB := FromContext(ctx)
	dest := &test.Item{
		ID:  id,
		Qty: num,
	}
	err := gormDB.Model(test.Item{}).Where("id=?", dest.ID).Updates(dest).Error
	if err != nil {
		return err
	}
	return nil
}

func update88(ctx context.Context) error {
	gormDB := FromContext(ctx)
	dest := &test.Item{
		ID:  88,
		Qty: 88,
	}
	err := gormDB.Model(test.Item{}).Where("id=?", dest.ID).Updates(dest).Error
	if err != nil {
		return err
	}
	return nil
}

func update99(ctx context.Context) error {
	gormDB := FromContext(ctx)
	dest := &test.Item{
		ID:  99,
		Qty: 99,
	}
	err := gormDB.Model(test.Item{}).Where("id=?", dest.ID).Updates(dest).Error
	if err != nil {
		return err
	}
	return nil
}
