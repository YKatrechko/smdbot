# SMD bot - A telegram bot written in Golang
Simple Telegram bot for finding SMD transistor by code

[![Build Status](https://travis-ci.org/YKatrechko/smdbot.svg?branch=master)](https://travis-ci.org/YKatrechko/smdbot)

## Demo bot: [@smdtransistorbot](https://telegram.me/smdtransistorbot)

# Installation

```bash
go build github.com/ykatrechko/smdbot
```

# Configuration

The bot reads a json file for configuration:

```json
{
        "Token" : "XXXXXXX",
        "DBFile": "smdbase.db",
        "LogFile": "output.log",
        "Debug": true
}
```

* **Token**: required - Telegram's Botfather will give you one
* **DBFile**: required -  sqlite3 database
* **LogFile**: optional -  log file name, default: "smdbot.log"

# Run

```bash
$ smdbot -c config-smdbot.json -l smdbot.log
```
