# Sammy the Botpher <img alt="Sammy" src="sammy.png" width="162,5">

This is the repository for Sammy source code.

If you want to enjoy Sammy go ahead >>> https://telegram.me/SammyGoBot

If you want to get your own sammy clone:

- Follow this [link](https://core.telegram.org/bots#3-how-do-i-create-a-bot) and create a Bot through BotFather
- Download the repo `go get github.com/MarcosSegovia/sammy-the-bot`
- Generate a `sammy_brain.toml` inside project root, at least like the following:
```toml
[welcome]
salutations = [
    "Hey, are you doing fine?",
    "Nice to hear from you, sir",
    "How are you feeling?",
    "Hi",
    "Time to make me work?",
    "It's been a while, it's all okay?"
]

[configuration]
telegram = "YOUR-TELEGRAM-BOT-TOKEN"
weather = "YOUR-OPEN-WEATHER-MAP-DOT-ORG-TOKEN"
```
- `go run main.go`
- Enjoy your own botpher !