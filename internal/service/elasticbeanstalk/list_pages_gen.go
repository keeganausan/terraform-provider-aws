// Code generated by "internal/generate/listpages/main.go -ListOps=DescribeEnvironments"; DO NOT EDIT.

package elasticbeanstalk

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk/elasticbeanstalkiface"
)

func describeEnvironmentsPages(ctx context.Context, conn elasticbeanstalkiface.ElasticBeanstalkAPI, input *elasticbeanstalk.DescribeEnvironmentsInput, fn func(*elasticbeanstalk.EnvironmentDescriptionsMessage, bool) bool) error {
	for {
		output, err := conn.DescribeEnvironmentsWithContext(ctx, input)
		if err != nil {
			return err
		}

		lastPage := aws.StringValue(output.NextToken) == ""
		if !fn(output, lastPage) || lastPage {
			break
		}

		input.NextToken = output.NextToken
	}
	return nil
}
