package models

// 超级节点算力表
type SuperForceTable struct {
	Id               int     `orm:"column(id);pk;auto" json:"id"`
	Level            string  `orm:"column(level)" json:"level" json:"level"`
	CoinNumberRule   int     `orm:"column(coin_number_rule)" json:"coin_number_rule"`
	BonusCalculation float64 `orm:"column(bonus_calculation)" json:"bonus_calculation"`
}

// 快速排序
func QuickSortSuperForce(arr []SuperForceTable, start, end int) {
	temp := arr[start]
	index := start
	i := start
	j := end

	for i <= j {
		for j >= index && arr[j].CoinNumberRule >= temp.CoinNumberRule {
			j--
		}
		if j > index {
			arr[index] = arr[j]
			index = j
		}
		for i <= index && arr[i].CoinNumberRule <= temp.CoinNumberRule {
			i++
		}
		if i <= index {
			arr[index] = arr[i]
			index = i
		}
	}
	arr[index] = temp
	if index-start > 1 {
		QuickSortSuperForce(arr, start, index-1)
	}
	if end-index > 1 {
		QuickSortSuperForce(arr, index+1, end)
	}
}
