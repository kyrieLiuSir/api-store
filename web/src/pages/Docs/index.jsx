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
//test

import React, { useMemo } from 'react';
import { Button, Card } from '@douyinfe/semi-ui';
import { useSearchParams } from 'react-router-dom';
import MarkdownRenderer from '../../components/common/markdown/MarkdownRenderer';

const DOC_ITEMS = [
  {
    key: 'quickstart',
    title: '创建API KEY',
    markdown: `

按照下面的步骤操作，即可完成账号注册、在线充值、兑换码充值和 API Key 创建。

---

## 1. 注册账号

- 登录注册地址：<https://apisstore.com/login>

---

## 2. 在线充值

- 登录后切换到控制台。
- 点击左侧导航栏的“钱包管理”。
- 选择对应的充值额度进行充值即可。

![在线充值示意图](/docs-redeem-topup.png)

---

## 3. 兑换码充值

- 登录后切换到控制台。
- 点击左侧导航栏的“钱包管理”。
- 输入兑换码后，点击“兑换额度”即可完成充值。

![兑换码充值示意图](/docs-redeem-code.png)

---

## 4. 创建API KEY

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
    title: 'codex 安装教程',
    markdown: `

## 1. 下载 cc switch

- 下载地址：
- <https://pan.baidu.com/s/1rNw_jVWEDd6Gz0HpDY7Quw?pwd=5588>

## 2. 安装 codex

- 下载安装nodejs: https://nodejs.org/zh-cn/download/
\`\`\`bash
npm i -g @openai/codex
\`\`\`

## 3. 配置 Key

- 切换到gpt页签,点击右侧加号

![Add Provider 示例](/img.png)

- 输入 供应商名称 ApiStore
- API Key：填写控制台令牌管理的秘钥
- API请求地址：\`https://apisstore.com/v1\`
- 其他可不填写,输入完成后点击右下角添加按钮

![配置完成示例](/img_1.png)
- 保存后即可在列表中查看
![Add Provider 示例](/codex-list.png)
- 打开终端,输入codex即可对话
![Add Provider 示例](/codex-terminal.png)

`,
  },
  {
    key: 'claudecode',
    title: 'claudecode 安装教程',
    markdown: `

按照下面步骤操作，即可在 Windows 或 macOS 上完成 Claude Code 安装与接入。

---

## 1. 下载 cc switch

- 下载地址：
- <https://pan.baidu.com/s/1rNw_jVWEDd6Gz0HpDY7Quw?pwd=5588>

---

## 2. 安装 Claude Code

下载安装nodejs: https://nodejs.org/zh-cn/download/
Windows 和 macOS 都可以在终端执行下面命令：

\`\`\`bash
npm install -g @anthropic-ai/claude-code
\`\`\`

如果提示没有 Node.js / npm，请先安装 Node.js 18 及以上版本，再重新执行命令。

## 3. 配置 Key

- 切换到claude页签,点击右侧加号

![Add Provider 示例](/claude-entry.png)

- 输入 供应商名称 ApiStore
- API Key：填写控制台令牌管理的秘钥
- API请求地址：\`https://apisstore.com\`
- 其他可不填写,输入完成后点击右下角添加按钮

![配置完成示例](/claude-form.png)
- 保存后即可在列表中查看
![Add Provider 示例](/claude-list.png)
- 打开终端,输入claude即可对话
![示例](/claude-terminal.png)
---

`,
  },
  {
    key: 'openclaw',
    title: 'OpenClaw 接入教程',
    markdown: `

按照下面步骤操作，即可完成 OpenClaw 安装、基础引导和自定义 Provider 配置。

---

## 1. 安装 OpenClaw

\`\`\`bash
npm install -g openclaw
\`\`\`

如果提示没有 Node.js / npm，请先安装 Node.js 18 及以上版本，再重新执行命令。

---

## 2. 执行初始化引导

安装完成后，在终端执行下面命令：

\`\`\`bash
openclaw onboard --install-daemon
\`\`\`

执行过程中根据提示完成基础设置。

---

## 3. 选择 Custom Provider

- 在模型提供商这里选择 \`Custom Provider\`
- 后续按下面参数填写

![选择 Custom Provider](/img_4.png)

---

## 4. 填写接入参数

- URL：\`https://apisstore.com/v1\`
- API KEY：填写你在站点后台创建的令牌 Key
- EndPoint-compatibility：选择 \`OpenAI-compatible\`
- Model ID：填写 \`gpt-5.4\`
- endpoint id：直接回车，使用默认值即可
- model alias：可选项，直接回车跳过即可

---

## 5. 完成配置并开始使用

- 保存配置后即可通过 OpenClaw 发起请求
- 如果是第一次接入，建议先发送一条简单测试消息，确认模型返回正常

## 常见检查项

- URL 必须填写为 \`https://apisstore.com/v1\`
- API KEY 前后不要带空格
- \`EndPoint-compatibility\` 需选择 \`OpenAI-compatible\`
- 模型名称请填写当前可用模型，例如 \`gpt-5.4\`
`,
  },
];

const Docs = () => {
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
