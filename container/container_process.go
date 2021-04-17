package container

import (
  log "github.com/sirupsen/logrus"
  "os"
  "os/exec"
  "syscall"
)

// 外部创建容器进程的command
func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
  // 通过管道实现父子进程通信
  // 管道也是文件的一种，但是管道有一个固定大小的缓冲区，大小一般是4KB。
  // 当管道被写满时，写进程就会被阻塞，直到有读进程把管道的内容读出来。
  // 同样地，如果管道的内容是空的，那么读进程同样会被阻塞，一直等到有写进程向管道内写数据。
  rPipe, wPipe, err := os.Pipe()
  if err != nil {
    log.Errorf("New pipe error %v", err)
    return nil, nil
  }

  cmd := exec.Command("/proc/self/exe", "init") // 自己调用自己
  cmd.SysProcAttr = &syscall.SysProcAttr{
    Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
      syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET,
  }
  // 为子进程传入管道读端，每个进程默认有3个文件描述符，stdin stdout stderr 这个传入的文件描述符就是第四个
  cmd.ExtraFiles = []*os.File{rPipe}

  if tty {
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
  }
  return cmd, wPipe
}
