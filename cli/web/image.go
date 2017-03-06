package web

import (
	"path/filepath"

	"github.com/bughou-go/xiaomei/cli/project/stack"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func init() {
	stack.RegisterImage(`web`, webImage{})
}

type webImage struct {
}

func (web webImage) Prepare() error {
	return nil
}

func (web webImage) BuildDir() string {
	return filepath.Join(config.Root(), `../img-web`)
}

func (web webImage) Dockerfile() string {
	return `Dockerfile`
}

func (web webImage) RunMapping() []string {
	return []string{
		filepath.Join(config.Root(), `../img-web/public`) + `:/var/www/` + config.Name(),
	}
}
