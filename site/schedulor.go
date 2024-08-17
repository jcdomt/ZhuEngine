// 负载均衡调度器
package site

import (
	"errors"
	"strings"
	"sync"
)

type Schedulor interface {
	Init(*Site) error
	Pick(*Site) *Site
}

// 轮询算法调度器
type RoundRobinSchedulor struct {
	m     sync.Mutex
	next  int
	items []string
}

func (s *RoundRobinSchedulor) Init(site *Site) error {
	ip_arr := strings.Split(site.Config.Server, ",")
	if len(ip_arr) == 0 {
		return errors.New("没有可用项目")
	}
	s.items = ip_arr
	return nil
}

func (s *RoundRobinSchedulor) Pick(site *Site) *Site {
	s.m.Lock()
	r := s.items[s.next]
	s.next = (s.next + 1) % len(s.items)
	s.m.Unlock()
	site.Server = r
	return site
}
