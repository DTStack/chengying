/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"text/template"
	"time"

	"easyagent/internal/sidecar/base"
	"github.com/containerd/cgroups"
)

type Cmd struct {
	*exec.Cmd

	cg  *Cgroup
	ctx context.Context

	pgid         int
	isSeniorKill bool
	runUser      string
}

type Cgroup struct {
	cgroups.Cgroup
	initStub string
}

func NewCgroup(cg cgroups.Cgroup, sid string) *Cgroup {
	paths := make([]string, 0)
	for _, s := range cg.Subsystems() {
		paths = append(paths, filepath.Join(s.(path).Path(sid), "cgroup.procs"))
	}

	var buf bytes.Buffer
	stubShell := `{{range $_, $path := .Paths}}echo $$ > {{$path}} && {{end}}exec "$@"`
	template.Must(template.New("test").Parse(stubShell)).Execute(&buf, struct{ Paths []string }{paths})

	return &Cgroup{cg, buf.String()}
}

type path interface {
	Path(path string) string
}

func (c *Cgroup) GetInitStub() string {
	return c.initStub
}

func CommandContext(ctx context.Context, user string, isSeniorKill bool, cg *Cgroup, name string, arg ...string) *Cmd {
	if ctx == nil {
		panic("nil Context")
	}
	var cmd *exec.Cmd
	if isSeniorKill && cg != nil && len(user) == 0 {
		stubArg := append([]string{"-c", cg.GetInitStub(), "", name}, arg...)
		cmd = exec.Command("sh", stubArg...)
	} else {
		cmd = exec.Command(name, arg...)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pdeathsig: syscall.SIGKILL}
	return &Cmd{Cmd: cmd, cg: cg, ctx: ctx, pgid: -1, isSeniorKill: isSeniorKill, runUser: user}
}

func (c *Cmd) Run() error {
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

func (c *Cmd) Start() error {
	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	default:
	}
	if err := c.Cmd.Start(); err != nil {
		return err
	}
	pgid, err := syscall.Getpgid(c.Process.Pid)
	if err != nil {
		base.Errorf("Getpgid error: %v", err)
		return nil
	}
	c.pgid = pgid
	return nil
}

func (c *Cmd) Wait() error {
	var err error
	waitDone := make(chan struct{})

	go func() {
		err = c.Cmd.Wait()
		close(waitDone)
	}()

	select {
	case <-waitDone:
		if c.cg != nil && c.isSeniorKill && len(c.runUser) == 0 {
			processes, _ := c.cg.Processes(cgroups.NetCLS, true)
			base.Debugf("kill processes: %v", processes)
			for _, proc := range processes {
				syscall.Kill(proc.Pid, syscall.SIGKILL)
			}
			for _, proc := range processes {
				for syscall.Kill(proc.Pid, 0) == nil {
					base.Debugf("waiting kill processes: %v", proc.Pid)
					time.Sleep(10 * time.Millisecond)
				}
			}
		} else if c.pgid != -1 {
			base.Debugf("kill pgid: %v", c.pgid)
			syscall.Kill(-c.pgid, syscall.SIGKILL)
			for syscall.Kill(-c.pgid, 0) == nil {
				base.Debugf("waiting kill pgid %v", c.pgid)
				time.Sleep(10 * time.Millisecond)
			}
		}
		return err
	case <-c.ctx.Done():
		// give the chance to cleanup children resources
		if c.cg != nil && c.isSeniorKill && len(c.runUser) == 0 {
			processes, _ := c.cg.Processes(cgroups.NetCLS, true)
			base.Debugf("kill processes: %v", processes)
			for _, proc := range processes {
				syscall.Kill(proc.Pid, syscall.SIGTERM)
			}

			time.Sleep(2 * time.Second)

			processes, _ = c.cg.Processes(cgroups.NetCLS, true)
			for _, proc := range processes {
				syscall.Kill(proc.Pid, syscall.SIGKILL)
			}
			for _, proc := range processes {
				for syscall.Kill(proc.Pid, 0) == nil {
					base.Debugf("waiting kill processes: %v", proc.Pid)
					time.Sleep(10 * time.Millisecond)
				}
			}
		} else if c.pgid != -1 {
			base.Debugf("kill pgid: %v", c.pgid)
			if c.isSeniorKill {
				syscall.Kill(-c.pgid, syscall.SIGTERM)
				time.Sleep(2 * time.Second)
			}
			syscall.Kill(-c.pgid, syscall.SIGKILL)
			for syscall.Kill(-c.pgid, 0) == nil {
				base.Debugf("waiting kill pgid %v", c.pgid)
				time.Sleep(10 * time.Millisecond)
			}
		} else {
			base.Debugf("Kill child process(%v) only!", c.Process.Pid)

			if c.isSeniorKill {
				c.Process.Signal(syscall.SIGTERM)
				time.Sleep(2 * time.Second)
			}
			c.Process.Signal(syscall.SIGKILL)
			for c.Process.Signal(syscall.Signal(0)) == nil {
				base.Debugf("waiting kill child process(%v) only!", c.Process.Pid)
				time.Sleep(10 * time.Millisecond)
			}
		}
		return c.ctx.Err()
	}
}

func CreateTempScript(content string, prefix string) (path string, err error) {
	f, err := ioutil.TempFile("", prefix)
	if err != nil {
		return "", err
	}

	if _, err = f.WriteString(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}
	if err = os.Chmod(f.Name(), 0500); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}

	f.Close()
	return f.Name(), nil
}
