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

package scripts

import (
	"bytes"
	"os"
	"path"
	"text/template"

	"github.com/pingcap-inc/tiem/tiup/embed"
)

// JaegerScript represent the data to generate Jaeger config
type JaegerScript struct {
	Host              string
	Port              int
	WebPort           int
	ZipkinCompactPort int
	JaegerCompactPort int
	JaegerBinaryPort  int
	JaegerHTTPPort    int
	CollectorHTTPPort int
	AdminHTTPPort     int
	CollectorGrpcPort int
	QueryGrpcPort     int
	DeployDir         string
	LogDir            string
}

// NewJaegerScript returns a JaegerScript with given arguments
func NewJaegerScript(ip, deployDir, logDir string) *JaegerScript {
	return &JaegerScript{
		Host:              ip,
		Port:              4114,
		WebPort:           4115,
		ZipkinCompactPort: 4116,
		JaegerCompactPort: 4117,
		JaegerBinaryPort:  4118,
		JaegerHTTPPort:    4119,
		CollectorHTTPPort: 4120,
		AdminHTTPPort:     4121,
		CollectorGrpcPort: 4122,
		QueryGrpcPort:     4123,
		DeployDir:         deployDir,
		LogDir:            logDir,
	}
}

// WithPort set Port field of JaegerScript
func (c *JaegerScript) WithPort(port int) *JaegerScript {
	c.Port = port
	return c
}

// WithWebPort set WebPort field of JaegerScript
func (c *JaegerScript) WithWebPort(port int) *JaegerScript {
	c.WebPort = port
	return c
}

// WithZipkinCompactPort set ZipkinCompactPort field of JaegerScript
func (c *JaegerScript) WithZipkinCompactPort(port int) *JaegerScript {
	c.ZipkinCompactPort = port
	return c
}

// WithJaegerCompactPort set JaegerCompactPort field of JaegerScript
func (c *JaegerScript) WithJaegerCompactPort(port int) *JaegerScript {
	c.JaegerCompactPort = port
	return c
}

// WithJaegerBinaryPort set JaegerBinaryPort field of JaegerScript
func (c *JaegerScript) WithJaegerBinaryPort(port int) *JaegerScript {
	c.JaegerBinaryPort = port
	return c
}

// WithJaegerHTTPPort set JaegerHTTPPort field of JaegerScript
func (c *JaegerScript) WithJaegerHTTPPort(port int) *JaegerScript {
	c.JaegerHTTPPort = port
	return c
}

// WithCollectorHTTPPort set CollectorHTTPPort field of JaegerScript
func (c *JaegerScript) WithCollectorHTTPPort(port int) *JaegerScript {
	c.CollectorHTTPPort = port
	return c
}

// WithAdminHTTPPort set AdminHTTPPort field of JaegerScript
func (c *JaegerScript) WithAdminHTTPPort(port int) *JaegerScript {
	c.AdminHTTPPort = port
	return c
}

// WithCollectorGrpcPort set CollectorGrpcPort field of JaegerScript
func (c *JaegerScript) WithCollectorGrpcPort(port int) *JaegerScript {
	c.CollectorGrpcPort = port
	return c
}

// WithQueryGrpcPort set QueryGrpcPort field of JaegerScript
func (c *JaegerScript) WithQueryGrpcPort(port int) *JaegerScript {
	c.QueryGrpcPort = port
	return c
}

// Script generate the config file data.
func (c *JaegerScript) Script() ([]byte, error) {
	fp := path.Join("templates", "scripts", "run_jaeger.sh.tpl")
	tpl, err := embed.ReadTemplate(fp)
	if err != nil {
		return nil, err
	}
	return c.ScriptWithTemplate(string(tpl))
}

// ScriptToFile write config content to specific path
func (c *JaegerScript) ScriptToFile(file string) error {
	config, err := c.Script()
	if err != nil {
		return err
	}
	return os.WriteFile(file, config, 0755)
}

// ScriptWithTemplate generate the Jaeger config content by tpl
func (c *JaegerScript) ScriptWithTemplate(tpl string) ([]byte, error) {
	tmpl, err := template.New("Jaeger").Parse(tpl)
	if err != nil {
		return nil, err
	}

	content := bytes.NewBufferString("")
	if err := tmpl.Execute(content, c); err != nil {
		return nil, err
	}

	return content.Bytes(), nil
}