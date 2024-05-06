package aliasing

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/dangermike/roku_toy/logging"
	"go.uber.org/zap"
)

type Alias struct {
	USN  string
	Name string
}

func Load(ctx context.Context) ([]Alias, error) {
	log := logging.FromContext(ctx)
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	targetPath := path.Join(home, ".config", "roku_toy", "aliases")

	f, err := os.Open(targetPath)
	if err != nil {
		log.Debug("failed to open aliases file", zap.String("path", targetPath), zap.Error(err))
		return nil, nil
	}
	defer f.Close()
	scn := bufio.NewScanner(f)
	var retval []Alias
	for scn.Scan() {
		if len(scn.Text()) == 0 {
			continue
		}
		kv := strings.SplitN(scn.Text(), ",", 2)
		retval = append(retval, Alias{kv[0], kv[1]})
	}

	return retval, nil
}

func uniqueify(aliases []Alias) []Alias {
	slices.Reverse(aliases)
	seenUSN := map[string]struct{}{}
	seenName := map[string]struct{}{}
	var dropped int

	for i := 0; i < len(aliases); i++ {
		if _, ok := seenUSN[aliases[i].USN]; ok {
			dropped++
			continue
		}
		if _, ok := seenName[aliases[i].Name]; ok {
			dropped++
			continue
		}
		seenUSN[aliases[i].USN] = struct{}{}
		seenName[aliases[i].Name] = struct{}{}
		if dropped > 0 {
			aliases[i-dropped] = aliases[i]
		}
	}
	aliases = aliases[:len(aliases)-dropped]
	slices.Reverse(aliases)
	return aliases
}

func Save(ctx context.Context, aliases []Alias) error {
	log := logging.FromContext(ctx)
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path.Join(home, ".config", "roku_toy")); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Join(home, ".config", "roku_toy"), 0o755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	targetPath := path.Join(home, ".config", "roku_toy", "aliases")

	f, err := os.Create(targetPath)
	if err != nil {
		log.Debug("failed to open aliases file", zap.String("path", targetPath), zap.Error(err))
		return nil
	}
	defer f.Close()
	for _, alias := range uniqueify(aliases) {
		if err := errors.Join(
			errOnly(fmt.Fprint(f, alias.USN)),
			errOnly(fmt.Fprint(f, ",")),
			errOnly(fmt.Fprint(f, alias.Name)),
			errOnly(fmt.Fprint(f, "\n")),
		); err != nil {
			return err
		}
	}

	return f.Close()
}

func errOnly(cnt int, err error) error {
	return err
}
