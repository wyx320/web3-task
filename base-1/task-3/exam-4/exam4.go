package exam4

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 添加这行
	"github.com/jmoiron/sqlx"
)

/*
题目四：
实现类型安全映射,
假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
要求 ：
定义一个 Book 结构体，包含与 books 表对应的字段。,
编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
*/

type BookEntity struct {
	Id     uint    `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

func (BookEntity) TableName() string {
	return "books"
}

var Db *sqlx.DB

func init() {
	host := "localhost"
	port := 3306
	user := "root"
	password := "1"
	database := "web3_task3"
	timeout := "10s"

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", user, password, host, port, database, timeout)

	var err error
	Db, err = sqlx.Open("mysql", connStr)
	if err != nil {
		fmt.Println("数据库连接失败：", err)
		return
	}

	// 连接成功
	fmt.Println(Db)
}

func QueryByPrice(price float64) []BookEntity {
	var books []BookEntity
	sql := "select * from books where price > ?"
	Db.Select(&books, sql, price)
	return books
}
