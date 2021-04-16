package container

import (
  "os"
  "os/exec"
  "syscall"
)

// 外部创建容器进程的command
func NewParentProcess(tty bool, command string) *exec.Cmd {
  args := []string{"init", command}
  cmd := exec.Command("/proc/self/exe", args...) // 自己调用自己
  cmd.SysProcAttr = &syscall.SysProcAttr{
    Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
      syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET,
  }

  if tty {
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
  }
  return cmd
}
