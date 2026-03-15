# Download File Activity

This activity allows you to download a file from a specified URL via HTTP and save it locally on your file system. It provides options for setting custom headers, customizing the HTTP method, and defining a timeout.

## Installation

```bash
go get github.com/postqode/qodeflow-component/activity/downloadFile
```

## Schema

Inputs and Outputs:

```json
{
    "input": [
        {
            "name": "url",
            "type": "string",
            "required": true,
            "description": "The URL of the file to download"
        },
        {
            "name": "destination",
            "type": "string",
            "required": true,
            "description": "The local file path or directory to save the downloaded file. If a directory is provided, the original filename from the URL will be used."
        },
        {
            "name": "method",
            "type": "string",
            "required": false,
            "description": "The HTTP method to use for the download request (default: GET)"
        },
        {
            "name": "headers",
            "type": "any",
            "required": false,
            "description": "Custom HTTP headers to include in the request as an object/map"
        },
        {
            "name": "timeout",
            "type": "integer",
            "required": false,
            "description": "Timeout in seconds for the download operation"
        },
        {
            "name": "append",
            "type": "boolean",
            "required": false,
            "description": "If true, appends to existing file. If false (default), overwrites existing file by deleting and recreating it."
        }
    ],
    "output": [
        {
            "name": "success",
            "type": "boolean",
            "description": "Whether the download was successful"
        },
        {
            "name": "error",
            "type": "string",
            "description": "Error message if the download failed"
        },
        {
            "name": "result",
            "type": "any",
            "description": "The result of the operation containing filePath and size in bytes"
        }
    ]
}
```

## Features and Usage

### Direct File Download
Downloads a remote file and saves it using a specific requested filename.
* **Input**: `url`="https://example.com/asset.zip", `destination`="/tmp/my-asset.zip"
* **Output Result**: `{"filePath": "/tmp/my-asset.zip", "size": 1048576}`

### Download to Directory
If the destination path points to an existing directory, the activity infers the filename directly from the URL.
* **Input**: `url`="https://example.com/images/photo.png", `destination`="/tmp/images_folder"
* **Output Result**: `{"filePath": "/tmp/images_folder/photo.png", "size": 2048}`

### Custom Headers
You can customize the HTTP request with specific headers, particularly useful for authentication.
* **Input**: `url`="https://api.example.com/data.csv", `destination`="/tmp/data.csv", `headers`={`Authorization`: `Bearer my-token`}

### File Append vs Overwrite
Control how the activity handles existing files with the same name.

**Overwrite Mode (Default - append=false):**
When a file with the same name exists and `append` is `false` (or not specified), the existing file is deleted and recreated with the downloaded content.
* **Input**: `url`="https://example.com/data.txt", `destination`="/tmp/data.txt", `append`=false
* **Behavior**: If `/tmp/data.txt` exists, it will be deleted and replaced with the newly downloaded content.

**Append Mode (append=true):**
When a file with the same name exists and `append` is `true`, the downloaded content is appended to the existing file.
* **Input**: `url`="https://example.com/data.txt", `destination`="/tmp/data.txt", `append`=true
* **Behavior**: If `/tmp/data.txt` exists, the newly downloaded content is appended to the end of the existing file. If the file doesn't exist, it will be created.

## Testing

To run the internal unit tests:

```bash
go test -v .
```
