# PRD: AI Contract Review SaaS

## 产品定位
AI-powered contract analysis platform for SMBs, freelancers, and legal professionals.

## 目标用户画像
- **中小企业主**：缺乏法务团队，需要快速审查合同风险
- **自由职业者**：接单前审查客户合同条款
- **初创公司**：融资协议、NDA、劳动合同审查

## MVP 功能范围

### 核心功能
1. **PDF/DOCX 上传** — 支持主流格式
2. **风险条款识别** — AI 检测不公平条款、隐藏风险
3. **智能摘要** — 关键条款提炼（金额、期限、违约条款）
4. **风险评估报告** — 风险等级 + 建议修改方案

### 技术架构
```
Go + SQLite + DeepSeek LLM (法律分析)
前端：Next.js 14 + Stripe 风格 UI
部署：Docker + Render Blueprint
```

## 定价策略
- **Free**: 3 次免费审查
- **Pro**: $19.9/月 或 $99/年
- **Team**: $49.9/月（多人协作 + 条款库）

## 技术依赖
- PDF 解析：`mohae/pdfreader` (纯Go) 或 `pdfcpu`
- DOCX 解析：`unioffice/docx`
- 法律 Prompt 工程：DeepSeek custom model

## 成功指标
- MRR: $500 以内验证 PMF
- 转化率: 免费→付费 > 5%
- 准确率: 风险条款召回率 > 80%

## 竞品分析
| 竞品 | 优势 | 劣势 |
|------|------|------|
| Harvey AI | 大模型能力 | 仅面向律所，$200+/月 |
| Clerky | 模板丰富 | 无 AI 审查 |
| LawDepot | 模板多 | 无智能分析 |
