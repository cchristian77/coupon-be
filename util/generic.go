package util

// GetPointerValue will try to get the underlying value from a pointer.
// If the pointer is nil, it will create a new Zero-Value based on the type.
func GetPointerValue[T any](ptr *T) T {
	if ptr == nil {
		var x interface{}
		ZeroValue, _ := x.(T)
		return ZeroValue
	}
	return *ptr
}

// ToPointerValue will convert a value to a pointer.
func ToPointerValue[T any](value T) *T {
	return &value
}

// Contains will check whether Target exists in list of Source.
func Contains[T comparable](source []T, target T) bool {
	for _, item := range source {
		if item == target {
			return true
		}
	}

	return false
}

// ToUniqueSlices will remove duplicates from a slice
func ToUniqueSlices[T comparable](source []T) []T {
	result := make([]T, 0, len(source))
	seen := make(map[T]bool)

	for _, item := range source {
		if !seen[item] {
			result = append(result, item)
			seen[item] = true
		}
	}

	return result
}
