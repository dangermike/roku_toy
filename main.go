package main

import (
	"errors"
	"fmt"
	"net"
	"time"
)

var ssdpHost = net.IPv4(239, 255, 255, 250)

const (
	ssdpPort = 1900
	ssdpBody = `M-SEARCH * HTTP/1.1
Host: 239.255.255.250:1900
Man: "ssdp:discover"
ST: roku:ecp
`

// ST: roku:ecp
)

var ssdpAddr = &net.UDPAddr{
	IP:   ssdpHost,
	Port: ssdpPort,
}

func main() {
	_, raddr, err := getLocalAddr()
	if err != nil {
		panic(err)
	}
	ua := &net.UDPAddr{
		IP:   raddr.IP,
		Port: 0,
		Zone: "",
	}

	lc, err := net.ListenUDP("udp", ua)
	if err != nil {
		panic(err)
	}
	if err := lc.SetReadBuffer(1 << 10); err != nil {
		panic(err)
	}
	ua = lc.LocalAddr().(*net.UDPAddr)

	scnt, err := lc.WriteTo([]byte(ssdpBody), &net.UDPAddr{
		IP: ssdpAddr.IP,
		Port: 1069,
	})
	if err != nil {
		panic(err)
	}

	_ = lc.SetDeadline(time.Now().Add(30 * time.Second))
	fmt.Printf("sent %d bytes %v -> %v\n", scnt, ua.String(), ssdpAddr.String())

	buf := make([]byte, 1<<20)
	for {
		_ = lc.SetDeadline(time.Now().Add(30 * time.Second))
		cnt, addr, err := lc.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}
		fmt.Println("-----------")
		fmt.Printf("read %d bytes from %s\n", cnt, addr.String())
		fmt.Println(string(buf[:cnt]))
	}
}

func roku() {
	_, raddr, err := getLocalAddr()
	if err != nil {
		panic(err)
	}
	ua := &net.UDPAddr{
		IP:   raddr.IP,
		Port: 0,
		Zone: "",
	}

	lc, err := net.ListenUDP("udp", ua)
	if err != nil {
		panic(err)
	}
	if err := lc.SetReadBuffer(1 << 10); err != nil {
		panic(err)
	}
	ua = lc.LocalAddr().(*net.UDPAddr)

	scnt, err := lc.WriteTo([]byte(ssdpBody), ssdpAddr)
	if err != nil {
		panic(err)
	}

	_ = lc.SetDeadline(time.Now().Add(30 * time.Second))
	fmt.Printf("sent %d bytes %v -> %v\n", scnt, ua.String(), ssdpAddr.String())

	buf := make([]byte, 1<<20)
	for {
		_ = lc.SetDeadline(time.Now().Add(30 * time.Second))
		cnt, addr, err := lc.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}
		fmt.Println("-----------")
		fmt.Printf("read %d bytes from %s\n", cnt, addr.String())
		fmt.Println(string(buf[:cnt]))
	}
}

func getLocalAddr() (*net.Interface, *net.IPNet, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}

	for _, iface := range ifaces {
		if 0 == (iface.Flags & net.FlagUp) {
			// fmt.Printf("skipping %s: not up\n", iface.Name)
			continue
		}
		if 0 == (iface.Flags & net.FlagBroadcast) {
			// fmt.Printf("skipping %s: not broadcast-enabled\n", iface.Name)
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, nil, err
		}
		for _, addr := range addrs {
			if a, ok := addr.(*net.IPNet); ok {
				return &iface, a, nil
			}
		}
	}
	return nil, nil, errors.New("Failed to find interface")
}
