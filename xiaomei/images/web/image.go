package web

import (
	"path/filepath"

	"github.com/bughou-go/xiaomei/xiaomei/release"
)

type Image struct {
}

func (i Image) PrepareForBuild() error {
	return nil
}

func (i Image) BuildDir() string {
	return filepath.Join(release.Root(), `img-web`)
}

func (i Image) Dockerfile() string {
	return `Dockerfile`
}

func (i Image) FilesForRun() []string {
	root := filepath.Join(release.Root(), `img-web`)
	return []string{
		root + `/site.conf:/etc/nginx/sites-enabled/` + release.Name(),
		root + `/public:/var/www/` + release.Name(),
	}
}

func (i Image) EnvForRun() []string {
	return []string{`GODEV=true`}
}

func (i Image) CmdForRun() []string {
	return []string{`sh`, `-c`, `nginx -t && nginx`}
}