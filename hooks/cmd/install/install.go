// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2021 Canonical Ltd
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 * SPDX-License-Identifier: Apache-2.0'
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	local "github.com/canonical/edgex-device-camera-go/hooks"
	hooks "github.com/canonical/edgex-snap-hooks"
)

var cli *hooks.CtlCli = hooks.NewSnapCtl()

// installProfiles copies the profile configuration.toml files from $SNAP to $SNAP_DATA.
func installConfig() error {
	var err error

	path := "/config/device-camera-go/res/configuration.toml"
	destFile := hooks.SnapData + path
	srcFile := hooks.Snap + path

	if err = os.MkdirAll(filepath.Dir(destFile), 0755); err != nil {
		return err
	}

	if err = hooks.CopyFile(srcFile, destFile); err != nil {
		return err
	}

	return nil
}

// TODO: merge into the above function...
func installDevProfiles() error {
	var err error

	profs := [...]string{"axis", "bosch"}

	for _, v := range profs {
		path := fmt.Sprintf("/config/device-camera-go/res/camera-%s.yaml", v)
		destFile := hooks.SnapData + path
		srcFile := hooks.Snap + path

		if err := os.MkdirAll(filepath.Dir(destFile), 0755); err != nil {
			return err
		}

		if err = hooks.CopyFile(srcFile, destFile); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	var debug = false
	var err error

	status, err := cli.Config("debug")
	if err != nil {
		fmt.Println(fmt.Sprintf("edgex-device-camera:install: can't read value of 'debug': %v", err))
		os.Exit(1)
	}
	if status == "true" {
		debug = true
	}

	if err = hooks.Init(debug, "edgex-device-camera"); err != nil {
		fmt.Println(fmt.Sprintf("edgex-device-camera::install: initialization failure: %v", err))
		os.Exit(1)

	}

	err = installConfig()
	if err != nil {
		hooks.Error(fmt.Sprintf("edgex-device-camera:install: %v", err))
		os.Exit(1)
	}

	err = installDevProfiles()
	if err != nil {
		hooks.Error(fmt.Sprintf("edgex-device-camera:install: %v", err))
		os.Exit(1)
	}

	cli := hooks.NewSnapCtl()
	svc := fmt.Sprintf("%s.device-camera-go", hooks.SnapInst)

	autostart, err := cli.Config(hooks.AutostartConfig)
	if err != nil {
		hooks.Error(fmt.Sprintf("Reading config 'autostart' failed: %v", err))
		os.Exit(1)
	}

	// TODO: move profile config before autostart, if profile=default, or
	// no configuration file exists for the profile, then ignore autostart

	switch strings.ToLower(autostart) {
	case "true":
	case "yes":
		break
	case "":
	case "no":
		// disable this service initially because it requires configuration
		// with a device profile that will be specific to each installation
		err = cli.Stop(svc, true)
		if err != nil {
			hooks.Error(fmt.Sprintf("Can't stop service - %v", err))
			os.Exit(1)
		}
	default:
		hooks.Error(fmt.Sprintf("Invalid value for 'autostart' : %s", autostart))
		os.Exit(1)
	}

	envJSON, err := cli.Config(hooks.EnvConfig)
	if err != nil {
		hooks.Error(fmt.Sprintf("Reading config 'env' failed: %v", err))
		os.Exit(1)
	}

	if envJSON != "" {
		hooks.Debug(fmt.Sprintf("edgex-device-camera:configure: envJSON: %s", envJSON))
		err = hooks.HandleEdgeXConfig("device-camera-go", envJSON, local.ConfToEnv)
		if err != nil {
			hooks.Error(fmt.Sprintf("HandleEdgeXConfig failed: %v", err))
			os.Exit(1)
		}
	}
}
