package my

type Pair struct {
	A, B interface{}
}

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

func Min(numbers map[interface{}]interface{}) interface{} {
	var minNumber interface{}
	for n := range numbers {
		minNumber = n
		break
	}
	for n := range numbers {
		if n < minNumber {
			minNumber = n
		}
	}
	return minNumber
}
