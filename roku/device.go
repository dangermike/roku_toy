package roku

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Device struct {
	Location          *url.URL
	USN               string
	DeviceGroup       string
	BroadcastInterval time.Duration
	Apps              []App
	AppNames          []string
}

type appsList struct {
	Apps []App `xml:"app"`
}

type App struct {
	Name    string `xml:",innerxml"`
	ID      int    `xml:"id,attr"`
	Version string `xml:"version,attr"`
}

func (rd *Device) QueryApps() ([]App, error) {
	resp, err := http.DefaultClient.Get(rd.Location.JoinPath("query", "apps").String())
	if err != nil {
		return nil, fmt.Errorf("failed to get apps from roku: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get apps from roku: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to get body from from roku apps response: %w", err)
	}
	return parseApps(body)
}

func (rd *Device) ActiveApp() (App, error) {
	var app App
	resp, err := http.DefaultClient.Get(rd.Location.JoinPath("query", "active-app").String())
	if err != nil {
		return app, fmt.Errorf("failed to get active app from roku: %w", err)
	}
	if resp.StatusCode != 200 {
		return app, fmt.Errorf("failed to get get active app from roku: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return app, fmt.Errorf("failed to get body from from roku active-app response: %w", err)
	}
	apps, err := parseApps(body)
	if err != nil {
		return app, fmt.Errorf("failed to parse roku active-app response: %w", err)
	}
	if len(apps) != 1 {
		return app, fmt.Errorf("failed to roku active-app response contained %d apps", len(apps))
	}

	return apps[0], nil
}

func (rd *Device) Launch(id int) error {
	if id == 0 {
		return rd.Home()
	}
	resp, err := http.DefaultClient.Post(rd.Location.JoinPath("launch", strconv.Itoa(id)).String(), "", nil)
	if err != nil {
		return fmt.Errorf("failed to launch app %d: %w", id, err)
	}
	// 200 successful channel change
	// 204 channel already set
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("failed to launch app %d: %s", id, resp.Status)
	}
	return nil
}

func (rd *Device) SetApps(apps []App) {
	rd.Apps = apps
	rd.AppNames = make([]string, len(rd.Apps))
	for i, a := range rd.Apps {
		rd.AppNames[i] = a.Name
	}
}

func (rd *Device) FindApp(name string) *App {
	ranks := fuzzy.RankFindFold(name, rd.AppNames)
	if len(ranks) == 0 {
		return nil
	}
	return &rd.Apps[ranks[0].OriginalIndex]
}

func (rd *Device) LaunchByName(name string) error {
	if strings.EqualFold(name, "home") {
		return rd.Home()
	}

	if len(rd.AppNames) == 0 {
		if len(rd.AppNames) == 0 {
			apps, err := rd.QueryApps()
			if err != nil {
				return err
			}
			rd.SetApps(apps)
		}
	}

	app := rd.FindApp(name)
	if app == nil {
		return fmt.Errorf("no applications found with a name like '%s'", name)
	}

	return rd.Launch(app.ID)
}

func (rd *Device) Home() error {
	resp, err := http.DefaultClient.Post(rd.Location.JoinPath("keypress", "home").String(), "", nil)
	if err != nil {
		return fmt.Errorf("failed to send 'home' keypress: %w", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to send 'home' keypress: %s", resp.Status)
	}
	return nil
}

func parseApps(data []byte) ([]App, error) {
	apps := appsList{}
	if err := xml.Unmarshal(data, &apps); err != nil {
		return nil, fmt.Errorf("failed to extract apps from xml response: %w", err)
	}
	return apps.Apps, nil
}
