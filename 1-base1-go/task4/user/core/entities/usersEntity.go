package entities

// 用户实体
type UserEntity struct {
	Id       uint64 `gorm:"primary_key;autoIncrement"`
	Username string
	Password string
	Salt     string
	Email    string
}

func (UserEntity) TableName() string {
	return "users"
}
