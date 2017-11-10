package main

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	// Using bastjan fork because sparrc has not applied
	// https://github.com/sparrc/go-ping/pull/15 and other useful fixes yet
	ping "github.com/bastjan/go-ping"
)

type NodePinger struct {
	Nodes        []string
	PingTimeout  time.Duration
	PingInterval time.Duration
}

func NewNodePinger(pingTimeout time.Duration, pingInterval time.Duration) *NodePinger {
	pinger := NodePinger{
		Nodes:        make([]string, 0),
		PingTimeout:  pingTimeout,
		PingInterval: pingInterval,
	}
	return &pinger
}

func (p *NodePinger) Start() {
	log.Info("NodePinger started")

	for {
		var wg sync.WaitGroup
		wg.Add(len(p.Nodes))

		for _, ip := range p.Nodes {
			go pingNode(&wg, ip, p.PingTimeout)
		}

		wg.Wait()

		time.Sleep(p.PingInterval)
	}
}

func pingNode(wg *sync.WaitGroup, ip string, timeout time.Duration) {
	defer wg.Done()

	log.Debugf("Pinging node %s", ip)

	ctx := context.TODO()
	pinger, err := ping.NewPinger(ctx, ip)
	if err != nil {
		log.Warnf("Unable to create pinger for %s, %s", ip, err)
		return
	}

	pinger.SetPrivileged(true)
	pinger.Count = 1
	pinger.Timeout = timeout
	pinger.Run()

	stats := pinger.Statistics()
	log.Debugf("Pinging node %s results: %+v", ip, stats)

	if stats.PacketLoss > 0 || stats.PacketsRecv == 0 || stats.PacketsSent == 0 {
		log.Warnf("Unable to reach %s, %+v", ip, stats)
	}
}
