package container

import (
  "github.com/sirupsen/logrus"
  "os"
  "syscall"
)

// 容器进程
func RunContainerInitProcess(command string, args []string) error {
  // MS_NOEXEC 在本文件系统中不允许运行其他程序
  // MS_NOSUID在本系统中运行程序的时候，不允许set-user-ID或set-group-ID
  // MS_NODEV这个参数是自从Linux 2.4以来，所有mount的系统都会默认设定的参数
  defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
  syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
  argv := []string{command}

  // 当前init进程是容器内PID为1的进程，但是在docker中，用户要运行的命令 的进程PID为1
  // 通过execve系统调用，可以执行一个新命令，同时会覆盖当前进程的镜像、数据和堆栈等信息，包括PID
  // 即替换掉init进程
  if err := syscall.Exec(command, argv, os.Environ()); err != nil {
    logrus.Errorf(err.Error())
  }

  return nil
}
