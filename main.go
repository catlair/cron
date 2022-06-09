package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	timeStr := flag.String("time", "08:08:08", "运行时间")
	configPath := flag.String("config", "", "配置文件路径")
	isStart := flag.Bool("start", false, "是否开机启动")
	isOnce := flag.Bool("once", false, "是否只执行一次")
	flag.Parse()
	if *isOnce && onceDaily(*configPath) {
		return
	}
	if *isStart {
		// 开机启动
		log.Println("开机启动")
		runBiliTools(*configPath)
	}
	s := gocron.NewScheduler(time.Local)
	s.Every(1).Day().At(*timeStr).Do(runBiliTools, *configPath)
	log.Println("设置定时任务：", *timeStr)
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
		log.Printf("cmd.Run() failed with %s\n", err)
	}
	log.Printf("combined out:\n%s\n", out)
}

/*
 * 每天只执行一次
 */
func onceDaily(configPath string) bool {
	// configPath 去掉最后的文件名
	configDir := configPath[:len(configPath)-len(filepath.Base(configPath))]
	// 检查 configPath 同目录下是否存在 bt_jobs.json
	jobsPath := filepath.Join(configDir, "bt_jobs.json")
	if _, err := os.Stat(jobsPath); err != nil {
		log.Println("没有找到 bt_jobs.json")
		return false
	}
	// 读取 bt_jobs.json
	jobsFile, err := os.Open(jobsPath)
	if err != nil {
		log.Println("打开 bt_jobs.json 失败：", err)
		return false
	}
	defer jobsFile.Close()
	// 获取 bt_jobs.json 中的 lastRun（float64 ）
	var jobs map[string]interface{}
	if err := json.NewDecoder(jobsFile).Decode(&jobs); err != nil {
		log.Println("解析 bt_jobs.json 失败：", err)
		return false
	}
	lastRun := jobs["lastRun"].(float64)
	// 获取今日时间，然后与 lastRun 进行年月日比较
	today := time.Now()
	lastRunTime, err := time.Parse("2006-01-02 15:04:05", time.UnixMilli(int64(lastRun)).Format("2006-01-02 15:04:05"))
	if err != nil {
		return false
	}
	// 查看 lastRun 是否是今天
	if today.Year() == lastRunTime.Year() && today.Month() == lastRunTime.Month() && today.Day() == lastRunTime.Day() {
		log.Println("今天已经执行过了")
		return true
	}
	return false
}

func init() {
	// 设置 log 写入文件 ./logs/bilitools.log ，不存在则创建
	if _, err := os.Stat("./logs"); err != nil {
		os.Mkdir("./logs", os.ModePerm)
	}
	logFile, err := os.OpenFile("./logs/bilitools.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("打开日志文件失败：", err)
	}
	log.SetOutput(logFile)
}
