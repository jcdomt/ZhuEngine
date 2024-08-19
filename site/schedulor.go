// 负载均衡调度器
package site

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

type Schedulor interface {
	Init(*Site) error
	Pick(*Site) string
}

// 轮询算法调度器
type RoundRobinSchedulor struct {
	m     sync.Mutex
	next  int
	items []string
}

func (s *RoundRobinSchedulor) Init(site *Site) error {
	ip_arr := strings.Split(site.Server, ",")
	if len(ip_arr) == 0 {
		return errors.New("没有可用项目")
	}
	for _, ip_weight := range ip_arr {
		a := strings.Split(ip_weight, "?")
		ip := a[0]
		s.items = append(s.items, ip)
	}
	return nil
}

func (s *RoundRobinSchedulor) Pick(site *Site) string {
	s.m.Lock()
	r := s.items[s.next]
	s.next = (s.next + 1) % len(s.items)
	s.m.Unlock()
	return r
}

// 随机算法调度器
type RandomSchedulor struct {
	m     sync.Mutex
	items []string
}

func (s *RandomSchedulor) Init(site *Site) error {
	ip_arr := strings.Split(site.Server, ",")
	if len(ip_arr) == 0 {
		return errors.New("没有可用项目")
	}

	for _, ip_weight := range ip_arr {
		a := strings.Split(ip_weight, "?")
		ip := a[0]
		weight, err := strconv.Atoi(a[1])
		if err != nil {
			return err
		}
		for i := 0; i < weight; i++ {
			s.items = append(s.items, ip)
		}
	}
	return nil
}

func (s *RandomSchedulor) Pick(site *Site) string {
	s.m.Lock()

	// fisher-yates 修正洗牌算法
	for i := len(s.items); i > 0; i-- {
		lastIdx := i - 1
		idx := rand.Intn(i)
		s.items[lastIdx], s.items[idx] = s.items[idx], s.items[lastIdx]
	}
	ret := s.items[0]

	s.m.Unlock()
	return ret
}
