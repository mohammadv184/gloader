package driver

import "fmt"

type Manager struct {
	drivers map[string]Driver
}

var driverManager *Manager

func init() {
	driverManager = newDriverManager()
}

func Register(driver Driver) error {
	return driverManager.Register(driver)
}
func GetDriver(name string) (Driver, error) {
	return driverManager.GetDriver(name)
}

func Drivers() map[string]Driver {
	return driverManager.Drivers()
}

func (dm *Manager) Register(driver Driver) error {
	dm.drivers[driver.GetDriverName()] = driver
	return nil
}

func (dm *Manager) GetDriver(driverName string) (Driver, error) {
	if driver, ok := dm.drivers[driverName]; ok {
		return driver, nil
	}
	return nil, fmt.Errorf("%v: %s", ErrDriverNotFound, driverName)
}

func (dm *Manager) Drivers() map[string]Driver {
	return dm.drivers
}
func newDriverManager() *Manager {
	return &Manager{drivers: make(map[string]Driver)}
}
