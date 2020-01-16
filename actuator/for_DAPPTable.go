package actuator

import (
	"ecology/models"
	"errors"
	"github.com/astaxie/beego/orm"
)

func SelectDAPP(o orm.Ormer, dapp_name, dapp_id, dapp_type string, page *models.Page) (dapp_list []models.DappTable, err error) {
	q := o.QueryTable("dapp_table")
	if dapp_name != "" {
		q = q.Filter("name", dapp_name)
	}
	if dapp_id != "" {
		q = q.Filter("id", dapp_id)
	}
	if dapp_type != "" {
		q = q.Filter("dapp_type", dapp_type)
	}
	_, err = q.All(&dapp_list)

	start, end := InitPage(page, len(dapp_list))

	if end > len(dapp_list) && start < len(dapp_list) {

		return dapp_list[start:], nil

	} else if start > len(dapp_list) {

		return []models.DappTable{}, nil

	} else if end < len(dapp_list) && start < len(dapp_list) {

		return dapp_list[start:end], nil

	}
	return dapp_list, errors.New("err")
}