package main

func Constant_ZoomArray() []int {
	return []int{64, 128, 256} //, 512}
}

func Constant_ZoomMax() int {
	arr := Constant_ZoomArray()
	return arr[len(arr)-1]
}
