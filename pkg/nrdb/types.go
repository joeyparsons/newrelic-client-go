// Code generated by tutone: DO NOT EDIT
package nrdb

import (
	"github.com/joeyparsons/newrelic-client-go/internal/serialization"
)

// EventAttributeDefinition - A human-readable definition of an NRDB Event Type Attribute
type EventAttributeDefinition struct {
	// This attribute's category
	Category string `json:"category"`
	// A short description of this attribute
	Definition string `json:"definition"`
	// The New Relic docs page for this attribute
	DocumentationUrl string `json:"documentationUrl"`
	// The human-friendly formatted name of the attribute
	Label string `json:"label"`
	// The name of the attribute
	Name string `json:"name"`
}

// EventDefinition - A human-readable definition of an NRDB Event Type
type EventDefinition struct {
	// A list of attribute definitions for this event type
	Attributes []EventAttributeDefinition `json:"attributes"`
	// A short description of this event
	Definition string `json:"definition"`
	// The human-friendly formatted name of the event
	Label string `json:"label"`
	// The name of the event
	Name string `json:"name"`
}

// NrdbMetadata - An object containing metadata about the query and result.
type NrdbMetadata struct {
	// A list of the event types that were queried.
	EventTypes []string `json:"eventTypes"`
	// A list of facets that were queried.
	Facets []string `json:"facets"`
	// Messages from NRDB included with the result.
	Messages []string `json:"messages"`
	// Details about the query time window.
	TimeWindow NrdbMetadataTimeWindow `json:"timeWindow"`
}

// NrdbMetadataTimeWindow - An object representing details about a query's time window.
type NrdbMetadataTimeWindow struct {
	// Timestamp marking the query begin time.
	Begin EpochMilliseconds `json:"begin"`
	// A clause representing the comparison time window.
	CompareWith string `json:"compareWith"`
	// Timestamp marking the query end time.
	End EpochMilliseconds `json:"end"`
	// SINCE clause resulting from the query
	Since string `json:"since"`
	// UNTIL clause resulting from the query
	Until string `json:"until"`
}

// NrdbResultContainer - A data structure that contains the results of the NRDB query along
// with other capabilities that enhance those results.
//
// Direct query results are available through `results`, `totalResult` and
// `otherResult`. The query you made is accessible through `nrql`, along with
// `metadata` about the query itself. Enhanced capabilities include
// `eventDefinitions`, `suggestedFacets` and more.
type NrdbResultContainer struct {
	// In a `COMPARE WITH` query, the `currentResults` contain the results for the current comparison time window.
	CurrentResults []NrdbResult `json:"currentResults"`
	// Generate a publicly sharable Embedded Chart URL for the NRQL query.
	//
	// For more details, see [our docs](https://docs.newrelic.com/docs/apis/graphql-api/tutorials/query-nrql-through-new-relic-graphql-api#embeddable-charts).
	EmbeddedChartUrl string `json:"embeddedChartUrl"`
	// Retrieve a list of event type definitions, providing descriptions
	// of the event types returned by this query, as well as details
	// of their attributes.
	EventDefinitions []EventDefinition `json:"eventDefinitions"`
	// Metadata about the query and result.
	Metadata NrdbMetadata `json:"metadata"`
	// The [NRQL](https://docs.newrelic.com/docs/insights/nrql-new-relic-query-language/nrql-resources/nrql-syntax-components-functions) query that was executed to yield these results.
	Nrql Nrql `json:"nrql"`
	// In a `FACET` query, the `otherResult` contains the aggregates representing the events _not_
	// contained in an individual `results` facet
	OtherResult NrdbResult `json:"otherResult"`
	// In a `COMPARE WITH` query, the `previousResults` contain the results for the previous comparison time window.
	PreviousResults []NrdbResult `json:"previousResults"`
	// The query results. This is a flat list of objects who's structure matches the query submitted.
	Results []NrdbResult `json:"results"`
	// Generate a publicly sharable static chart URL for the NRQL query.
	StaticChartUrl string `json:"staticChartUrl"`
	// Retrieve a list of suggested NRQL facets for this NRDB query, to be used with
	// the `FACET` keyword in NRQL.
	//
	// Results are based on historical query behaviors.
	//
	// If the query already has a `FACET` clause, it will be ignored for the purposes
	// of suggesting facets.
	//
	// For more details, see [our docs](https://docs.newrelic.com/docs/apis/graphql-api/tutorials/nerdgraph-graphiql-nrql-tutorial#suggest-facets).
	SuggestedFacets []NrqlFacetSuggestion `json:"suggestedFacets"`
	// Suggested queries that could help explain an anomaly in your timeseries based on either statistical differences in the data or historical usage.
	//
	// If no `anomalyTimeWindow` is supplied, we will attempt to detect a spike in the NRQL results. If no spike is found, the suggested query results will be empty.
	//
	// Input NRQL must be a TIMESERIES query and must have exactly one result.
	SuggestedQueries SuggestedNrqlQueryResponse `json:"suggestedQueries"`
	// In a `FACET` query, the `totalResult` contains the aggregates representing _all_ the events,
	// whether or not they are contained in an individual `results` facet
	TotalResult NrdbResult `json:"totalResult"`
}

// NrqlFacetSuggestion - A suggested NRQL facet. Facet suggestions may be either a single attribute, or
// a list of attributes in the case of multi-attribute facet suggestions.
type NrqlFacetSuggestion struct {
	// A list of attribute names comprising the suggested facet.
	//
	// Raw attribute names will be returned here. Attribute names may need to be
	// backtick-quoted before inclusion in a NRQL query.
	Attributes []string `json:"attributes"`
	// A modified version of the input NRQL, with a `FACET ...` clause appended.
	// If the original NRQL had a `FACET` clause already, it will be replaced.
	Nrql Nrql `json:"nrql"`
}

// NrqlHistoricalQuery - An NRQL query executed in the past.
type NrqlHistoricalQuery struct {
	// The Account ID queried.
	AccountID int `json:"accountId"`
	// The NRQL query executed.
	Nrql Nrql `json:"nrql"`
	// The time the query was executed.
	Timestamp EpochSeconds `json:"timestamp"`
}

// SuggestedNrqlQueryResponse - result type encapsulating suggested queries
type SuggestedNrqlQueryResponse struct {
	// List of suggested queries.
	Suggestions []SuggestedNrqlQuery `json:"suggestions"`
}

// SuggestedNrqlQuery - Interface type representing a query suggestion.
type SuggestedNrqlQueryType struct {
	// The NRQL string to run for the suggested query
	Nrql string `json:"nrql"`
	// A human-readable title describing what the query shows
	Title string `json:"title"`
}

func (x *SuggestedNrqlQueryType) ImplementsSuggestedNrqlQuery() {}

// EpochMilliseconds - The `EpochMilliseconds` scalar represents the number of milliseconds since the Unix epoch
type EpochMilliseconds serialization.EpochTime

// EpochSeconds - The `EpochSeconds` scalar represents the number of seconds since the Unix epoch
type EpochSeconds serialization.EpochTime

// NrdbResult - This scalar represents a NRDB Result. It is a `Map` of `String` keys to values.
//
// The shape of these objects reflect the query used to generate them, the contents
// of the objects is not part of the GraphQL schema.
type NrdbResult map[string]interface{}

// Nrql - This scalar represents a NRQL query string.
//
// See the [NRQL Docs](https://docs.newrelic.com/docs/insights/nrql-new-relic-query-language/nrql-resources/nrql-syntax-components-functions) for more information about NRQL syntax.
type Nrql string

// SuggestedNrqlQuery - Interface type representing a query suggestion.
type SuggestedNrqlQuery interface {
	ImplementsSuggestedNrqlQuery()
}
