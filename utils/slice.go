package utils

import (
	"errors"
)

var ErrorIndexOutOfRange = errors.New("小标超出")

func DeleteAt[T any](src []T, index int) ([]T, error) {
	length := len(src)
	if index < 0 || index >= length {
		return nil, ErrorIndexOutOfRange
	}
	// 实现把下标前面的加入到新数组，后面的再加入
	//result := src[:index]
	//result = append(result, src[index+1:]...)
	//return result, nil

	// 实现直接把index后的值往前排列重新赋值，然后删除最后一位
	for i := index; i+1 < length; i++ {
		src[i] = src[i+1]
	}

	return src[:length-1], nil
}

// Shrink[T any] 缩容
func Shrink[T any](src []T) []T {
	c, l := cap(src), len(src)

	n, changed := calCapacity(c, l)

	if !changed {
		return src
	}

	s := make([]T, 0, n)
	s = append(s, src...)
	return s
}

func calCapacity(cap, length int) (int, bool) {
	// 小于64 不需要缩容
	if cap < 64 {
		return cap, false
	}

	// cap大于2048的时候, 并且src元素不足1半
	// 降低到比一半多一点，大概 5/8的样子
	if cap > 2048 && (cap/length) > 2 {
		return int(float64(cap) * 0.625), true
	}

	// 缩容1半， src元素不足4分之1,
	if cap <= 2048 && (cap/length) >= 4 {
		return cap / 2, true
	}

	return cap, false
}
