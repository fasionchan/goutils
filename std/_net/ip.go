package _net

import (
	"net"

	"github.com/fasionchan/goutils/stl"
)

var Ipv4Intranets = IpNets{
	{
		IP:   net.ParseIP("10.0.0.0"),
		Mask: net.CIDRMask(8, 32),
	},
	{
		IP:   net.ParseIP("172.16.0.0"),
		Mask: net.CIDRMask(12, 32),
	},
	{
		IP:   net.ParseIP("192.168.0.0"),
		Mask: net.CIDRMask(16, 32),
	},
}

type IpFilter = func(net.IP) bool

type Ips []net.IP

func (ips Ips) Append(more ...net.IP) Ips {
	return append(ips, more...)
}

func (ips Ips) Filter(f IpFilter) Ips {
	return stl.Filter(ips, f)
}

func (ips Ips) AllMatch(f IpFilter) bool {
	return stl.AllMatch(ips, f)
}

func (ips Ips) AllInIpv4Intranets() bool {
	return ips.AllMatch(Ipv4Intranets.Contains)
}

type IpNets []*net.IPNet

func (nets IpNets) Contains(ip net.IP) bool {
	return stl.AnyMatchUnary(nets, (*net.IPNet).Contains, ip)
}

func GetLocalIpv4Addr() (string, error) {
	srcAddr := net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 0,
	}
	dstAddr := net.UDPAddr{
		IP:   net.ParseIP("8.8.8.8"),
		Port: 8,
	}

	conn, err := net.DialUDP("udp", &srcAddr, &dstAddr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr()
	host, _, err := net.SplitHostPort(localAddr.String())
	if err != nil {
		return "", err
	}

	return host, nil
}