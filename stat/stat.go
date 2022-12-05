package stat

import (
	"checker/message"
	"time"
)

type Stat struct {
	check        float64
	speed        int64
	averageSpeed float64
	speedch      int64
	averagech    float64
	dtime        time.Time
	time         time.Time
}

func NewStat() *Stat {
	return &Stat{time: time.Now(), dtime: time.Now()}
}

func (s *Stat) Success() {
	s.check++

}
func (s *Stat) Speed() {
	s.averageSpeedCalc()
	s.speedCalc()
}

func (s *Stat) speedCalc() {
	s.speedch++
	if time.Now().Sub(s.dtime).Milliseconds() >= time.Second.Milliseconds() {
		s.speed = s.speedch
		s.dtime = time.Now()
		s.speedch = 0
	}

}

func (s *Stat) averageSpeedCalc() {
	s.averagech++
	if time.Now().Sub(s.time).Seconds() > 0 {
		s.averageSpeed = s.averagech / time.Now().Sub(s.time).Seconds()
	}
}

func (s *Stat) GetStat() message.Stat {
	stat := message.Stat{Success: s.check, Speed: s.speed, AverageSpeed: s.averageSpeed}
	return stat
}
func (s *Stat) Clear() {
	s.averagech = 0
	s.speedch = 0
	s.dtime = time.Now()
	s.time = time.Now()
	s.speed = 0
	s.averageSpeed = 0
}
