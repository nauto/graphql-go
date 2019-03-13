package graphql_test

import (
	"fmt"
	"testing"

	graphql "github.com/nauto/graphql-go"
	qerrors "github.com/nauto/graphql-go/errors"
	"github.com/nauto/graphql-go/gqltesting"
)

type enumResolver struct{}

func (self *enumResolver) Greet(args struct{ Mood string }) string {
	return fmt.Sprintf("Hi, %s!", args.Mood)
}

func (self *enumResolver) Leave(args struct{ Moods *[]*string }) string {
	retVal := "Bye"
	for _, s := range *args.Moods {
		retVal += ", " + *s
	}
	return retVal + "!"
}

func (self *enumResolver) Grasp(args struct{ None *string }) string {
	return fmt.Sprintf("None, %s.", *args.None)
}

func (self *enumResolver) Crash(args struct{ Why string }) string {
	return fmt.Sprintf("Why, %s!", args.Why)
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
			// 1. misspelled scalar enum literal
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: `
			query {
				greet(mood: WRONG)
			}`,
			ExpectedErrors: []*qerrors.QueryError{{
				Message:   "Argument \"mood\" has invalid value WRONG.\nExpected type \"Mood\", found WRONG.",
				Locations: []qerrors.Location{{Line: 3, Column: 17}},
				Rule:      "ArgumentsOfCorrectType",
			}},
		},
		{
			// 2. correct scalar enum literal
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: `
			query {
				greet(mood: WRUNG)
			}`,
			ExpectedResult: `{ "greet": "Hi, WRUNG!" }`,
		},
		{
			// 3. misspelled list-of-enum literal
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: `
			query {
				leave(moods: [WRONG])
			}`,
			ExpectedErrors: []*qerrors.QueryError{{
				Message:   "Argument \"moods\" has invalid value [WRONG].\nIn element #0: Expected type \"Mood\", found WRONG.",
				Locations: []qerrors.Location{{Line: 3, Column: 18}},
				Rule:      "ArgumentsOfCorrectType",
			}},
		},
		{
			// 4. correct list-of-enum literal
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: `
			query {
				leave(moods: [WRUNG])
			}`,
			ExpectedResult: `{ "leave": "Bye, WRUNG!" }`,
		},
		{
			// 5. misspelled scalar enum variable
			Schema:    graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query:     varScalar,
			Variables: map[string]interface{}{"wrong": "WRONG"},
			ExpectedErrors: []*qerrors.QueryError{{
				Message:   "Expected type \"Mood\", found WRONG.",
				Locations: []qerrors.Location{{Line: 2, Column: 8}, {Line: 3, Column: 15}},
				Rule:      "ArgumentsOfCorrectType",
			}},
		},
		{
			// 6. correct scalar enum variable
			Schema:         graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query:          varScalar,
			Variables:      map[string]interface{}{"wrong": "WRUNG"},
			ExpectedResult: `{ "greet": "Hi, WRUNG!" }`,
		},
		{
			// 7. misspelled list-of-enum variable
			Schema:    graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query:     varList,
			Variables: map[string]interface{}{"wrong": `[WRONG]`},
			ExpectedErrors: []*qerrors.QueryError{{
				Message:   "In element #0: Expected type \"Mood\", found WRONG.",
				Locations: []qerrors.Location{{Line: 2, Column: 8}, {Line: 3, Column: 16}},
				Rule:      "ArgumentsOfCorrectType",
			}},
		},
		{
			// 8. correct list-of-enum variable
			Schema:         graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query:          varList,
			Variables:      map[string]interface{}{"wrong": `[WRUNG]`},
			ExpectedResult: `{ "leave": "Bye, [WRUNG]!" }`,
		},
		{
			// 9. misspelled again scalar enum literal
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: `
			query {
				greet(mood: WRU)
			}`,
			ExpectedErrors: []*qerrors.QueryError{{
				Message:   "Argument \"mood\" has invalid value WRU.\nExpected type \"Mood\", found WRU.",
				Locations: []qerrors.Location{{Line: 3, Column: 17}},
				Rule:      "ArgumentsOfCorrectType",
			}},
		},
		{
			// 10. misspelled again scalar enum variable
			Schema:    graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query:     varScalar,
			Variables: map[string]interface{}{"wrong": "WRU"},
			ExpectedErrors: []*qerrors.QueryError{{
				Message:   "Expected type \"Mood\", found WRU.",
				Locations: []qerrors.Location{{Line: 2, Column: 8}, {Line: 3, Column: 15}},
				Rule:      "ArgumentsOfCorrectType",
			}},
		},
		{
			// 11. spelled empty enum literal
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: `
			query {
				grasp(none: NOTHING)
			}`,
			ExpectedErrors: []*qerrors.QueryError{{
				Message:   "Argument \"none\" has invalid value NOTHING.\nExpected type \"Nothing\", found NOTHING.",
				Locations: []qerrors.Location{{Line: 3, Column: 17}},
				Rule:      "ArgumentsOfCorrectType",
			}},
		},
		{
			// 12. empty string literal
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: `
			query {
				crash(why: "")
			}`,
			ExpectedResult: `{ "crash": "Why, !" }`,
		},
		{
			// 13. empty string variable
			Schema: graphql.MustParseSchema(rightSchema, &enumResolver{}),
			Query: `
			query($justso: String!) {
				crash(why: $justso)
			}`,
			Variables: map[string]interface{}{"justso": ""},
			ExpectedResult: `{ "crash": "When, !" }`,
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

	enum Nothing {
	}

	type Query {
		greet(mood: Mood!): String!
		leave(moods: [Mood]): String!
		grasp(none: Nothing): String!
		crash(why: String!): String!
	}
`
