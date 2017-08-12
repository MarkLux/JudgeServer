package server

import (
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func Ping(c *gin.Context) {

	hostname, err := os.Hostname()

	if err != nil {
		hostname = ""
		log.Println("Failed to get HostName")
	}

	cpuPercent, _ := cpu.Percent(0, false)
	vmem, _ := mem.VirtualMemory()

	c.JSON(http.StatusOK, gin.H{
		"judger_version": "0.1.0",
		"hostname":       hostname,
		"cpu_core":       runtime.NumCPU(),
		"cpu":            cpuPercent,
		"memory":         vmem.UsedPercent,
		"action":         "pong",
	})
}

func Judge(c *gin.Context) {
	
}
