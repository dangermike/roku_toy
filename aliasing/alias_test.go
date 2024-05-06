package aliasing

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUniqueify(t *testing.T) {
	for i, test := range []struct {
		src []Alias
		exp []Alias
	}{
		{},
		{
			src: []Alias{{"u_a", "n_a"}, {"u_b", "n_b"}},
			exp: []Alias{{"u_a", "n_a"}, {"u_b", "n_b"}},
		},
		{
			src: []Alias{{"u_a", "n_a"}, {"u_a", "n_c"}, {"u_b", "n_b"}, {"u_b", "n_d"}},
			exp: []Alias{{"u_a", "n_c"}, {"u_b", "n_d"}},
		},
		{
			src: []Alias{{"u_a", "n_a"}, {"u_a", "n_c"}, {"u_b", "n_b"}, {"u_b", "n_c"}},
			exp: []Alias{{"u_a", "n_a"}, {"u_b", "n_c"}},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			act := uniqueify(test.src)
			require.Equal(t, test.exp, act)
		})
	}
}
