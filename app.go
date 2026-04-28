package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx       context.Context
	config    *AIConfig
	outputDir string
	styles    []StyleSkill
}

type AIConfig struct {
	Provider string `json:"provider"`
	BaseURL  string `json:"base_url"`
	Model    string `json:"model"`
	APIKey   string `json:"api_key"`
}

type StyleSkill struct {
	Style       string   `json:"style"`
	StyleName   string   `json:"styleName"`
	Aliases     []string `json:"aliases"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Source      string   `json:"source"`
	StyleSkill  string   `json:"styleSkill"`
}

type SlideContent struct {
	Title   string   `json:"title"`
	Content []string `json:"content"`
}

type StyleCatalogItem struct {
	ID          string `json:"id"`
	StyleKey    string `json:"styleKey"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Source      string `json:"source"`
}

func NewApp() *App {
	return &App{
		config: &AIConfig{
			Provider: "openai",
			BaseURL:  "https://api.openai.com/v1",
			Model:    "gpt-4o",
		},
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	home, _ := os.UserHomeDir()
	a.outputDir = filepath.Join(home, "OhMyPPT", "outputs")
	os.MkdirAll(a.outputDir, 0755)

	a.loadStyles()
	runtime.LogInfo(a.ctx, fmt.Sprintf("App started: %s, loaded %d styles", a.outputDir, len(a.styles)))
}

func (a *App) loadStyles() {
	execPath, _ := os.Executable()
	baseDir := filepath.Dir(execPath)

	stylesPath := filepath.Join(baseDir, "resources", "styles.json")

	if _, err := os.Stat(stylesPath); os.IsNotExist(err) {
		devPath := filepath.Join("resources", "styles.json")
		if _, err := os.Stat(devPath); err == nil {
			stylesPath = devPath
		}
	}

	data, err := os.ReadFile(stylesPath)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to load styles: %s", err.Error()))
		a.styles = getDefaultStyles()
		return
	}

	var styles []StyleSkill
	if err := json.Unmarshal(data, &styles); err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to parse styles: %s", err.Error()))
		a.styles = getDefaultStyles()
		return
	}

	a.styles = styles
	if len(a.styles) == 0 {
		a.styles = getDefaultStyles()
	}
}

func getDefaultStyles() []StyleSkill {
	return []StyleSkill{
		{
			Style:      "minimal-white",
			StyleName:  "极简白",
			Aliases:    []string{"minimal", "light", "简约"},
			Description: "极简白，克制高级",
			Category:   "浅色 · 沉静",
			Source:     "builtin",
			StyleSkill: "Use minimal-white style with clean white background and blue accents.",
		},
		{
			Style:      "cyberpunk-neon",
			StyleName:  "赛博霓虹",
			Aliases:    []string{"cyberpunk", "neon", "赛博"},
			Description: "纯黑+霓虹粉青黄+发光",
			Category:   "效果 · 戏剧",
			Source:     "builtin",
			StyleSkill: "Use cyberpunk-neon style with dark background and neon pink/cyan/yellow accents.",
		},
	}
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) SaveConfig(config AIConfig) (string, error) {
	a.config = &config
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, "OhMyPPT", "config.json")
	os.MkdirAll(filepath.Dir(configPath), 0755)
	data, _ := json.MarshalIndent(config, "", "  ")
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return "", err
	}
	return configPath, nil
}

func (a *App) LoadConfig() (AIConfig, error) {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, "OhMyPPT", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return AIConfig{
			Provider: "openai",
			BaseURL:  "https://api.openai.com/v1",
			Model:    "gpt-4o",
			APIKey:   "",
		}, nil
	}
	var config AIConfig
	json.Unmarshal(data, &config)
	a.config = &config
	return config, nil
}

func (a *App) ListStyles() ([]StyleCatalogItem, error) {
	items := make([]StyleCatalogItem, len(a.styles))
	for i, s := range a.styles {
		items[i] = StyleCatalogItem{
			ID:          s.Style,
			StyleKey:    s.Style,
			Label:       s.StyleName,
			Description: s.Description,
			Category:    s.Category,
			Source:      s.Source,
		}
	}
	return items, nil
}

func (a *App) GetStyleDetail(styleId string) (StyleSkill, error) {
	for _, s := range a.styles {
		if s.Style == styleId {
			return s, nil
		}
	}
	return StyleSkill{}, fmt.Errorf("风格不存在: %s", styleId)
}

func (a *App) GeneratePPT(prompt string, theme string) (string, error) {
	runtime.LogInfo(a.ctx, fmt.Sprintf("Generating PPT: theme=%s", theme))

	if a.config.APIKey == "" {
		return "", fmt.Errorf("请先配置API密钥")
	}

	style := a.getStyleById(theme)

	slides, err := a.callAIForSlides(prompt, style)
	if err != nil {
		return "", err
	}

	htmlPath, err := a.saveHTML(slides, style)
	if err != nil {
		return "", err
	}

	runtime.LogInfo(a.ctx, "Generated: "+htmlPath)
	return htmlPath, nil
}

func (a *App) getStyleById(styleId string) StyleSkill {
	for _, s := range a.styles {
		if s.Style == styleId {
			return s
		}
	}
	if len(a.styles) > 0 {
		return a.styles[0]
	}
	return StyleSkill{
		Style:     "minimal-white",
		StyleName: "极简白",
		StyleSkill: "Use minimal-white style.",
	}
}

func (a *App) callAIForSlides(prompt string, style StyleSkill) ([]SlideContent, error) {
	systemPrompt := fmt.Sprintf(`你是一个专业的PPT设计师。请根据用户需求生成PPT大纲。

风格: %s
风格描述: %s

风格详细指南:
%s

请生成5-8张幻灯片，每张包含标题和2-4个要点。
JSON格式:
{"slides":[{"title":"标题","content":["要点1","要点2"]},...]}

只返回JSON。`, style.StyleName, style.Description, style.StyleSkill)

	messages := []map[string]string{
		{"role": "system", "content": systemPrompt},
		{"role": "user", "content": prompt},
	}

	body, _ := json.Marshal(map[string]interface{}{
		"model":    a.config.Model,
		"messages": messages,
	})

	req, _ := http.NewRequest("POST",
		fmt.Sprintf("%s/chat/completions", strings.TrimSuffix(a.config.BaseURL, "/")),
		bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.config.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API请求失败: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode != 200 {
		errMsg, _ := json.Marshal(result["error"])
		return nil, fmt.Errorf("API错误: %s", string(errMsg))
	}

	choices := result["choices"].([]interface{})
	msg := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	content := msg["content"].(string)

	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var data struct {
		Slides []SlideContent `json:"slides"`
	}
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}

	return data.Slides, nil
}

func parseGradient(gradient string) string {
	if strings.HasPrefix(gradient, "linear-gradient") || strings.HasPrefix(gradient, "radial-gradient") {
		return gradient
	}
	return "linear-gradient(135deg, #1a1a2e, #16213e)"
}

func (a *App) saveHTML(slides []SlideContent, style StyleSkill) (string, error) {
	bg := parseGradient("linear-gradient(135deg, #ffffff, #f8fafc)")
	color := "#0f172a"
	cardBg := "rgba(255,255,255,0.94)"

	switch style.Style {
	case "bauhaus":
		bg = "linear-gradient(160deg, #fef3c7 0%, #fef9c3 50%, #f0f9ff 100%)"
		color = "#991b1b"
		cardBg = "rgba(255,255,255,0.94)"
	case "memphis-pop":
		bg = "linear-gradient(135deg, #fef3c7 0%, #f0abfc 45%, #818cf8 100%)"
		color = "#0f172a"
		cardBg = "rgba(255,255,255,0.92)"
	case "sharp-mono":
		bg = "linear-gradient(160deg, #ffffff 0%, #f5f5f5 55%, #e5e5e5 100%)"
		color = "#000000"
		cardBg = "rgba(255,255,255,0.96)"
	case "swiss-grid":
		bg = "linear-gradient(160deg, #ffffff 0%, #f8fafc 65%, #f1f5f9 100%)"
		color = "#0f172a"
		cardBg = "rgba(255,255,255,0.96)"
	case "neo-brutalism":
		bg = "linear-gradient(160deg, #fff7ed 0%, #ffedd5 50%, #fde68a 100%)"
		color = "#0f172a"
		cardBg = "rgba(255,255,255,0.96)"
	case "retro-tv":
		bg = "linear-gradient(145deg, #fef3c7 0%, #fde68a 50%, #fcd34d 100%)"
		color = "#78350f"
		cardBg = "rgba(254,243,199,0.85)"
	case "midcentury":
		bg = "linear-gradient(160deg, #fef3c7 0%, #fef9c3 50%, #fef3c7 100%)"
		color = "#78350f"
		cardBg = "rgba(255,255,255,0.9)"
	case "news-broadcast":
		bg = "linear-gradient(145deg, #ffffff 0%, #fef2f2 55%, #fee2e2 100%)"
		color = "#7f1d1d"
		cardBg = "rgba(255,255,255,0.94)"
	case "magazine-bold":
		bg = "linear-gradient(160deg, #fef3c7 0%, #fef9c3 50%, #fefce8 100%)"
		color = "#451a03"
		cardBg = "rgba(255,255,255,0.92)"
	case "arctic-cool":
		bg = "linear-gradient(145deg, #f0f9ff 0%, #e0f2fe 50%, #bae6fd 100%)"
		color = "#0c4a6e"
		cardBg = "rgba(255,255,255,0.9)"
	case "nord":
		bg = "radial-gradient(circle at 20% 0%, #2e3440 0%, #3b4252 50%, #434c5e 100%)"
		color = "#eceff4"
		cardBg = "rgba(46,52,64,0.75)"
	case "tokyo-night":
		bg = "radial-gradient(circle at 20% 0%, #1e293b 0%, #0f172a 50%, #020617 100%)"
		color = "#e2e8f0"
		cardBg = "rgba(15,23,42,0.72)"
	case "rose-pine":
		bg = "radial-gradient(circle at 20% 0%, #191724 0%, #1f1d2e 50%, #26233a 100%)"
		color = "#e0def4"
		cardBg = "rgba(25,23,36,0.75)"
	case "catppuccin-mocha":
		bg = "radial-gradient(circle at 20% 0%, #303446 0%, #24273a 50%, #1e2030 100%)"
		color = "#c6d0f5"
		cardBg = "rgba(48,52,70,0.75)"
	case "dracula":
		bg = "radial-gradient(circle at 20% 0%, #282a36 0%, #21222c 50%, #191a21 100%)"
		color = "#f8f8f2"
		cardBg = "rgba(40,42,54,0.8)"
	case "gruvbox-dark":
		bg = "radial-gradient(circle at 20% 0%, #282828 0%, #1d2021 50%, #16191a 100%)"
		color = "#ebdbb2"
		cardBg = "rgba(40,40,40,0.8)"
	case "sunset-warm":
		bg = "linear-gradient(135deg, #fed7aa 0%, #fdba74 45%, #fb923c 100%)"
		color = "#7c2d12"
		cardBg = "rgba(255,255,255,0.88)"
	case "minimal-white":
		bg = "linear-gradient(145deg, #ffffff 0%, #f8fafc 55%, #f1f5f9 100%)"
		color = "#0f172a"
		cardBg = "rgba(255,255,255,0.94)"
	case "solarized-light":
		bg = "linear-gradient(145deg, #fdf6e3 0%, #eee8d5 55%, #e8e0cc 100%)"
		color = "#073642"
		cardBg = "rgba(253,246,227,0.92)"
	case "soft-pastel":
		bg = "linear-gradient(135deg, #fef3c7 0%, #fce7f3 45%, #dbeafe 100%)"
		color = "#831843"
		cardBg = "rgba(255,255,255,0.88)"
	case "xiaohongshu-white":
		bg = "linear-gradient(160deg, #ffffff 0%, #fff5f5 55%, #fff1f0 100%)"
		color = "#7c2d12"
		cardBg = "rgba(255,255,255,0.95)"
	case "editorial-serif":
		bg = "linear-gradient(160deg, #fef7ed 0%, #fff7ed 52%, #fffbeb 100%)"
		color = "#431407"
		cardBg = "rgba(255,250,245,0.9)"
	case "catppuccin-latte":
		bg = "linear-gradient(145deg, #eff1f5 0%, #e6e9ef 55%, #ccd0da 100%)"
		color = "#4c4f69"
		cardBg = "rgba(255,255,255,0.85)"
	case "engineering-whiteprint":
		bg = "linear-gradient(145deg, #ffffff 0%, #fafafa 55%, #f5f5f5 100%)"
		color = "#1e3a5f"
		cardBg = "rgba(255,255,255,0.94)"
	case "corporate-clean":
		bg = "linear-gradient(145deg, #ffffff 0%, #f8fafc 55%, #f1f5f9 100%)"
		color = "#0f172a"
		cardBg = "rgba(255,255,255,0.96)"
	case "japanese-minimal":
		bg = "linear-gradient(145deg, #faf8f5 0%, #f5f0e8 55%, #ede5d8 100%)"
		color = "#1c1c1c"
		cardBg = "rgba(255,255,255,0.9)"
	case "pitch-deck-vc":
		bg = "linear-gradient(145deg, #ffffff 0%, #f5f3ff 50%, #ede9fe 100%)"
		color = "#1e1b4b"
		cardBg = "rgba(255,255,255,0.92)"
	case "academic-paper":
		bg = "linear-gradient(145deg, #fafafa 0%, #f5f5f5 55%, #e5e5e5 100%)"
		color = "#171717"
		cardBg = "rgba(255,255,255,0.95)"
	case "cyberpunk-neon":
		bg = "radial-gradient(circle at 30% 10%, #1a1a2e 0%, #0f0f1a 50%, #000000 100%)"
		color = "#f8fafc"
		cardBg = "rgba(15,15,26,0.85)"
	case "vaporwave":
		bg = "radial-gradient(circle at 30% 0%, #2e1065 0%, #7c3aed 40%, #06b6d4 80%, #0f172a 100%)"
		color = "#f8fafc"
		cardBg = "rgba(46,16,101,0.6)"
	case "y2k-chrome":
		bg = "linear-gradient(135deg, #e5e5e5 0%, #f5f5f5 30%, #d4d4d4 60%, #a8a8a8 100%)"
		color = "#171717"
		cardBg = "rgba(255,255,255,0.7)"
	case "rainbow-gradient":
		bg = "linear-gradient(90deg, #fecaca 0%, #fef3c7 17%, #fef9c3 33%, #dcfce7 50%, #dbeafe 67%, #e0e7ff 83%, #fae8ff 100%)"
		color = "#1e3a5f"
		cardBg = "rgba(255,255,255,0.9)"
	case "aurora":
		bg = "linear-gradient(135deg, #a7f3d0 0%, #6ee7b7 30%, #6366f1 60%, #a855f7 100%)"
		color = "#1e1b4b"
		cardBg = "rgba(255,255,255,0.55)"
	case "blueprint":
		bg = "linear-gradient(145deg, #1e3a5f 0%, #1e40af 50%, #1d4ed8 100%)"
		color = "#dbeafe"
		cardBg = "rgba(30,58,95,0.7)"
	case "glassmorphism":
		bg = "linear-gradient(140deg, #c7d2fe 0%, #e0e7ff 45%, #cffafe 100%)"
		color = "#0f172a"
		cardBg = "rgba(255,255,255,0.46)"
	case "terminal-green":
		bg = "radial-gradient(circle at 20% 0%, #0a1a0a 0%, #0d280d 50%, #0f330f 100%)"
		color = "#4ade80"
		cardBg = "rgba(10,26,10,0.85)"
	}

	var sb strings.Builder

	sb.WriteString("<!DOCTYPE html>\n")
	sb.WriteString("<html lang=\"zh-CN\">\n<head>\n")
	sb.WriteString("<meta charset=\"UTF-8\">\n")
	sb.WriteString("<meta name=\"viewport\" content=\"width=device-width,initial-scale=1\">\n")
	sb.WriteString(fmt.Sprintf("<title>%s - PPT Preview</title>\n", escapeHTML(style.StyleName)))
	sb.WriteString("<style>\n")
	sb.WriteString("*{margin:0;padding:0;box-sizing:border-box}\n")
	sb.WriteString("body{font-family:'Segoe UI','PingFang SC',sans-serif;overflow:hidden}\n")
	sb.WriteString(fmt.Sprintf(".slide{width:100vw;height:100vh;display:flex;flex-direction:column;justify-content:center;align-items:center;padding:80px;background:%s;color:%s;transition:opacity 0.5s}\n", bg, color))
	sb.WriteString(".slide h1{font-size:3.5em;margin-bottom:50px;text-shadow:3px3px15px rgba(0,0,0,0.4);text-align:center}\n")
	sb.WriteString(".slide ul{list-style:none;max-width:950px;width:100%}\n")
	sb.WriteString(fmt.Sprintf(".slide li{font-size:1.7em;padding:18px 35px;margin:12px 0;background:%s;border-radius:12px;backdrop-filter:blur(10px)}\n", cardBg))
	sb.WriteString(".slide li::before{content:'▸ ';opacity:0.7}\n")
	sb.WriteString(".nav{position:fixed;bottom:30px;left:50%;transform:translateX(-50%);display:flex;gap:20px;z-index:100}\n")
	sb.WriteString(".nav button{padding:14px 35px;font-size:1.1em;border:none;border-radius:30px;cursor:pointer;background:rgba(255,255,255,0.2);color:white;backdrop-filter:blur(15px);transition:all 0.3s}\n")
	sb.WriteString(".nav button:hover{background:rgba(255,255,255,0.4);transform:scale(1.05)}\n")
	sb.WriteString(".page{position:fixed;bottom:35px;right:40px;font-size:1.1em;color:rgba(255,255,255,0.6)}\n")
	sb.WriteString(".fullscreen{position:fixed;top:20px;right:20px;padding:10px 20px;font-size:0.9em}\n")
	sb.WriteString(".style-tag{position:fixed;top:20px;left:20px;padding:8px 16px;background:rgba(255,255,255,0.15);border-radius:20px;font-size:0.9em;color:rgba(255,255,255,0.7)}\n")
	sb.WriteString("</style>\n")
	sb.WriteString("</head>\n<body>\n")

	sb.WriteString(fmt.Sprintf("<div class=\"style-tag\">风格: %s</div>\n", escapeHTML(style.StyleName)))

	for i, s := range slides {
		var items strings.Builder
		for _, c := range s.Content {
			items.WriteString(fmt.Sprintf("<li>%s</li>", escapeHTML(c)))
		}
		sb.WriteString(fmt.Sprintf("<div class=\"slide\" id=\"s%d\" style=\"display:flex\"><h1>%s</h1><ul>%s</ul></div>\n",
			i, escapeHTML(s.Title), items.String()))
	}

	sb.WriteString("<div class=\"nav\">\n")
	sb.WriteString("<button onclick=\"prev()\">◀</button>\n")
	sb.WriteString("<button onclick=\"next()\">▶</button>\n")
	sb.WriteString("<button class=\"fullscreen\" onclick=\"toggleFS()\">⛶</button>\n")
	sb.WriteString("</div>\n")
	sb.WriteString(fmt.Sprintf("<div class=\"page\"><span id=\"cur\">1</span> / %d</div>\n", len(slides)))
	sb.WriteString("<script>\n")
	sb.WriteString(fmt.Sprintf("let i=0,total=%d;\n", len(slides)))
	sb.WriteString("const slides=document.querySelectorAll('.slide');\n")
	sb.WriteString("function show(n){slides.forEach(s=>s.style.display='none');i=(n+total)%total;slides[i].style.display='flex';document.getElementById('cur').textContent=i+1}\n")
	sb.WriteString("function next(){show(i+1)}\n")
	sb.WriteString("function prev(){show(i-1)}\n")
	sb.WriteString("function toggleFS(){document.fullscreenElement?document.exitFullscreen():document.documentElement.requestFullscreen()}\n")
	sb.WriteString("document.addEventListener('keydown',e=>{if(e.key==='ArrowRight'||e.key===' ')next();if(e.key==='ArrowLeft')prev();if(e.key==='Escape')document.exitFullscreen()});\n")
	sb.WriteString("show(0);\n")
	sb.WriteString("</script>\n")
	sb.WriteString("</body>\n</html>")

	filename := fmt.Sprintf("presentation_%d.html", time.Now().UnixNano()/1e6)
	path := filepath.Join(a.outputDir, filename)
	os.WriteFile(path, []byte(sb.String()), 0644)
	return path, nil
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
