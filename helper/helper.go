package helper

func Abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func Swap(x, y *int) {
	temp := y
	y = x
	x = temp
}

func ConvertMapToArray[T comparable](inputMap map[T]T) []T {
	array := make([]T, 0)
	for k := range inputMap {
		array = append(array, k)
	}

	return array
}
