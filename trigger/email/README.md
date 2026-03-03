<!--
title: Email
weight: 4622
-->

# Email Trigger

This trigger listens for incoming emails via IMAP. It polls a configured mailbox folder for new (unseen) messages and triggers the associated handler when new emails arrive.

## Installation

### Qodeflow CLI

```bash
qodeflow install github.com/postqode/qodeflow-component/trigger/email
```

## Configuration

### Settings:

| Name     | Type    | Required | Description                                              |
|----------|---------|----------|----------------------------------------------------------|
| host     | string  | true     | The IMAP server host (e.g., imap.gmail.com)              |
| port     | int     | true     | The IMAP server port (e.g., 993 for TLS)                 |
| username | string  | true     | The email account username                               |
| password | string  | true     | The email account password or app-specific password      |
| useTLS   | boolean | false    | Use TLS for the IMAP connection (default: true for port 993) |

### Handler Settings:

| Name         | Type   | Description                                                  |
|--------------|--------|--------------------------------------------------------------|
| folder       | string | The mailbox folder to monitor (default: INBOX)               |
| pollInterval | string | The polling interval (e.g., 30s, 1m, 5m), defaults to 1m    |

### Output:

| Name    | Type   | Description                      |
|---------|--------|----------------------------------|
| from    | string | The sender email address         |
| to      | string | The recipient email address(es)  |
| subject | string | The email subject                |
| body    | string | The email body content           |
| date    | string | The date the email was received  |

## Example Configuration

```json
{
  "id": "email-trigger",
  "ref": "github.com/postqode/qodeflow-component/trigger/email",
  "settings": {
    "host": "imap.gmail.com",
    "port": 993,
    "username": "your-email@gmail.com",
    "password": "your-app-password",
    "useTLS": true
  },
  "handlers": [
    {
      "settings": {
        "folder": "INBOX",
        "pollInterval": "1m"
      },
      "actions": [
        {
          "id": "log-email",
          "ref": "github.com/postqode/qodeflow-component/activity/log",
          "input": {
            "message": "=string.concat('Received email from: ', $.from, ' Subject: ', $.subject)"
          }
        }
      ]
    }
  ]
}
```

## Notes

- For Gmail, you will need to use an [App Password](https://support.google.com/accounts/answer/185833) rather than your regular account password.
- The trigger marks fetched messages as "seen" by the IMAP server when they are read (standard IMAP behavior).
- Each poll cycle creates a new IMAP connection, fetches unseen messages, and then disconnects.
- Multiple handlers can be configured to monitor different folders with different polling intervals.
