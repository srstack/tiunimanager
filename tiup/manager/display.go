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

package manager

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	operator "github.com/pingcap-inc/tiem/tiup/operation"
	"github.com/pingcap-inc/tiem/tiup/spec"
	perrs "github.com/pingcap/errors"
	"github.com/pingcap/tiup/pkg/cluster/clusterutil"
	"github.com/pingcap/tiup/pkg/cluster/ctxt"
	"github.com/pingcap/tiup/pkg/cluster/executor"
	"github.com/pingcap/tiup/pkg/logger/log"
	"github.com/pingcap/tiup/pkg/meta"
	"github.com/pingcap/tiup/pkg/set"
	"github.com/pingcap/tiup/pkg/tui"
	"github.com/pingcap/tiup/pkg/utils"
)

// InstInfo represents an instance info
type InstInfo struct {
	ID        string `json:"id"`
	Role      string `json:"role"`
	Host      string `json:"host"`
	Ports     string `json:"ports"`
	OsArch    string `json:"os_arch"`
	Status    string `json:"status"`
	Since     string `json:"since"`
	DataDir   string `json:"data_dir"`
	DeployDir string `json:"deploy_dir"`

	ComponentName string
	Port          int
}

// ClusterMetaInfo hold the structure for the JSON output of the dashboard info
type ClusterMetaInfo struct {
	ClusterType    string `json:"cluster_type"`
	ClusterName    string `json:"cluster_name"`
	ClusterVersion string `json:"cluster_version"`
	DeployUser     string `json:"deploy_user"`
	SSHType        string `json:"ssh_type"`
}

// JSONOutput holds the structure for the JSON output of `tiup cluster display --json`
type JSONOutput struct {
	ClusterMetaInfo ClusterMetaInfo `json:"cluster_meta"`
	InstanceInfos   []InstInfo      `json:"instances"`
}

// Display cluster meta and topology.
func (m *Manager) Display(name string, opt operator.Options) error {
	if err := clusterutil.ValidateClusterNameOrError(name); err != nil {
		return err
	}

	clusterInstInfos, err := m.GetClusterTopology(name, opt)
	if err != nil {
		return err
	}

	metadata, _ := m.meta(name)
	topo := metadata.GetTopology()
	base := metadata.GetBaseMeta()
	cyan := color.New(color.FgCyan, color.Bold)
	// display cluster meta
	var j *JSONOutput
	if opt.JSON {
		j = &JSONOutput{
			ClusterMetaInfo: ClusterMetaInfo{
				m.sysName,
				name,
				base.Version,
				topo.BaseTopo().GlobalOptions.User,
				string(topo.BaseTopo().GlobalOptions.SSHType),
			},
			InstanceInfos: clusterInstInfos,
		}
	} else {
		fmt.Printf("Cluster type:       %s\n", cyan.Sprint(m.sysName))
		fmt.Printf("Cluster name:       %s\n", cyan.Sprint(name))
		fmt.Printf("Cluster version:    %s\n", cyan.Sprint(base.Version))
		fmt.Printf("Deploy user:        %s\n", cyan.Sprint(topo.BaseTopo().GlobalOptions.User))
		fmt.Printf("SSH type:           %s\n", cyan.Sprint(topo.BaseTopo().GlobalOptions.SSHType))
		fmt.Printf("WebServer URL:      %s\n", cyan.Sprint(formatWebServerURL(topo.BaseTopo().WebServers)))
	}

	// display topology
	var clusterTable [][]string
	if opt.ShowUptime {
		clusterTable = append(clusterTable, []string{"ID", "Role", "Host", "Ports", "OS/Arch", "Status", "Since", "Data Dir", "Deploy Dir"})
	} else {
		clusterTable = append(clusterTable, []string{"ID", "Role", "Host", "Ports", "OS/Arch", "Status", "Data Dir", "Deploy Dir"})
	}

	masterActive := make([]string, 0)
	for _, v := range clusterInstInfos {
		row := []string{
			color.CyanString(v.ID),
			v.Role,
			v.Host,
			v.Ports,
			v.OsArch,
			formatInstanceStatus(v.Status),
		}
		if opt.ShowUptime {
			row = append(row, v.Since)
		}
		row = append(row, v.DataDir, v.DeployDir)
		clusterTable = append(clusterTable, row)

		if strings.HasPrefix(v.Status, "Up") || strings.HasPrefix(v.Status, "Healthy") {
			instAddr := fmt.Sprintf("%s:%d", v.Host, v.Port)
			masterActive = append(masterActive, instAddr)
		}
	}

	if opt.JSON {
		d, err := json.MarshalIndent(j, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(d))
		return nil
	}

	tui.PrintTable(clusterTable, true)
	fmt.Printf("Total nodes: %d\n", len(clusterTable)-1)

	return nil
}

// GetClusterTopology get the topology of the cluster.
func (m *Manager) GetClusterTopology(name string, opt operator.Options) ([]InstInfo, error) {
	ctx := ctxt.New(context.Background(), opt.Concurrency)
	metadata, err := m.meta(name)
	if err != nil && !errors.Is(perrs.Cause(err), meta.ErrValidate) {
		return nil, err
	}

	topo := metadata.GetTopology()
	base := metadata.GetBaseMeta()

	err = SetSSHKeySet(ctx, m.specManager.Path(name, "ssh", "id_rsa"), m.specManager.Path(name, "ssh", "id_rsa.pub"))
	if err != nil {
		return nil, err
	}

	err = SetClusterSSH(ctx, topo, base.User, opt.SSHTimeout, opt.SSHType, topo.BaseTopo().GlobalOptions.SSHType)
	if err != nil {
		return nil, err
	}

	filterRoles := set.NewStringSet(opt.Roles...)
	filterNodes := set.NewStringSet(opt.Nodes...)
	masterList := topo.BaseTopo().MasterList
	tlsCfg := &tls.Config{} // not implemented for tiem

	masterActive := make([]string, 0)
	masterStatus := make(map[string]string)

	topo.IterInstance(func(ins spec.Instance) {
		status := ins.Status(tlsCfg, masterList...)
		if strings.HasPrefix(status, "Up") || strings.HasPrefix(status, "Healthy") {
			instAddr := fmt.Sprintf("%s:%d", ins.GetHost(), ins.GetPort())
			masterActive = append(masterActive, instAddr)
		}
		masterStatus[ins.ID()] = status
	})

	clusterInstInfos := []InstInfo{}

	topo.IterInstance(func(ins spec.Instance) {
		// apply role filter
		if len(filterRoles) > 0 && !filterRoles.Exist(ins.Role()) {
			return
		}
		// apply node filter
		if len(filterNodes) > 0 && !filterNodes.Exist(ins.ID()) {
			return
		}

		dataDir := "-"
		insDirs := ins.UsedDirs()
		deployDir := insDirs[0]
		if len(insDirs) > 1 {
			dataDir = insDirs[1]
		}

		var status string
		switch ins.ComponentName() {
		default:
			status = ins.Status(tlsCfg, masterActive...)
		}

		since := "-"
		if opt.ShowUptime {
			since = formatInstanceSince(ins.Uptime(tlsCfg))
		}

		// Query the service status and uptime
		if status == "-" || (opt.ShowUptime && since == "-") {
			e, found := ctxt.GetInner(ctx).GetExecutor(ins.GetHost())
			if found {
				active, _ := operator.GetServiceStatus(ctx, e, ins.ServiceName())
				if status == "-" {
					if parts := strings.Split(strings.TrimSpace(active), " "); len(parts) > 2 {
						if parts[1] == "active" {
							status = "Up"
						} else {
							status = parts[1]
						}
					}
				}
				if opt.ShowUptime && since == "-" {
					since = formatInstanceSince(parseSystemctlSince(active))
				}
			}
		}

		// check if the role is patched
		roleName := ins.Role()
		if ins.IsPatched() {
			roleName += " (patched)"
		}
		clusterInstInfos = append(clusterInstInfos, InstInfo{
			ID:            ins.ID(),
			Role:          roleName,
			Host:          ins.GetHost(),
			Ports:         utils.JoinInt(ins.UsedPorts(), "/"),
			OsArch:        tui.OsArch(ins.OS(), ins.Arch()),
			Status:        status,
			DataDir:       dataDir,
			DeployDir:     deployDir,
			ComponentName: ins.ComponentName(),
			Port:          ins.GetPort(),
			Since:         since,
		})
	})

	// Sort by role,host,ports
	sort.Slice(clusterInstInfos, func(i, j int) bool {
		lhs, rhs := clusterInstInfos[i], clusterInstInfos[j]
		if lhs.Role != rhs.Role {
			return lhs.Role < rhs.Role
		}
		if lhs.Host != rhs.Host {
			return lhs.Host < rhs.Host
		}
		return lhs.Ports < rhs.Ports
	})

	return clusterInstInfos, nil
}

func formatInstanceStatus(status string) string {
	lowercaseStatus := strings.ToLower(status)

	startsWith := func(prefixs ...string) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(lowercaseStatus, prefix) {
				return true
			}
		}
		return false
	}

	switch {
	case startsWith("up|l", "healthy|l"): // up|l, up|l|ui, healthy|l
		return color.HiGreenString(status)
	case startsWith("up", "healthy", "free"):
		return color.GreenString(status)
	case startsWith("down", "err", "inactive"): // down, down|ui
		return color.RedString(status)
	case startsWith("tombstone", "disconnected", "n/a"), strings.Contains(status, "offline"):
		return color.YellowString(status)
	default:
		return status
	}
}

func formatWebServerURL(webServerSpec []*spec.WebServerSpec) string {
	urls := make([]string, 0)
	for _, spec := range webServerSpec {
		urls = append(urls, "http://"+spec.Host+":"+strconv.Itoa(spec.Port))
	}
	return strings.Join(urls, ",")
}

func formatInstanceSince(uptime time.Duration) string {
	if uptime == 0 {
		return "-"
	}

	d := int64(uptime.Hours() / 24)
	h := int64(math.Mod(uptime.Hours(), 24))
	m := int64(math.Mod(uptime.Minutes(), 60))
	s := int64(math.Mod(uptime.Seconds(), 60))

	chunks := []struct {
		unit  string
		value int64
	}{
		{"d", d},
		{"h", h},
		{"m", m},
		{"s", s},
	}

	parts := []string{}

	for _, chunk := range chunks {
		switch chunk.value {
		case 0:
			continue
		default:
			parts = append(parts, fmt.Sprintf("%d%s", chunk.value, chunk.unit))
		}
	}

	return strings.Join(parts, "")
}

// `systemctl status xxx.service` returns as below
// Active: active (running) since Sat 2021-03-27 10:51:11 CST; 41min ago
func parseSystemctlSince(str string) (dur time.Duration) {
	// if service is not found or other error, don't need to parse it
	if str == "" {
		return 0
	}
	defer func() {
		if dur == 0 {
			log.Warnf("failed to parse systemctl since '%s'", str)
		}
	}()
	parts := strings.Split(str, ";")
	if len(parts) != 2 {
		return
	}
	parts = strings.Split(parts[0], " ")
	if len(parts) < 3 {
		return
	}

	dateStr := strings.Join(parts[len(parts)-3:], " ")

	tm, err := time.Parse("2006-01-02 15:04:05 MST", dateStr)
	if err != nil {
		return
	}

	return time.Since(tm)
}

// SetSSHKeySet set ssh key set.
func SetSSHKeySet(ctx context.Context, privateKeyPath string, publicKeyPath string) error {
	ctxt.GetInner(ctx).PrivateKeyPath = privateKeyPath
	ctxt.GetInner(ctx).PublicKeyPath = publicKeyPath
	return nil
}

// SetClusterSSH set cluster user ssh executor in context.
func SetClusterSSH(ctx context.Context, topo spec.Topology, deployUser string, sshTimeout uint64, sshType, defaultSSHType executor.SSHType) error {
	if sshType == "" {
		sshType = defaultSSHType
	}
	if len(ctxt.GetInner(ctx).PrivateKeyPath) == 0 {
		return perrs.Errorf("context has no PrivateKeyPath")
	}

	for _, com := range topo.ComponentsByStartOrder() {
		for _, in := range com.Instances() {
			cf := executor.SSHConfig{
				Host:    in.GetHost(),
				Port:    in.GetSSHPort(),
				KeyFile: ctxt.GetInner(ctx).PrivateKeyPath,
				User:    deployUser,
				Timeout: time.Second * time.Duration(sshTimeout),
			}

			e, err := executor.New(sshType, false, cf)
			if err != nil {
				return err
			}
			ctxt.GetInner(ctx).SetExecutor(in.GetHost(), e)
		}
	}

	return nil
}