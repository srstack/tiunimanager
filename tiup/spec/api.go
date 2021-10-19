// Copyright 2021 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package spec

import (
	"context"
	"crypto/tls"
	"fmt"
	"path/filepath"
	"time"

	"github.com/pingcap-inc/tiem/tiup/templates/config"

	"github.com/pingcap-inc/tiem/tiup/templates/scripts"
	"github.com/pingcap/tiup/pkg/cluster/ctxt"
	"github.com/pingcap/tiup/pkg/cluster/spec"
	"github.com/pingcap/tiup/pkg/logger/log"
	"github.com/pingcap/tiup/pkg/meta"
	"github.com/pingcap/tiup/pkg/utils"
)

// APIServerSpec represents the Master topology specification in topology.yaml
type APIServerSpec struct {
	Host            string                 `yaml:"host"`
	SSHPort         int                    `yaml:"ssh_port,omitempty" validate:"ssh_port:editable"`
	Port            int                    `yaml:"port,omitempty" default:"4116"`
	MetricsPort     int                    `yaml:"metrics_port,omitempty" default:"4123"`
	DeployDir       string                 `yaml:"deploy_dir,omitempty"`
	DataDir         string                 `yaml:"data_dir,omitempty"`
	LogDir          string                 `yaml:"log_dir,omitempty"`
	Config          map[string]interface{} `yaml:"config,omitempty" validate:"config:ignore"`
	Arch            string                 `yaml:"arch,omitempty"`
	OS              string                 `yaml:"os,omitempty"`
	LogLevel        string                 `yaml:"log_level,omitempty" default:"info" validate:"log_level:editable"`
	EnableHttps     string                 `yaml:"enable_https,omitempty" default:"true"`
	ResourceControl meta.ResourceControl   `yaml:"resource_control,omitempty" validate:"resource_control:editable"`
}

// Status queries current status of the instance
func (s *APIServerSpec) Status(tlsCfg *tls.Config, _ ...string) string {
	client := utils.NewHTTPClient(statusQueryTimeout, tlsCfg)

	path := "/system/check"
	protocol := "https"
	if s.EnableHttps == "false" {
		protocol = "http"
	}
	url := fmt.Sprintf("%s://%s:%d%s", protocol, s.Host, s.Port, path)

	// body doesn't have any status section needed
	body, err := client.Get(context.TODO(), url)
	if err != nil || body == nil {
		return "Down"
	}
	return "Up"
}

// Role returns the component role of the instance
func (s *APIServerSpec) Role() string {
	return ComponentTiEMAPIServer
}

// SSH returns the host and SSH port of the instance
func (s *APIServerSpec) SSH() (string, int) {
	return s.Host, s.SSHPort
}

// GetMainPort returns the main port of the instance
func (s *APIServerSpec) GetMainPort() int {
	return s.Port
}

// IsImported implements the instance interface, not needed for tiem
func (s *APIServerSpec) IsImported() bool {
	return false
}

// IgnoreMonitorAgent returns if the node does not have monitor agents available
func (s *APIServerSpec) IgnoreMonitorAgent() bool {
	return false
}

// APIServerComponent represents TiEM component.
type APIServerComponent struct{ Topology *Specification }

// Name implements Component interface.
func (c *APIServerComponent) Name() string {
	return ComponentTiEMAPIServer
}

// Role implements Component interface.
func (c *APIServerComponent) Role() string {
	return RoleTiEMAPI
}

// Instances implements Component interface.
func (c *APIServerComponent) Instances() []Instance {
	ins := make([]Instance, 0)
	for _, s := range c.Topology.APIServers {
		s := s
		ins = append(ins, &APIServerInstance{
			BaseInstance: BaseInstance{
				InstanceSpec: s,
				Name:         c.Name(),
				Host:         s.Host,
				Port:         s.Port,
				SSHP:         s.SSHPort,

				Ports: []int{
					s.Port,
					s.MetricsPort,
				},
				Dirs: []string{
					s.DeployDir,
					s.DataDir,
					s.LogDir,
				},
				StatusFn: s.Status,
				UptimeFn: func(tlsCfg *tls.Config) time.Duration {
					return spec.UptimeByHost(s.Host, s.Port, tlsCfg)
				},
			},
			topo: c.Topology,
		})
	}
	return ins
}

// APIServerInstance represent the TiEM instance
type APIServerInstance struct {
	BaseInstance
	topo *Specification
}

// InitConfig implement Instance interface
func (i *APIServerInstance) InitConfig(
	ctx context.Context,
	e ctxt.Executor,
	clusterName,
	clusterVersion,
	deployUser string,
	paths meta.DirPaths,
) error {
	if err := i.BaseInstance.InitConfig(ctx, e, i.topo.GlobalOptions, deployUser, paths); err != nil {
		return err
	}

	spec := i.InstanceSpec.(*APIServerSpec)

	cfg := config.NewAPIServerConfig().
		WithPrometheusAddress(i.topo.MonitorEndpoints()).
		WithGrafanaAddress(i.topo.GrafanaEndpoints()).
		WithAlertManagerAddress(i.topo.AlertManagerEndpoints()).
		WithKibanaAddress(i.topo.KibanaEndpoints()).
		WithJaegerAddress(i.topo.TracerEndpoints()).
		WithElasticsearchAddress(i.topo.ElasticSearchEndpoints())
	fp := filepath.Join(paths.Cache, fmt.Sprintf("openapi_%s_%d.yml", i.GetHost(), i.GetPort()))
	if err := cfg.ConfigToFile(fp); err != nil {
		return err
	}
	dst := filepath.Join(paths.Deploy, "conf", "env.yml")
	if err := e.Transfer(ctx, fp, dst, false, 0); err != nil {
		return err
	}

	scpt := scripts.NewTiEMAPIServerScript(
		i.GetHost(),
		paths.Deploy,
		paths.Data[0],
		paths.Log,
		spec.LogLevel,
	).
		WithPort(spec.Port).
		WithMetricsPort(spec.MetricsPort).
		WithRegistry(i.topo.RegistryEndpoints()).
		WithTracer(i.topo.TracerEndpoints()).
		WithElasticsearch(i.topo.ElasticSearchEndpoints()).
		WithEnableHttps(spec.EnableHttps)

	fp = filepath.Join(paths.Cache, fmt.Sprintf("run_openapi-server_%s_%d.sh", i.GetHost(), i.GetPort()))
	if err := scpt.ScriptToFile(fp); err != nil {
		return err
	}
	dst = filepath.Join(paths.Deploy, "scripts", "run_openapi-server.sh")
	if err := e.Transfer(ctx, fp, dst, false, 0); err != nil {
		return err
	}
	if _, _, err := e.Execute(ctx, "chmod +x "+dst, false); err != nil {
		return err
	}

	// TODO: support user specified certificates
	if _, _, err := e.Execute(ctx,
		fmt.Sprintf("cp -r %s/bin/cert %s/", paths.Deploy, paths.Deploy),
		false); err != nil {
		return err
	}

	// no config file needed
	return nil
}

// ScaleConfig deploy temporary config on scaling
func (i *APIServerInstance) ScaleConfig(
	ctx context.Context,
	e ctxt.Executor,
	topo Topology,
	clusterName,
	clusterVersion,
	deployUser string,
	paths meta.DirPaths,
) error {
	if err := i.InitConfig(ctx, e, clusterName, clusterVersion, deployUser, paths); err != nil {
		return err
	}

	spec := i.InstanceSpec.(*APIServerSpec)
	cfg := config.NewAPIServerConfig().
		WithPrometheusAddress(i.topo.MonitorEndpoints()).
		WithGrafanaAddress(i.topo.GrafanaEndpoints()).
		WithAlertManagerAddress(i.topo.AlertManagerEndpoints()).
		WithKibanaAddress(i.topo.KibanaEndpoints()).
		WithJaegerAddress(i.topo.TracerEndpoints()).
		WithElasticsearchAddress(i.topo.ElasticSearchEndpoints())
	fp := filepath.Join(paths.Cache, fmt.Sprintf("openapi_%s_%d.yml", i.GetHost(), i.GetPort()))
	if err := cfg.ConfigToFile(fp); err != nil {
		return err
	}
	dst := filepath.Join(paths.Deploy, "conf", "env.yml")
	if err := e.Transfer(ctx, fp, dst, false, 0); err != nil {
		return err
	}

	scpt := scripts.NewTiEMAPIServerScript(
		i.GetHost(),
		paths.Deploy,
		paths.Data[0],
		paths.Log,
		spec.LogLevel,
	).
		WithPort(spec.Port).
		WithMetricsPort(spec.MetricsPort).
		WithRegistry(i.topo.RegistryEndpoints()).
		WithTracer(i.topo.TracerEndpoints()).
		WithElasticsearch(i.topo.ElasticSearchEndpoints()).
		WithEnableHttps(spec.EnableHttps)

	fp = filepath.Join(paths.Cache, fmt.Sprintf("run_openapi-server_%s_%d.sh", i.GetHost(), i.GetPort()))
	log.Infof("script path: %s", fp)
	if err := scpt.ScriptToFile(fp); err != nil {
		return err
	}

	dst = filepath.Join(paths.Deploy, "scripts", "run_openapi-server.sh")
	if err := e.Transfer(ctx, fp, dst, false, 0); err != nil {
		return err
	}
	if _, _, err := e.Execute(ctx, "chmod +x "+dst, false); err != nil {
		return err
	}

	return nil
}
