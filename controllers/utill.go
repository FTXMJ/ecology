package controllers

import "ecology/models"

func SelectPeerABounsList(page models.Page, user_name string) ([]models.PeerAbouns, models.Page, error) {
	peer_a_bouns := []models.TxIdList{}
	o := models.NewOrm()
	switch user_name {
	case "":
		o.Raw("select * from tx_id_list order by create_time limit ?,?", page.Count, page.PageSize).QueryRows(&peer_a_bouns)
	default:
		users := []models.User{}
		o.Raw("select * from user where user_name=?", user_name).QueryRows(&users)
		for _, v := range users {
			peers := []models.TxIdList{}
			o.Raw("select * from tx_id_list where user_id=?", v.UserId).QueryRows(&peers)
			for _, vv := range peers {
				peer_a_bouns = append(peer_a_bouns, vv)
			}
		}
	}

	models.QuickSortPeerABouns(peer_a_bouns, 0, len(peer_a_bouns)-1)

	page.Count = len(peer_a_bouns)
	if page.PageSize < 10 {
		page.PageSize = 10
	}
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}
	start := (page.CurrentPage - 1) * page.PageSize
	end := start + page.PageSize

	listle := []models.PeerAbouns{}

	if start > len(peer_a_bouns) {
		return []models.PeerAbouns{}, page, nil
	} else if end > len(peer_a_bouns) {
		end = len(peer_a_bouns)
	}

	for _, v := range peer_a_bouns[start:end] {
		u := models.User{
			UserId: v.UserId,
		}
		o.Read(&u)
		_, level, tfors, err_tfor := ReturnSuperPeerLevel(v.UserId)
		if err_tfor != nil {
			return []models.PeerAbouns{}, page, err_tfor
		}
		p := models.PeerAbouns{
			Id:       v.Id,
			UserName: u.UserName,
			Level:    level,
			Tfors:    tfors,
			Time:     v.CreateTime,
		}
		listle = append(listle, p)
	}
	return listle, page, nil
}
