<!--
title: AWS SES
weight: 4621
-->

# AWS SES
This activity allows you to send emails via AWS SES (Simple Email Service).

## Installation

### Qodeflow CLI
```bash
qodeflow install github.com/postqode/qodeflow-component/activity/awsses
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
| Name    | Type   | Required | Description
|:---     | :---   | :---     | :---
| from    | string | Yes      | The sender email address (must be verified in SES)
| to      | array  | Yes      | List of recipient email addresses
| subject | string | No       | The email subject
| body    | string | No       | The email body (HTML supported)
| files   | any    | No       | File attachments (Upcoming)

### Output:
| Name    | Type    | Description
|:---     | :---    | :---
| success | boolean | Whether the email was sent successfully
| error   | string  | Error message if the email failed to send

## Examples
### Simple AWS SES Email
The below example sends a test email:

```json
{
  "id": "send_ses_email",
  "name": "Send SES Email",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/awsses",
    "settings": {
        "region": "us-east-1"
    },
    "input": {
      "from": "sender@yourdomain.com",
      "to": ["recipient@example.com"],
      "subject": "Test Subject",
      "body": "<h1>Hello from AWS SES</h1><p>This is a test email.</p>"
    }
  }
}
```
