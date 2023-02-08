package driver

import "fmt"

// Manager is a manager for drivers.
type Manager struct {
	drivers map[string]Driver
}

var driverManager *Manager

func init() {
	driverManager = newDriverManager()
}

// Register registers a driver.
func Register(driver Driver) error {
	return driverManager.Register(driver)
}

// GetDriver returns a driver by name.
func GetDriver(name string) (Driver, error) {
	return driverManager.GetDriver(name)
}

// Drivers returns all drivers.
func Drivers() map[string]Driver {
	return driverManager.Drivers()
}

// Register registers a driver.
func (dm *Manager) Register(driver Driver) error {
	dm.drivers[driver.GetDriverName()] = driver
	return nil
}

// GetDriver returns a driver by name.
func (dm *Manager) GetDriver(driverName string) (Driver, error) {
	if driver, ok := dm.drivers[driverName]; ok {
		return driver, nil
	}
	return nil, fmt.Errorf("%v: %s", ErrDriverNotFound, driverName)
}

// Drivers returns all drivers.
func (dm *Manager) Drivers() map[string]Driver {
	return dm.drivers
}

// newDriverManager returns a new driver manager.
func newDriverManager() *Manager {
	return &Manager{drivers: make(map[string]Driver)}
}
