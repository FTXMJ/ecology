package actuator

import (
	db "ecology/db"
	"ecology/models"
	"errors"
	"github.com/jinzhu/gorm"

	"github.com/astaxie/beego/orm"

	"fmt"
	"strconv"
)

// 根据条件  进行数据查询
func GeneratedSQLAndExec(o orm.Ormer, table_name string, p models.FindObj) (blos []interface{}, err error) {
	err = errors.New("")
	us := []models.User{}
	blo := []models.BlockedDetail{}
	acc := []models.AccountDetail{}
	ac := []models.Account{}

	q_user := o.QueryTable("user")

	if p.UserName != "" {
		q_user = q_user.Filter("user_name", p.UserName)
	}
	if p.UserId != "" {
		q_user = q_user.Filter("user_id", p.UserId)
	}
	q_user.All(&us)
	user_ids := []string{}
	for _, v := range us {
		user_ids = append(user_ids, v.UserId)
	}

	q_blos := o.QueryTable(table_name).Filter("user_id__in", user_ids)

	if p.StartTime != "" && p.EndTime != "" {
		q_blos = q_blos.Filter("create_date__gte", p.StartTime).Filter("create_date__lte", p.EndTime)
	}
	if table_name == "blocked_detail" {
		q_blos.OrderBy("-create_date").All(&blo)
		for _, v := range blo {
			blos = append(blos, v)
		}
	} else if table_name == "account_detail" {
		q_blos.OrderBy("-create_date").All(&acc)
		for _, v := range acc {
			blos = append(blos, v)
		}
	} else if table_name == "account" {
		q_blos.OrderBy("-create_date").All(&ac)
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
	blo := orm.ParamsList{}
	sql_str := "SELECT sum(current_outlay) from blocked_detail where user_id=? and create_date>=? and comment=? "
	_, err := db.NewEcologyOrm().Raw(sql_str, user_id, time, "直推收益").ValuesFlat(&blo)
	if err != nil {
		return 0, err
	}
	var zhitui float64
	if len(blo) > 0 && blo[0] != nil {
		z, _ := strconv.ParseFloat(blo[0].(string), 64)
		zhitui = z
	}
	zhit, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", zhitui), 64)
	return zhit, nil
}

// 直推算力的计算　　　－－　　　任意天
func RecommendReturnRateEveryDay(user_id, time_start, time_end string) (float64, error) {
	blo := orm.ParamsList{}
	sql_str := "SELECT sum(current_outlay) from blocked_detail where user_id=? and create_date>=? and create_date<=? and comment=? "
	_, err := db.NewEcologyOrm().Raw(sql_str, user_id, time_start, time_end, "直推收益").ValuesFlat(&blo)
	if err != nil {
		return 0, err
	}
	var zhitui float64
	if len(blo) > 0 && blo[0] != nil {
		z, _ := strconv.ParseFloat(blo[0].(string), 64)
		zhitui = z
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
		err_acc := o.QueryTable("account").Filter("user_id", v).One(&account)
		if err_acc != nil {
			if err_acc.Error() != "<QuerySeter> no row found" {
				return 0, err_acc
			}
		}
		// 拿到生态项目对应的算力表
		formula := models.Formula{}
		err_for := o.QueryTable("formula").Filter("ecology_id", account.Id).One(&formula)
		if err_for != nil {
			if err_for.Error() != "<QuerySeter> no row found" {
				return 0, err_for
			}
		}
		coin_abouns += formula.HoldReturnRate * account.Balance
	}
	return coin_abouns, nil
}
