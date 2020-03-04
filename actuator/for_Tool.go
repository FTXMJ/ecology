package actuator

import (
	db "ecology/db"
	"ecology/models"
	"errors"

	"github.com/jinzhu/gorm"

	"fmt"
	"strconv"
)

// 根据条件  进行数据查询
func GeneratedSQLAndExec(o *gorm.DB, table_name string, p models.FindObj) (blos []interface{}, err error) {
	err = errors.New("")
	us := []models.User{}
	blo := []models.BlockedDetail{}
	acc := []models.AccountDetail{}
	ac := []models.Account{}

	q_user := o.Table("user")

	if p.UserName != "" {
		q_user = q_user.Where("user_name = ?", p.UserName)
	}
	if p.UserId != "" {
		q_user = q_user.Where("user_id = ?", p.UserId)
	}
	q_user.Find(&us)
	user_ids := []string{}
	for _, v := range us {
		user_ids = append(user_ids, v.UserId)
	}

	q_blos := o.Table(table_name).Where("user_id in (?)", user_ids)

	if p.StartTime != "" && p.EndTime != "" {
		q_blos = q_blos.Where("create_date >= ?", p.StartTime).Where("create_date <= ?", p.EndTime)
	}
	if table_name == "blocked_detail" {
		q_blos.Order("create_date desc").Find(&blo)
		for _, v := range blo {
			blos = append(blos, v)
		}
	} else if table_name == "account_detail" {
		q_blos.Order("create_date desc").Find(&acc)
		for _, v := range acc {
			blos = append(blos, v)
		}
	} else if table_name == "account" {
		q_blos.Order("create_date desc").Find(&ac)
		for _, v := range ac {
			blos = append(blos, v)
		}
	}

	return blos, nil
}

// page 初始化
func InitPage(page *models.Page, count int) (start, end int) {
	page.Count = count
	if page.PageSize < 5 {
		page.PageSize = 5
	}
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}
	start = (page.CurrentPage - 1) * page.PageSize
	end = start + page.PageSize

	page.TotalPage = (page.Count / page.PageSize) + 1 //总页数
	if page.Count <= 5 {
		page.CurrentPage = 1
	}
	return
}

// 数据截取
func ListLimit(list []interface{}, start, end int) []interface{} {
	var listle []interface{}
	if end > len(list) && start < len(list) {
		for _, v := range list[start:] {
			listle = append(listle, v)
		}
	} else if start > len(list) {

	} else if end < len(list) && start < len(list) {
		for _, v := range list[start:end] {
			listle = append(listle, v)
		}
	}
	return listle
}

// 直推算力的计算　　　－－　　　当天
func RecommendReturnRate(user_id, time string) (float64, error) {
	blo := []models.BlockedDetail{}
	sql_str := "SELECT * from blocked_detail where user_id=? and create_date>=? and comment=? "
	er := db.NewEcologyOrm().Raw(sql_str, user_id, time, "直推收益").Find(&blo)
	if er.Error != nil {
		return 0, er.Error
	}
	var zhitui float64
	if len(blo) > 0 {
		for _, v := range blo {
			zhitui += v.CurrentOutlay
		}
	}
	zhit, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", zhitui), 64)
	return zhit, nil
}

// 直推算力的计算　　　－－　　　任意天
func RecommendReturnRateEveryDay(user_id, time_start, time_end string) (float64, error) {
	blo := []models.BlockedDetail{}
	sql_str := "SELECT * from blocked_detail where user_id=? and create_date>=? and create_date<=? and comment=? "
	er := db.NewEcologyOrm().Raw(sql_str, user_id, time_start, time_end, "直推收益").Find(&blo)
	if er.Error != nil {
		return 0, er.Error
	}
	var zhitui float64
	if len(blo) > 0 {
		for _, v := range blo {
			zhitui += v.CurrentOutlay
		}
	}
	return zhitui, nil
}

// 查看用户有史以来所有的收益
func AddAllSum(o *gorm.DB, user_id string) float64 {
	var blos []models.BlockedDetail
	o.Raw("select * from blocked_detail where user_id=? and comment!=?", user_id, "直推收益").Find(&blos)
	var zhitui float64
	if len(blos) > 0 {
		for _, v := range blos {
			zhitui += v.CurrentRevenue
		}
	}
	return zhitui
}

// 处理器，计算所有用户的收益  并发布任务和 分红记录
func HandlerOperation(users []string, user_id string) (float64, error) {
	o := db.NewEcologyOrm()
	var coin_abouns float64
	for _, v := range users {
		// 拿到生态项目实例
		account := models.Account{}
		err_acc := o.Table("account").Where("user_id = ?", v).First(&account)
		if err_acc.Error != nil {
			if err_acc.Error.Error() != "<QuerySeter> no row found" {
				return 0, err_acc.Error
			}
		}
		// 拿到生态项目对应的算力表
		formula := models.Formula{}
		err_for := o.Table("formula").Where("ecology_id = ?", account.Id).First(&formula)
		if err_for.Error != nil {
			if err_for.Error.Error() != "<QuerySeter> no row found" {
				return 0, err_for.Error
			}
		}
		coin_abouns += formula.HoldReturnRate * account.Balance
	}
	return coin_abouns, nil
}
