package parsing

import (
	"github.com/davenfroberg/wireyak/metrics"
	"github.com/shirou/gopsutil/v3/process"
)

type ProcessParser struct {
	Metrics *metrics.Metrics
}

func NewProcessParser(m *metrics.Metrics) *ProcessParser {
	return &ProcessParser{
		Metrics: m,
	}
}

func (p *ProcessParser) ParseProcesses(procs []*process.Process) {
	p.Metrics.ProcessCount.Set(float64(len(procs)))
	for _, proc := range procs {
		p.parseProcess(proc)
	}
}

func (p *ProcessParser) parseProcess(proc *process.Process) {
	// name, _ := proc.Name()
	// cpu, _ := proc.CPUPercent()
	// mem, _ := proc.MemoryPercent()
	// par, err := proc.Parent()
	// parent := "None"

	// if err == nil {
	// 	parent, _ = par.Name()
	// 	if parent == "" {
	// 		parent = "None"
	// 	}
	// }

	// if cpu > 0 || mem > 0 {
	// 	fmt.Printf("PID: %d | Name: %s | CPU: %.2f%% | RAM: %.2f%%| Spawned By: %s\n", proc.Pid, name, cpu, mem, parent)
	// }
}
