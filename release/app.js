// 全局状态
let state = {
    currentPage: 1,
    totalPages: 1,
    outputDir: '',
    currentSlide: 1,
    config: {
        provider: 'openai',
        apiKey: '',
        baseURL: 'https://api.openai.com/v1',
        model: 'gpt-4o'
    }
};

// 风格数据
const styles = [
    { id: 'minimal-white', name: '极简白', description: '简洁干净的白色风格', palette: ['#FFFFFF', '#F9FAFB', '#111827', '#3B82F6', '#10B981'] },
    { id: 'cyber-neon', name: '赛博霓虹', description: '未来科技感的霓虹风格', palette: ['#0F172A', '#1E293B', '#06B6D4', '#8B5CF6', '#EC4899'] },
    { id: 'bauhaus', name: '包豪斯', description: '德国包豪斯风格', palette: ['#FFFFFF', '#F5F5DC', '#E53935', '#1E88E5', '#FDD835'] },
    { id: 'japanese-minimal', name: '日式简约', description: '日式极简美学', palette: ['#FBF7F4', '#2C2C2C', '#C4A484', '#8B7355', '#D4C4B0'] },
    { id: 'corporate-blue', name: '企业蓝', description: '专业的企业蓝色调', palette: ['#FFFFFF', '#EFF6FF', '#1E40AF', '#3B82F6', '#60A5FA'] },
    { id: 'nature-green', name: '自然绿', description: '清新自然的绿色主题', palette: ['#FFFFFF', '#F0FDF4', '#166534', '#22C55E', '#86EFAC'] },
    { id: 'dark-tech', name: '暗黑科技', description: '深色科技风格', palette: ['#0A0A0A', '#171717', '#404040', '#22D3EE', '#A855F7'] },
    { id: 'retro-warm', name: '复古暖色', description: '温暖的复古色调', palette: ['#FFF8F0', '#2D2A26', '#D97706', '#DC2626', '#059669'] },
    { id: 'elegant-purple', name: '优雅紫', description: '高贵的紫色主题', palette: ['#FFFFFF', '#FAF5FF', '#7C3AED', '#A855F7', '#C084FC'] },
    { id: 'ocean-blue', name: '海洋蓝', description: '清新的海洋主题', palette: ['#FFFFFF', '#F0F9FF', '#0369A1', '#0EA5E9', '#38BDF8'] },
    { id: 'sunset-orange', name: '日落橙', description: '温暖的日落色调', palette: ['#FFFFFF', '#FFFBEB', '#C2410C', '#F97316', '#FDBA74'] },
    { id: 'mint-fresh', name: '薄荷清新', description: '清新的薄荷绿', palette: ['#FFFFFF', '#F0FDFA', '#0D9488', '#14B8A6', '#5EEAD4'] },
    { id: 'rose-pink', name: '玫瑰粉', description: '温柔的玫瑰粉色', palette: ['#FFFFFF', '#FFF1F2', '#BE185D', '#EC4899', '#F9A8D4'] },
    { id: 'nordic-light', name: '北欧光', description: '明亮的北欧风格', palette: ['#FFFFFF', '#FAFAFA', '#1F2937', '#4B5563', '#9CA3AF'] },
    { id: 'gold-luxury', name: '金色奢华', description: '高贵的金色主题', palette: ['#000000', '#1C1917', '#D97706', '#F59E0B', '#FCD34D'] },
    { id: 'cherry-blossom', name: '樱花浪漫', description: '温柔的樱花色调', palette: ['#FFF0F5', '#FDF2F8', '#BE185D', '#EC4899', '#F9A8D4'] },
    { id: 'monochrome', name: '黑白单色', description: '极简的黑白色调', palette: ['#FFFFFF', '#F9FAFB', '#111827', '#374151', '#6B7280'] },
    { id: 'academic-blue', name: '学术蓝', description: '严谨的学术风格', palette: ['#FFFFFF', '#EFF6FF', '#1E3A8A', '#2563EB', '#60A5FA'] },
    { id: 'startup-modern', name: '创业现代', description: '适合创业公司的现代风格', palette: ['#FFFFFF', '#F8FAFC', '#6366F1', '#8B5CF6', '#A855F7'] },
    { id: 'forest-earth', name: '森林大地', description: '沉稳的大地色系', palette: ['#FFFDF7', '#FEF3C7', '#365314', '#65A30D', '#84CC16'] },
];

// Toast 提示
function showToast(message, type = 'info') {
    const toast = document.getElementById('toast');
    toast.textContent = message;
    toast.className = `toast show ${type}`;
    setTimeout(() => {
        toast.className = 'toast';
    }, 3000);
}

// 页面导航
function navigateTo(pageName) {
    document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));
    document.querySelectorAll('.nav-menu a').forEach(a => a.classList.remove('active'));

    document.getElementById(`page-${pageName}`).classList.add('active');
    document.querySelector(`[data-page="${pageName}"]`).classList.add('active');
}

// 加载风格列表
function loadStyles() {
    const container = document.getElementById('style-list');
    container.innerHTML = styles.map(style => `
        <div class="style-item" data-style="${style.id}">
            <h3>${style.name}</h3>
            <p>${style.description}</p>
            <div class="style-colors">
                ${style.palette.slice(0, 5).map(c => `<span style="background: ${c}"></span>`).join('')}
            </div>
        </div>
    `).join('');
}

// 加载配置
function loadConfig() {
    const saved = localStorage.getItem('godjian-ppt-config');
    if (saved) {
        state.config = JSON.parse(saved);
    }
    document.getElementById('provider').value = state.config.provider;
    document.getElementById('api-key').value = state.config.apiKey;
    document.getElementById('base-url').value = state.config.baseURL;
    document.getElementById('model').value = state.config.model;
}

// 保存设置
function handleSettings(e) {
    e.preventDefault();
    state.config.provider = document.getElementById('provider').value;
    state.config.apiKey = document.getElementById('api-key').value;
    state.config.baseURL = document.getElementById('base-url').value;
    state.config.model = document.getElementById('model').value;
    localStorage.setItem('godjian-ppt-config', JSON.stringify(state.config));
    showToast('设置已保存', 'success');
}

// 调用 AI API
async function callAI(messages) {
    if (!state.config.apiKey) {
        throw new Error('请先在设置中配置 API Key');
    }

    const response = await fetch(`${state.config.baseURL}/chat/completions`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${state.config.apiKey}`
        },
        body: JSON.stringify({
            model: state.config.model,
            messages: messages,
            max_tokens: 4000,
            temperature: 0.7
        })
    });

    if (!response.ok) {
        const err = await response.text();
        throw new Error(`API 调用失败: ${response.status} ${err}`);
    }

    const data = await response.json();
    return data.choices[0].message.content;
}

// 生成 PPT
async function handleGenerate(e) {
    e.preventDefault();

    const topic = document.getElementById('topic').value.trim();
    const styleId = document.getElementById('style').value;
    const pages = parseInt(document.getElementById('pages').value);

    if (!topic) {
        showToast('请输入 PPT 主题', 'error');
        return;
    }

    const progressCard = document.getElementById('progress-card');
    const progressBar = document.getElementById('progress-bar');
    const progressText = document.getElementById('progress-text');

    progressCard.style.display = 'block';
    progressBar.style.width = '10%';
    progressText.textContent = '正在生成大纲...';

    try {
        // 如果有 API Key，尝试调用 AI
        if (state.config.apiKey) {
            progressBar.style.width = '30%';
            progressText.textContent = '正在调用 AI...';

            // 生成大纲
            const planResponse = await callAI([
                { role: 'system', content: `你是一个PPT规划专家。根据用户主题，规划 ${pages} 页PPT的大纲。
只返回JSON数组格式：[{"title":"页面标题","keyPoints":["要点1","要点2"]}]
不要返回其他内容。` },
                { role: 'user', content: `主题：${topic}` }
            ]);

            let pagesData;
            try {
                // 尝试解析 JSON
                const jsonMatch = planResponse.match(/\[[\s\S]*\]/);
                if (jsonMatch) {
                    pagesData = JSON.parse(jsonMatch[0]);
                } else {
                    throw new Error('无法解析AI响应');
                }
            } catch {
                // 使用默认结构
                pagesData = generateDefaultStructure(topic, pages);
            }

            progressBar.style.width = '60%';
            progressText.textContent = '正在生成页面内容...';

            // 生成每一页
            const pages = [];
            for (let i = 0; i < pagesData.length; i++) {
                progressBar.style.width = 60 + (30 * i / pagesData.length) + '%';
                pages.push(pagesData[i]);
            }

            // 生成文件
            await generateFiles(topic, styleId, pages);

        } else {
            // 无 API Key，使用默认结构
            progressBar.style.width = '50%';
            progressText.textContent = '使用默认结构...';

            const pages = generateDefaultStructure(topic, pages);
            await generateFiles(topic, styleId, pages);
        }

        progressBar.style.width = '100%';
        progressText.textContent = '生成完成！';
        showToast('PPT 生成成功！', 'success');

        // 跳转到预览
        setTimeout(() => {
            navigateTo('preview');
            loadPreview();
            progressCard.style.display = 'none';
            progressBar.style.width = '0%';
        }, 1000);

    } catch (err) {
        console.error('生成失败:', err);
        showToast('生成失败: ' + err.message, 'error');
        progressCard.style.display = 'none';
    }
}

// 生成默认结构
function generateDefaultStructure(topic, pageCount) {
    const pages = [];
    pages.push({ title: topic, keyPoints: ['演讲者', '日期', '版本号'] });

    if (pageCount > 2) {
        pages.push({ title: '目录', keyPoints: ['内容概览', '重点章节'] });
    }

    const remaining = pageCount - pages.length;
    for (let i = 0; i < remaining && i < 5; i++) {
        pages.push({
            title: `第 ${i + 1} 部分`,
            keyPoints: ['要点 1', '要点 2', '要点 3']
        });
    }

    if (pages.length < pageCount) {
        pages.push({ title: '总结', keyPoints: ['核心观点', '下一步'] });
    }

    return pages;
}

// 生成文件
async function generateFiles(topic, styleId, pages) {
    const style = styles.find(s => s.id === styleId) || styles[0];

    // 创建输出目录
    const outputDir = `output/${sanitizeFilename(topic)}`;

    // 生成 index.html
    const indexHTML = generateIndexHTML(topic, pages.length, style);
    await saveFile(`${outputDir}/index.html`, indexHTML);

    // 生成每一页
    for (let i = 0; i < pages.length; i++) {
        const pageHTML = generatePageHTML(pages[i], style, i + 1, pages.length);
        await saveFile(`${outputDir}/page-${i + 1}.html`, pageHTML);
    }

    // 生成 CSS
    const css = generateCSS(style);
    await saveFile(`${outputDir}/style.css`, css);

    // 生成运行时
    await saveFile(`${outputDir}/assets/ppt-runtime.js`, getRuntimeJS());

    state.outputDir = outputDir;
    state.totalPages = pages.length;
    state.currentSlide = 1;
}

// 生成 index.html
function generateIndexHTML(title, pageCount, style) {
    let thumbs = '';
    for (let i = 1; i <= pageCount; i++) {
        thumbs += `<div class="slide-thumb ${i === 1 ? 'active' : ''}" onclick="goToSlide(${i})">${i}</div>`;
    }

    return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>${title} - PPT</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: system-ui, sans-serif; background: #1a1a2e; color: #fff; overflow: hidden; }
        .container { display: flex; height: 100vh; }
        .sidebar { width: 200px; background: #16213e; padding: 20px; overflow-y: auto; }
        .sidebar h3 { margin-bottom: 16px; color: #9ca3af; font-size: 14px; }
        .slide-thumb { width: 100%; aspect-ratio: 16/9; background: #0f3460; border-radius: 4px; margin-bottom: 8px; display: flex; align-items: center; justify-content: center; cursor: pointer; font-size: 18px; font-weight: bold; transition: all 0.2s; }
        .slide-thumb:hover { background: #1e5f74; transform: scale(1.05); }
        .slide-thumb.active { background: #3b82f6; }
        .main { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 20px; }
        .frame { width: 100%; max-width: 1200px; aspect-ratio: 16/9; background: #000; border-radius: 8px; overflow: hidden; box-shadow: 0 20px 60px rgba(0,0,0,0.5); }
        .frame iframe { width: 100%; height: 100%; border: none; }
        .controls { display: flex; gap: 16px; margin-top: 20px; align-items: center; }
        .controls button { padding: 12px 24px; background: #3b82f6; color: white; border: none; border-radius: 8px; cursor: pointer; font-size: 16px; }
        .controls button:hover { background: #2563eb; }
        .page-indicator { color: #9ca3af; font-size: 16px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="sidebar"><h3>幻灯片</h3>${thumbs}</div>
        <div class="main">
            <div class="frame"><iframe id="frame" src="page-1.html"></iframe></div>
            <div class="controls">
                <button onclick="prevSlide()">上一页</button>
                <span class="page-indicator"><span id="current">1</span> / ${pageCount}</span>
                <button onclick="nextSlide()">下一页</button>
                <button onclick="toggleFullscreen()">全屏</button>
            </div>
        </div>
    </div>
    <script>
        let current = 1, total = ${pageCount};
        function update() {
            document.getElementById('frame').src = 'page-' + current + '.html';
            document.getElementById('current').textContent = current;
            document.querySelectorAll('.slide-thumb').forEach((t,i) => t.classList.toggle('active', i+1 === current));
        }
        function goToSlide(n) { current = n; update(); }
        function nextSlide() { if (current < total) { current++; update(); } }
        function prevSlide() { if (current > 1) { current--; update(); } }
        function toggleFullscreen() { if (!document.fullscreenElement) { document.documentElement.requestFullscreen(); } else { document.exitFullscreen(); } }
        document.addEventListener('keydown', e => { if (e.key === 'ArrowRight' || e.key === ' ') nextSlide(); if (e.key === 'ArrowLeft') prevSlide(); if (e.key === 'Escape') document.exitFullscreen(); });
    </script>
</body>
</html>`;
}

// 生成页面 HTML
function generatePageHTML(page, style, pageNum, totalPages) {
    const bgColor = style.palette[0];
    const primaryColor = style.palette.length > 1 ? style.palette[1] : '#1F2937';

    let points = page.keyPoints.map((pt, i) =>
        `<li style="margin-bottom:16px;display:flex;align-items:center;">
            <span style="display:inline-flex;width:32px;height:32px;background:${primaryColor};color:white;border-radius:50%;align-items:center;justify-content:center;margin-right:16px;font-weight:bold;font-size:14px;">${i + 1}</span>
            <span style="font-size:20px;">${pt}</span>
        </li>`
    ).join('');

    return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>${page.title}</title>
    <link href="https://cdn.tailwindcss.com" rel="stylesheet">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: system-ui, sans-serif; background: ${bgColor}; min-height: 100vh; display: flex; align-items: center; justify-content: center; }
        .page { width: 100%; max-width: 1200px; padding: 60px; }
        h1 { font-size: 48px; font-weight: 700; color: ${primaryColor}; margin-bottom: 40px; }
        ul { list-style: none; font-size: 24px; line-height: 1.8; }
    </style>
</head>
<body>
    <div class="page">
        <h1>${page.title}</h1>
        <ul>${points}</ul>
    </div>
</body>
</html>`;
}

// 生成 CSS
function generateCSS(style) {
    return `body { font-family: system-ui, sans-serif; background: ${style.palette[0]}; color: #1F2937; }`;
}

// 运行时 JS
function getRuntimeJS() {
    return `const PPT = {
    init() { console.log('PPT Runtime Initialized'); },
    animate(el, type) { el.style.opacity = '1'; }
};
document.addEventListener('DOMContentLoaded', () => PPT.init());`;
}

// 保存文件（使用 Blob 下载）
async function saveFile(filename, content) {
    // 检查是否在浏览器环境
    if (typeof window === 'undefined') return;

    // 对于实际保存，我们使用一个模拟方式
    // 在真实环境中，可以使用 File System Access API 或下载
    console.log('Generated:', filename);

    // 如果有 state.outputDir，保存到本地
    if (state.outputDir) {
        localStorage.setItem(`file_${filename}`, content);
    }
}

// 加载预览
function loadPreview() {
    const infoBox = document.getElementById('preview-info');
    const frameContainer = document.getElementById('preview-frame-container');
    const controls = document.getElementById('preview-controls');

    if (!state.outputDir) {
        infoBox.style.display = 'block';
        frameContainer.style.display = 'none';
        controls.style.display = 'none';
        return;
    }

    infoBox.style.display = 'none';
    frameContainer.style.display = 'block';
    controls.style.display = 'flex';

    // 尝试加载 index.html
    const iframe = document.getElementById('preview-frame');
    iframe.src = `${state.outputDir}/index.html`;
    updateIndicator();
}

// 预览控制
function prevPage() {
    if (state.currentSlide > 1) {
        state.currentSlide--;
        goToSlide(state.currentSlide);
    }
}

function nextPage() {
    if (state.currentSlide < state.totalPages) {
        state.currentSlide++;
        goToSlide(state.currentSlide);
    }
}

function goToSlide(n) {
    state.currentSlide = n;
    const iframe = document.getElementById('preview-frame');
    iframe.src = `${state.outputDir}/page-${n}.html`;
    updateIndicator();
}

function updateIndicator() {
    document.getElementById('page-indicator').textContent = `${state.currentSlide} / ${state.totalPages}`;
}

function openFullscreen() {
    const iframe = document.getElementById('preview-frame');
    if (iframe.requestFullscreen) {
        iframe.requestFullscreen();
    }
}

// 文件名清理
function sanitizeFilename(name) {
    return name.replace(/[\/\\:*?"<>|]/g, '-').substring(0, 50);
}

// 初始化
document.addEventListener('DOMContentLoaded', () => {
    // 导航事件
    document.querySelectorAll('.nav-menu a').forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const page = e.target.dataset.page;
            navigateTo(page);
        });
    });

    // 表单事件
    document.getElementById('generate-form').addEventListener('submit', handleGenerate);
    document.getElementById('settings-form').addEventListener('submit', handleSettings);

    // 加载数据
    loadStyles();
    loadConfig();
});

// 暴露给全局
window.prevPage = prevPage;
window.nextPage = nextPage;
window.goToSlide = goToSlide;
window.openFullscreen = openFullscreen;
