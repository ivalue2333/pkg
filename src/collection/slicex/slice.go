package slicex

func Range(start, end int) []int {
	return RangeInternal(start, end, 1)
}

func RangeInternal(start, end int, internal int) []int {

	length := lengthFn(start, end, internal)
	if length == -1 {
		return nil
	}

	sls := make([]int, length)
	index := 0
	for i := start; i < end && index < length; i += internal {
		sls[index] = i
		index += 1
	}
	return sls
}

func lengthFn(start, end int, internal int) int {
	if internal <= 0 {
		return -1
	}

	var remainder int

	if (end-start)%internal != 0 {
		remainder = 1
	}

	return ((end - start) / internal) + remainder
}
