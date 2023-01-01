package drivers

type Driver interface {
	GetDriverName() string
}
type WritableDriver interface {
}

type ReadableDriver interface {
}
