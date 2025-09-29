# Excelize JSON to Excel API

This project provides a simple HTTP API to convert a JSON file into an Excel (.xlsx) file. It is built with Go and uses the [excelize](https://github.com/xuri/excelize) library for efficient Excel generation.

## Features
- Accepts a POST request with a JSON file at the `/excel` endpoint
- Converts the JSON data to an Excel file and returns it as a download
- Handles large files efficiently using streaming

## JSON Structure
The JSON file should be an object where each key is a row number (as a string), and each value is an array of objects mapping column letters to cell values. Example:

```json
{
  "1": [
    { "A": "Name", "B": "Age", "C": "Address" }
  ],
  "2": [
    { "A": "Alice", "B": 30, "C": "123 Main St" }
  ],
  "3": [
    { "A": "Bob", "B": 25, "C": "456 Oak Ave" }
  ]
}
```
- Each key (e.g., "1", "2", "3") is the row number in Excel.
- Each value is an array of objects, where each object maps column letters (A, B, C, ...) to cell values.
- The first row typically contains headers.

## API Usage

### Endpoint
`POST /excel`

### Request
- Content-Type: `multipart/form-data`
- Field name: `file`
- Value: The JSON file to convert

#### Example using `curl`:
```sh
curl -X POST http://localhost:3000/excel \
  -F "file=@my.json" \
  -o result.xlsx
```

### Response
- Content-Type: `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
- Content-Disposition: `attachment; filename="excel_TIMESTAMP.xlsx"`
- Body: The generated Excel file

## Example

#### Request JSON (`my.json`):
```json
{
  "1": [ { "A": "Name", "B": "Age", "C": "Address" } ],
  "2": [ { "A": "Alice", "B": 30, "C": "123 Main St" } ],
  "3": [ { "A": "Bob", "B": 25, "C": "456 Oak Ave" } ]
}
```

#### Response
- Downloaded file: `excel_TIMESTAMP.xlsx`
- The Excel file will have:

| Name  | Age | Address      |
|-------|-----|--------------|
| Alice | 30  | 123 Main St  |
| Bob   | 25  | 456 Oak Ave  |

## Running the Project

1. Build and run the server:
   ```sh
   go run main.go
   ```
2. Send a POST request as shown above.

## Dummy JSON Generator Endpoint

### Endpoint
`GET /dummy`

### Description
This endpoint generates and serves a large dummy JSON file (up to 500,000 rows and 10 columns) for testing purposes. The file is saved in the `temp/` folder and returned as a download.

> **Caution:**
> - The generated file can be very large (hundreds of MBs or more).
> - Downloading or processing this file may consume significant disk, memory, and network resources.
> - Use only for testing and benchmarking. Not recommended for production or low-resource environments.

### Usage Example
```sh
curl -O http://localhost:3000/dummy
```

The downloaded file will be `temp/dummy.json` and can be used as input for the `/excel` endpoint to test large-scale Excel generation.

## License
MIT
