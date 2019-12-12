package models

//用户表
type User struct {
	Id       int    `orm:"column(id);pk;auto"`
	Name     string `orm:column(name)`      // 对应 monggodb 的fallname
	UserId   string `orm:column(user_id)`   //对应 monggodb 的user_id
	FatherId string `orm:column(father_id)` //父亲id
}

func (this *User) TableName() string {
	return "user"
}

func (this *User) Insert() error {
	_, err := NewOrm().Insert(this)
	return err
}

func (this *User) Update() (err error) {
	_, err = NewOrm().Update(this)
	return err
}
