package actuator

import "ecology/models"

// QuickSort
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
