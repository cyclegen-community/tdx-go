package core

import "sync"

// todo: 链接池，无感切换
type Hub struct {
	lock    sync.Mutex
	Clients []Client
}

func (hub *Hub) Do() error {
	return nil
}
