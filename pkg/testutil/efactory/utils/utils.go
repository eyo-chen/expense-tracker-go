package utils

func CvtToAnysWithOW[T any](i int, ow *T) []interface{} {
	res := make([]interface{}, i)
	for k := 0; k < i; k++ {
		if ow != nil {
			res[k] = ow
		} else {
			var v T
			ptr := &v
			res[k] = ptr
		}
	}

	return res
}

func CvtToT[T any](vals []interface{}) []T {
	res := make([]T, len(vals))
	for k, v := range vals {
		res[k] = *v.(*T)
	}

	return res
}
