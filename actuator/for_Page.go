package actuator

import (
	db "ecology/db"
	"ecology/models"
	"errors"

	"github.com/astaxie/beego/orm"
)

// user
func SelectHostery(ecology_id int, page models.Page) ([]models.HostryValues, models.Page, error) {
	o := db.NewEcologyOrm()

	var acc_list []models.AccountDetail
	_, acc_read_err := o.QueryTable("account_detail").Filter("account", ecology_id).All(&acc_list)
	index_values := append_acc_to_public(acc_list)
	if acc_read_err != nil {
		return nil, page, acc_read_err
	}

	var blo_list []models.BlockedDetail
	_, blo_read_err := o.QueryTable("blocked_detail").Filter("account", ecology_id).All(&blo_list)
	if blo_read_err != nil {
		return nil, page, blo_read_err
	}
	last_values := append_blo_to_public(blo_list, index_values)
	if len(last_values) == 0 {
		return []models.HostryValues{}, page, errors.New("没有历史交易记录!")
	}
	QuickSortAgreement(last_values, 0, len(last_values)-1)
	page.Count = len(last_values)
	if page.PageSize < 10 {
		page.PageSize = 10
	}
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}
	start := (page.CurrentPage - 1) * page.PageSize
	end := start + page.PageSize
	listle := []models.HostryValues{}
	if end > len(last_values) && start < len(last_values) {
		for _, v := range last_values[start:] {
			listle = append(listle, v)
		}
	} else if start > len(last_values) {

	} else if end <= len(last_values) && start < len(last_values) {
		for _, v := range last_values[start:end] {
			listle = append(listle, v)
		}
	}
	page.TotalPage = (page.Count / page.PageSize) + 1 //总页数
	if page.Count <= 5 {
		page.CurrentPage = 1
	}
	return listle, page, nil
}

//root
func SelectHosteryRoot(o orm.Ormer, page models.Page) ([]models.HostryValues, models.Page, error) {
	var acc_list []models.AccountDetail
	o.Raw("select * from account_detail").QueryRows(&acc_list)
	index_values := append_acc_to_public(acc_list)

	var blo_list []models.BlockedDetail
	o.Raw("select * from blocked_detail").QueryRows(&blo_list)
	last_values := append_blo_to_public(blo_list, index_values)

	if len(last_values) == 0 {
		return []models.HostryValues{}, page, errors.New("没有历史交易记录!")
	}

	QuickSortAgreement(last_values, 0, len(last_values)-1)
	page.Count = len(last_values)
	if page.PageSize < 10 {
		page.PageSize = 10
	}
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}

	start := (page.CurrentPage - 1) * page.PageSize
	end := start + page.PageSize
	listle := make([]models.HostryValues, 0)

	if end > len(last_values) && start < len(last_values) {
		for _, v := range last_values[start:] {
			listle = append(listle, v)
		}
	} else if start > len(last_values) {

	} else if end < len(last_values) && start < len(last_values) {
		for _, v := range last_values[start:end] {
			listle = append(listle, v)
		}
	}
	page.TotalPage = (page.Count / page.PageSize) + 1 //总页数
	if page.Count <= 5 {
		page.CurrentPage = 1
	}
	return listle, page, nil
}

func append_acc_to_public(acc []models.AccountDetail) []models.HostryValues {
	var hostry_values []models.HostryValues
	for _, v := range acc {
		hos_va := models.HostryValues{
			Id:             v.Id,
			UserId:         v.UserId,
			CurrentRevenue: v.CurrentRevenue,
			CurrentOutlay:  v.CurrentOutlay,
			OpeningBalance: v.OpeningBalance,
			CurrentBalance: v.CurrentBalance,
			CreateDate:     v.CreateDate,
			Comment:        v.Comment,
			TxId:           v.TxId,
			Account:        v.Account,
		}
		hostry_values = append(hostry_values, hos_va)
	}
	return hostry_values
}

func append_blo_to_public(blo []models.BlockedDetail, hostry_values []models.HostryValues) []models.HostryValues {
	for _, v := range blo {
		hos_va := models.HostryValues{
			Id:             v.Id,
			UserId:         v.UserId,
			CurrentRevenue: v.CurrentRevenue,
			CurrentOutlay:  v.CurrentOutlay,
			OpeningBalance: v.OpeningBalance,
			CurrentBalance: v.CurrentBalance,
			CreateDate:     v.CreateDate,
			Comment:        v.Comment,
			TxId:           v.TxId,
			Account:        v.Account,
		}
		hostry_values = append(hostry_values, hos_va)
	}
	return hostry_values
}

//page
func PageS(peer_user_list []models.PeerUser, page models.Page) ([]models.PeerUser, models.Page) {
	page.Count = len(peer_user_list)
	if page.PageSize < 5 {
		page.PageSize = 5
	}
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}
	start := (page.CurrentPage - 1) * page.PageSize
	end := start + page.PageSize
	page.TotalPage = (page.Count / page.PageSize) + 1 //总页数
	if page.Count <= 5 {
		page.CurrentPage = 1
	}

	if end > len(peer_user_list) && start < len(peer_user_list) {

		return peer_user_list[start:], page

	} else if start > len(peer_user_list) {

		return []models.PeerUser{}, page

	} else if end < len(peer_user_list) && start < len(peer_user_list) {

		return peer_user_list[start:end], page

	}
	return nil, page
}

//page peer_history
func PageHistory(peer_user_list []models.PeerHistory, page models.Page) ([]models.PeerHistory, models.Page) {
	page.Count = len(peer_user_list)
	if page.PageSize < 5 {
		page.PageSize = 5
	}
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}
	start := (page.CurrentPage - 1) * page.PageSize
	end := start + page.PageSize
	page.TotalPage = (page.Count / page.PageSize) + 1 //总页数
	if page.Count <= 5 {
		page.CurrentPage = 1
	}

	if end > len(peer_user_list) && start < len(peer_user_list) {

		return peer_user_list[start:], page

	} else if start > len(peer_user_list) {

		return []models.PeerHistory{}, page

	} else if end < len(peer_user_list) && start < len(peer_user_list) {

		return peer_user_list[start:end], page

	}
	return []models.PeerHistory{}, page
}
