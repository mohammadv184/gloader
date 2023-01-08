package driver

func IsReadableConnection(Conn Connection) bool {
	_, ok := Conn.(ReadableConnection)
	return ok
}

func IsWritableConnection(Conn Connection) bool {
	_, ok := Conn.(WritableConnection)
	return ok
}
