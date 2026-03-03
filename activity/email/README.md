<!--
title: Send Email
weight: 4620
-->

# Send Email
This activity allows you to send emails via SMTP.

## Installation

### Qodeflow CLI
```bash
qodeflow install github.com/postqode/qodeflow-component/activity/email
```

## Configuration

### Settings:
| Name     | Type   | Description
|:---      | :---   | :---
| host     | string | The SMTP server host
| port     | string | The SMTP server port (e.g., 587)
| username | string | The SMTP server username
| password | string | The SMTP server password

### Input:
| Name    | Type   | Description
|:---     | :---   | :---
| to      | array  | List of recipient email addresses
| subject | string | The email subject
| body    | string | The email body (HTML supported)
| files   | any    | File attachments (supports single file or array of files)

### Output:
| Name    | Type    | Description
|:---     | :---    | :---
| success | boolean | Whether the email was sent successfully
| error   | string  | Error message if the email failed to send

## Examples
### Simple HTML Email
The below example sends a test email:

```json
{
  "id": "send_email",
  "name": "Send Email",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/email",
    "input": {
      "to": ["test@example.com"],
      "subject": "Test Subject",
      "body": "<h1>Hello</h1><p>This is a test email.</p>"
    }
  }
}
```

### Email with Attachment
The below example sends an email with an attachment:

```json
{
  "id": "send_email_attachment",
  "name": "Send Email with Attachment",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/email",
    "input": {
      "to": ["test@example.com"],
      "subject": "Report",
      "body": "<p>Please find the report attached.</p>",
      "files": {
        "filename": "report.pdf",
        "data": "...",
        "mimeType": "application/pdf"
      }
    }
  }
}
```

