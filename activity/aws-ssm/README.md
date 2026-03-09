# AWS SSM Activity

This activity allows you to perform common Parameter Store operations in AWS Systems Manager (SSM), including Getting, Putting, Deleting, and Listing parameters.

## Installation

### Qodeflow CLI
```bash
qodeflow install github.com/postqode/qodeflow-component/activity/aws-ssm
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
| Name            | Type    | Required | Description
|:---             | :---    | :---     | :---
| method          | string  | Yes      | The SSM operation to perform (`GetParameter`, `PutParameter`, `DeleteParameter`, `GetParametersByPath`, `DescribeParameters`)
| parameterName   | string  | No       | The name of the parameter (Required for Get, Put, and Delete)
| parameterValue  | string  | No       | The value of the parameter (Required for `PutParameter`)
| parameterType   | string  | No       | The type of parameter (`String`, `StringList`, `SecureString`) (Optional for `PutParameter`)
| path            | string  | No       | The path to parameters (Required for `GetParametersByPath`)
| recursive       | boolean | No       | Whether to fetch parameters recursively (Optional for `GetParametersByPath`)
| withDecryption  | boolean | No       | Whether to decrypt `SecureString` parameters (Optional for Get and GetByPath)
| overwrite       | boolean | No       | Whether to overwrite an existing parameter (Optional for `PutParameter`)

### Output:
| Name    | Type    | Description
|:---     | :---    | :---
| success | boolean | Whether the operation was successful
| error   | string  | Error message if the operation failed
| result  | any     | The result of the operation (details of fetched, created, or listed parameters)

## Supported Methods

### GetParameter
Retrieves information about a specific parameter.
- **Input**: `method: "GetParameter"`, `parameterName`, `withDecryption` (optional).
- **Output**: Returns the parameter `name`, `value`, and `type` in the `result` field.

### PutParameter
Creates a new parameter or updates an existing one.
- **Input**: `method: "PutParameter"`, `parameterName`, `parameterValue`, `parameterType` (optional, default `String`), `overwrite` (optional).
- **Output**: Returns the new `version` of the parameter in the `result` field.

### DeleteParameter
Deletes a specified parameter.
- **Input**: `method: "DeleteParameter"`, `parameterName`.

### GetParametersByPath
Retrieves all parameters in a hierarchy.
- **Input**: `method: "GetParametersByPath"`, `path`, `recursive` (optional), `withDecryption` (optional).
- **Output**: Returns a list of parameters, each containing `name`, `value`, and `type`.

### DescribeParameters
Returns a list of all parameters.
- **Input**: `method: "DescribeParameters"`.
- **Output**: Returns a list of parameter names, their `type`, and `description`.

## Examples

### Get Parameter
```json
{
  "id": "get_ssm_parameter",
  "name": "Get SSM Parameter",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/aws-ssm",
    "settings": {
        "region": "us-east-1"
    },
    "input": {
      "method": "GetParameter",
      "parameterName": "/my/app/config",
      "withDecryption": true
    }
  }
}
```

### Put Parameter
```json
{
  "id": "put_ssm_parameter",
  "name": "Put SSM Parameter",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/aws-ssm",
    "settings": {
        "region": "us-east-1"
    },
    "input": {
      "method": "PutParameter",
      "parameterName": "/my/app/config",
      "parameterValue": "new_value",
      "parameterType": "String",
      "overwrite": true
    }
  }
}
```
