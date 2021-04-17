package subsystems

import (
  "fmt"
  log "github.com/sirupsen/logrus"
  "io/ioutil"
  "os"
  "path"
  "strconv"
)

type MemorySubSystem struct {

}

func (m *MemorySubSystem) Name() string {
  return "memory"
}

func (m *MemorySubSystem) Set(cgouppath string, res *ResourceConfig) (err error) {
  if res.MemoryLimit == "" {
    return
  }

  var subsysCgroupPath string
  if subsysCgroupPath, err = GetCgroupPath(m.Name(), cgouppath, true); err == nil {
    log.Printf("memory limit is %v byte", res.MemoryLimit)
    if err = ioutil.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"),
      []byte(res.MemoryLimit), 0644); err != nil {
      return fmt.Errorf("set cgroup memory fail %v", err)
    }
  } else {
    log.Errorf("memory error %v", err)
  }

  return
}

func (m *MemorySubSystem) Apply(cgouppath string, pid int) (err error) {
  var subsysCgroupPath string
  if subsysCgroupPath, err = GetCgroupPath(m.Name(), cgouppath, false); err == nil {
    if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
      return fmt.Errorf("set cgroup proc fail %v", err)
    }
  }

  return
}

func (m *MemorySubSystem) Remove(cgouppath string) error {
  if subsysCgroupPath, err := GetCgroupPath(m.Name(), cgouppath, false); err == nil {
    return os.Remove(subsysCgroupPath)
  } else {
    return err
  }
}

