// keep alive with datanode

package alive

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	timeout     = time.Second * 1 //
	backupN int = 3
)

// record alived datanode
type Alive struct {
	rw sync.RWMutex

	// string records the address of datanode
	// it stored last := time datanode heartbeat to namenode, when I need to use

	// this datanode, I need to check whether the := time is expired or not.
	mp map[string]time.Time //

	balance []string
}

func init() {
}

func NewAlive() *Alive {
	rand.Seed(time.Now().UnixMicro())
	return &Alive{
		mp:      make(map[string]time.Time),
		balance: make([]string, 0),
	}
}

// datanode online
func (a *Alive) Update(key string) {
	a.rw.Lock()
	defer a.rw.Unlock()
	_, ok := a.mp[key]
	if !ok { // not register before
		a.balance = append(a.balance, key)
		a.mp[key] = time.Now()
	}

	a.mp[key] = time.Now()
}

// check datanode is alive or not
// if datanode is not alive, then kick it out
func (a *Alive) IsAlive(key string) bool {
	a.rw.RLock()
	defer a.rw.RUnlock()

	t, ok := a.mp[key]
	if !ok { // datanode haven't register
		return false
	}

	expireTime := time.Since(t)
	return expireTime <= timeout
}

// backup needs to think about the number of datanode and backups,
// because, stored the same backups in one datanode is useless,
// so choose Backup datanode should ompare the number of servers and the number of backups which is smaller
func (a *Alive) Backup() ([]string, error) {
	length := len(a.balance)
	if length == 0 {
		return nil, fmt.Errorf("there is no alived address")
	}

	str := make([]string, 0)

	startIndex := rand.Intn(length) // return [0, len)

	var n int
	if backupN < length {
		n = backupN
	} else {
		n = length
	}

	for i := 0; i < n; i++ {
		index := (i + startIndex) % length
		key := a.balance[index]
		if ok := a.IsAlive(key); ok {
			str = append(str, key)
		}
	}

	if len(str) == 0 {
		return nil, fmt.Errorf("there is no alived address")
	}

	return str, nil
}