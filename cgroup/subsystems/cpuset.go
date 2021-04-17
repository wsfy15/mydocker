package subsystems

import (
  "fmt"
  "io/ioutil"
  "os"
  "path"
  "strconv"
)

type CpusetSubSystem struct {
  
}

func (c *CpusetSubSystem) Name() string {
  return "cpuset"
}

func (c *CpusetSubSystem) Set(cgouppath string, res *ResourceConfig) (err error) {
  if res.CpuSet == "" {
    return
  }

  var subsysCgroupPath string
  if subsysCgroupPath, err = GetCgroupPath(c.Name(), cgouppath, true); err == nil {
    if err = ioutil.WriteFile(path.Join(subsysCgroupPath, "cpuset.cpus"),
      []byte(res.CpuSet), 0644); err != nil {
      return fmt.Errorf("set cgroup cpu share fail %v", err)
    }
  }

  return
}

func (c *CpusetSubSystem) Apply(cgouppath string, pid int) (err error) {
  var subsysCgroupPath string
  if subsysCgroupPath, err = GetCgroupPath(c.Name(), cgouppath, false); err == nil {
    if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
      return fmt.Errorf("set cgroup proc fail %v", err)
    }
  }

  return
}

func (c *CpusetSubSystem) Remove(cgouppath string) error {
  if subsysCgroupPath, err := GetCgroupPath(c.Name(), cgouppath, false); err == nil {
    return os.Remove(subsysCgroupPath)
  } else {
    return err
  }
}

