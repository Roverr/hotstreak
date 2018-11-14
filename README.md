
# <img src="./hotstreak.png"/>
[ ![Codeship Status for Roverr/hotstreak](https://app.codeship.com/projects/498c53d0-c99b-0136-5bad-7e8852079539/status?branch=master)](https://app.codeship.com/projects/314997)
[![Go Report Card](https://goreportcard.com/badge/github.com/Roverr/hotstreak)](https://goreportcard.com/report/github.com/Roverr/hotstreak)


Hotstreak is lightweight library for creating a certain type of rate limiting solution.

It provides the tools needed for rate limiting, but does not try to solve everything.

# How

Hotstreak uses 2 terms `Active` and `Hot`. </br>
While it's active, you can call `Hit` to increase the inner counter.</br>
After you call `Hit` for a configurable amount of times, the streak will become `Hot`.</br>
`Hot` means that only deactivation can stop the service from being `Active` for a configurable amount of time. (It also helps with handling hits as fast as possible)</br>
After a configurable time, if no `Hit` were made at all, it deactivates. 

### Config
* `Limit`      - int           - _Describes how many times we have to hit before a streak becomes hot_
* `HotWait`    - time.Duration - _Describes the amount of time we are waiting before declaring a cool down_
* `ActiveWait` - time.Duration - _Describes the amount of time we are waiting to check on a streak being active_

### Chainability
Most commands are chainable to allow easier handling.
```go
    streak := hotstreak.New(hotstreak.Config{
        Limit: 20, // Hit 20 times before getting hot
        HotWait: time.Minute * 5, // Wait 5 minutes before cooling down
        ActiveWait:  time.Minute * 10, // Wait 10 minutes before deactivation
    })
    streak.Activate().Hit()

    // do things

    if streak.Hit().IsHot() {
        // Hit and do other things if the streak became hot
    }
```

See [docs for more info](https://godoc.org/github.com/Roverr/hotstreak).

# Example

Make certain number of requests in given time period
```go
    streak := hotstreak.New(hotstreak.Config{
        Limit: 20,
        HotWait: time.Minute * 5,
        ActiveWait:  time.Minute * 10,
    })

    streak.Activate()
    for _, request := range requests {
        streak.Hit()
        // If we are hitting it too hard, try slowing down
        if streak.IsHot() {
            <-time.After(time.Second * 5)
        }
        // .. logic
    }
```
