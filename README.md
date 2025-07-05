# `chipi-bot` &mdash; A personal financial telegram bot

---

## A story behiend this bot
Hi, `chipi-bot` is an expence-tracking bot, that I've made for personal use with my partner. I was getting tired of doing an accounting work and I'm learning `Golang`, so I decided to automate a bit expence-tracking and learn a thing or two from it 😄.

So, bot can recieve a new recipt from a user (each item must be hand-writen: name, price, owner), calculate how much each of users payed for each other. Every check 'lives' inside of a session. When session is being closed, bot notifies both of us who and how much money owes to another. Sessions are meant to be closed, when it's time for justice to be done 🤓😆. After closure, session result is saved in `db`, to keep the history (litteraly basic accounting).

So, bot operates with only four basic entities:
- Sessions
- Checks
- Items
- Totals


---

## My experience developing this thing (a.k.a. dev-blog)
### Starting point
I thing I should first say something about my background. I can't say I've never programmed before, I have already made several chat-bots using `Python`, but I did with defferent approach: apps themselves were just means to meet the ends (part of project for Uni, or just to play musik on my discord server). With this one I have decided I want to learn as much as possible about software-dev: ideas, concepts, practicies, hard-skills and just learn `Golang`. 

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

