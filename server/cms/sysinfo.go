package cms

import (
	"bytes"
	"fmt"
	"herefriend/lib"
	"html/template"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type process struct {
	pid int
	cpu float64
}

func getCPUUsage() float64 {
	var out bytes.Buffer

	var processes []*process
	cmd := exec.Command("ps", "aux")

	cmd.Stdout = &out
	cmd.Run()

	for {
		line, err := out.ReadString('\n')
		if err != nil {
			break
		}

		tokens := strings.Split(line, " ")
		var ft []string
		for _, t := range tokens {
			if t != "" && t != "\t" {
				ft = append(ft, t)
			}
		}
		pid, err := strconv.Atoi(ft[1])
		if err != nil {
			continue
		}
		cpu, err := strconv.ParseFloat(ft[2], 64)
		if err == nil {
			processes = append(processes, &process{pid, cpu})
		}
	}

	usage := float64(0)
	for _, p := range processes {
		usage += p.cpu
	}

	if 100.0001 < usage {
		usage = 100.0
	}

	return usage
}

// SystemInfo 获取最新的系统信息
func SystemInfo(c *gin.Context) {
	meminfo, _ := mem.VirtualMemory()
	cpuinfo, _ := cpu.CPUInfo()
	diskinfo, _ := disk.DiskUsage("/")
	hostinfo, _ := host.HostInfo()

	info := cmsSystemSummary{
		OSDescribe:  fmt.Sprintf("%s %s", hostinfo.OS, hostinfo.PlatformVersion),
		CPUDescribe: fmt.Sprintf("%s %d Cores", cpuinfo[0].ModelName, cpuinfo[0].Cores),
		MemTotal:    meminfo.Total / 1024 / 1024,
		MemUsed:     meminfo.Used / 1024 / 1024,
		MemUsage:    lib.TruncFloat(meminfo.UsedPercent, 1),
		HDUsage:     lib.TruncFloat(diskinfo.UsedPercent, 1),
		HDTotal:     diskinfo.Total / 1024 / 1024 / 1024,
		HDUsed:      diskinfo.Used / 1024 / 1024 / 1024,
	}

	c.JSON(http.StatusOK, info)
}

// CPUInfo .
func CPUInfo(c *gin.Context) {
	info := cmsCPUInfo{
		CPUUsage: lib.TruncFloat(getCPUUsage(), 1),
	}

	c.JSON(http.StatusOK, info)
}

// Log .
func Log(c *gin.Context) {
	t, err := template.ParseFiles("/var/log/herefriend.log")
	if err != nil {
		c.Status(http.StatusOK)
		return
	}

	t.Execute(c.Writer, nil)
}
