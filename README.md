# JSON to HTML Report Generator

A Go script that converts JSON security testing reports into beautiful HTML reports with Snyk-themed styling.

## Features

- 🎨 Beautiful Snyk-themed design with gradient backgrounds
- 📊 Summary dashboard showing severity counts
- 🔍 Detailed view of each security finding
- 📱 Responsive design that works on all devices
- 🎯 Easy-to-read conversation turns with scrollable sections
- 🔗 Clickable links to target URLs

## Usage

### Basic Usage

```bash
go run json-to-html.go <input.json>
```

This will generate a `report.html` file in the current directory.

### Custom Output File

```bash
go run json-to-html.go <input.json> <output.html>
```

### Example

```bash
go run json-to-html.go test.json report.html
```

## Building

To build a standalone executable:

```bash
go build -o json-to-html json-to-html.go
```

Then run:

```bash
./json-to-html test.json report.html
```

## JSON Format

The script expects JSON files in the following format:

```json
{
  "id": "report-id",
  "results": [
    {
      "id": "result-id",
      "definition": {
        "id": "definition-id",
        "name": "Finding Name",
        "description": "Finding description"
      },
      "severity": "high",
      "url": "http://example.com",
      "turns": [
        {
          "request": "Request text",
          "response": "Response text"
        }
      ],
      "evidence": {
        "type": "json",
        "content": {
          "reason": "Evidence reason"
        }
      }
    }
  ]
}
```

## Design

The HTML report features:
- Snyk's signature purple gradient background
- Clean white cards with subtle shadows
- Color-coded severity badges (High: Red, Medium: Orange, Low: Blue)
- Smooth hover effects and transitions
- Professional typography with system fonts
- Scrollable sections for long content

## Requirements

- Go 1.16 or higher

## License

MIT

