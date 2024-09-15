package initializer

import (
	"context"
	"log"
	"sync"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

var (
	Client   influxdb2.Client
	WriteAPI api.WriteAPIBlocking
	QueryAPI api.QueryAPI
	once     sync.Once
)

func Initialize(url, token, org, bucket string) error {
	var onceErr error
	once.Do(func() {
		Client = influxdb2.NewClient(url, token)

		orgInfo, err := Client.OrganizationsAPI().FindOrganizationByName(context.Background(), org)
		if err != nil {
			orgInfo, err = Client.OrganizationsAPI().CreateOrganizationWithName(context.Background(), org)
			if err != nil {
				log.Fatalf("failed to create organization: %v", err)
				onceErr = err
				return
			}
			log.Printf("Organization created: %s\n", org)
		} else {
			log.Printf("Organization found: %s\n", orgInfo.Name)
		}

		bucketInfo, err := Client.BucketsAPI().FindBucketByName(context.Background(), bucket)
		if err != nil {
			bucketInfo, err = Client.BucketsAPI().CreateBucketWithName(context.Background(), orgInfo, bucket, domain.RetentionRule{})
			if err != nil {
				log.Fatalf("failed to create bucket: %v", err)
				onceErr = err
				return
			}
			log.Printf("Bucket created: %s\n", bucket)
		} else {
			log.Printf("Bucket found: %s\n", bucketInfo.Name)
		}
		WriteAPI = Client.WriteAPIBlocking(org, bucket)
		QueryAPI = Client.QueryAPI(org)
	})
	return onceErr
}
