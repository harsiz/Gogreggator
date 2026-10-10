# GOGreggator

This project shouldn't be taken seriously.

It is an RSS blog (and static podcast) aggregator created in Go.

Do not use this. It's just a Boot.dev project I added (and removed out of sheer laziness) extra features to.

To get:

```go
go install github.com/harsiz/gogreggator
```

You'd probably also want to use Postgresql for the sql queries, just to stay on the safe side.

I might rewrite for sqlite in the future (probably will actually.)

You can run these below. Some require extra arguments:

- "login" (requires username)
- "register" (requires username)
- "reset"
- "users"
- "addfeed" (requires name, and url)
- "agg" (requires url) (also, this aggregator command currently just blurts out a struct since i am lazy. in the future it will be formatted - hopefully)
- "feeds"
- "follow" (requires url)
- "following" ()
- "unfollow" (requires url)
