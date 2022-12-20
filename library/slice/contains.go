package slice

type Primitive interface {
	string | int
}

func Contains[T Primitive](value T, slice *[]T) bool {
	if slice == nil {
		return false
	}

	for i := range *slice {
		if (*slice)[i] == value {
			return true
		}
	}

	return false
}
