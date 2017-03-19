[![Go Report Card](https://goreportcard.com/badge/github.com/marcossegovia/sammy-the-bot)](https://goreportcard.com/report/github.com/marcossegovia/sammy-the-bot)

# Sammy the Botpher <img alt="Sammy" src="sammy.png" width="162,5">

This is the repository for Sammy source code.

If you want to enjoy Sammy go ahead >>> https://telegram.me/SammyGoBot

If you want to get your own sammy clone:

- Follow this [link](https://core.telegram.org/bots#3-how-do-i-create-a-bot) and create a Bot through BotFather
- Download the repo `go get github.com/marcossegovia/sammy-the-bot`
- Generate a `sammy_brain.toml` inside project root, at least like the following:

```toml
[configuration]
# Telegram bot token to be able to send/retrieve messages
telegram = "YOUR-TELEGRAM-BOT-TOKEN"
# http://openweathermap.org/ token to retrieve weather on your bot
weather = "YOUR-OPEN-WEATHER-MAP-DOT-ORG-TOKEN"
# Managing Entities/Entities of your api.ai account
api_ai_developer = "YOUR-API-DOT-AI-DEVELOPER-TOKEN"
# Making queries to your api.ai account
api_ai_client = "YOUR-API-DOT-AI-CLIENT-TOKEN"
```
- `go run main.go`
- Enjoy your own botpher !

## Contribution
If you want to stay tuned about the new features Sammy is going to get, we encourage you to follow https://trello.com/b/U6k3ThJe/sammy
And of course if you find something useful, get in touch by [opening an issue](https://github.com/marcossegovia/sammy-the-bot/issues/new) or [sending an email](mailto:velozmarkdrea@gmail.com)