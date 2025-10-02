// SPDX-FileCopyrightText: 2024 Intel Corporation
// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/server"
	"github.com/urfave/cli/v3"
)

var SSM = server.SsmServer

func main() {
	logger.AppLog.Infoln("SSM is starting wait some seconds while the configs are loading")
	app := &cli.Command{}
	app.Name = "ssm"
	logger.AppLog.Infoln(app.Name)
	app.Usage = "Access & Mobility Management function"
	app.UsageText = "ssm -cfg <ssm_config_file.conf>"
	app.Action = action
	app.Flags = SSM.GetCliCmd()
	if err := app.Run(context.Background(), os.Args); err != nil {
		logger.AppLog.Fatalf("SSM run error: %v", err)
	}
}

func action(ctx context.Context, c *cli.Command) error {
	if err := SSM.Initialize(c); err != nil {
		logger.CfgLog.Errorf("%+v", err)
		return fmt.Errorf("failed to initialize")
	}

	SSM.Start()

	return nil
}
