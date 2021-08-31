package mgql

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nrfta/go-log"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Map map[string]interface{}

func JSON(input interface{}) []byte {
	b, err := json.Marshal(input)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

func writeJsonErrorf(w io.Writer, format string, args ...interface{}) {
	writeJson(w, &graphql.Response{Errors: gqlerror.List{{Message: fmt.Sprintf(format, args...)}}})
}

func writeJson(w io.Writer, response *graphql.Response) {
	b, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	w.Write(b)
}

func jsonDecode(r io.Reader, val interface{}) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(val)
}
