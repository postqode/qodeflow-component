# WriteFile Activity

This activity allows you to perform basic file system write operations. It supports writing new files, appending to existing files, deleting files, and creating directories.

## Installation

```bash
go get github.com/postqode/qodeflow-component/activity/writeFile
```

## Schema
Inputs and Outputs:

```json
{
    "input": [
        {
            "name": "method",
            "type": "string",
            "required": true,
            "description": "The file operation to perform (WriteFile, AppendFile, DeleteFile, CreateDirectory)"
        },
        {
            "name": "filePath",
            "type": "string",
            "required": true,
            "description": "The path to the file or directory"
        },
        {
            "name": "content",
            "type": "any",
            "required": false,
            "description": "The content to write or append to the file"
        }
    ],
    "output": [
        {
            "name": "success",
            "type": "boolean",
            "description": "Whether the operation was successful"
        },
        {
            "name": "error",
            "type": "string",
            "description": "Error message if the operation failed"
        },
        {
            "name": "result",
            "type": "any",
            "description": "The result of the operation"
        }
    ]
}
```

## Supported Operations

### WriteFile
Creates a new file or overwrites an existing file with the provided `content`.
* **Input**: `method`="WriteFile", `filePath`="/path/to/file.txt", `content`="Data to write"
* **Output Result**: `{"filePath": "/path/to/file.txt", "size": 13}`

### AppendFile
Appends the provided `content` to an existing file, creating it if it doesn't exist.
* **Input**: `method`="AppendFile", `filePath`="/path/to/file.txt", `content`="\nMore data"
* **Output Result**: `{"filePath": "/path/to/file.txt", "bytesWritten": 10}`

### DeleteFile
Deletes the file at the specified `filePath`.
* **Input**: `method`="DeleteFile", `filePath`="/path/to/file.txt"
* **Output Result**: `{"filePath": "/path/to/file.txt", "deleted": true}`

### CreateDirectory
Creates a directory and any necessary parents at the specified `filePath`.
* **Input**: `method`="CreateDirectory", `filePath`="/path/to/new_dir"
* **Output Result**: `{"dirPath": "/path/to/new_dir", "created": true}`

## Testing

To run the unit tests:
```bash
go test -v .
```

To run integration tests against actual files, use environment variables:
```bash
WRITE_FILE_PATH=/tmp/qftest.txt WRITE_DIR_PATH=/tmp/qfdir go test -v -run TestWriteFileActivity_Integration
```
