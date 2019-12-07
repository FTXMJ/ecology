package models

// Page 分页参数  ---  历史信息
type HostryPageInfo struct {
	Items []HostryValues //数据列表
	Page  Page       //分页信息
}

type Page struct {
	TotalPage   int //总页数
	CurrentPage int //当前页数
	PageSize    int //每页数据条数
	Count       int //总数据量
}

type HostryValues struct {
	Id             int
	UserId         string
	CurrentRevenue float64 //上期支出
	CurrentOutlay  float64 //本期支出
	OpeningBalance float64 //上期余额
	CurrentBalance float64 //本期余额
	CreateDate     string  //创建时间
	Comment        string  //评论
	TxId           string  //任务id
	Account        int     //生态仓库id
}

func SelectHostery(ecology_id int, page Page) ([]HostryValues, Page, error) {
	o := NewOrm()

	var acc_list []AccountDetail
	_, acc_read_err := o.QueryTable("account_detail").Filter("account", ecology_id).All(&acc_list)
	index_values := append_acc_to_public(acc_list)
	if acc_read_err != nil {
		return nil, page, acc_read_err
	}

	var blo_list []BlockedDetail
	_, blo_read_err := o.QueryTable("account_detail").Filter("account", ecology_id).All(&blo_list)
	if blo_read_err != nil {
		return nil, page, blo_read_err
	}
	last_values := append_blo_to_public(blo_list, index_values)
	QuickSortAgreement(last_values, 0, len(last_values)-1)
	page.Count = len(last_values)
	if page.PageSize < 5 {
		page.PageSize = 5
	}
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}
	//listle, _ := o.Limit(page.PageSize, (page.PageNo-1)*page.PageSize).OrderBy("-createtime").All(&list)
	start := (page.CurrentPage-1)*page.PageSize - 1
	end := page.PageSize + 1
	listle := last_values[start:end]
	page.CurrentPage = (page.Count / page.PageSize) + 1 //总页数
	if page.Count <= 5 {
		page.CurrentPage = 1
	}
	return listle, page, nil
}

func append_acc_to_public(acc []AccountDetail) []HostryValues {
	var hostry_values []HostryValues
	for _, v := range acc {
		hos_va := HostryValues{
			Id:             v.Id,
			UserId:         v.UserId,
			CurrentRevenue: v.CurrentBalance,
			CurrentOutlay:  v.CurrentOutlay,
			OpeningBalance: v.OpeningBalance,
			CurrentBalance: v.CurrentBalance,
			CreateDate:     v.CreateDate.String(),
			Comment:        v.Comment,
			TxId:           v.TxId,
			Account:        v.Account,
		}
		hostry_values = append(hostry_values, hos_va)
	}
	return hostry_values
}

func append_blo_to_public(blo []BlockedDetail, hostry_values []HostryValues) []HostryValues {
	for _, v := range blo {
		hos_va := HostryValues{
			Id:             v.Id,
			UserId:         v.UserId,
			CurrentRevenue: v.CurrentBalance,
			CurrentOutlay:  v.CurrentOutlay,
			OpeningBalance: v.OpeningBalance,
			CurrentBalance: v.CurrentBalance,
			CreateDate:     v.CreateDate.String(),
			Comment:        v.Comment,
			TxId:           v.TxId,
			Account:        v.Account,
		}
		hostry_values = append(hostry_values, hos_va)
	}
	return hostry_values
}

// 快速排序
func QuickSortAgreement(arr []HostryValues, start, end int) {
	temp := arr[start]
	index := start
	i := start
	j := end

	for i <= j {
		for j >= index && arr[j].CreateDate <= temp.CreateDate {
			j--
		}
		if j > index {
			arr[index] = arr[j]
			index = j
		}
		for i <= index && arr[i].CreateDate >= temp.CreateDate {
			i++
		}
		if i <= index {
			arr[index] = arr[i]
			index = i
		}
	}
	arr[index] = temp
	if index-start > 1 {
		QuickSortAgreement(arr, start, index-1)
	}
	if end-index > 1 {
		QuickSortAgreement(arr, index+1, end)
	}
}
