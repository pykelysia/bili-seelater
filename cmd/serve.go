package cmd

import (
	"fmt"
	"log"

	"bili-seelater/internal/bilibili"
	"bili-seelater/internal/config"
	"bili-seelater/internal/email"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动定时推送服务",
	Run:   runServe,
}

func runServe(cmd *cobra.Command, args []string) {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	schedule := cfg.Schedule
	if schedule == "" {
		schedule = "0 9 * * *"
	}

	c := cron.New()
	c.AddFunc(schedule, func() {
		runJob(cfg)
	})

	fmt.Printf("定时任务已启动，计划: %s\n按 Ctrl+C 退出\n", schedule)
	c.Run()
}

func runJob(cfg *config.Config) {
	fmt.Println("开始执行任务...")

	biliClient := bilibili.NewClient(
		cfg.Bilibili.SESSDATA,
		cfg.Bilibili.BiliJct,
		cfg.Bilibili.Buvid3,
	)

	videos, err := biliClient.GetToViewList()
	if err != nil {
		log.Printf("获取稍后再看列表失败: %v", err)
		return
	}

	if len(videos) == 0 {
		fmt.Println("稍后再看列表为空")
		return
	}

	fmt.Printf("获取到 %d 个视频\n", len(videos))

	emailSender := email.NewSender(&cfg.Email)
	if err := emailSender.SendPlainText(videos); err != nil {
		log.Printf("发送邮件失败: %v", err)
		return
	}

	fmt.Println("邮件发送成功")
}
