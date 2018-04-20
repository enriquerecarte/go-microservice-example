package messaging

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

const (
	QueueUrl    = "https://sqs.ap-northeast-1.amazonaws.com/3**********2/my-queue"
	Region      = "eu-west-1"
	CredPath    = "/Users/home/.aws/credentials"
	CredProfile = "aws-cred-profile"
	Endpoint    = "http://localhost:9324"
)

func SendMessage() {
	sqsSession := newSession()

	sqsConnection := sqs.New(sqsSession)

	queue, err := queueUrl(sqsConnection, "local-bacs_gateway-paymentsubmitted")

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("Success", *queue.QueueUrl)

	// Send message
	messageInput := &sqs.SendMessageInput{
		MessageBody:  aws.String("message body"),
		QueueUrl:     queue.QueueUrl,
		DelaySeconds: aws.Int64(3),
	}

	messageOutput, err := sqsConnection.SendMessage(messageInput)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("[Send message] \n%v \n\n", messageOutput)

}
func queueUrl(svc *sqs.SQS, queueName string) (*sqs.GetQueueUrlOutput, error) {
	return svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
}
func newSession() *session.Session {
	return session.New(&aws.Config{
		Region:   aws.String(Region),
		Endpoint: aws.String(Endpoint),
		Credentials: credentials.NewStaticCredentials(CredPath, CredProfile),
		MaxRetries: aws.Int(5),
	})
}

func ReceiveMessage() {
	// Receive message
	sqsSession := newSession()

	sqsConnection := sqs.New(sqsSession)

	queue, err := queueUrl(sqsConnection, "local-bacs_gateway-paymentsubmitted")

	receiveParams := &sqs.ReceiveMessageInput{
		QueueUrl:            queue.QueueUrl,
		MaxNumberOfMessages: aws.Int64(3),
		VisibilityTimeout:   aws.Int64(30),
		WaitTimeSeconds:     aws.Int64(5),
	}
	receiveResponse, err := sqsConnection.ReceiveMessage(receiveParams)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("[Receive message] \n%v \n\n", receiveResponse)

	// Delete message
	for _, message := range receiveResponse.Messages {
		deleteParams := &sqs.DeleteMessageInput{
			QueueUrl:      queue.QueueUrl,
			ReceiptHandle: message.ReceiptHandle,
		}
		_, err := sqsConnection.DeleteMessage(deleteParams) // No response returned when successed.
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("[Delete message] \nMessage ID: %s has beed deleted.\n\n", *message.MessageId)
	}

}
