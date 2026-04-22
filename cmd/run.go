package cmd

import (
	"fmt"
	"log"

	"bili-seelater/internal/bilibili"
	"bili-seelater/internal/config"
	"bili-seelater/internal/email"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "立即执行一次获取并发送",
	Run:   runRun,
}

func runRun(cmd *cobra.Command, args []string) {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	biliClient := bilibili.NewClient(
		cfg.Bilibili.SESSDATA,
		cfg.Bilibili.BiliJct,
		cfg.Bilibili.Buvid3,
	)

	videos, err := biliClient.GetToViewList()
	if err != nil {
		log.Fatalf("获取稍后再看列表失败: %v", err)
	}

	if len(videos) == 0 {
		fmt.Println("稍后再看列表为空")
		return
	}

	fmt.Printf("获取到 %d 个视频\n", len(videos))

	emailSender := email.NewSender(&cfg.Email)
	if err := emailSender.SendPlainText(videos); err != nil {
		fmt.Printf("发送邮件失败: %v\n", err)
		fmt.Printf("SMTP配置: host=%s, port=%d, user=%s\n", cfg.Email.SMTPHost, cfg.Email.SMTPPort, cfg.Email.Username)
		log.Fatalf("发送邮件失败: %v", err)
	}

	fmt.Println("邮件发送成功")
}
