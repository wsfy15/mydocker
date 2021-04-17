package cgroup

import (
  "github.com/sirupsen/logrus"
  "mydocker/cgroup/subsystems"
)

type CgroupManager struct {
  // 相对于各root cgroup 目录的路径
  Path string
  Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
  return &CgroupManager{Path: path}
}

func (m *CgroupManager) Apply(pid int) error {
  for _, subSysIns := range subsystems.SubsystemsIns {
    subSysIns.Apply(m.Path, pid)
  }
  return nil
}

func (m *CgroupManager) Set(res *subsystems.ResourceConfig) error  {
  for _, subSysIns := range subsystems.SubsystemsIns {
    subSysIns.Set(m.Path, res)
  }
  return nil
}

func (m *CgroupManager) Destroy() error {
  for _, subSysIns := range subsystems.SubsystemsIns {
    if err := subSysIns.Remove(m.Path); err != nil {
      logrus.Warnf("remove cgroup fail %v", err)
    }
  }
  return nil
}
