package subsystems

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string // CPU 时间片权重
	CpuSet      string // CPU 核心数
}

type Subsystem interface {
  // eg, cpu or memory
  Name() string
  // 设置某个cgroup的资源限制
  Set(path string, res *ResourceConfig) error
  // 将进程添加到cgroup
  Apply(path string, pid int) error
  // 删除cgroup
  Remove(path string) error
}

var (
  SubsystemsIns = []Subsystem{
    &MemorySubSystem{},
    &CpusetSubSystem{},
    &CpuSubSystem{},
  }
)
