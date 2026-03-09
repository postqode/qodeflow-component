# Read File Activity

This activity allows you to perform common file system operations, including reading file content, checking file existence, listing directory entries, and retrieving file metadata.

## Installation

### Qodeflow CLI
```bash
qodeflow install github.com/postqode/qodeflow-component/activity/readFile
```

## Configuration

This activity requires no settings — it operates directly on the local file system.

### Input:
| Name     | Type   | Required | Description                                                                             |
|:---      | :---   | :---     | :---                                                                                    |
| method   | string | Yes      | The file operation to perform (`ReadFile`, `FileExists`, `ListDirectory`, `GetFileInfo`) |
| filePath | string | Yes      | The path to the file or directory                                                        |
| encoding | string | No       | The file encoding to use when reading (default: `utf-8`)                                 |

### Output:
| Name    | Type    | Description                                              |
|:---     | :---    | :---                                                     |
| success | boolean | Whether the operation was successful                     |
| error   | string  | Error message if the operation failed                    |
| result  | any     | The result of the operation (varies by method)           |

## Supported Methods

### ReadFile
Reads the entire content of a file.
- **Input**: `method: "ReadFile"`, `filePath`
- **Output**: Returns `content` (string), `filePath`, and `size` (bytes) in the `result` field.

### FileExists
Checks whether a file or directory exists at the given path.
- **Input**: `method: "FileExists"`, `filePath`
- **Output**: Returns `exists` (boolean) and `filePath` in the `result` field. Does **not** fail if the path does not exist.

### ListDirectory
Lists all entries inside a directory.
- **Input**: `method: "ListDirectory"`, `filePath` (must be a directory path)
- **Output**: Returns `dirPath`, `count` (number of entries), and `entries` (array of objects with `name`, `isDir`, `size`, `modTime`).

### GetFileInfo
Retrieves metadata for a file or directory.
- **Input**: `method: "GetFileInfo"`, `filePath`
- **Output**: Returns `name`, `size`, `isDir`, `modTime` (RFC3339), and `mode` in the `result` field.

## Examples

### Read a File
```json
{
  "id": "read_config_file",
  "name": "Read Config File",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/readFile",
    "input": {
      "method": "ReadFile",
      "filePath": "/etc/app/config.json"
    }
  }
}
```

### Check If a File Exists
```json
{
  "id": "check_file_exists",
  "name": "Check File Exists",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/readFile",
    "input": {
      "method": "FileExists",
      "filePath": "/var/data/report.csv"
    }
  }
}
```

### List a Directory
```json
{
  "id": "list_logs_dir",
  "name": "List Logs Directory",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/readFile",
    "input": {
      "method": "ListDirectory",
      "filePath": "/var/log/myapp"
    }
  }
}
```

### Get File Info
```json
{
  "id": "get_file_info",
  "name": "Get File Info",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/readFile",
    "input": {
      "method": "GetFileInfo",
      "filePath": "/var/data/report.csv"
    }
  }
}
```
