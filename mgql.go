package mgql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

type ResponseResolver func(ctx Context)

func New() *MGQL {
	m := &MGQL{
		queries:   make(map[string]ResponseResolver),
		mutations: make(map[string]ResponseResolver),

		querySpies:    make(map[string]ResponseResolver),
		mutationSpies: make(map[string]ResponseResolver),
	}

	return m
}

type MGQL struct {
	queries   map[string]ResponseResolver
	mutations map[string]ResponseResolver

	querySpies    map[string]ResponseResolver
	mutationSpies map[string]ResponseResolver
}

func (m *MGQL) Reset() {
	m.queries = make(map[string]ResponseResolver)
	m.mutations = make(map[string]ResponseResolver)

	m.querySpies = make(map[string]ResponseResolver)
	m.mutationSpies = make(map[string]ResponseResolver)
}

func (m *MGQL) Query(operationName string, responseResolver ResponseResolver) {
	m.queries[operationName] = responseResolver
}

func (m *MGQL) Mutation(operationName string, responseResolver ResponseResolver) {
	m.mutations[operationName] = responseResolver
}

func (m *MGQL) SpyQuery(operationName string, responseResolver ResponseResolver) {
	m.querySpies[operationName] = responseResolver
}

func (m *MGQL) SpyMutation(operationName string, responseResolver ResponseResolver) {
	m.mutationSpies[operationName] = responseResolver
}

func (m *MGQL) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var params *graphql.RawParams
		if err := jsonDecode(r.Body, &params); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJsonErrorf(w, "json body could not be decoded: "+err.Error())
			return
		}

		doc, err := parser.ParseQuery(&ast.Source{Input: params.Query})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJsonErrorf(w, "unable to parse query: "+err.Error())
			return
		}

		operation := doc.Operations.ForName(params.OperationName)
		var resolvers map[string]ResponseResolver
		var spies map[string]ResponseResolver

		switch operation.Operation {
		case ast.Query:
			resolvers = m.queries
			spies = m.querySpies
		case ast.Mutation:
			resolvers = m.mutations
			spies = m.mutationSpies
		case ast.Subscription:
			writeJsonErrorf(w, "subscription are not current supported")
			return
		default:
			writeJsonErrorf(w, "unsupported GraphQL operation")
			return
		}

		ctx := Context{
			ctx:           r.Context(),
			w:             w,
			operationName: operation.Name,
			Variables:     params.Variables,
		}

		spy, ok := spies[operation.Name]
		if ok {
			spy(ctx)
		}

		resolverFn, ok := resolvers[operation.Name]
		if ok {
			resolverFn(ctx)
			return
		}
		writeJsonErrorf(w, "Mock ResponseResolver not found: "+operation.Name)
	}
}
