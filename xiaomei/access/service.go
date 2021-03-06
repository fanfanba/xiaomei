package access

import (
	"errors"
	"strings"

	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/release"
)

type service struct {
	Env     string
	svcName string
	addrs   []string
}

func newService(svcName, env string) *service {
	if conf.HasService(svcName, env) {
		return &service{Env: env, svcName: svcName}
	} else {
		return nil
	}
}

func (s *service) Addrs() ([]string, error) {
	if s == nil {
		return nil, nil
	}
	if s.addrs == nil {
		addrs := []string{}
		instances := conf.GetService(s.svcName, s.Env).Instances()
		for _, node := range s.Nodes() {
			for _, instance := range instances {
				addrs = append(addrs, node.GetListenAddr()+`:`+instance)
			}
		}
		s.addrs = addrs
		if len(addrs) == 0 {
			return nil, errors.New(`no instance defined for: ` + s.svcName)
		}
	}
	return s.addrs, nil
}

func (s *service) Nodes() (nodes []cluster.Node) {
	if s == nil {
		return nil
	}
	nodesCondition := conf.GetService(s.svcName, s.Env).Nodes
	for _, node := range cluster.Get(s.Env).GetNodes(``) {
		if node.Match(nodesCondition) {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (s *service) DeployName() string {
	if s == nil {
		return ``
	}
	return release.AppConf(s.Env).DeployName()
}

func (s *service) Domain() string {
	if s == nil {
		return ``
	}
	domain := release.AppConf(s.Env).Domain
	parts := strings.SplitN(domain, `.`, 2)
	if len(parts) == 2 {
		return parts[0] + `-` + s.svcName + `.` + parts[1]
	} else {
		return domain + `-` + s.svcName
	}
}
