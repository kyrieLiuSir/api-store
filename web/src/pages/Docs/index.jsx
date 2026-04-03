/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React, { useMemo } from 'react';
import { Button, Card, Typography } from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router-dom';
import MarkdownRenderer from '../../components/common/markdown/MarkdownRenderer';

const { Title, Text } = Typography;

const DOC_ITEMS = [
  {
    key: 'quickstart',
    title: '创建API KEY',
    markdown: `
# 创建API KEY

按照下面的步骤操作，即可完成账号注册、兑换码充值和 API Key 创建。

---

## 1. 注册账号

- 登录注册地址：<https://apisstore.com/login>

---

## 2. 兑换码充值

- 登录后切换到控制台。
- 点击左侧导航栏的“钱包管理”。
- 输入兑换码后，点击“兑换额度”即可完成充值。

![兑换码充值示意图](/docs-redeem-code.png)

---

## 3. 创建API KEY

### 进入令牌管理

- 点击顶部“控制台”。
- 点击左侧导航栏“令牌管理”。
- 点击“添加令牌”。

![进入令牌管理](/docs-token-manage.png)

### 填写令牌信息

- 令牌名称：输入便于自己识别的名称。
- 令牌分组：选择 \`default\`。
- 填写完成后点击“提交”，即可生成令牌。

![创建令牌](/docs-create-api-key.png)

### 复制 API Key

- 创建完成后，可以在令牌列表中看到刚生成的令牌。
- 点击复制按钮，即可复制 API Key。

![复制 API Key](/docs-copy-api-key.png)
`,
  },
  {
    key: 'gpt-codex',
    title: 'GPT-codex 安装教程',
    markdown: `
# GPT-codex 安装教程

## 1. 准备 API 信息

- Base URL：\`https://apisstore.com/v1\`
- API Key：使用你在站点后台生成的密钥

## 2. 常见配置思路

支持自定义 OpenAI 兼容地址的工具，一般只需要填写：

- API Key：\`YOUR_API_KEY\`
- Base URL：\`https://apisstore.com/v1\`
- Model：按站点已开通模型填写

## 3. 环境变量示例

\`\`\`bash
export OPENAI_API_KEY=YOUR_API_KEY
export OPENAI_BASE_URL=https://apisstore.com/v1
\`\`\`

## 4. 验证是否连通

\`\`\`bash
curl https://apisstore.com/v1/models \\
  -H "Authorization: Bearer YOUR_API_KEY"
\`\`\`

如果能正常返回模型列表，说明接入完成。
`,
  },
  {
    key: 'claudecode',
    title: 'claudecode 安装教程（Win/Mac）',
    markdown: `
# claudecode 安装教程（Win/Mac）

## Windows / macOS 通用配置

1. 安装客户端或命令行工具。
2. 在配置页中找到自定义 API / OpenAI Compatible / Base URL 配置项。
3. 填入以下信息：

- Base URL：\`https://apisstore.com/v1\`
- API Key：\`YOUR_API_KEY\`
- Model：按站点支持的模型填写

## 推荐检查项

- 确认系统代理没有拦截请求。
- 确认 API Key 没有多余空格。
- 确认请求地址使用的是 \`https://apisstore.com/v1\`，不是旧域名。

## 测试建议

先发送一个最简单的对话请求，确认响应正常后，再导入你自己的完整配置。
`,
  },
  {
    key: 'openclaw',
    title: 'OpenClaw 接入 API 教程',
    markdown: `
# OpenClaw 接入 API 教程

## 接入参数

- API Key：站点后台生成的密钥
- Base URL：\`https://apisstore.com/v1\`
- 协议类型：OpenAI Compatible

## 配置流程

1. 打开 OpenClaw 设置。
2. 新增一个自定义提供商。
3. 将接口地址设置为 \`https://apisstore.com/v1\`。
4. 填入 API Key。
5. 选择站点支持的模型并保存。

## 排错建议

- 报 401：检查 API Key 是否正确。
- 报 404：检查 Base URL 是否写成了 \`https://apisstore.com/v1\`。
- 报模型不存在：改为站点当前可用模型名称。
`,
  },
];

const Docs = () => {
  const { t } = useTranslation();
  const [searchParams, setSearchParams] = useSearchParams();
  const currentDoc = searchParams.get('doc') || DOC_ITEMS[0].key;

  const activeDoc = useMemo(
    () => DOC_ITEMS.find((item) => item.key === currentDoc) || DOC_ITEMS[0],
    [currentDoc],
  );

  const handleDocChange = (docKey) => {
    setSearchParams({ doc: docKey });
  };

  return (
    <div className='mt-[60px] px-4 py-6 md:px-6 lg:px-8'>
      <div className='mx-auto max-w-7xl'>

        <div className='grid grid-cols-1 gap-6 lg:grid-cols-[280px_minmax(0,1fr)]'>
          <Card bodyStyle={{ padding: 16 }} className='h-fit'>
            <div className='flex flex-col gap-2'>
              {DOC_ITEMS.map((item) => {
                const isActive = item.key === activeDoc.key;
                return (
                  <Button
                    key={item.key}
                    block
                    theme={isActive ? 'solid' : 'light'}
                    type={isActive ? 'primary' : 'tertiary'}
                    className='!justify-start !h-auto !py-3 !px-4 !whitespace-normal text-left'
                    onClick={() => handleDocChange(item.key)}
                  >
                    {item.title}
                  </Button>
                );
              })}
            </div>
          </Card>

          <Card
            title={activeDoc.title}
            bodyStyle={{ padding: 24 }}
            className='min-w-0'
          >
            <MarkdownRenderer content={activeDoc.markdown} />
          </Card>
        </div>
      </div>
    </div>
  );
};

export default Docs;
