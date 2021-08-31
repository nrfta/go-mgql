# Go Mock Graphql (go-mgql)

Library to mock graphql requests.


## Usage

### Create a new mgql instance: 

```go
m := mgql.New()
m.Query("GetAccount", func(ctx mgql.Context) {
	ctx.Data(
		mgql.Map{
			"account": mgql.Map{
				"id": ctx.Variables["id"],
			},
		},
	)
})
```

### Pass the HTTP handler to a test server:

```go
testServer := httptest.NewServer(m.Handler())
```

### Reset mgql registry before each test:

```go
ginkgo.BeforeEach(func() {
	m.Reset()
})
```

### Overwrite default resolvers at your test:

```go
It("creates service appointment when onboarding", func() {
  // ...
	m.Query("GetAccount", func(ctx mgql.Context) {
    // ..
	})
  // ... 
})
```

## API

* func New() *MGQL
* func (m *MGQL) Handler() http.HandlerFunc
* func (m *MGQL) Mutation(operationName string, responseResolver ResponseResolver)
* func (m *MGQL) Query(operationName string, responseResolver ResponseResolver)
* func (m *MGQL) Reset()
* func (m *MGQL) SpyMutation(operationName string, responseResolver ResponseResolver)
* func (m *MGQL) SpyQuery(operationName string, responseResolver ResponseResolver)
