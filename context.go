package mgql

import (
	"context"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nrfta/go-log"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Context struct {
	Variables     map[string]interface{}
	w             io.Writer
	ctx           context.Context
	operationName string
	wrote         bool
}

func (c *Context) check() {
	if c.wrote {
		log.Fatal("MGQL: You tried to write a response more then once in " + c.operationName)
	}
}

func (c *Context) Context() context.Context {
	return c.ctx
}

func (c *Context) Data(data interface{}) {
	c.check()
	c.wrote = true
	writeJson(
		c.w,
		&graphql.Response{
			Data: JSON(data),
		},
	)
}

func (c *Context) GraphqlError(err ...*gqlerror.Error) {
	c.check()
	c.wrote = true
	writeJson(c.w, &graphql.Response{Errors: err})
}

func (c *Context) Error(err ...error) {
	c.check()
	c.wrote = true
	gqlErrors := gqlerror.List{}

	for _, e := range err {
		gqlErrors = append(gqlErrors, &gqlerror.Error{Message: e.Error()})
	}

	writeJson(c.w, &graphql.Response{Errors: gqlErrors})
}
