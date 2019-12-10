package models

type ForceTable struct {
	Id                  int     `orm:"column(id);pk;auto"`
	Level               string  `orm:column(level)`
	LowHold             int     `orm:column(low_hold)`              //低位
	HighHold            int     `orm:column(high_hold)`             //高位
	ReturnMultiple      float64 `orm:column(return_multiple)`       //杠杆
	HoldReturnRate      float64 `orm:column(hold_return_rate)`      //本金自由算力
	RecommendReturnRate float64 `orm:column(recommend_return_rate)` //直推算力
	TeamReturnRate      float64 `orm:column(team_return_rate)`      //动态算力
}


// 快速排序
func QuickSortForce(arr []ForceTable, start, end int) {
	temp := arr[start]
	index := start
	i := start
	j := end

	for i <= j {
		for j >= index && arr[j].LowHold >= temp.LowHold {
			j--
		}
		if j > index {
			arr[index] = arr[j]
			index = j
		}
		for i <= index && arr[i].LowHold <= temp.LowHold {
			i++
		}
		if i <= index {
			arr[index] = arr[i]
			index = i
		}
	}
	arr[index] = temp
	if index-start > 1 {
		QuickSortForce(arr, start, index-1)
	}
	if end-index > 1 {
		QuickSortForce(arr, index+1, end)
	}
}


type ForceTable_test struct {
	Id__id                       int     `orm:"column(id);pk;auto"`
	Level___等级                   string  `orm:column(level)`
	LowHold___充值或者升级的低位          int     `orm:column(low_hold)`              //低位
	HighHold___高位                int     `orm:column(high_hold)`             //高位
	ReturnMultiple___杠杆          float64 `orm:column(return_multiple)`       //杠杆
	HoldReturnRate____自由算力       float64 `orm:column(hold_return_rate)`      //本金自由算力
	RecommendReturnRate____直推算力  float64 `orm:column(recommend_return_rate)` //直推算力
	TeamReturnRate____动态算力__团队算力 float64 `orm:column(team_return_rate)`      //动态算力
}