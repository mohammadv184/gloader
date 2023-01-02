package driver

func IsReadableDriver(driver Driver) bool {
	_, ok := driver.(ReadableDriver)
	return ok
}

func IsWritableDriver(driver Driver) bool {
	_, ok := driver.(WritableDriver)
	return ok
}
