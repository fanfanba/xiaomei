package develop

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Spec(t string) bool {
	if err := os.Chdir(filepath.Join(config.Root(), `..`)); err != nil {
		panic(err)
	}

	targets := specTargets()
	if t != `all` {
		targets = specChangedTargets(targets)
	}
	if len(targets) == 0 {
		return true
	}

	if !cmd.Ok(cmd.O{NoStdout: true}, `which`, `gospec`) {
		cmd.Run(cmd.O{Panic: true}, `go`, `get`, `-u`, `github.com/bughou-go/spec/gospec`)
	}

	return cmd.Ok(cmd.O{}, `gospec`, targets...)
}

func specTargets() []string {
	matches, err := filepath.Glob(`*`)
	if err != nil {
		panic(err)
	}

	targets := []string{}
	for _, v := range matches {
		if v != `release` && v != `vendor` {
			targets = append(targets, v)
		}
	}
	return targets
}

func specChangedTargets(targets []string) []string {
	output, _ := cmd.Run(cmd.O{Output: true}, `git`,
		append([]string{`diff`, `--name-only`, `--diff-filter=AMR`, `--`}, targets...)...,
	)
	lines := strings.Split(output, "\n")
	results := []string{}
	for _, v := range lines {
		v = strings.TrimSpace(v)
		if len(v) > 0 {
			results = append(results, v)
		}
	}
	return results
}
