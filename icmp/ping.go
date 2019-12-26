package icmp

import (
	"errors"
	"github.com/gy-games-libs/golang/x/net/icmp"
	"github.com/gy-games-libs/golang/x/net/ipv4"
	"math/rand"
	"net"
	"time"
)

func RunPing(Addr string, maxrtt time.Duration, maxttl int, seq int) (float64, error) {
	var res pkg
	var err error
	res.dest, err = net.ResolveIPAddr("ip", Addr)
	if err != nil {
		return 0, errors.New("Unable to resolve destination host")
	}
	res.maxrtt = maxrtt
	//res.id = rand.Int() % 0x7fff
	res.id = rand.Intn(65535)
	res.seq = seq
	res.msg = icmp.Message{Type: ipv4.ICMPTypeEcho, Code: 0, Body: &icmp.Echo{ID: res.id, Seq: res.seq}}
	res.netmsg, err = res.msg.Marshal(nil)
	if nil != err {
		return 0, err
	}
	pingRsult := res.Send(maxttl)
	return float64(pingRsult.RTT.Nanoseconds()) / 1e6, pingRsult.Error
}

func Ping(host string, count int, interval int) (int, float64) {
	var AvgDelay float64
	var SendPk, RevcPk, LossPk int
	for i := 0; i <= count; i++ {
		delay, err := RunPing(host, 1*time.Second, 64, i)
		if err == nil {
			AvgDelay = AvgDelay + delay
			RevcPk = RevcPk + 1
		} else {
			LossPk = LossPk + 1
		}
		SendPk = SendPk + 1
		time.Sleep(time.Millisecond * time.Duration(interval))
	}

	LossRate := int((float64(LossPk) / float64(SendPk)) * 100)
	if RevcPk > 0 {
		AvgDelay = AvgDelay / float64(RevcPk)
	} else {
		AvgDelay = 0.0
	}
	return LossRate, AvgDelay
}
