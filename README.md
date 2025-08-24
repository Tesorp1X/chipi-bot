# `chipi-bot` &mdash; A personal financial telegram bot

---

# Overview

### What does it do?

Basic accounting for me and my partner. Under the hood there is state-management, db requests (raw sql), and string-processing.

### What's in use?

`Go 1.25`, `SQLite`, [telebot]("gopkg.in/telebot.v4"), [fsm-telebot](github.com/vitaliy-ukiru/fsm-telebot), `testing`.

### Can it be accessed?
No, it's a private bot only for both of us 💜. Simple design and absolutely no room for scaling was intentional. Too much complexity must have a solid reason, not my case.

---

## A story behind this bot
Hi, `chipi-bot` is an expense-tracking bot, that I've made for personal use with my partner. I was getting tired of doing an accounting work and I'm learning `Golang`, so I decided to automate a bit expense-tracking and learn a thing or two from it 😄.

So, bot can receive a new receipt from a user (each item must be hand-written: name, price, owner), calculate how much each of users payed for each other. Every check 'lives' inside of a session. When session is being closed, bot notifies both of us who and how much money owes to another. Sessions are meant to be closed, when it's time for justice to be done 🤓😆. After closure, session result is saved in `db`, to keep the history (literally basic accounting).

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

As I said previously, I've had some experience building telegram-bots before and it was done with [aiogram](https://github.com/aiogram/aiogram). It's very powerful and easy to use, so I wanted something similar but in Go. I've stumbled upon two frameworks: [go-telegram](github.com/go-telegram/bot) and [telebot]("gopkg.in/telebot.v4"). I've tried both and I found [go-telegram](github.com/go-telegram/bot) harder to work with and its FSM (Finite State Machine). So, I decided to stick with [telebot]("gopkg.in/telebot.v4"). I also needed a FSM, and I found a [repo](github.com/vitaliy-ukiru/fsm-telebot/) from [vitaliy-ukiru](github.com/vitaliy-ukiru/). It's exactly what I needed and was familiar with. 
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


### Testing

It was problematic from the start. Main things I needed to test are:
- Message handlers;
- DB operations;
- Utility functions.

No problem with `util` package, because every function returns something. For example:

```Go
func TestExtractAdminsIDs(t *testing.T) {
    t.Run("line like [id, id]", func(t *testing.T) {
		s := "[123, 234, 456]"
		got := util.ExtractAdminsIDs(s)
		want := []int64{123, 234, 456}
		if !slices.Equal(got, want) {
			t.Fatalf("got %v want %v", got, want)
		}
	})
}
```

But I was kinda stuck with testing db requests or message-handlers. 

First I cracked testing `db` package! One particular feature of `SQLite3` helped me out: in memory mode. As path you can specify `:memory:` and db will be created in RAM, so it's super-fast and easy to clean up. Here is an example of how I test `AlterItem` function:

```Go
func TestAlterItem(t *testing.T) {
	db := makeInMemoryDB(t) 
	defer db.Close()

	populateItemsDB(t, db, []models.Item{
		{Id: 1, CheckId: 1, Name: "Item 1", Owner: "Owner 1", Price: 100},
		{Id: 2, CheckId: 1, Name: "Item 2", Owner: "Owner 1", Price: 200},
		{Id: 3, CheckId: 2, Name: "Item 3", Owner: "Owner 2", Price: 300},
	})

	t.Run("update item name", func(t *testing.T) {
		item := &models.Item{Id: 1, Name: "Updated Item 1"}

		if err := alterItem(db, item); err != nil {
			t.Fatalf("expected no error, but got %v", err)
		}

		var nameGot string
		err := db.QueryRow("SELECT Name FROM items WHERE id = ?", item.Id).Scan(&nameGot)
		if err != nil {
			t.Fatalf("failed to query updated item: %v", err)
		}

		if nameGot != item.Name {
			t.Fatalf("expected item name to be '%s', got '%s'", item.Name, nameGot)
		}
	})
}
```

Last thing to beat "message-handler testing". I thought of mocking everything from the start, but `tele.API` has like 100 methods, and I was not thrilled about it. But nothing came to mind so I went down the road of implementing every method of `tele.Context`, `fsm.Context` and `tele.API`. Only by the last one (and the largest one) I remembered about AI (yeah, not vibe coding, I like to suffer authentically). I prompted claude to implement all methods with zero-values as return values. I only altered some of the methods, such as: `Send`, `Reply`, `Edit` and `Response`. 

So, how do I test things? I made a struct:

```Go
// what was sent to a user
type HandlerResponse struct {
	// Text of a displayed message or Text field of [tele.tele.CallbackResponse].
	Text string
	// In which way message was sent (Send, Reply, Edit, EditOrReply, Respond).
	// Supported options are defined as iota-constants.
	Type int
	// Which [SendOptions] were used with response.
	SendOptions *tele.SendOptions
}
```

Every mocked type has a field of `HandlerResponse` and altered methods, that I mentioned earlier (Send, Reply etc.) just populating that field with response-data. For example, here is implementation of `Send` method from `MockContext`:

```Go
// MockContext mocks the original tele.Context for testing purposes.
type MockContext struct {
	bot     tele.API
	update  tele.Update
	storage *MockStorage

	response *HandlerResponse
}

func (c *MockContext) Send(what interface{}, opts ...interface{}) error {
	text, ok := what.(string)
	if !ok {
		return errors.New("expected what of type string")
	}

	c.response.Text = text
	c.response.Type = ResponseTypeSend
	c.response.SendOptions = extractOptions(opts)

	return nil
}
```

In test set-up I just inject `HandlerResponse` object and after handler done its work, I can look into that injected object and assert everything. Here is how a simple test looks like:

```Go
func TestHelloHandler(t *testing.T) {
	response := mocks.HandlerResponse{}
	bot := mocks.NewMockBot(&response)
	botStorage := mocks.NewMockStorage()
	fsmStorage := mocks.NewMockStorage()

	update := makeUpdateWithMessageText("hello")

	teleCtx := mocks.NewMockContext(bot, update, botStorage, &response)
	stateCtx := mocks.NewMockFsmContext(fsmStorage, models.StateDefault)

	expextedResponse := mocks.HandlerResponse{
		Text: "Hello, 1",
		Type: mocks.ResponseTypeSend,
	}

	if err := stateCtx.SetState(context.Background(), models.StateStart); err != nil {
		t.Fatalf("couldn't change state to %s: %v", models.StateStart, err)
	}

	handlerErr := HelloHandler(teleCtx, stateCtx)

	assertHandlerError(t, false, errEmpty, handlerErr)
	assertHandlerResponse(t, &expextedResponse, &response)
}
```

Another things I needed to assert are `fsm.State` and storage after execution of the handler. State is easy, just need to compare what state `fsm.Context` mock has with the expected one. Storage is bit harder, because it's a `map[string]any` and how to compare any values? Right, reflection! Here's what I did:

```Go

// Fails a test if storage is missing any (key, value) tuple from expected,
// or if expected and got values are not deeply equal (must have the same type).
func assertStorage(t testing.TB, expected *map[string]any, storage *mocks.MockStorage) {
	t.Helper()
	for k, v := range *expected {
		storageVal := storage.Get(k)
		if storageVal == nil {
			t.Fatalf("in storage expected (key, value): (%s, %v), but instead got nil", k, v)
		}

		expectedReflectValue := reflect.ValueOf(v)
		gotReflectValue := reflect.ValueOf(storageVal)

		if expectedReflectValue.Type() != gotReflectValue.Type() {
			t.Fatalf("in storage for key %s expected value type of %v, but insted got %v", k, expectedReflectValue.Type(), gotReflectValue.Type())
		}

		if !reflect.DeepEqual(v, storageVal) {
			t.Fatalf("in storage for for key %s expected value %v, but instaed got %v", k, expectedReflectValue, gotReflectValue)
		}
	}
}
```

To conclude... I've managed to come up with tests for all key components, but I still has problems with testing, because of architecture. For example, I can't properly test handlers, that has to interact with db. Therefore my work continues &mdash; going to rewrite it all with gained experience. 

---

### How I improved item price validation and fixed some bugs I didn't know existed

It all started with a proposal from my partner to add ability to write down floats not with a dot, but with a coma (`45,5` instead of `45.5`). Why? In russian math tradition floats are written down with a coma. So, it's just for the sake of convenience. I said yeah, why not. Should be easy: just swap coma with a dot and it's a done deal. So I did exactly so, and decided to test the thing. And that's where I discovered some unwanted behaviors.


First things first, how it should work:
- Bot sends a message asking for price of an item
- The answer can be an int, a float or \<multiplier>*\<`int` or `float`>

How it was implemented:
```Go
msgText := c.Text()
msgText = strings.ReplaceAll(msgText, " ", "")
msgText = strings.ReplaceAll(msgText, ",", ".")
```
And the validation and conversion was implemented that way:
```Go
var (
	price  float64
	amount int
	err    error
)
if strings.Contains(msgText, "*") {
	tokens := strings.Split(msgText, "*")
	if len(tokens) != 2 {
		return c.Send(models.AmountOfItemsHelpMsg)
	}
	if amount, err = strconv.Atoi(tokens[0]); err != nil {
		return c.Send(models.ErrorAmountOfItemsMustBeANumberMsg)
	}
	if price, err = strconv.ParseFloat(tokens[1], 64); err != nil {
		return c.Send(models.ErrorItemPriceMustBeANumberMsg)
	}
	price *= float64(amount)
} else {
	if price, err = strconv.ParseFloat(msgText, 64); err != nil {
		return c.Send(models.ErrorItemPriceMustBeANumberMsg)
	}
}
```

With this in mind I made `TestItemPriceResponseHandler` in `handlers_test.go` file. And "good" scenarios there was no problem. By good I mean:
- 56.5
- 56,5
- 3*56.5
- 3 * 56.5

But what's the point to test only good cases, right? Here are some more test, that should trigger an error:
- hi <3
- 55 56
- 55.6 56.5
- 55 * 56.5 *

Test with line `55 56` must fail but it doesn't, instead its being transformed into number `5556` and it becomes a price. That moment I came up with a solution &mdash; regexp :). And I decided to move all line parsing and verification into a different function. It should simplify the handler function itself and improve readability. So, I made three functions: 
- `ParsePrice` &mdash; to parse and verify the price-containing string. It returns a float value &mdash; a price multiplied by multiplier (if provided) and an error (if verification failed);
- `verifyPrice` &mdash; to verify price token. Returns true if price matches `^\s*[0-9]+([.,][0-9]+)?\s*$` pattern;
- `verifyPriceMultiplier` &mdash; to verify multiplier token. Returns true if multiplier matches `^\s*[0-9]+\s*$` pattern.

Now, part of the handler function code, dedicated to price is much simpler:
```Go
msgText := c.Text()
price, err := util.ParsePrice(msgText)
if err != nil {
	switch err {
	case models.ErrItemPriceMultiplierNotSingleInt:
		return c.Send(models.ErrorAmountOfItemsMustBeANumberMsg)
	case models.ErrItemPriceNotSingleIntOrFloat:
		return c.Send(models.ErrorItemPriceMustBeANumberMsg)
	case models.ErrItemPriceWrongFormat:
		return c.Send(models.AmountOfItemsHelpMsg)
	}
}
``` 
It also simplified the testing it seems, if it would be made that way from the start. First is the verification of user's response, make sure it works it was intended to work, than I could focus on making a handler behave the right way. Well, lesson learned (I hope so)😅.

---
To be continued...