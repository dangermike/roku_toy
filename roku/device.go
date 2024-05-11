package roku

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dangermike/roku_toy/logging"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"go.uber.org/zap"
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
	ID      string `xml:"id,attr"`
	Version string `xml:"version,attr"`
}

type ErrApplicationNotFound string

func (e ErrApplicationNotFound) Error() string {
	return fmt.Sprintf("no applications found with a name like '%s'", string(e))
}

func (rd *Device) QueryApps(ctx context.Context) ([]App, error) {
	log := logging.FromContext(ctx)
	log.Debug("getting channels")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rd.Location.JoinPath("query", "apps").String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
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
	log.Debug("got channels")
	return parseApps(body)
}

func (rd *Device) ActiveApp(ctx context.Context) (App, error) {
	log := logging.FromContext(ctx)
	log.Debug("getting channel")
	var app App
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rd.Location.JoinPath("query", "active-app").String(), nil)
	if err != nil {
		return app, err
	}
	resp, err := http.DefaultClient.Do(req)
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

	log.Debug("got channel")

	return apps[0], nil
}

func (rd *Device) Launch(ctx context.Context, id string) error {
	if id == "0" {
		return rd.Home(ctx)
	}
	log := logging.FromContext(ctx)
	log.Debug("setting channel", zap.String("channel_id", id))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rd.Location.JoinPath("launch", id).String(), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to launch app %s: %w", id, err)
	}
	// 200 successful channel change
	// 204 channel already set
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("failed to launch app %s: %s", id, resp.Status)
	}
	log.Debug("set channel", zap.String("channel_id", id))
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

func (rd *Device) LaunchByName(ctx context.Context, name string) error {
	if strings.EqualFold(name, "home") {
		return rd.Home(ctx)
	}

	log := logging.FromContext(ctx)
	log.Debug("setting channel", zap.String("channel", name))

	if len(rd.AppNames) == 0 {
		if len(rd.AppNames) == 0 {
			apps, err := rd.QueryApps(ctx)
			if err != nil {
				return err
			}
			rd.SetApps(apps)
		}
	}

	app := rd.FindApp(name)
	if app == nil {
		return ErrApplicationNotFound(name)
	}

	return rd.Launch(ctx, app.ID)
}

func (rd *Device) Home(ctx context.Context) error {
	log := logging.FromContext(ctx)
	log.Debug("setting channel home")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rd.Location.JoinPath("keypress", "home").String(), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send 'home' keypress: %w", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to send 'home' keypress: %s", resp.Status)
	}
	log.Debug("set channel home")
	return nil
}

func parseApps(data []byte) ([]App, error) {
	apps := appsList{}
	if err := xml.Unmarshal(data, &apps); err != nil {
		return nil, fmt.Errorf("failed to extract apps from xml response: %w", err)
	}
	return apps.Apps, nil
}
