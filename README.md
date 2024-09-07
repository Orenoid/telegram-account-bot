# Telegram 记账机器人

这是一个 Telegram 记账机器人，可以帮助你记录支出和收入。

## 部署

1. 安装 Docker 和 Docker Compose
2. 克隆本仓库到本地
3. 将 env 目录中的 .example 文件都重命名为 .env 文件，并按需修改数据库名称、密码等配置
4. 在重命名后的 `bot.env` 文件中修改 `TELEBOT_TOKEN` 为你自己的 Telegram Bot Token
5. 在终端中进入项目目录，运行 `docker-compose up -d` 启动容器，如果你不需要 OpenAPI 功能，则执行 `docker-compose up -d bot` 即可
6. 打开 Telegram，搜索你的 Bot，开始使用

## 使用

以下是目前可用的命令：

    /start - 开始使用

    /day - 查看当日账单

    /month - 查看当月账单

    /set_keyboard - 设置快捷键盘

    /cancel - 取消当前操作

    /set_balance - 设置余额

    /balance - 查询余额

    /create_token 创建用于 OpenAPI 的 token

    /disable_all_tokens 废弃所有 token

## TODO
- [x] 自动设置机器人 Commands
- [x] Open API
- [ ] 多语言（不一定做）
- [ ] ~~定时自动记账~~（暂不考虑，未来可自行使用 Open API 实现）
