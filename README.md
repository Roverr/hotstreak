
# <img src="./hotstreak.png"/>

Hotstreak is lightweight library for creating a certain type of rate limiting solution.

It provides the tools needed for rate limiting, but does not try to solve everything.

# How

Hotstreak uses 2 terms `Active` and `Hot`. 
While it's active, you can call `Hit` to increase the inner counter
After you call `Hit` for a configurable amount of times, the streak will become `Hot`. 
`Hot` means that only deactivation can stop the service from being `Active` for a configurable amount of time.

# Example

```go
    streak := hotstreak.New(hotstreak.Config{
        Limit: 20, // Hit 20 times before getting hot
        HotWait: time.Minute * 5, // Wait 5 minutes before cooling down
        ActiveWait:  time.Minute * 10, // Wait 10 minutes before deactivation
    })

    streak.Activate()
    for _, request := range requests {
        streak.Hit()
        // If we are hitting it too hard, try slowing down
        if streak.IsHot() {
            <-time.After(time.Second * 5)
        }
    }
```