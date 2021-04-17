package main

import (
  "fmt"
  log "github.com/sirupsen/logrus"
  "github.com/urfave/cli/v2"
  "mydocker/cgroup/subsystems"
  "mydocker/container"
)

var runCommand = cli.Command{
  Name: "run",
  Usage: `Create a container with namespace and cgroups limit
			    mydocker run -ti [command]`,
  Flags: []cli.Flag{
    &cli.BoolFlag{
      Name:  "ti",
      Usage: "enable tty",
    },
    &cli.StringFlag{
      Name:  "m",
      Usage: "memory limit in bytes",
    },
    &cli.StringFlag{
      Name:  "cpushare",
      Usage: "cpushare limit",
    },
    &cli.StringFlag{
      Name:  "cpuset",
      Usage: "cpuset limit",
    },
  },
  Action: func(context *cli.Context) error {
    if context.Args().Len() < 1 {
      return fmt.Errorf("Missing container command")
    }

    var cmdArray = context.Args().Slice()
    tty := context.Bool("ti")
    resConf := &subsystems.ResourceConfig{
      MemoryLimit: context.String("m"),
      CpuShare:    context.String("cpushare"),
      CpuSet:      context.String("cpuset"),
    }

    Run(tty, cmdArray, resConf)
    return nil
  },
}

var initCommand = cli.Command{
  Name:  "init",
  Usage: "Init container process run user's process in container. Do not call it outside",
  Action: func(context *cli.Context) error {
    log.Infof("init come on")
    cmd := context.Args().Get(0)
    log.Infof("command %s", cmd)
    err := container.RunContainerInitProcess()
    return err
  },
}
