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

