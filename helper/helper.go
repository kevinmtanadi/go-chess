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
