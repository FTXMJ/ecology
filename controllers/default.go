package controllers

import (
	"ecology1/models"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
)

type FirstController struct {
	beego.Controller
}

func (this *FirstController) CreateUserAbout() error {
	user_id := this.GetString("user_id")
	user_name := this.GetString("user_name")
	fmt.Println("User_id", user_id)
	fmt.Println("User_name", user_name)
	users := []models.User{}
	models.NewOrm().QueryTable("user").Filter("user_id",user_id).All(&users)
	if len(users) == 0{
		user := models.User{
			Name:   user_name,
			UserId: user_id,
		}
		erruser := user.Insert()
		if erruser != nil {
			return erruser
		}
	}else if len(users) == 1{
		return nil
	}

	/*formula := models.Formula{
		EcologyId:            ,
		Level:               "",
		LowHold:             0,
		HighHold:            0,
		ReturnMultiple:      1,
		HoldReturnRate:      0,
		RecommendReturnRate: 0,
		TeamReturnRate:      0,
	}
	errformula := formula.Insert()
	if errformula != nil {
		return errformula
	}*/
	return errors.New("错误")
}
