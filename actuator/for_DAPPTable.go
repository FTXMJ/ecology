package actuator

import (
	"ecology/models"
	"fmt"
	"github.com/jinzhu/gorm"
)

func SelectDAPP(o *gorm.DB, dapp_name, dapp_id, dapp_type string, page *models.Page) (dapp_list []models.DappTable, err error) {
	q := o.Table("dapp_table")
	if dapp_name != "" {
		q = q.Where("name = ?", dapp_name)
	}
	if dapp_id != "" {
		q = q.Where("id = ?", dapp_id)
	}
	if dapp_type != "" {
		q = q.Where("agreement_type = ?", dapp_type)
	}
	e := q.Find(&dapp_list).Error
	fmt.Println(e)

	start, end := InitPage(page, len(dapp_list))

	if end > len(dapp_list) && start < len(dapp_list) {

		return dapp_list[start:], nil

	} else if start > len(dapp_list) {

		return make([]models.DappTable, 0), nil

	} else if end < len(dapp_list) && start < len(dapp_list) {

		return dapp_list[start:end], nil

	}
	return dapp_list, nil
}
