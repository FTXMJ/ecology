package actuator

import (
	"ecology/models"
	"github.com/jinzhu/gorm"

	"time"
)

// 生成充值表的借贷记录
func FindLimitOneAndSaveAcc_d(o *gorm.DB, user_id, comment, tx_id string, money_out, money_in float64, account_id int) error {
	account_old := models.AccountDetail{}
	o.Raw("select * from account_detail where user_id=? and account=? order by create_date desc,id desc limit 1", user_id, account_id).First(&account_old)
	if account_old.Id == 0 {
		account_old.CurrentBalance = 0
	}
	account_new := models.AccountDetail{
		UserId:         user_id,
		CurrentRevenue: money_in,
		CurrentOutlay:  money_out,
		OpeningBalance: account_old.CurrentBalance,
		CurrentBalance: account_old.CurrentBalance + money_in - money_out,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        comment,
		TxId:           tx_id,
		Account:        account_id,
		CoinType:       "USDD",
	}
	if account_new.CurrentBalance < 0 {
		account_new.CurrentBalance = 0
	}
	err_acc := o.Create(&account_new)
	if err_acc.Error != nil {
		return err_acc.Error
	}

	err_up := o.Model(&models.Account{}).Where("id = ?", account_id).Update("balance", account_new.CurrentBalance)
	if err_up.Error != nil {
		return err_acc.Error
	}
	return nil
}

// 创建第一条充值表的借贷记录
func NewCreateAndSaveAcc_d(o *gorm.DB, user_id, comment, tx_id string, money_out, money_in float64, account_id int) error {
	account_new := models.AccountDetail{
		UserId:         user_id,
		CurrentRevenue: money_in,
		CurrentOutlay:  money_out,
		OpeningBalance: 0,
		CurrentBalance: money_in,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        comment,
		TxId:           tx_id,
		Account:        account_id,
		CoinType:       "USDD",
	}

	err_acc := o.Create(&account_new)
	if err_acc.Error != nil {
		return err_acc.Error
	}

	err_up := o.Model(&models.Account{}).Where("id = ?", account_id).Update("balance", account_new.CurrentBalance)
	if err_up.Error != nil {
		return err_acc.Error
	}
	return nil
}
