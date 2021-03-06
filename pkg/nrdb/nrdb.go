// Package nrdb provides a programmatic API for interacting with NRDB, New Relic's Datastore
package nrdb

import (
	"github.com/joeyparsons/newrelic-client-go/internal/http"
	"github.com/joeyparsons/newrelic-client-go/internal/logging"
	"github.com/joeyparsons/newrelic-client-go/pkg/config"
)

// Nrdb is used to communicate with the New Relic's Datastore, NRDB.
type Nrdb struct {
	client http.Client
	logger logging.Logger
}

// New returns a new GraphQL client for interacting with New Relic's Datastore
func New(config config.Config) Nrdb {
	return Nrdb{
		client: http.NewClient(config),
		logger: config.GetLogger(),
	}
}
