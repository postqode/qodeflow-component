# AWS SQS Activity

This activity allows you to perform common queue operations in AWS SQS, including Sending messages, Receiving messages, Deleting messages, and Listing queues.

## Installation

### Qodeflow CLI
```bash
qodeflow install github.com/postqode/qodeflow-component/activity/aws-sqs
```

## Configuration

### Settings:
| Name            | Type   | Required | Description
|:---             | :---   | :---     | :---
| region          | string | Yes      | The AWS region (e.g., us-east-1)
| accessKeyID     | string | No       | The AWS Access Key ID (Optional if using AWS shared credentials)
| secretAccessKey | string | No       | The AWS Secret Access Key (Optional if using AWS shared credentials)
| sessionToken    | string | No       | The AWS Session Token (Optional)

### Input:
| Name                | Type    | Required | Description
|:---                 | :---    | :---     | :---
| method              | string  | Yes      | The SQS operation to perform (`SendMessage`, `ReceiveMessage`, `DeleteMessage`, `ListQueues`)
| queueURL            | string  | No       | The URL of the SQS queue (Required for most operations)
| messageBody         | string  | No       | The body of the message to send (Required for `SendMessage`)
| receiptHandle       | string  | No       | The receipt handle of the message to delete (Required for `DeleteMessage`)
| maxNumberOfMessages | integer | No       | Maximum number of messages to receive (1 to 10, default 1)
| waitTimeSeconds     | integer | No       | The duration (in seconds) for which the call waits for a message to arrive

### Output:
| Name    | Type    | Description
|:---     | :---    | :---
| success | boolean | Whether the operation was successful
| error   | string  | Error message if the operation failed
| result  | any     | The result of the operation (e.g., messageId for Send, list of messages for Receive, list of URLs for List)

## Supported Methods

### SendMessage
Sends a message to the specified SQS queue.
- **Input**: `method: "SendMessage"`, `queueURL`, `messageBody`.
- **Output**: Returns the `messageId` in the `result` field.

### ReceiveMessage
Retrieves one or more messages from the specified queue.
- **Input**: `method: "ReceiveMessage"`, `queueURL`, `maxNumberOfMessages` (optional), `waitTimeSeconds` (optional).
- **Output**: Returns a list of message objects, each containing `messageId`, `body`, and `receiptHandle`.

### DeleteMessage
Deletes a specified message from the queue.
- **Input**: `method: "DeleteMessage"`, `queueURL`, `receiptHandle`.

### ListQueues
Returns a list of your queues in the current region.
- **Input**: `method: "ListQueues"`.
- **Output**: Returns a list of queue URLs in the `result` field.

## Examples

### Send Message
```json
{
  "id": "send_sqs_message",
  "name": "Send SQS Message",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/aws-sqs",
    "settings": {
        "region": "us-east-1"
    },
    "input": {
      "method": "SendMessage",
      "queueURL": "https://sqs.us-east-1.amazonaws.com/123456789012/MyQueue",
      "messageBody": "Hello from Qodeflow!"
    }
  }
}
```

### Receive Messages
```json
{
  "id": "receive_sqs_messages",
  "name": "Receive SQS Messages",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/aws-sqs",
    "settings": {
        "region": "us-east-1"
    },
    "input": {
      "method": "ReceiveMessage",
      "queueURL": "https://sqs.us-east-1.amazonaws.com/123456789012/MyQueue",
      "maxNumberOfMessages": 5
    }
  }
}
```
