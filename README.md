# go-account

This is a Go Client library to access fake account API. I've written this without any prior experience with Go so the decisions I've made are that of Go newbie

## Usage

```import github.com/clmoni/go-account/account```

Initialise a new account service and use the methods on the object to interact with the downstream account api. See example below;

```hp := infrastructure.NewHTTP(ctx, cancel, httpClient, url, userAgent)```
```accountService := account.NewService(hp)```
```account, err := accountService.GetByID("626e880a-e719-11ea-8eaa-8c85903c0c20")```

if you pass ```nil``` to the NewService function, then a set of defualt values are used to initialise the client.

## Decisions

Injecting HTTP with client into service to give flexibility or default
