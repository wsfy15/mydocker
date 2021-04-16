package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "os/exec"
  "path"
  "strconv"
  "syscall"
)

const cgroupMemoryHierachyMount = "/sys/fs/cgroup/memory"

func useNamespaceAndCgroup() {
  if os.Args[0] == "/proc/self/exe" {
    // 容器进程
    fmt.Println("current pid", syscall.Getpid());
    cmd := exec.Command("sh", "-c", "stress --vm-bytes 200m --vm-keep -m 1")
    cmd.SysProcAttr =&syscall.SysProcAttr{}
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout

    if err := cmd.Run(); err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
  }

  // /proc/self/exe 表示当前进程执行的可执行文件
  cmd := exec.Command("/proc/self/exe")
  cmd.SysProcAttr = &syscall.SysProcAttr{
    Cloneflags: syscall.CLONE_NEWUTS | // hostname domainname
      syscall.CLONE_NEWIPC | // System V IPC、POSIX message queue
      syscall.CLONE_NEWPID |
      syscall.CLONE_NEWNS | // mount   mount -t proc proc /proc
      syscall.CLONE_NEWUSER |
      syscall.CLONE_NEWNET,
  }
  //cmd.SysProcAttr.Credential = &syscall.Credential{
  // Uid: uint32(1),
  // Gid: uint32(1),
  //}
  cmd.Stderr = os.Stderr
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout

  if err := cmd.Start(); err != nil {
    log.Fatal(err)
  } else {
    // fork 的进程在外部namespace 的pid
    fmt.Println(cmd.Process.Pid)

    os.Mkdir(path.Join(cgroupMemoryHierachyMount, "testmemorylimit"), 0755)

    ioutil.WriteFile(path.Join(cgroupMemoryHierachyMount, "testmemorylimit", "tasks"),
      []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
    ioutil.WriteFile(path.Join(cgroupMemoryHierachyMount, "testmemorylimit", "memory.limit_in_bytes"),
      []byte("100m"), 0644)
  }
  cmd.Process.Wait()
}
