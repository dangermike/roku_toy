package roku_test

import (
	"testing"

	"github.com/dangermike/roku_toy/roku"

	"github.com/stretchr/testify/require"
)

func TestFuzzyMatch(t *testing.T) {
	rd := &roku.Device{}
	rd.SetApps([]roku.App{
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
	})

	for _, test := range []struct {
		input string
		expID int
	}{
		{"Hulu", 2285},
		{"acorn", 14295},
		{"roku", 151908},
		{"youtube", 837},
		{"plex", 13535},
		{"spotfy", 22297},
		{"xyzzy", 0},
	} {
		t.Run(test.input, func(t *testing.T) {
			app := rd.FindApp(test.input)
			if test.expID == 0 {
				require.Nil(t, app)
				return
			}
			require.NotNil(t, app)
			require.Equal(t, test.expID, app.ID)
		})
	}
}
