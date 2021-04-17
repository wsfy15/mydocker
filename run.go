package main

import (
  log "github.com/sirupsen/logrus"
  "mydocker/cgroup"
  "mydocker/cgroup/subsystems"
  "mydocker/container"
  "os"
  "strings"
)

func Run(tty bool, command []string, res *subsystems.ResourceConfig) {
  parent, writePipe := container.NewParentProcess(tty)
  if err := parent.Start(); err != nil {
    log.Fatal(err)
  }

  log.Printf("container pid is %v", parent.Process.Pid)

  cgroupManager := cgroup.NewCgroupManager("mydocker-cgroup")
  cgroupManager.Set(res)
  cgroupManager.Apply(parent.Process.Pid)
  sendInitCommand(command, writePipe)
  parent.Wait()
}

func sendInitCommand(comArray []string, writePipe *os.File) {
  command := strings.Join(comArray, " ")
  log.Infof("command all is %s", command)
  writePipe.WriteString(command)
  writePipe.Close()
}
