# AWS S3 Activity

This activity allows you to perform common object storage operations in AWS S3, including Upload, Download, List, and Delete.

## Installation

### Qodeflow CLI
```bash
qodeflow install github.com/postqode/qodeflow-component/activity/aws-s3
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
| Name      | Type   | Required | Description
|:---       | :---   | :---     | :---
| method    | string | Yes      | The S3 operation to perform (`Upload`, `Download`, `List`, `Delete`)
| bucket    | string | Yes      | The S3 bucket name
| key       | string | Yes      | The S3 object key (or prefix for `List` operation)
| data      | any    | No       | The data to upload (required for `Upload`)
| localPath | string | No       | The local file path to save to (optional for `Download`)

### Output:
| Name    | Type    | Description
|:---     | :---    | :---
| success | boolean | Whether the operation was successful
| error   | string  | Error message if the operation failed
| result  | any     | The result of the operation (e.g., list of keys for `List`, or file content for `Download` if `localPath` is not provided)

## Supported Methods

### Upload
Uploads data to the specified bucket and key.
- **Input**: `method: "Upload"`, `bucket`, `key`, `data` (string, []byte, or io.Reader).

### Download
Downloads an object from the specified bucket and key.
- **Input**: `method: "Download"`, `bucket`, `key`, `localPath` (optional).
- **Output**: If `localPath` is provided, the file is saved locally. If not, the content is returned in the `result` field as a string.

### List
Lists objects in the specified bucket starting with the prefix provided in the `key` field.
- **Input**: `method: "List"`, `bucket`, `key` (used as prefix).
- **Output**: A list of strings containing the keys of the objects found.

### Delete
Deletes an object from the specified bucket and key.
- **Input**: `method: "Delete"`, `bucket`, `key`.

## Examples

### Upload File to S3
```json
{
  "id": "upload_s3_file",
  "name": "Upload S3 File",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/aws-s3",
    "settings": {
        "region": "us-east-1"
    },
    "input": {
      "method": "Upload",
      "bucket": "my-cool-bucket",
      "key": "documents/hello.txt",
      "data": "Hello World from Qodeflow!"
    }
  }
}
```

### List Files in Prefix
```json
{
  "id": "list_s3_files",
  "name": "List S3 Files",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/aws-s3",
    "settings": {
        "region": "us-east-1"
    },
    "input": {
      "method": "List",
      "bucket": "my-cool-bucket",
      "key": "documents/"
    }
  }
}
```
 Salmon
