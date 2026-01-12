package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	autoSign    bool
	autoComment bool
)

// topfeelCmd 用于 Topfeel 社区
var topfeelCmd = &cobra.Command{
	Use:   "topfeel",
	Short: "用于 Topfeel 社区的自动脚本",
	Long:  "支持自动签到、自动评论等功能，可组合使用",
	Run: func(cmd *cobra.Command, args []string) {
		if !autoSign && !autoComment {
			cmd.Help()
			return
		}

		if autoSign {
			runAutoSign()
		}

		if autoComment {
			runAutoComment()
		}
	},
}

func init() {
	rootCmd.AddCommand(topfeelCmd)

	topfeelCmd.PersistentFlags().BoolVarP(&autoSign, "sign", "c", false, "自动签到")
	topfeelCmd.PersistentFlags().BoolVarP(&autoComment, "comment", "m", false, "自动评论")
}

// SignPayload 签到请求参数结构
type SignPayload struct {
	Oldtime int64 `json:"oldtime"`
	Newtime int64 `json:"newtime"`
}

// runAutoSign 运行自动签到逻辑
func runAutoSign() {
	sugar.Info("开始自动签到...")

	token := viper.GetString("topfeel.token")
	if token == "" {
		sugar.Error("未配置 Topfeel 访问 Token")
		return
	}

	const (
		baseURL    = "https://bbs.topfeel.com"
		signInPath = "/api/gift/day_sign"
		userAgent  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36"
		referer    = "https://bbs.topfeel.com/h5/"
	)

	now := time.Now().UnixMilli()
	randomDiff := rand.Int63n(5000) + 3000 // 3~8秒随机差值（毫秒）

	payload := SignPayload{
		Oldtime: now,
		Newtime: now + randomDiff,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		sugar.Error(fmt.Sprintf("生成请求体失败: %v\n", err))
		return
	}

	req, err := http.NewRequest("POST", baseURL+signInPath, bytes.NewBuffer(jsonData))
	if err != nil {
		sugar.Error(fmt.Sprintf("创建请求失败: %v\n", err))
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="143", "Chromium";v="143", "Not A(Brand";v="24"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("token", token)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		sugar.Error(fmt.Sprintf("请求失败: %v\n", err))
		return
	}
	defer resp.Body.Close()

	sugar.Infof("签到状态码: %d", resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		sugar.Info("签到请求发送成功（状态码 200）")
	} else {
		sugar.Warn("签到可能失败，建议检查响应状态码")
	}

	body, _ := io.ReadAll(resp.Body)
	sugar.Infof("响应内容: %s\n", string(body)) // {"code":200,"msg":"签到成功","time":1768187726,"data":[]}
}

// runAutoComment 运行自动评论逻辑
func runAutoComment() {
	sugar.Info("开始自动评论...")
	// TODO: 评论逻辑
}
