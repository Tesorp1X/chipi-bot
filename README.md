# `chipi-bot` &mdash; A personal financial telegram bot

---

# Overview

### What does it do?

Basic accounting for me and my partner. Under the hood there is state-manegment, db requests (raw sql), and string-processing.

### What's in use?

`Go 1.25`, `SQLite`, [telebot]("gopkg.in/telebot.v4"), [fsm-telebot](github.com/vitaliy-ukiru/fsm-telebot), `testing`.

### Can it be accessed?
No, it's a private bot only for both of us 💜. Simple design and absolutely no room for scaling was intentional. Too much complexity must have a solid reason, not my case.

---

## A story behind this bot
Hi, `chipi-bot` is an expence-tracking bot, that I've made for personal use with my partner. I was getting tired of doing an accounting work and I'm learning `Golang`, so I decided to automate a bit expence-tracking and learn a thing or two from it 😄.

So, bot can receive a new receipt from a user (each item must be hand-writen: name, price, owner), calculate how much each of users payed for each other. Every check 'lives' inside of a session. When session is being closed, bot notifies both of us who and how much money owes to another. Sessions are meant to be closed, when it's time for justice to be done 🤓😆. After closure, session result is saved in `db`, to keep the history (litterally basic accounting).

So, bot operates with only four basic entities:
- Sessions
- Checks
- Items
- Totals


---

## My experience developing this thing (a.k.a. dev-blog)
### Starting point
I thing I should first say something about my background. I can't say I've never programmed before, I have already made several chat-bots using `Python`, but I did with different approach: apps themselves were just means to meet the ends (part of project for Uni, or just to play music on my discord server). With this one I have decided I want to learn as much as possible about software-dev: ideas, concepts, practices, hard-skills and just learn `Golang`. 

So, let's make a check-point on my abilities at the moment of beginning this project:
- Previous experience with telegram/discord bots dev
- Familiarity with Go-syntax, telegram-api
- Basic knowledge of algorithms & data-structures, SQL, Unit-testing 
- Some git knowledge
- Very vague understanding of how to structure Go project, or how it should be built architecture wise
- No idea about concurrency in Go and patterns related to it, logging, docker, CI/CD

My goals for `ver 1.0`:
- Fully functioning bot with core functionality
- Some amount of testing of different components (message-responders, utility-tools, db-tools)
- Implement basic logging
- Deploy

### The bot-framework 
Bot itself is just a HTTP-server, so I could've made it from scratch (for educational purposes), but I thought it would add so much complexity, that I might never accomplish anything with project. Maybe one day... (`ver 3.0` goals???)

As I said previously, I've had some experience building telegram-bots before and it was done with [aiogram](https://github.com/aiogram/aiogram). It's very powerful and easy to use, so I wanted something similar but in Go. I've stumbled upon two frameworks: [go-telegram](github.com/go-telegram/bot) and [telebot]("gopkg.in/telebot.v4"). I've tried both and I found [go-telegram](github.com/go-telegram/bot) harder to work with and its FSM (Finite State Machine). So, I decided to stick with [telebot]("gopkg.in/telebot.v4"). I also needed a FSM, and I found a [repo](github.com/vitaliy-ukiru/fsm-telebot/) from [vitaliy-ukiru](github.com/vitaliy-ukiru/). It's exactly what I needded and was familiar with. 
Finite State Machine for telebot. Based on aiogram FSM version. This is what said on the `README.md`:

> It not a full implementation FSM. It just states manager for telegram bots.

### Storage
Cold storage was no-brainer for me &mdash; I chose `SQLite`. I didn't need any overhead with db, and... only two users. I guess it's a great choice in my case 😃. For temporary storage of any data I'm just using `context` and `FSM-Storage` from [fsm-telebot](github.com/vitaliy-ukiru/fsm-telebot/). But maybe, I need to use some Key/Value db, just to be sure 👀.

Why no ORM? It's simple, I wanted to learn SQL and how to handle raw-queries. I kinda love to know "how it works on the inside" about everything, so why the hell not learn SQL 😃.

### Architecture and planning
Well, I guess my project structure and architecture is pretty chaotic 😅. Well, that was my idea &mdash; I wanted to see how it will end up. The result is expected: it's hard to add functionality, hard to test and no way of scaling. At least its seems so.

To be honest, I encountered challenges with my ingenious design the moment I've tried to db-storage to my bot. For the moment for each request app opens db, queries and closes it. 

Another problem I've encountered is logging. I wanted logging in different files, but how can I supply it to a message-handler function? No way in my implementation.

And of course, it's kinda hard to write unit-tests for this thing, but more about that in Testing block.

For the `ver 2.0` I already know what I'm going to change:

- Bot object should live with db and logging services in one happy struct.
```Go
type App struct {
    bot *tele.Bot
    log *services.Log
    db *services.DB
}
```
And any handler must be a method of `App`, so it has access to logging and db of this app.

- All three of these services will be running in their own goroutines, to optimize performance (I will try to benchmark it).

- On app architecture I'm still undecided. Will talk about in `ver 2.0` dev-blog 😉👀.



---
To be continued...