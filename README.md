Simple Tg Bot
====

Simple Telegram Bot framework use to easily create easy interaction with bots commands.

## Demo -- Zhizhang sentence maker
Below is an interactive demo

![simple-tg-bot.gif](https://ooo.0o0.ooo/2017/04/08/58e7d8e29e4c3.gif)

## Getting Started

#### First `go get -v -u github.com/VOID001/simple-tg-bot`
#### Then create config.toml File
```
dsn = "MYSQL DSN Required"
token = "Your-Awesome-Bot-Token"
root = []
```
#### Create a MySQL database and run `github.com/VOID001/simple-tg-bot/init.sql` in it
 this will create storage for session (Future will support multiple storage type: maybe redis, sqlite, mongodb, etc)

#### Now you have setup the environment, now create a main file for your bot

There are simply five steps to go:

* Parse Config file
* Init DB Connection
* Init Bot API
* Register commands
* Run it

Sample file:

```go
package main

import (
	"flag"
	"time"

	"github.com/VOID001/simple-tg-bot/command"
	"github.com/VOID001/simple-tg-bot/model"
	"github.com/VOID001/simple-tg-bot/module"
	_ "github.com/go-sql-driver/mysql"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "config.toml", "Select the config file to use ")
	flag.Parse()
}

func main() {
	cfg := new(module.Config)
	cmdHandler := new(module.CommandHandler)

	// Read the config file
	err := cfg.Parse(configPath)

	// Initialize DB
	model.DB, err = sqlx.Connect("mysql", cfg.DSN)

	// Initialize bot
	bot, err := tg.NewBotAPI(cfg.Token)
	cmdHandler.Bot = bot

  // Register the Command
	YourCommand := new(command.YourCommand)
	cmdHandler.Register("yourcommand", YourCommand)

  // Run it!
	cmdHandler.Run()
}
```


#### Implement the command

Then implement the command  just create a file and import `github.com/VOID001/simple-tg-bot/module/session`

And implement the Command Interface

```go

package command

// Demo to show Dialog Basic usage

import (
	"fmt"

	session "github.com/Wheeeel/todobot/module/session"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

// You do not need to contain any export member here
type Ping struct {
	msg   *tg.Message
	sess  session.Session
}

func (p *Ping) Init(sess *session.Session, args ...interface{}) {
	p.sess = *sess
}

func (p *Ping) Dialog(state int, data string) (reply tg.Chattable, nextState int, err error) {
	return
}

func (p *Ping) Info() (reply tg.Chattable, err error) {
	return
}

func (p *Ping) Run() (reply tg.Chattable, err error) {
	return
}

func (p *Ping) ParseMessage(msg *tg.Message) (errMsg tg.Chattable, err error) {
	return
}

func (p *Ping) IsFinalState(state int) (ok bool) {
	return
}
```

#### Then All things aer done!

`go build && ./your-awesome-bot -c config.toml`


## Examples

You can find examples in [command/](https://github.com/VOID001/simple-tg-bot/tree/master/command)

* [Dialog State Switch & Input example] (https://github.com/VOID001/simple-tg-bot/tree/master/command/demo.go)
* [Simple Dialog Example](https://github.com/VOID001/simple-tg-bot/tree/master/command/ping.go)
