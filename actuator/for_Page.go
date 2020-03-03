package actuator

import (
	db "ecology/db"
	"ecology/models"
	"errors"
	"github.com/jinzhu/gorm"
)

// user
func SelectHostery(ecology_id int, page models.Page) ([]models.HostryValues, models.Page, error) {
	o := db.NewEcologyOrm()

	acc_list := make([]models.AccountDetail, 0)
	acc_read_err := o.Table("account_detail").Where("account = ?", ecology_id).Find(&acc_list)
	index_values := append_acc_to_public(acc_list)
	if acc_read_err.Error != nil {
		return nil, page, acc_read_err.Error
	}

	blo_list := make([]models.BlockedDetail, 0)
	blo_read_err := o.Table("blocked_detail").Where("account = ?", ecology_id).Find(&blo_list)
	if blo_read_err.Error != nil {
		return nil, page, blo_read_err.Error
	}
	last_values := append_blo_to_public(blo_list, index_values)
	if len(last_values) == 0 {
		return make([]models.HostryValues, 0), page, errors.New("没有历史交易记录!")
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
func SelectHosteryRoot(o *gorm.DB, page models.Page) ([]models.HostryValues, models.Page, error) {
	acc_list := make([]models.AccountDetail, 0)

	o.Raw("select * from account_detail").Find(&acc_list)
	index_values := append_acc_to_public(acc_list)

	blo_list := make([]models.BlockedDetail, 0)
	o.Raw("select * from blocked_detail").Find(&blo_list)
	last_values := append_blo_to_public(blo_list, index_values)

	if len(last_values) == 0 {
		return make([]models.HostryValues, 0), page, errors.New("没有历史交易记录!")
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
	hostry_values := make([]models.HostryValues, 0)
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

		return make([]models.PeerUser, 0), page

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

		return make([]models.PeerHistory, 0), page

	} else if end < len(peer_user_list) && start < len(peer_user_list) {

		return peer_user_list[start:end], page

	}
	return make([]models.PeerHistory, 0), page
}
