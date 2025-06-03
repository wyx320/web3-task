package exam3

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 添加这行
	"github.com/jmoiron/sqlx"
)

/*
题目三：
使用SQL扩展库进行查询,
假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
要求 ：
编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。,
编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
*/

// 按部门名称查询
func QueryByDepartment(deptName string) []EmployeeEntity {
	var employees []EmployeeEntity
	sql := "select * from employees where department = ?"
	Db.Select(&employees, sql, deptName)
	return employees
}

// 查询工资最高的员工
func QueryMaxSalary() EmployeeEntity {
	var employee EmployeeEntity
	sql := "select * from employees where salary = (select max(salary) from employees)"
	Db.Get(&employee, sql)
	return employee
}

type EmployeeEntity struct {
	Id         uint    `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

func (EmployeeEntity) TableName() string {
	return "employees"
}

var Db *sqlx.DB

func init() {

	host := "localhost"
	port := 3306
	userName := "root"
	password := "1"
	dbName := "web3_task3"
	timeout := "10s"

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", userName, password, host, port, dbName, timeout)

	var err error
	Db, err = sqlx.Connect("mysql", connStr)
	if err != nil {
		fmt.Println("数据库连接失败")
	}

	// 连接成功
	fmt.Println(Db)
}
