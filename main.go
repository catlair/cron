package main

import (
	"flag"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	timeStr := flag.String("time", "08:00", "运行时间")
	configPath := flag.String("config", "", "配置文件路径")
	isStart := flag.Bool("start", false, "是否开机启动")
	flag.Parse()
	if *isStart {
		runBiliTools(*configPath)
	}
	s := gocron.NewScheduler(time.Local)
	s.Every(1).Day().At(*timeStr).Do(runBiliTools, *configPath)
	s.StartBlocking()
}

func runBiliTools(configPath string) {
	var cmd *exec.Cmd
	shell := "bash"
	cArg := "-c"
	if runtime.GOOS == "windows" {
		shell = "cmd"
		cArg = "/c"
	}
	installCmd := exec.Command(shell, cArg, "npm", "install", "-g", "@catlair/bilitools")
	cmd = exec.Command(shell, cArg, "bilitools", "-c", configPath)
	// 取消 cmd 窗口
	if runtime.GOOS == "windows" {
		prepareBackgroundCommand(cmd)
		prepareBackgroundCommand(installCmd)
	}

	installCmd.Run()
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", out)
}
