package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/godjian/myppt-app/internal/config"
)

// Client AI 客户端
type Client struct {
	client   *http.Client
	baseURL  string
	model    string
	apiKey   string
	provider string
}

// Message 对话消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Choice 选择
type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage 使用量
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Response AI 响应
type Response struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// PPTPlan PPT 大纲计划
type PPTPlan struct {
	Title     string     `json:"title"`
	Subtitle  string     `json:"subtitle"`
	Pages     []PagePlan `json:"pages"`
	Theme     string     `json:"theme"`
	Palette   []string   `json:"palette"`
}

// PagePlan 页面计划
type PagePlan struct {
	Title      string   `json:"title"`
	KeyPoints  []string `json:"key_points"`
	LayoutType string   `json:"layout_type"` // "cover", "content", "chart", "summary"
}

// DesignContract 设计契约
type DesignContract struct {
	Theme        string   `json:"theme"`
	Background   string   `json:"background"`
	Palette      []string `json:"palette"`
	TitleStyle   string   `json:"title_style"`
	LayoutMotif  string   `json:"layout_motif"`
	ChartStyle   string   `json:"chart_style"`
	ShapeLanguage string  `json:"shape_language"`
}

// NewClient 创建 AI 客户端
func NewClient(cfg *config.Config) *Client {
	return &Client{
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
		baseURL:  cfg.BaseURL,
		model:    cfg.Model,
		apiKey:   cfg.APIKey,
		provider: cfg.Provider,
	}
}

// ChatCompletion 聊天完成
func (c *Client) ChatCompletion(messages []Message) (string, error) {
	url := fmt.Sprintf("%s/chat/completions", c.baseURL)

	payload := map[string]interface{}{
		"model": c.model,
		"messages": messages,
		"max_tokens": 4000,
		"temperature": 0.7,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 返回错误 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("API 未返回任何选择")
	}

	return result.Choices[0].Message.Content, nil
}

// GeneratePPTPlan 生成 PPT 大纲
func (c *Client) GeneratePPTPlan(topic string, pageCount int) (*PPTPlan, error) {
	systemPrompt := fmt.Sprintf(`你是一位PPT结构规划专家。根据用户的主题和需求，规划出每页的标题和关键点。

## 强制约束（最高优先级）
- 你必须恰好返回 %d 页的规划结果
- 无论主题内容多少，都不允许返回少于或多于 %d 项
- 如果内容不够分 %d 页，请合理拆分或补充过渡页

## 规则：
- 标题应简洁、有层次、能体现叙事逻辑
- 首页通常是封面，末页通常是总结或致谢
- 关键点必须短句化：每页只给 1-6 个关键点
- 每个关键点尽量控制在 8-20 个字

只返回 JSON 数组，不要返回任何额外说明。
格式示例：
[{"title":"封面","key_points":["项目名称","演讲者","日期"]},{"title":"目录","key_points":["内容概览"]},...]`, pageCount, pageCount, pageCount)

	userPrompt := fmt.Sprintf("主题：%s\n\n请规划一个 %d 页的 PPT 大纲。", topic, pageCount)

	resp, err := c.ChatCompletion([]Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	})
	if err != nil {
		return nil, err
	}

	// 解析 JSON
	var plan []map[string]interface{}
	if err := json.Unmarshal([]byte(resp), &plan); err != nil {
		return nil, fmt.Errorf("解析 PPT 计划失败: %w", err)
	}

	pptPlan := &PPTPlan{
		Title:  topic,
		Pages:  make([]PagePlan, 0, len(plan)),
	}

	for _, p := range plan {
		page := PagePlan{
			Title: getString(p, "title"),
		}
		
		if kp, ok := p["key_points"].([]interface{}); ok {
			for _, k := range kp {
				if s, ok := k.(string); ok {
					page.KeyPoints = append(page.KeyPoints, s)
				}
			}
		}
		pptPlan.Pages = append(pptPlan.Pages, page)
	}

	return pptPlan, nil
}

// GenerateDesignContract 生成设计契约
func (c *Client) GenerateDesignContract(topic, style string) (*DesignContract, error) {
	systemPrompt := fmt.Sprintf(`你是一位PPT视觉系统设计师。根据主题、风格和大纲，生成一份设计契约。

## 风格约束
你必须严格遵循风格规范来生成设计契约。

## 字段语义：
- theme 是视觉气质/设计方向
- background 背景描述
- palette 配色方案 (3-6 个颜色)
- titleStyle 标题样式 (Tailwind CSS)
- layoutMotif 布局特征
- chartStyle 图表风格
- shapeLanguage 形状语言

只返回 JSON 对象，不要返回任何额外说明。
格式示例：
{"theme":"calm editorial","background":"warm white","palette":["#f7f3e8","#5f7550"],"title_style":"text-5xl font-semibold","layout_motif":"spacious grids","chart_style":"muted lines","shape_language":"8px radius"}`)

	userPrompt := fmt.Sprintf("主题：%s\n风格：%s\n\n请生成设计契约。", topic, style)

	resp, err := c.ChatCompletion([]Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	})
	if err != nil {
		return nil, err
	}

	var contract DesignContract
	if err := json.Unmarshal([]byte(resp), &contract); err != nil {
		// 尝试清理 JSON
		resp = cleanJSON(resp)
		if err := json.Unmarshal([]byte(resp), &contract); err != nil {
			return nil, fmt.Errorf("解析设计契约失败: %w", err)
		}
	}

	return &contract, nil
}

// GeneratePageContent 生成页面内容
func (c *Client) GeneratePageContent(page PagePlan, contract *DesignContract, style string) (string, error) {
	systemPrompt := fmt.Sprintf(`你是PPT生成专家，负责将页面大纲落地为 HTML 内容。

## 风格约束
风格：%s

## 设计契约
- 主题：%s
- 背景：%s
- 配色：%v
- 标题样式：%s
- 布局：%s
- 图表：%s
- 形状：%s

## 画布约束
- 页面固定按 16:9 比例 (1600×900 像素)
- 根容器使用 p-8
- 不要使用固定像素宽度
- 使用 Tailwind CSS 类

## 内容规则
- 只返回 <main> 内部的内容 HTML
- 不要返回 <!doctype>、<html>、<head>、<body>
- 禁止输出 <script src=...>
- 所有标签必须成对闭合
- 最低内容密度：每页至少 3-5 个内容块
- 使用 grid/flex 布局

页面标题：%s
关键点：%v`, style, contract.Theme, contract.Background, contract.Palette,
		contract.TitleStyle, contract.LayoutMotif, contract.ChartStyle, contract.ShapeLanguage,
		page.Title, page.KeyPoints)

	resp, err := c.ChatCompletion([]Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: fmt.Sprintf("请生成页面内容：%s", page.Title)},
	})
	if err != nil {
		return "", err
	}

	return resp, nil
}

// 辅助函数
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func cleanJSON(s string) string {
	// 移除可能的 markdown 代码块
	s = removePrefix(s, "```json")
	s = removePrefix(s, "```")
	s = removeSuffix(s, "```")
	
	// 移除首尾空白
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\n' || s[0] == '\r' || s[0] == '\t') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\n' || s[len(s)-1] == '\r' || s[len(s)-1] == '\t') {
		s = s[:len(s)-1]
	}
	
	return s
}

func removePrefix(s, prefix string) string {
	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}

func removeSuffix(s, suffix string) string {
	if len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix {
		return s[:len(s)-len(suffix)]
	}
	return s
}
