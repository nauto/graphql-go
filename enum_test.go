package graphql_test

import (
	//"context"
	//"errors"
	"fmt"
	"testing"

	graphql "github.com/nauto/graphql-go"
	qerrors "github.com/nauto/graphql-go/errors"
	"github.com/nauto/graphql-go/gqltesting"
)

type enumResolver struct {}

func (self *enumResolver) Greet(args struct { Mood string }) string {
	return fmt.Sprintf("Hi, %s!", args.Mood)
}

func TestInvalidEnum(t *testing.T) {
	gqltesting.RunTests(t, []*gqltesting.Test{
		{
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: `
			{
				greet(mood: WRONG)
			}`,
			ExpectedErrors: []*qerrors.QueryError{{
				Message: "Argument \"mood\" has invalid value WRONG.\nExpected type \"Mood\", found WRONG.",
				Locations: []qerrors.Location{{Line: 3, Column: 17}},
				Rule: "ArgumentsOfCorrectType",
			}},
		},
		{
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: `
			query($wrong: Mood!) {
				greet(mood: $wrong)
			}`,
			Variables: map[string]interface{}{ "wrong": "WRONG" },
			ExpectedErrors: []*qerrors.QueryError{{
				Message: "Argument \"mood\" has invalid value WRONG.\nExpected type \"Mood\", found WRONG.",
				Locations: []qerrors.Location{{Line: 3, Column: 17}},
				Rule: "ArgumentsOfCorrectType",
			}},
		},
	})
}

const rightSchema = `
	schema {
		query: Query
	}

	enum Mood {
		RIGHT
		WRUNG
	}

	type Query {
		greet(mood: Mood!): String!
	}
`
