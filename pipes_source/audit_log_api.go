package pipes_source

import (
	"context"
	"fmt"
	"github.com/turbot/pipes-sdk-go"
	"github.com/turbot/tailpipe-plugin-pipes/pipes_collection"
	"github.com/turbot/tailpipe-plugin-sdk/grpc/proto"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/source"
)

// AuditLogAPI source is responsible for collecting audit logs from Turbot Pipes API
type AuditLogAPI struct {
	source.Base
	Config pipes_collection.AuditLogConfig
}

func NewAuditLogAPISource(config pipes_collection.AuditLogConfig) plugin.Source {
	return &AuditLogAPI{
		Config: config,
	}

}

func (s *AuditLogAPI) Collect(ctx context.Context, req *proto.CollectRequest) error {

	// Create a default configuration
	configuration := pipes.NewConfiguration()

	// Add your Turbot Pipes user token as an auth header
	configuration.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", s.Config.Token))

	// Create a client
	client := pipes.NewAPIClient(configuration)

	nextToken := ""

	// TODO - fix me
	orgHandle := "turbot-ops"

	conn := client.GetConfig().Host
	if conn == "" {
		conn = "pipes.turbot.com"
	}
	conn = conn + ":" + orgHandle

	for {
		listReq := client.Orgs.ListAuditLogs(ctx, orgHandle)
		if nextToken != "" {
			listReq = listReq.NextToken(nextToken)
		}

		fmt.Println("Request with NextToken: ", nextToken)

		listReq = listReq.Limit(100)

		response, _, err := listReq.Execute()
		if err != nil {
			// Do something with the error
			panic(err)
		}

		if response.HasItems() {

			fmt.Printf("Response item count: %d\n", len(*response.Items))

			for _, item := range *response.Items {
				s.OnRow(req, conn, item)
			}
		}

		if response.HasNextToken() {
			nextToken = *response.NextToken
		} else {
			break
		}

	}

	fmt.Printf("Done!\n")

	return nil

}
