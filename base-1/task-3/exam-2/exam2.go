package exam2

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
题目2：事务语句,
假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
要求 ：
编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。
在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
*/

// 创建账户
func CreateAccount(balance float64) *AccountEntity {
	account := AccountEntity{
		Balance: balance,
	}
	Db.Create(&account)

	return &account
}

// 转账
func Transfer(fromAccountId uint, toAccountId uint, amount float64) {
	Db.Transaction(func(tx *gorm.DB) error {
		// 查询账户余额
		var fromAccount AccountEntity
		tx.Take(&fromAccount, fromAccountId)
		if fromAccount.Balance < amount {
			return fmt.Errorf("账户%v余额不足", fromAccountId)
		}

		// 扣除金额
		fromAccount.Balance -= amount
		Db.Save(&fromAccount)

		// 增加金额
		var toAccount AccountEntity
		tx.Take(&toAccount, toAccountId)
		if toAccount.Id == 0 {
			return fmt.Errorf("账户%v不存在", toAccountId)
		}
		toAccount.Balance += amount
		Db.Save(&toAccount)

		return nil

	})
}

type AccountEntity struct {
	Id      uint    `grom:"primaryKey;autoIncrement;column:id"`
	Balance float64 `gorm:"column:balance"`
}

func (AccountEntity) TableName() string {
	return "accounts"
}

type TransactionEntity struct {
	Id            uint    `grom:"primaryKey;autoIncrement;column:id"`
	FromAccountId uint    `grom:"column:from_account_id"`
	ToAccountId   uint    `grom:"column:to_account_id"`
	Amount        float64 `grom:"column:amount"`
}

func (TransactionEntity) TableName() string {
	return "transactions"
}

var Db *gorm.DB

func init() {
	host := "localhost"
	port := 3306
	userName := "root"
	password := "1"
	timeout := "10s"

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", userName, password, host, port, DbName, timeout)

	var err error
	Db, err = gorm.Open(mysql.Open(connStr))
	if err != nil {
		fmt.Println("数据库连接失败：", err)
		return
	}

	// 连接成功
	fmt.Println(Db)

	// 添加数据表注释
	Db.Exec("ALTER TABLE accounts COMMENT '账户表'")
}
