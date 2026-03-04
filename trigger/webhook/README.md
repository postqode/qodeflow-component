# Webhook Trigger

The Webhook trigger allows you to trigger Qodeflow flows via HTTP requests. It starts an HTTP server on a specified port and routes incoming requests to the appropriate handlers based on the configured method and path.

## Installation

### Qodeflow CLI

```bash
qodeflow install github.com/postqode/qodeflow-component/trigger/webhook
```

## Configuration

### Settings

| Name | Type | Required | Description |
|------|------|----------|-------------|
| port | int  | true     | The port to listen on for incoming HTTP requests |

### Handler Settings

| Name   | Type   | Required | Description |
|--------|--------|----------|-------------|
| method | string | true     | The HTTP method (GET, POST, PUT, PATCH, DELETE) |
| path   | string | true     | The resource path (e.g., `/api/webhook` or `/user/{id}`) |

### Output

| Name        | Type   | Description |
|-------------|--------|-------------|
| pathParams  | params | The parameters extracted from the path (e.g., `id` from `/user/{id}`) |
| queryParams | params | The query parameters from the URL |
| headers     | params | The HTTP headers from the request |
| method      | string | The HTTP method used for the request |
| content     | any    | The body content of the request (automatically decoded if JSON) |

### Reply (Handler Results)

The handler can return a reply that the trigger uses to respond to the HTTP request:

| Name | Type | Description |
|------|------|-------------|
| code | int  | The HTTP status code to return (e.g., 200, 201) |
| data | any  | The data to return in the response body (sent as JSON) |

## Example Configuration

```json
{
  "id": "webhook-trigger",
  "ref": "github.com/postqode/qodeflow-component/trigger/webhook",
  "settings": {
    "port": 8080
  },
  "handlers": [
    {
      "settings": {
        "method": "POST",
        "path": "/webhook"
      },
      "actions": [
        {
          "id": "log-request",
          "ref": "github.com/postqode/qodeflow-component/activity/log",
          "input": {
            "message": "=string.concat('Received webhook with content: ', $.content)"
          }
        }
      ]
    }
  ]
}
```

## Path Parameters

The Webhook trigger uses `gorilla/mux` for routing, allowing you to define path parameters using the `{name}` syntax. For example, a path of `/user/{id}` will populate `$.pathParams.id` in the trigger output.
