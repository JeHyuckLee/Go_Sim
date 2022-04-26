package executor

import (
	"math/rand"

	"github.com/gammazero/deque"
)

//슬라이스에서 특정한 값을 찾아 리턴
func Slice_Find(slice []interface{}, val interface{}) (int, bool) {

	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func Slice_Find_string(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

//맵에서 특정값을 찾음
func Map_Find(m map[interface{}]interface{}, val interface{}) (interface{}, bool) {
	for k, v := range m {
		if v == val {
			return k, true
		}
	}
	return -1, false
}

func Custom_Sorted(list *deque.Deque) {
	var A []*BehaviorModelExecutor
	length := list.Len()
	for i := 1; i <= length; i++ {
		A = append(A, list.PopFront().(*BehaviorModelExecutor))
	}

	quickSort(A)

	for i := 0; i < length-1; i++ {
		list.PushBack(A[i])
	}
}

func quickSort(arr []*BehaviorModelExecutor) []*BehaviorModelExecutor {

	if len(arr) <= 1 {
		return arr
	}

	median := arr[rand.Intn(len(arr))].Get_req_time()

	lowPart := make([]*BehaviorModelExecutor, 0, len(arr))
	highPart := make([]*BehaviorModelExecutor, 0, len(arr))
	middlePart := make([]*BehaviorModelExecutor, 0, len(arr))

	for _, item := range arr {
		switch {
		case item.Get_req_time() < median:
			lowPart = append(lowPart, item)
		case item.Get_req_time() == median:
			middlePart = append(middlePart, item)
		case item.Get_req_time() > median:
			highPart = append(highPart, item)
		}
	}

	lowPart = quickSort(lowPart)
	highPart = quickSort(highPart)

	lowPart = append(lowPart, middlePart...)
	lowPart = append(lowPart, highPart...)

	return lowPart
}
