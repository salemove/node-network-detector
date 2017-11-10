package main

import (
	"reflect"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	kubeapi "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func MonitorNodes(client *kubernetes.Clientset, pinger *NodePinger, nodeFetchInterval time.Duration) {
	for {
		log.Debug("Fetching nodes list")

		nodes, err := client.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			log.Warnf("Fetching nodes failed: %s", err)
			continue
		}
		ips := getIpsForNodes(nodes.Items)
		sort.Strings(ips)

		log.Debugf("Found nodes: %v", ips)

		if !reflect.DeepEqual(ips, pinger.Nodes) {
			log.Infof("Nodes changed, previously: %v, now: %v", pinger.Nodes, ips)
			pinger.Nodes = ips
		}

		time.Sleep(nodeFetchInterval)
	}
}

func getIpsForNodes(nodes []kubeapi.Node) []string {
	ips := make([]string, 0)
	for _, node := range nodes {
		if ip := getNodeInternalIp(node); ip != "" {
			ips = append(ips, ip)
		}
	}
	return ips
}

func getNodeInternalIp(node kubeapi.Node) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == kubeapi.NodeInternalIP && addr.Address != "" {
			return addr.Address
		}
	}

	return ""
}
