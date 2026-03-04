<!--
title: MongoDB
weight: 4630
-->

# MongoDB
This activity integrates with MongoDB to perform common database operations like GET, INSERT, UPDATE, DELETE, and REPLACE.

## Installation

### Qodeflow CLI
```bash
qodeflow install github.com/postqode/qodeflow-component/activity/mongodb
```

## Configuration

### Settings:
| Name       | Type   | Description
|:---        | :---   | :---
| uri        | string | The MongoDB connection URI (e.g., mongodb://localhost:27017)
| dbName     | string | The name of the database
| collection | string | The name of the collection

### Input:
| Name     | Type   | Description
|:---      | :---   | :---
| method   | string | The operation to perform (GET, DELETE, INSERT, REPLACE, UPDATE)
| keyName  | string | The name of the key field for filters (e.g., _id or email)
| keyValue | string | The value of the key field for filters
| data     | any    | The data object for INSERT, REPLACE, or UPDATE operations

### Output:
| Name   | Type    | Description
|:---    | :---    | :---
| output | any     | The result of the operation (document for GET, ID for INSERT)
| count  | integer | The number of documents affected (for DELETE, UPDATE, REPLACE)

## Examples

### Get a Document
```json
{
  "id": "get_user",
  "name": "Get User",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/mongodb",
    "input": {
      "method": "GET",
      "keyName": "email",
      "keyValue": "user@example.com"
    }
  }
}
```

### Insert a Document
```json
{
  "id": "insert_log",
  "name": "Insert Log",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/mongodb",
    "input": {
      "method": "INSERT",
      "data": {
        "event": "login",
        "timestamp": "2024-03-04T00:00:00Z",
        "user": "admin"
      }
    }
  }
}
```

### Update a Document
```json
{
  "id": "update_status",
  "name": "Update Status",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/mongodb",
    "input": {
      "method": "UPDATE",
      "keyName": "_id",
      "keyValue": "60d5ec491234567890abcdef",
      "data": {
        "status": "active"
      }
    }
  }
}
```

## Testing

This activity includes functional integration tests that require a running MongoDB instance.

### Run Tests
To run the tests, provide the `MONGODB_URI` environment variable:

```bash
MONGODB_URI="mongodb://localhost:27017" go test -v ./activity/mongodb/...
```

### Local Setup (Docker)
You can quickly start a local MongoDB instance for testing using Docker:

```bash
docker run --name mongodb-test -p 27017:27017 -d mongodb/mongodb-community-server:latest
```
