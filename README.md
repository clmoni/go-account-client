# go-account

This is a Go Client library to make accessing account API easier. I've written this without any prior experience with Go so the decisions I've made are that of a Go newbie writing it with ideas I've developed working with other langauges such as C#, swift, java, javascript & some C++ (whilst at university). I've found the language to be very pleasant to write too.

## Usage

```import account```

Initialise a new account service and use the methods on the object to interact with the downstream account api. See example below;

```hp := infrastructure.NewHTTP(ctx, cancel, httpClient, url, userAgent)```
```accountService := account.NewService(hp)```
```account, err := accountService.GetByID("626e880a-e719-11ea-8eaa-8c85903c0c20")```

if you pass ```nil``` to the NewService function, then a set of defualt values are used to initialise the client.

## What drove my decisions?
 Reading https://golang.org/doc/effective_go.html, the way they want you to be concise is a bit different in that it promotes naming without stutter. My initial instincts was to have very descriptive & expressive names but I learned that Go way can be just as expressive to the reader. So, I have tried to adopt this in my code. 
 
 I wrote this with the following at the forefront of my mind;

1. reusability
2. readability
3. maintainability
4. testability

Injecting the HTTP struct with the client into service to give flexibility, I think this is a good approach because the user can create their own as per their needs or use the default http properties when ```nil``` is passed to the ```NewService``` function. Here I'm trying to follow the notion of dependency injection. To achieve this I separated the bits of code that is solely concerned with doing http communication from the bits of code that is responsible for knowing that there is a downstream account API to communicate with.

Learning Golang & looking at various Go repositories on GitHub, I noticed people write quite big/long methods that do multiple things. I've moved away from this because it makes code hard to reason about and follow.

Blocking IO calls should be cancleable or timeout after a set duration. The user can choose to provide a timeout to the http client in the HTTP struct or even better provide a context with timout that has a cancelfunc that can be called with defer & executed upon exit of the calling method. A better approach with context would be to pass this context in on a per request basis (to the methods on the Service struct & passed down to the client request in HTTP where the call is made). This would mean you can cancel each request independently. 

I wrote component tests that don't make real calls out, they just test that the account service and all its elements are working correctly. 

The integration tests call out to the real downtream account service, the integration tests exist to ensure that if the implementation of the downstream account API changes, it will be caught. The component tests wont catch these sort of errors since they work with mock responses from a test server.

In the integration tests, I provide a mechanism to get the account api base url from an environment variable ```ACCOUNT_API_BASE_URL```. I did this so the base url for the test can be passed in for account APIs running in different environments to the one running locally, perhaps a staging or pre-production environment that is an exact replica of production. It might be the case where we have a fake stub of the API running locally to gain some confidence whilst developing but to gain even more confidence we run the integration tests against the real thing in a staging environment. Environment variables give us the flexibility to do this.

