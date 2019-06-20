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

// func NewDatasetsService(s *Service) *DatasetsService {
// 	rs := &DatasetsService{s: s}
// 	return rs
// }

// Query: Runs a BigQuery SQL query synchronously and returns query
// results if the query completes within a specified timeout.
// func (r *JobsService) Query(projectId string, queryrequest *QueryRequest) *JobsQueryCall {
// 	c := &JobsQueryCall{s: r.s, urlParams_: make(gensupport.URLParams)}
// 	c.projectId = projectId
// 	c.queryrequest = queryrequest
// 	return c
// }

// q := client.Query(
// 	"SELECT name FROM `bigquery-public-data.usa_names.usa_1910_2013` " +
// 			"WHERE state = \"TX\" " +
// 			"LIMIT 100")
// // Location must match that of the dataset(s) referenced in the query.
// q.Location = "US"
// job, err := q.Run(ctx)
// if err != nil {
// 	return err
// }
// status, err := job.Wait(ctx)
// if err != nil {
// 	return err
// }
// if err := status.Err(); err != nil {
// 	return err
// }
// it, err := job.Read(ctx)
// for {
// 	var row []bigquery.Value
// 	err := it.Next(&row)
// 	if err == iterator.Done {
// 			break
// 	}
// 	if err != nil {
// 			return err
// 	}
// 	fmt.Println(row)
// }

// To run this sample, you will need to create (or reuse) a context and
// an instance of the bigquery client.  For example:
// import "cloud.google.com/go/bigquery"
// ctx := context.Background()
// client, err := bigquery.NewClient(ctx, "your-project-id")
// u := client.Dataset(datasetID).Table(tableID).Uploader()
// items := []*Item{
//         // Item implements the ValueSaver interface.
//         {Name: "Phred Phlyntstone", Age: 32},
//         {Name: "Wylma Phlyntstone", Age: 29},
// }
// if err := u.Put(ctx, items); err != nil {
//         return err
// }
