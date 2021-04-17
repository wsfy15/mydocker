package subsystems

import (
  "bufio"
  "fmt"
  log "github.com/sirupsen/logrus"
  "os"
  "path"
  "strings"
)

func FindCgroupMountpoint(subsystem string) string {
  f, err := os.Open("/proc/self/mountinfo")
  if err != nil {
    return ""
  }
  defer f.Close()

  scanner := bufio.NewScanner(f)
  for scanner.Scan() {
    text := scanner.Text()
    fields := strings.Split(text, " ")
    for _, opt := range strings.Split(fields[len(fields) - 1], ",") {
      if opt == subsystem {
        return fields[4]
      }
    }
  }

  return ""
}

func GetCgroupPath(subsystem, cgroupPath string, autoCreate bool) (string, error) {
  var err error
  cgroupRoot := FindCgroupMountpoint(subsystem)
  if _, err = os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil ||
      (autoCreate && os.IsNotExist(err)) {
    if os.IsNotExist(err) {
      if err = os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err != nil {
        return "", fmt.Errorf("error create cgroup %v", err)
      }
    }
    return path.Join(cgroupRoot, cgroupPath), nil
  }

  log.Errorf("GetCgroupPath %v", err)
  return "", fmt.Errorf("error path cgroup %v", err)
}

