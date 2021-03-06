package app

import (
	"fmt"
	"path"
	"strings"

	"github.com/lovego/cmd"
	"github.com/lovego/fs"
	"github.com/lovego/slice"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func depsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `deps`,
		Short: "list dependence packages.",
		Run: func(c *cobra.Command, args []string) {
			listDeps()
		},
	}
	return cmd
}

func listDeps() {
	deps := getAllDeps()
	inVendor := []string{}
	notInVendor := []string{}
	for _, dep := range deps {
		i := strings.Index(dep, `vendor/`)
		if i > -1 {
			inVendor = append(inVendor, dep[i+7:]) // 7 = len(`vendor/`)
		} else {
			notInVendor = append(notInVendor, dep)
		}
	}
	fmt.Printf("<dependence in vendor>:\n%s\n", strings.Join(inVendor, "\n"))
	fmt.Printf("<dependence not in vendor>:\n%s\n", strings.Join(notInVendor, "\n"))
}

func getAllDeps() []string {
	projectDir := path.Join(release.Root(), `../`)
	return getDirDeps(projectDir)
}

var already = make(map[string]bool)

func getDirDeps(dir string) []string {
	if already[dir] {
		return []string{}
	}
	goSrcPath, err := fs.GetGoSrcPath()
	if err != nil {
		panic(err)
	}
	result, _ := cmd.Run(
		cmd.O{Output: true, Dir: dir}, `go`, `list`, `-e`, `-f`, `{{join .Imports "\n"}}`,
	)
	already[dir] = true
	deps := filterStandard(strings.Split(result, "\n"))
	for _, depPath := range deps {
		if strings.HasPrefix(depPath, release.Path()) {
			childDeps := filterStandard(getDirDeps(path.Join(goSrcPath, depPath)))
			for _, childDep := range childDeps {
				if !slice.ContainsString(deps, childDep) {
					deps = append(deps, childDep)
				}
			}
		}
	}
	return filterDeps(deps)
}

// 过滤xiaomei和项目内的包
func filterDeps(deps []string) []string {
	pkgs := []string{}
	for _, dep := range deps {
		if strings.HasPrefix(dep, path.Join(release.Path(), `vendor`)) ||
			!strings.HasPrefix(dep, `github.com/lovego/xiaomei`) &&
				!strings.HasPrefix(dep, release.Path()) {
			pkgs = append(pkgs, dep)
		}
	}
	return pkgs
}

// 过滤标准库的包
func filterStandard(deps []string) []string {
	pkgs := []string{}
	for _, dep := range deps {
		if strings.Contains(dep, `.`) {
			pkgs = append(pkgs, dep)
		}
	}
	return pkgs
}
