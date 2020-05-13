package main

import (
	"applike/pkg"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os"
	"strconv"
)

func main() {
	// sqs
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Endpoint: aws.String(os.Getenv("SQS_ENDPOINT")),
	})
	svc := sqs.New(sess)
	_, err := svc.SetQueueAttributes(&sqs.SetQueueAttributesInput{
		QueueUrl: aws.String(os.Getenv("SQS_QUEUE_URL")),
		Attributes: aws.StringMap(map[string]string{
			"ReceiveMessageWaitTimeSeconds": strconv.Itoa(3),
		}),
	})
	if err != nil {
		_ = fmt.Errorf("Unable to update queue %q, %v.", os.Getenv("SQS_QUEUE_NAME"), err)
	}
	queue := pkg.NewSqsQueue(svc, os.Getenv("SQS_QUEUE_URL"))

	for {
		item, err := queue.ReceiveMessage()
		if err != nil {
			log.Print(err)
		}
		log.Printf("Message id=%d description=%s time=%v:", item.Id, item.Description, item.DueDate)
	}
}
