package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

type SqsQueue struct {
	*sqs.SQS
	url string
}

func NewSqsQueue(q *sqs.SQS, url string) *SqsQueue {
	return &SqsQueue{
		SQS: q,
		url: url,
	}
}

func (q *SqsQueue) PublishMessage(item TodoItem) error {
	bytes, err := json.Marshal(item)
	if err != nil {
		return err
	}
	sendParams := &sqs.SendMessageInput{
		MessageBody: aws.String(string(bytes)),
		QueueUrl:    aws.String(q.url),
	}
	result, err := q.SendMessage(sendParams)

	if err != nil {
		return err
	}

	fmt.Println("Success", *result.MessageId)
	return nil
}
// should check this
// https://medium.com/@marcioghiraldelli/elegant-use-of-golang-channels-with-aws-sqs-dad20cd59f34
func (q *SqsQueue) ReceiveMessage() (*TodoItem, error) {
	receiveParams := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(q.url),
	}
	receiveResp, err := q.SQS.ReceiveMessage(receiveParams)
	if err != nil {
		log.Println(err)
	}

	// Delete message
	if len(receiveResp.Messages) > 0 {
		message := receiveResp.Messages[0]
		item := &TodoItem{}
		err = json.NewDecoder(bytes.NewBufferString(*message.Body)).Decode(item)
		if err != nil {
			return nil, err
		}
		// delete message
		deleteParams := &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(q.url),
			ReceiptHandle: message.ReceiptHandle,
		}
		_, err := q.DeleteMessage(deleteParams)
		if err != nil {
			log.Println(err)
		}
		return item, nil
	}
	return nil, errors.New("message is empty")
}
