package roku

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppParse(t *testing.T) {
	appXML := []byte(`<?xml version="1.0" encoding="UTF-8" ?>
	<app id="2285" type="appl" version="6.81.0">Hulu</app>`)

	app := App{}
	require.NoError(t, xml.Unmarshal(appXML, &app))
	require.Equal(t, App{ID: 2285, Version: "6.81.0", Name: "Hulu"}, app)
}

func TestAppsParse(t *testing.T) {
	appsXML := `<?xml version="1.0" encoding="UTF-8" ?>
<apps>
	<app id="2285" type="appl" version="6.81.0">Hulu</app>
	<app id="12" type="appl" version="4.2.100079005">Netflix</app>
	<app id="13535" type="appl" version="7.18.10">Plex - Free Movies &amp; TV</app>
	<app id="837" type="appl" version="2.20.110005159">YouTube</app>
	<app id="551012" type="appl" version="14.2.89">Apple TV</app>
	<app id="1980" type="appl" version="4.2.2036">Vimeo</app>
	<app id="151908" type="appl" version="9.3.10">The Roku Channel</app>
	<app id="23353" type="appl" version="5.6.1">PBS</app>
	<app id="13" type="appl" version="15.1.2024030812">Prime Video</app>
	<app id="164003" type="appl" version="2.16.306230007">Cartoon Network</app>
	<app id="143088" type="appl" version="3.94.3">BritBox</app>
	<app id="14295" type="appl" version="4.23.240318">Acorn TV</app>
	<app id="593099" type="appl" version="5.5.21">Peacock TV</app>
	<app id="22297" type="appl" version="2.11.67">Spotify Music</app>
	<app id="23048" type="appl" version="12.2.0">Spectrum TV</app>
	<app id="636527" type="appl" version="1.2.49">AMC+</app>
	<app id="683311" type="appl" version="10.3.17">Live TV Guide</app>
</apps>`

	apps, err := parseApps([]byte(appsXML))
	require.NoError(t, err)
	require.Equal(t, []App{
		{ID: 2285, Version: "6.81.0", Name: "Hulu"},
		{ID: 12, Version: "4.2.100079005", Name: "Netflix"},
		{ID: 13535, Version: "7.18.10", Name: "Plex - Free Movies &amp; TV"},
		{ID: 837, Version: "2.20.110005159", Name: "YouTube"},
		{ID: 551012, Version: "14.2.89", Name: "Apple TV"},
		{ID: 1980, Version: "4.2.2036", Name: "Vimeo"},
		{ID: 151908, Version: "9.3.10", Name: "The Roku Channel"},
		{ID: 23353, Version: "5.6.1", Name: "PBS"},
		{ID: 13, Version: "15.1.2024030812", Name: "Prime Video"},
		{ID: 164003, Version: "2.16.306230007", Name: "Cartoon Network"},
		{ID: 143088, Version: "3.94.3", Name: "BritBox"},
		{ID: 14295, Version: "4.23.240318", Name: "Acorn TV"},
		{ID: 593099, Version: "5.5.21", Name: "Peacock TV"},
		{ID: 22297, Version: "2.11.67", Name: "Spotify Music"},
		{ID: 23048, Version: "12.2.0", Name: "Spectrum TV"},
		{ID: 636527, Version: "1.2.49", Name: "AMC+"},
		{ID: 683311, Version: "10.3.17", Name: "Live TV Guide"},
	}, apps)
}
