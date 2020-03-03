package actuator

import (
	"ecology/common"
	"ecology/models"
	"github.com/jinzhu/gorm"

	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func SendHttpPost(urls string, api string, data map[string]string, token string) (error, *common.ResponseData) {
	form := url.Values{}
	for k, v := range data {
		form.Set(k, v)
	}
	u, err := url.ParseRequestURI("http://" + urls)
	if err != nil {
		//log.Log.Error(err)
	}
	u.Path = api
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(form.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if token != "" {
		r.Header.Set("Authorization", token)
	}
	resp, err := client.Do(r)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var rep common.ResponseData
	json.Unmarshal(body, &rep)
	if rep.Code != 200 {
		return errors.New(rep.Msg), nil
	}
	return nil, &rep
}

// 判断是否达到超级节点的要求 --- 页面显示
func SuperLevelSet(o *gorm.DB, user_id string, ec_obj *models.Ecology_index_obj, tfor float64) {
	s_f_t := make([]models.SuperForceTable, 0)
	o.Table("super_force_table").Find(&s_f_t)

	for i := 0; i < len(s_f_t); i++ {
		for j := 1; j < len(s_f_t)-i-1; j++ {
			if s_f_t[i].CoinNumberRule > s_f_t[j].CoinNumberRule {
				s_f_t[i], s_f_t[j] = s_f_t[j], s_f_t[i]
			}
		}
	}
	index := make([]int, 0)
	for i, v := range s_f_t {
		if tfor >= float64(v.CoinNumberRule) {
			index = append(index, i)
		}
	}
	blo := models.TxIdList{}
	start_time := time.Now().Format("2006-01-02") + " 00:00:00"
	end_time := time.Now().Format("2006-01-02") + " 59:59:59"
	o.Raw("select * from tx_id_list where user_id=? and comment=? and create_time>=? and create_time<=? order by create_time desc limit 1", user_id, "节点分红", start_time, end_time).Find(&blo)
	if blo.Id < 1 {
		blo.Expenditure = 0
	}
	if len(index) > 0 {
		ec_obj.Super_peer_bool = true
		ec_obj.Super_peer.Level = s_f_t[index[len(index)-1]].Level
		ec_obj.Super_peer.Usdd = tfor
		ec_obj.Super_peer.TodayABouns = blo.Expenditure
	}
}
