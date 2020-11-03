package snowflake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

var (
	Epoch		int64 = 1288834974657
	NodeBits	uint8 = 10
	StepBits	uint8 = 12

	mu			sync.Mutex

	nodeMax		int64 = -1 ^ (-1 << NodeBits)
	nodeMask		  = nodeMax << StepBits
	stepMask	int64 = -1 ^ (-1 << StepBits)

	timeShift 		  = NodeBits + StepBits
	nodeShift		  = StepBits
)

type Node struct {
	mu			sync.Mutex

	epoch		time.Time
	time		int64
	node 		int64
	step		int64

	nodeMax   	int64
	nodeMask  	int64
	stepMask  	int64
	timeShift 	uint8
	nodeShift 	uint8
}

type ID int64

func NewNode(node int64) (*Node, error) {
	n := Node{}
	n.node = node
	n.nodeMax = -1 ^ (-1 << NodeBits)
	n.nodeMask = n.nodeMax << StepBits
	n.stepMask = -1 ^ (-1 << StepBits)
	n.timeShift = NodeBits + StepBits
	n.nodeShift = StepBits

	if n.node < 0 || n.node > n.nodeMax {
		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	}
	var curTime = time.Now()

	n.epoch = curTime.Add(time.Unix(Epoch/1000, (Epoch%1000)*1000000).Sub(curTime))

	return &n, nil
}

func (n *Node) Generate () ID {
	n.mu.Lock()

	now := time.Since(n.epoch).Nanoseconds()

	if now == n.time {
		n.step = (n.step + 1) & n.stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Nanoseconds() / 1000000
			}
		}
	} else {
		n.step = 0
	}

	if now < n.time {
		panic("Clock is moving backwards, rejecting requests")
	}

	n.time = now

	r := ID((now)<<n.timeShift |
		(n.node << n.nodeShift) |
		(n.step),
	)

	n.mu.Unlock()
	return r
}

func (f ID) Int64() int64 {
	return int64(f)
}

func ParseInt64(id int64) ID {
	return ID(id)
}



