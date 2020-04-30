package server

import (
	"sync/atomic"

	"github.com/mason-leap-lab/infinicache/common/logger"
	"github.com/mason-leap-lab/infinicache/proxy/global"
)

type Scaler struct {
	log     logger.ILogger
	proxy   *Proxy
	Signal  chan struct{}
	ready   chan struct{}
	counter int32
}

func NewScaler() *Scaler {
	s := &Scaler{
		log: &logger.ColorLogger{
			Prefix: "Scaler ",
			Level:  global.Log.GetLevel(),
			Color:  true,
		},
		Signal:  make(chan struct{}, 1),
		ready:   make(chan struct{}),
		counter: 0,
	}
	return s
}

// check the cluster usage information periodically
func (s *Scaler) Daemon() {
	for {
		s.log.Debug("in scaler Daemon, Group len is %v", s.proxy.group.Len())
		select {
		// receive scaling out signal
		case <-s.Signal:

			// get current bucket
			bucket := s.proxy.movingWindow.getCurrentBucket()
			tmpGroup := NewGroup(NumLambdaClusters)

			for i := range tmpGroup.All {
				node := scheduler.GetForGroup(tmpGroup, i)
				node.Meta.Capacity = InstanceCapacity
				node.Meta.IncreaseSize(InstanceOverhead)
				//s.log.Debug("[scaling lambda instance %v, size %v]", node.Name(), node.Size())

				go func() {
					node.WarmUp()
					if atomic.AddInt32(&s.counter, 1) == int32(len(tmpGroup.All)) {
						s.log.Info("[scale out is ready]")
					}
				}()

				// Begin handle requests
				go node.HandleRequests()
			}

			// reset counter
			s.counter = 0

			// update bucket and placer info

			// append tmpGroup to current bucket group
			bucket.append(tmpGroup)

			// add tmp group and update proxy group pointer
			s.proxy.movingWindow.appendToGroup(tmpGroup)
			atomic.AddInt32(&s.proxy.placer.pointer, NumLambdaClusters)

			//scale out phase finished
			s.proxy.placer.scaling = false
			s.log.Debug("scale out finish")
		}
	}
}

func (s *Scaler) Ready() chan struct{} {
	return s.ready
}
