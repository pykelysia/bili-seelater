package cmd

import (
	"errors"
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
		handleBilibiliErrorServe(err, cfg)
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

func handleBilibiliErrorServe(err error, cfg *config.Config) {
	if errors.Is(err, bilibili.ErrSessionExpired) || errors.Is(err, bilibili.ErrAuthFailed) {
		emailSender := email.NewSender(&cfg.Email)
		alertErr := emailSender.SendAlert(
			"B站 Cookie 已过期",
			"您的 B站 SESSDATA 已过期，请重新登录 B站 获取新的 Cookie。\n\n请更新配置文件 config.yaml 中的 sessdata、bili_jct、buvid3。",
		)
		if alertErr != nil {
			log.Printf("发送过期提醒邮件失败: %v", alertErr)
		}
		log.Printf("B站认证失败，请更新 Cookie: %v", err)
		return
	}
	log.Printf("获取稍后再看列表失败: %v", err)
}
