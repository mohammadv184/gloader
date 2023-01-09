package data

func IsCompatibleWithType(t1, t2 ValueType) bool {
	if _, err := t1.To(t2); err != nil {
		return false
	}
	return true
}
