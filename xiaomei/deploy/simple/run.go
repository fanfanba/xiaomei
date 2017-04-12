package simple

import (
	"fmt"
)

func (d driver) FlagsForRun(svcName string) ([]string, error) {
	flags := []string{`--network=host`}
	portEnv := portEnvName(svcName)
	ports := portsOf(svcName)
	if portEnv != `` && len(ports) > 0 {
		flags = append(flags, fmt.Sprintf(`-e=%s=%s`, portEnv, ports[0]))
	}
	return flags, nil
}

func portEnvName(svcName string) string {
	switch svcName {
	case `app`:
		return `GOPORT`
	case `web`, `access`:
		return `NGPORT`
	default:
		return ``
	}
}