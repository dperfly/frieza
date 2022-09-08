package master

import "errors"

type Slave struct {
	Ch     chan int32
	Ip     string
	Status int32
}

var (
	ErrNotFound = errors.New("not found; unregistered")
	ErrExist    = errors.New("exist")
)

type Slaves interface {
	Create(uuid string) error
	Get(uuid string) (chan int32, error)
	IsExist(uuid string) bool
	UpdateStatus(uuid string, status int32) error
	UpdateIp(uuid string, ip string) error
	Delete(name string) error

	GetAllSlavesSlice() []chan int32 //获取所有slaves slice

	//GetAllSlaves() (map[string]*Slave, error)

	GetAllCanRunIP() []string // 获取所有状态是Idle的IP地址
}
