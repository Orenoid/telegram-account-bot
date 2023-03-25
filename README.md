# Telegram 记账机器人

这是一个 Telegram 记账机器人，可以帮助你记录支出和收入。

## 部署

1. 安装 Docker 和 Docker Compose
2. 克隆本仓库到本地
3. 在 `docker-compose.yml` 文件中修改 `TELEBOT_TOKEN` 为你自己的 Telegram Bot Token
4. 在终端中进入项目目录，运行 `docker-compose up -d` 启动容器
5. 打开 Telegram，搜索你的 Bot，开始使用

## 使用

以下是目前可用的命令：

    /start - 开始使用

    /day - 查看当日账单

    /month - 查看当月账单

    /set_keyboard - 设置快捷键盘

    /cancel - 取消当前操作

    /set_balance - 设置余额

    /balance - 查询余额

## TODO
- [ ] 自动设置机器人 Commands
- [ ] Open API
- [ ] 多语言（不一定做）
- [ ] ~~支持定时自动记账~~（暂不考虑，未来可自行使用 Open API 实现）
