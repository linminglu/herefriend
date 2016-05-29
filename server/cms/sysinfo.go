package cms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"herefriend/lib"
	"html/template"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type Process struct {
	pid int
	cpu float64
}

func getCpuUsage() float64 {
	var out bytes.Buffer

	processes := make([]*Process, 0)
	cmd := exec.Command("ps", "aux")

	cmd.Stdout = &out
	cmd.Run()

	for {
		line, err := out.ReadString('\n')
		if err != nil {
			break
		}

		tokens := strings.Split(line, " ")
		ft := make([]string, 0)
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
			processes = append(processes, &Process{pid, cpu})
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

/*
 |    Function: SystemInfo
 |      Author: Mr.Sancho
 |        Date: 2016-02-12
 |   Arguments:
 |      Return:
 | Description: 获取最新的系统信息
 |
*/
func SystemInfo(r *http.Request) string {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.CPUInfo()
	d, _ := disk.DiskUsage("/")
	n, _ := host.HostInfo()

	info := cmsSystemSummary{
		OSDescribe:  fmt.Sprintf("%s %s", n.OS, n.PlatformVersion),
		CpuDescribe: fmt.Sprintf("%s %d Cores", c[0].ModelName, c[0].Cores),
		MemTotal:    v.Total / 1024 / 1024,
		MemUsed:     v.Used / 1024 / 1024,
		MemUsage:    lib.TruncFloat(v.UsedPercent, 1),
		HDTotal:     d.Total / 1024 / 1024 / 1024,
		HDUsed:      d.Used / 1024 / 1024 / 1024,
		HDUsage:     lib.TruncFloat(d.UsedPercent, 1),
	}

	jsonRlt, _ := json.Marshal(info)
	return string(jsonRlt)
}

/*
 |    Function: CpuInfo
 |      Author: Mr.Sancho
 |        Date: 2016-02-28
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func CpuInfo(r *http.Request) string {
	info := cmsCpuInfo{
		CpuUsage: lib.TruncFloat(getCpuUsage(), 1),
	}

	jsonRlt, _ := json.Marshal(info)
	return string(jsonRlt)
}

/*
 |    Function: Log
 |      Author: Mr.Sancho
 |        Date: 2016-05-08
 | Description:
 |      Return:
 |
*/
func Log(w http.ResponseWriter) {
	t, err := template.ParseFiles("/var/log/herefriend.log")
	if nil != err {
		w.WriteHeader(200)
		return
	}

	t.Execute(w, nil)
}
