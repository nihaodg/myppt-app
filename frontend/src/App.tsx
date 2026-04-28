import {useState, useEffect} from 'react';
import {Settings, Sparkles, FileText, Download, ChevronRight, Palette, Layers, X} from 'lucide-react';
import './style.css';
import {GeneratePPT, ListStyles, LoadConfig, SaveConfig, GetStyleDetail} from '../wailsjs/go/main/App';
import {main} from '../wailsjs/go/models';

function App() {
  const [showSettings, setShowSettings] = useState(false);
  const [showStylePicker, setShowStylePicker] = useState(false);
  const [prompt, setPrompt] = useState('');
  const [theme, setTheme] = useState('minimal-white');
  const [themeName, setThemeName] = useState('极简白');
  const [generating, setGenerating] = useState(false);
  const [htmlPath, setHtmlPath] = useState('');
  const [error, setError] = useState('');
  const [config, setConfig] = useState<main.AIConfig>({
    provider: 'openai',
    base_url: 'https://api.openai.com/v1',
    model: 'gpt-4o',
    api_key: '',
  });
  const [styles, setStyles] = useState<main.StyleCatalogItem[]>([]);
  const [categories, setCategories] = useState<Record<string, main.StyleCatalogItem[]>>({});
  const [selectedStyle, setSelectedStyle] = useState<main.StyleSkill | null>(null);

  useEffect(() => {
    LoadConfig().then(c => {
      setConfig(c as main.AIConfig);
    }).catch(() => {});

    ListStyles().then(s => {
      setStyles(s as main.StyleCatalogItem[]);
      const cats: Record<string, main.StyleCatalogItem[]> = {};
      (s as main.StyleCatalogItem[]).forEach((style: main.StyleCatalogItem) => {
        const cat = style.category || '其他';
        if (!cats[cat]) cats[cat] = [];
        cats[cat].push(style);
      });
      setCategories(cats);
    }).catch(() => {});

    GetStyleDetail('minimal-white').then(s => {
      setSelectedStyle(s as main.StyleSkill);
      setThemeName((s as main.StyleSkill).styleName);
    }).catch(() => {});
  }, []);

  const handleSaveConfig = async () => {
    try {
      await SaveConfig(config);
      setShowSettings(false);
    } catch (e) {
      console.error(e);
    }
  };

  const handleGenerate = async () => {
    if (!prompt.trim()) {
      setError('请输入PPT主题或内容描述');
      return;
    }
    if (!config.api_key) {
      setError('请先配置API密钥');
      setShowSettings(true);
      return;
    }

    setGenerating(true);
    setError('');
    try {
      const path = await GeneratePPT(prompt, theme);
      setHtmlPath(path);
    } catch (e: any) {
      setError(e.message || '生成失败');
    } finally {
      setGenerating(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) {
      handleGenerate();
    }
  };

  const selectStyle = async (styleId: string) => {
    try {
      const s = await GetStyleDetail(styleId);
      setSelectedStyle(s as main.StyleSkill);
      setTheme(styleId);
      setThemeName((s as main.StyleSkill).styleName);
      setShowStylePicker(false);
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <div className="app-container">
      <header className="header">
        <div className="logo">
          <Sparkles className="logo-icon" size={28} />
          <span>Oh My PPT</span>
        </div>
        <div className="header-actions">
          <button className="style-preview-btn" onClick={() => setShowStylePicker(true)}>
            <Palette size={18} />
            <span>{themeName}</span>
          </button>
          <button className="settings-btn" onClick={() => setShowSettings(true)}>
            <Settings size={20} />
          </button>
        </div>
      </header>

      <main className="main-content">
        <div className="input-section">
          <h2>描述你的PPT需求</h2>
          <p className="subtitle">输入主题或详细描述，AI将为你生成专业PPT</p>

          <div className="textarea-wrapper">
            <textarea
              className="prompt-input"
              placeholder="例如：帮我做一个关于人工智能在教育领域应用的产品介绍PPT，包含背景、核心功能、优势和展望..."
              value={prompt}
              onChange={e => setPrompt(e.target.value)}
              onKeyDown={handleKeyDown}
              rows={5}
            />
          </div>

          {error && <div className="error-msg">{error}</div>}

          <button
            className={`generate-btn ${generating ? 'loading' : ''}`}
            onClick={handleGenerate}
            disabled={generating}
          >
            {generating ? (
              <>
                <span className="spinner"></span>
                生成中...
              </>
            ) : (
              <>
                <Sparkles size={20} />
                生成PPT
                <ChevronRight size={18} />
              </>
            )}
          </button>
        </div>

        {htmlPath && (
          <div className="preview-section">
            <div className="preview-header">
              <FileText size={20} />
              <span>预览文件已生成</span>
            </div>
            <p className="preview-path">{htmlPath}</p>
            <div className="preview-actions">
              <a href={`file://${htmlPath}`} target="_blank" className="preview-btn primary">
                <Download size={18} />
                在浏览器中打开
              </a>
              <button onClick={() => setHtmlPath('')} className="preview-btn secondary">
                新建
              </button>
            </div>
          </div>
        )}
      </main>

      {showStylePicker && (
        <div className="modal-overlay" onClick={() => setShowStylePicker(false)}>
          <div className="modal style-picker-modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3><Palette size={20} /> 选择风格</h3>
              <button className="close-btn" onClick={() => setShowStylePicker(false)}>×</button>
            </div>
            <div className="modal-body style-picker-body">
              {Object.entries(categories).map(([cat, items]) => (
                <div key={cat} className="style-category">
                  <h4 className="category-title">{cat}</h4>
                  <div className="style-grid">
                    {items.map(item => (
                      <button
                        key={item.id}
                        className={`style-card ${theme === item.id ? 'active' : ''}`}
                        onClick={() => selectStyle(item.id)}
                      >
                        <div className="style-name">{item.label}</div>
                        <div className="style-desc">{item.description}</div>
                      </button>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {showSettings && (
        <div className="modal-overlay" onClick={() => setShowSettings(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>API 设置</h3>
              <button className="close-btn" onClick={() => setShowSettings(false)}>×</button>
            </div>
            <div className="modal-body">
              <div className="form-group">
                <label>Provider</label>
                <select value={config.provider} onChange={e => setConfig({...config, provider: e.target.value})}>
                  <option value="openai">OpenAI</option>
                  <option value="custom">Custom</option>
                </select>
              </div>
              <div className="form-group">
                <label>Base URL</label>
                <input
                  type="text"
                  value={config.base_url}
                  onChange={e => setConfig({...config, base_url: e.target.value})}
                  placeholder="https://api.openai.com/v1"
                />
              </div>
              <div className="form-group">
                <label>Model</label>
                <input
                  type="text"
                  value={config.model}
                  onChange={e => setConfig({...config, model: e.target.value})}
                  placeholder="gpt-4o"
                />
              </div>
              <div className="form-group">
                <label>API Key</label>
                <input
                  type="password"
                  value={config.api_key}
                  onChange={e => setConfig({...config, api_key: e.target.value})}
                  placeholder="sk-..."
                />
              </div>
              <div className="form-tip">
                支持 OpenAI 兼容 API，包括本地 Ollama
              </div>
            </div>
            <div className="modal-footer">
              <button className="btn-secondary" onClick={() => setShowSettings(false)}>取消</button>
              <button className="btn-primary" onClick={handleSaveConfig}>保存</button>
            </div>
          </div>
        </div>
      )}

      <footer className="footer">
        <span>Powered by AI • Windows Desktop App • {styles.length} 种风格可用</span>
      </footer>
    </div>
  );
}

export default App;
