package main

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
)

func GetPing(target string, count int) (string, error) {
	var result string = ""
	pinger, err := ping.NewPinger(target)
	if err != nil {
		return result, err
	}
	pinger.SetPrivileged(true)
	go func() {
		for i := 0; i < count; i++ {
			time.Sleep(time.Second * 1)
		}
		pinger.Stop()
	}()
	pinger.Count = count
	/*
	   pinger.OnRecv = func(pkt *ping.Packet) {
	       fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
	           pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	   }

	   pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
	       fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
	           pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
	   }
	*/
	pinger.OnFinish = func(stat *ping.Statistics) {
		min := float64(stat.MinRtt.Microseconds()) / 1000
		avg := float64(stat.AvgRtt.Microseconds()) / 1000
		max := float64(stat.MaxRtt.Microseconds()) / 1000
		minstr := fmt.Sprintf("# TYPE dns_lookup_mintime_microseconds gauge\ndevPing{cost=\"min\"} %f\n", min)
		avgstr := fmt.Sprintf("# TYPE dns_lookup_avgtime_microseconds gauge\ndevPing{cost=\"avg\"} %f\n", avg)
		maxstr := fmt.Sprintf("# TYPE dns_lookup_maxtime_microseconds gauge\ndevPing{cost=\"max\"} %f\n", max)
		transtr := fmt.Sprintf("# TYPE dns_lookup_package_transmmit gauge\ndevPingPackage{type=\"transmmit\"} %d\n", stat.PacketsSent)
		receivestr := fmt.Sprintf("# TYPE dns_lookup_package_received gauge\ndevPingPackage{type=\"received\"} %d\n", stat.PacketsRecv)
		lossstr := fmt.Sprintf("# TYPE dns_lookup_package_loss_percente gauge\ndevPingPackageLossPercente %v\n", stat.PacketLoss)
		var updown int
		if stat.PacketLoss == 100 {
			updown = 0
		} else {
			updown = 1
		}
		updownStr := fmt.Sprintf("devPingDeviceUp %d\n", updown)
		result = fmt.Sprintf("%s%s%s%s%s%s%s", minstr, avgstr, maxstr, transtr, receivestr, lossstr, updownStr)
	}

	err = pinger.Run()
	if err != nil {
		return result, err
	}
	return result, nil
}
