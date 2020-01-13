package actuator

import "ecology/models"

// 快速排序
func QuickSortSuperForce(arr []models.SuperForceTable, start, end int) {
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
