package main

import (
	"context"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/bigquery/v2"
	"google.golang.org/api/dns/v1"
)

// BigQueryClient is the interface for connecting to bigquery
type BigQueryClient interface {
}

type bigQueryClientImpl struct {
	bqService *bigquery.Service
}

// NewBigQueryClient returns new BigQueryClient
func NewBigQueryClient() (BigQueryClient, error) {

	ctx := context.Background()
	googleClient, err := google.DefaultClient(ctx, dns.NdevClouddnsReadwriteScope)
	if err != nil {
		return nil, err
	}

	bigqueryService, err := bigquery.New(googleClient)
	if err != nil {
		return nil, err
	}

	return &bigQueryClientImpl{
		bqService: bigqueryService,
	}, nil
}
