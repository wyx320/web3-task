package exam1

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
假设有一个名为 students 的表，
包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
要求 ：
编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。,
编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。,
编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。,
编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
*/

type StudentEntity struct {
	Id    uint   `gorm:"primaryKey;autoIncrement;column:id;comment:'主键'"`
	Name  string `gorm:"column:name;comment:'姓名'"`
	Age   int    `gorm:"column:age;comment:'年龄'"`
	Grade string `gorm:"column:grade;comment:'年级'"`
}

// 指定表名
func (StudentEntity) TableName() string {
	return "students"
}

var Db *gorm.DB

func init() {
	host := "localhost"
	port := 3306
	userName := "root"
	password := "1"
	DbName := "web3_task3"
	timeout := "10s"

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", userName, password, host, port, DbName, timeout)

	var err error
	Db, err = gorm.Open(mysql.Open(connStr))
	if err != nil {
		fmt.Println("数据库链接失败：", err)
		return
	}
	// 连接成功
	fmt.Println(Db)

	Db.AutoMigrate(&StudentEntity{})
}

func CreateNewOne() {
	student := StudentEntity{
		Name:  "张三",
		Age:   20,
		Grade: "三年级",
	}
	Db.Create(&student)
}

func QueryByAge() []StudentEntity {
	var students []StudentEntity
	Db.Find(&students, "age > ? ", 18)
	return students
}

func UpdateGradeByName(name string, grade string) {
	Db.Model(&StudentEntity{}).Where("name = ?", name).Update("grade", grade)
}

func DeleteLessThanByAge() {
	Db.Where("age < ?", 15).Delete(&StudentEntity{})
}
