# SuperMuxer

Super **useful** Go package to configure your HTTP routes using only the **standard library**. Define routes, middlewares, groups, and subgroups effortlessly!

This package acts like a **Swiss Army Knife**: It is **tiny** and **compact**, providing everything you need in just **one** file with **less than 200 lines of code**.

### SuperMuxer is for you if:

- You want to **declaratively** define your HTTP routes while using **only** the standard library.
- You want to define **middlewares** for your routes, groups, and subgroups while still relying on the standard library.
- You **don’t** want to use **third-party libraries** bloated with excessive functionalities that you might never use.


### How to Use

Simply **clone** the file into your project or **import** it using the following command:

```shell
go get github.com/DBarbosaDev/supermuxer@v0.1.0
```


## Functionalities

### Declarative route creation

```go

serverMux := http.NewServeMux()
superRouter := supermuxer.New(serverMux)

superRouter.Get("/users", handler).Put("/users/{id}", handler)
superRouter.Post("/login", handler)

```

### Middleware association

```go

serverMux := http.NewServeMux()
superRouter := supermuxer.New(serverMux)

superRouter.AddMiddlewares(middleware1, middleware2)
superRouter.Get("/users", handler).Put("/users/{id}", handler)
superRouter.Post("/login", handler)

```

### Groups
Create a **group of routes** from a **copy** of the original router **without middlewares**. This allows you to define separate route structures while maintaining the original route path.

```go

serverMux := http.NewServeMux()
superRouter := supermuxer.New(serverMux)
superRouter.AddMiddlewares(middleware1, middleware2)

// Route "POST /login" with middleware1 and middleware2
superRouter.Post("/login", handler)

// Creates the users group
userRouter := superRouter.Group("/users")
// Group Route "GET /users" without middlewares
userRouter.Get("", handler)

userRouter.AddMiddlewares(middleware3)
// Group Route "Put /users/{id}" with middleware3
userRouter.Put("/{id}", handler)


```

### Subgroups
Works the **same way** as the **Group** function, **but** the **subgroup reuses** the middlewares from the original router. This makes it easy to maintain middleware consistency across related routes.
```go

serverMux := http.NewServeMux()
superRouter := supermuxer.New(serverMux)
superRouter.AddMiddlewares(middleware1, middleware2)

// Route "POST /login" with middleware1 and middleware2
superRouter.Post("/login", handler)

// Creates the users group
userRouter := superRouter.SubGroup("/users")
// Group Route "GET /users" with middlewares1 and middlewares2
userRouter.Get("", handler)

userRouter.AddMiddlewares(middleware3)
// Group Route "Put /users/{id}" with middlewares1, middlewares2 and middlewares3
userRouter.Put("/{id}", handler)

```
