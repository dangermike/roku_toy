package roku

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "embed"

	"github.com/dangermike/roku_toy/logging"
	"go.uber.org/zap"
)

const (
	ssdpPort = 1900
)

//go:embed ssdp_req.txt
var ssdpBody string

var (
	ErrNotRoku = errors.New("SSPD ST header is not roku:ecp")

	rxCacheControl = regexp.MustCompile("^max-age=([0-9]+)$")

	ssdpHost = net.IPv4(239, 255, 255, 250)
	ssdpAddr = &net.UDPAddr{
		IP:   ssdpHost,
		Port: ssdpPort,
	}
)

func SSDP(ctx context.Context, cb func(*Device) error) error {
	log := logging.FromContext(ctx)
	_, raddr, err := getLocalAddr(ctx)
	if err != nil {
		return fmt.Errorf("failed to get local address: %w", err)
	}
	ua := &net.UDPAddr{
		IP:   raddr.IP,
		Port: 0,
		Zone: "",
	}

	lc, err := net.ListenUDP("udp", ua)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	if err := lc.SetReadBuffer(1 << 10); err != nil {
		return fmt.Errorf("failed to set read buffer: %w", err)
	}
	ua = lc.LocalAddr().(*net.UDPAddr)

	sent, err := lc.WriteTo([]byte(ssdpBody), ssdpAddr)
	if err != nil {
		return fmt.Errorf("failed to send SSDP request: %w", err)
	}

	log.Debug("sent ssdp request", zap.Int("bytes", sent), zap.String("from", ua.String()), zap.String("to", ssdpAddr.String()))

	buf := make([]byte, 1<<20)
	if err := lc.SetDeadline(time.Now().Add(2 * time.Second)); err != nil {
		return fmt.Errorf("failed to set read deadline: %w", err)
	}
	start := time.Now()
	for {
		cnt, addr, err := lc.ReadFromUDP(buf)
		if os.IsTimeout(err) {
			break
		}
		if err != nil {
			return fmt.Errorf("failed upp (SSDP) read: %w", err)
		}
		log.Debug("read udp bytes", zap.Int("bytes", cnt), zap.String("source", addr.String()), zap.Duration("time", time.Since(start)))
		rokuDev, err := handleSSDPResponse(rokuDevFromHTTP, buf[:cnt])
		if err == ErrNotRoku {
			continue
		}
		if err != nil {
			return fmt.Errorf("failed to get roku device from SSDP response: %w", err)
		}

		if err := cb(rokuDev); err != nil {
			return fmt.Errorf("SSDP callback returned error: %w", err)
		}
	}
	return nil
}

func getLocalAddr(ctx context.Context) (*net.Interface, *net.IPNet, error) {
	log := logging.FromContext(ctx)
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}

	for _, iface := range ifaces {
		if ok, a, err := getIfaceAddr(ctx, &iface); ok {
			log.Debug("using interface", zap.String("name", iface.Name), zap.String("address", a.IP.String()))
			return &iface, a, nil
		} else if err != nil {
			return nil, nil, err
		}
	}
	return nil, nil, errors.New("Failed to find interface")
}

func getIfaceAddr(ctx context.Context, iface *net.Interface) (bool, *net.IPNet, error) {
	log := logging.FromContext(ctx)
	if iface == nil {
		return false, nil, nil
	}
	if 0 == (iface.Flags & net.FlagUp) {
		log.Debug("skipping interface", zap.String("name", iface.Name), zap.String("reason", "down"))
		return false, nil, nil
	}
	if 0 == (iface.Flags & net.FlagBroadcast) {
		log.Debug("skipping interface", zap.String("name", iface.Name), zap.String("reason", "not broadast enabled"))
		return false, nil, nil
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return false, nil, err
	}
	for _, addr := range addrs {
		if a, ok := addr.(*net.IPNet); ok && len(a.Mask) == 4 {
			return true, a, nil
		}
	}
	return false, nil, nil
}

func parseCacheControl(cc string) (time.Duration, error) {
	m := rxCacheControl.FindStringSubmatch(cc)
	if len(m) < 2 {
		return 0, fmt.Errorf("Cannot parse '%s' as max-age header", cc)
	}
	secs, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Cannot parse '%s' as max-age header: %w", cc, err)
	}
	return time.Duration(secs) * time.Second, nil
}

func handleSSDPResponse(parser func(*http.Response) (*Device, error), ssdp []byte) (*Device, error) {
	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(ssdp)), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSDP response: %w", err)
	}
	return parser(resp)
}

func rokuDevFromHTTP(resp *http.Response) (*Device, error) {
	if resp.StatusCode != 200 || resp.Header.Get("ST") != "roku:ecp" {
		return nil, ErrNotRoku
	}

	rokuUrl, err := url.Parse(resp.Header.Get("location"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url '%s': %w", resp.Header.Get("location"), err)
	}
	broadcastInterval, err := parseCacheControl(resp.Header.Get("Cache-Control"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cache-control '%s': %w", resp.Header.Get("Cache-Control"), err)
	}

	rokuDev := Device{
		Location:          rokuUrl,
		USN:               strings.TrimPrefix(resp.Header.Get("USN"), "uuid:roku:ecp:"),
		DeviceGroup:       resp.Header.Get("device-group.roku.com"),
		BroadcastInterval: broadcastInterval,
	}

	return &rokuDev, nil
}
