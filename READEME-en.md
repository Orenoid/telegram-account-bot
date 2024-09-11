English | [简体中文](./README.md)

# Telegram Accounting Bot

This is a Telegram accounting bot that can help you record expenses and income.

## Deployment

1. Install Docker and Docker Compose.
2. Clone this repository to your local machine.
3. Rename all `.example` files in the `env` directory to `.env` files and modify configurations such as the database name and password as needed.
4. In the renamed `bot.env` file, change the `TELEBOT_TOKEN` to your own Telegram Bot Token.
5. Navigate to the project directory in the terminal and run `docker-compose up -d` to start the containers. If you don't need the OpenAPI feature, you can run `docker-compose up -d bot` instead.
6. Open Telegram, search for your bot, and start using it.

## Usage

Here are the currently available commands:

- `/start` - Start using the bot
- `/day` - View today's bill
- `/month` - View this month's bill
- `/set_keyboard` - Set quick access keyboard
- `/cancel` - Cancel the current operation
- `/set_balance` - Set the balance
- `/balance` - Check the balance
- `/create_token` - Create a token for OpenAPI
- `/disable_all_tokens` - Disable all tokens

## TODO
- [x] Automatically update Bot Commands
- [x] Open API
- [ ] Multilingual
- [ ] Natural language interface
