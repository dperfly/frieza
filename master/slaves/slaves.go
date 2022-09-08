package slaves

import (
	"Frieza/constant"
	"Frieza/master"
	"fmt"
	"strings"
	"sync"
)

type Slaves struct {
	sync.RWMutex
	SlavesMap map[string]*master.Slave
}

func (c *Slaves) IsExist(uuid string) bool {
	c.RLock()
	defer c.RUnlock()
	if _, ok := c.SlavesMap[uuid]; ok {
		return true
	}
	return false
}
func (c *Slaves) Create(uuid string) error {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.SlavesMap[uuid]; ok {
		return master.ErrExist
	}

	Ch := make(chan int32, 1)
	c.SlavesMap[uuid] = &master.Slave{Ch: Ch}
	return nil
}

func (c *Slaves) Get(uuid string) (chan int32, error) {
	c.RLock()
	defer c.RUnlock()
	slave, ok := c.SlavesMap[uuid]
	if ok {
		return slave.Ch, nil
	}
	return nil, master.ErrExist
}

func (c *Slaves) Delete(uuid string) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.SlavesMap[uuid]; !ok {
		return master.ErrExist
	}
	delete(c.SlavesMap, uuid)
	return nil
}

func (c *Slaves) GetAllSlavesSlice() []chan int32 {
	c.RLock()
	defer c.RUnlock()
	allSlave := make([]chan int32, 0)
	for _, ch := range c.SlavesMap {
		allSlave = append(allSlave, ch.Ch)
	}
	return allSlave
}

//func (c *Slaves) GetAllSlaves() (map[string]*master.Slave, error) {
//	c.RLock()
//	defer c.RUnlock()
//
//	return c.SlavesMap, nil
//}

func (c *Slaves) UpdateStatus(uuid string, status int32) error {
	c.Lock()
	defer c.Unlock()
	// check stream is Exist
	if ch, ok := c.SlavesMap[uuid]; ok {
		ch.Status = status
		return nil
	} else {
		return master.ErrExist
	}

}

func (c *Slaves) UpdateIp(uuid string, ip string) error {
	c.Lock()
	defer c.Unlock()
	if ch, ok := c.SlavesMap[uuid]; ok {
		ch.Ip = ip
		return nil
	} else {
		return master.ErrExist
	}
}

func (c *Slaves) GetAllCanRunIP() []string {
	c.RLock()
	defer c.RUnlock()
	ips := make([]string, 0)
	for _, v := range c.SlavesMap {
		if v.Status == constant.Idle {
			ips = append(ips, v.Ip)
		}
	}
	if len(ips) == 0 {
		return nil
	}
	return ips
}

func (c *Slaves) String() string {
	c.RLock()
	defer c.RUnlock()
	res := strings.Builder{}
	if len(c.SlavesMap) == 0 {
		res.WriteString("not found Frieza slave was run")
	} else {
		res.WriteString("nums\tIp\t\tstatus\n")
	}
	num := 0
	for _, value := range c.SlavesMap {
		if value.Ip != "" {
			num++
			res.WriteString(fmt.Sprintf("%d\t%s\t", num, value.Ip))

			// todo 这里有些重复了
			switch value.Status {
			case constant.Starting:
				res.WriteString(fmt.Sprintf("Starting"))
			case constant.Idle:
				res.WriteString(fmt.Sprintf("Idle"))
			case constant.Failed:
				res.WriteString(fmt.Sprintf("Failed"))
			case constant.Stopped:
				res.WriteString(fmt.Sprintf("Stopped"))
			case constant.Running:
				res.WriteString(fmt.Sprintf("Running"))
			}
			res.WriteString("\n")
		}

	}

	return res.String()
}
