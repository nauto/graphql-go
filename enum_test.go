package graphql_test

import (
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

func (self *enumResolver) Leave(args struct { Moods *[]*string }) string {
	retVal := "Bye";
	for _, s := range *args.Moods {
		retVal += ", " + *s
	}
	return retVal + "!"
}

func TestInvalidEnum(t *testing.T) {
	varScalar := `
	query($wrong: Mood!) {
		greet(mood: $wrong)
	}`
	varList := `
	query($wrong: [Mood]) {
		leave(moods: $wrong)
	}`
	gqltesting.RunTests(t, []*gqltesting.Test{
		{
			// misspelled scalar enum literal
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: `
			query {
				greet(mood: WRONG)
			}`,
			ExpectedErrors: []*qerrors.QueryError{{
				Message: "Argument \"mood\" has invalid value WRONG.\nExpected type \"Mood\", found WRONG.",
				Locations: []qerrors.Location{{Line: 3, Column: 17}},
				Rule: "ArgumentsOfCorrectType",
			}},
		},
		{
			// correct scalar enum literal
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: `
			query {
				greet(mood: WRUNG)
			}`,
			ExpectedResult: `{ "greet": "Hi, WRUNG!" }`,
		},
		{
			// misspelled list-of-enum literal
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: `
			query {
				leave(moods: [WRONG])
			}`,
			ExpectedErrors: []*qerrors.QueryError{{
				Message: "Argument \"moods\" has invalid value [WRONG].\nIn element #0: Expected type \"Mood\", found WRONG.",
				Locations: []qerrors.Location{{Line: 3, Column: 18}},
				Rule: "ArgumentsOfCorrectType",
			}},
		},
		{
			// correct list-of-enum literal
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: `
			query {
				leave(moods: [WRUNG])
			}`,
			ExpectedResult: `{ "leave": "Bye, WRUNG!" }`,
		},
		{
			// misspelled scalar enum variable
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: varScalar,
			Variables: map[string]interface{}{ "wrong": "WRONG" },
			ExpectedErrors: []*qerrors.QueryError{{
				Message: "Argument \"mood\" has invalid value $wrong.\nExpected type \"Mood\", found WRONG.",
				Locations: []qerrors.Location{{Line: 3, Column: 15}},
				Rule: "ArgumentsOfCorrectType",
			}},
		},
		{
			// correct scalar enum variable
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: varScalar,
			Variables: map[string]interface{}{ "wrong": "WRUNG" },
			ExpectedResult: `{ "greet": "Hi, WRUNG!" }`,
		},
		{
			// misspelled list-of-enum variable
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: varList,
			Variables: map[string]interface{}{ "wrong": `[WRONG]` },
			ExpectedErrors: []*qerrors.QueryError{{
				Message: "Argument \"moods\" has invalid value $wrong.\nIn element #0: Expected type \"Mood\", found WRONG.",
				Locations: []qerrors.Location{{Line: 3, Column: 16}},
				Rule: "ArgumentsOfCorrectType",
			}},
		},
		{
			// correct list-of-enum variable
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: varList,
			Variables: map[string]interface{}{ "wrong": `[WRUNG]` },
			ExpectedResult: `{ "leave": "Bye, [WRUNG]!" }`,
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
		leave(moods: [Mood]): String!
	}
`
