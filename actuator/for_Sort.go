package actuator

import "ecology/models"

// 快速排序   ForceTable_Sort
func QuickSortForce(arr []models.ForceTable, start, end int) {
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

// 快速排序   SuperForceTable_Sort
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

// 快速排序     TxIdList_Sort
func QuickSortPeerABouns(arr []models.TxIdList, start, end int) {
	temp := arr[start]
	index := start
	i := start
	j := end

	for i <= j {
		for j >= index && arr[j].CreateTime <= temp.CreateTime {
			j--
		}
		if j > index {
			arr[index] = arr[j]
			index = j
		}
		for i <= index && arr[i].CreateTime >= temp.CreateTime {
			i++
		}
		if i <= index {
			arr[index] = arr[i]
			index = i
		}
	}
	arr[index] = temp
	if index-start > 1 {
		QuickSortPeerABouns(arr, start, index-1)
	}
	if end-index > 1 {
		QuickSortPeerABouns(arr, index+1, end)
	}
}

// 快速排序      HostryValues_Sort
func QuickSortAgreement(arr []models.HostryValues, start, end int) {
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

// 快速排序     BlockedDetail_Sort
func QuickSortBlockedDetail(arr []models.BlockedDetail, start, end int) {
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
		QuickSortBlockedDetail(arr, start, index-1)
	}
	if end-index > 1 {
		QuickSortBlockedDetail(arr, index+1, end)
	}
}
