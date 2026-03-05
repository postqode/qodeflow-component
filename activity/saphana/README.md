<!--
title: SAP HANA
weight: 4640
-->

# SAP HANA

This activity integrates with SAP HANA to execute SQL queries, DML statements (INSERT, UPDATE, DELETE), and stored procedure calls.

## Installation

### Qodeflow CLI
```bash
qodeflow install github.com/postqode/qodeflow-component/activity/saphana
```

## Configuration

### Settings:
| Name     | Type    | Required | Description |
|:---      | :---    | :---     | :--- |
| dsn      | string  | no¹      | Full SAP HANA DSN (e.g. `hdb://user:password@host:39017`). If provided, the individual fields below are ignored. |
| host     | string  | no¹      | SAP HANA server hostname or IP address |
| port     | integer | no       | SAP HANA server port (default: `39017`) |
| user     | string  | no¹      | Database username |
| password | string  | no       | Database password |

> ¹ Either `dsn` **or** both `host` and `user` must be provided.

### Input:
| Name   | Type   | Required | Description |
|:---    | :---   | :---     | :--- |
| method | string | yes      | Operation to perform: `QUERY`, `EXEC`, or `CALL` |
| query  | string | yes      | SQL statement or stored procedure call |
| args   | array  | no       | Positional arguments bound to `?` placeholders in the query |

#### Method reference

| Method  | Use for | Typical SQL |
|:---     | :---    | :--- |
| `QUERY` | Read rows | `SELECT ...` |
| `EXEC`  | Modify data / DDL | `INSERT`, `UPDATE`, `DELETE`, `CREATE TABLE`, etc. |
| `CALL`  | Stored procedures | `CALL my_procedure(?, ?)` |

### Output:
| Name         | Type    | Description |
|:---          | :---    | :--- |
| result       | array   | Rows returned by `QUERY` or `CALL` (each row is a key-value map keyed by column name). Empty for `EXEC`. |
| rowsAffected | integer | Number of rows affected (`EXEC`) or number of rows returned (`QUERY`/`CALL`). |

## Examples

### Query rows
```json
{
  "id": "get_orders",
  "name": "Get Orders",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/saphana",
    "input": {
      "method": "QUERY",
      "query": "SELECT ORDER_ID, STATUS FROM ORDERS WHERE CUSTOMER_ID = ?",
      "args": ["C001"]
    }
  }
}
```

### Insert a row
```json
{
  "id": "insert_event",
  "name": "Insert Event",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/saphana",
    "input": {
      "method": "EXEC",
      "query": "INSERT INTO AUDIT_LOG (EVENT, TS, USER) VALUES (?, NOW(), ?)",
      "args": ["login", "admin"]
    }
  }
}
```

### Update a row
```json
{
  "id": "update_status",
  "name": "Update Status",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/saphana",
    "input": {
      "method": "EXEC",
      "query": "UPDATE ORDERS SET STATUS = ? WHERE ORDER_ID = ?",
      "args": ["SHIPPED", "ORD-42"]
    }
  }
}
```

### Call a stored procedure
```json
{
  "id": "call_proc",
  "name": "Call Procedure",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/saphana",
    "input": {
      "method": "CALL",
      "query": "CALL GET_CUSTOMER_ORDERS(?)",
      "args": ["C001"]
    }
  }
}
```

## Testing

This activity includes unit tests (no database required) and integration tests (require a running SAP HANA instance).

### Unit tests
```bash
cd activity/saphana
go test -v -run "TestActivity_Metadata|TestActivity_New|TestActivity_Eval_MissingQuery" ./...
```

### Integration tests
Provide the `SAPHANA_DSN` environment variable pointing to your SAP HANA instance:

```bash
SAPHANA_DSN="hdb://SYSTEM:YourPassword@localhost:39017" go test -v ./activity/saphana/...
```

### Local setup (SAP HANA Express via Docker)
You can run a local SAP HANA Express edition for testing:

```bash
docker pull saplabs/hanaexpress:latest

docker run -d \
  --name hana-express \
  -p 39013:39013 -p 39015:39015 -p 39017:39017 -p 39041-39045:39041-39045 \
  -e AGREE_TO_SAP_LICENSE=true \
  -e MASTER_PASSWORD=MyStr0ngPass \
  saplabs/hanaexpress:latest

# Wait ~2 minutes for HANA to start, then run tests:
SAPHANA_DSN="hdb://SYSTEM:MyStr0ngPass@localhost:39017" go test -v ./activity/saphana/...
