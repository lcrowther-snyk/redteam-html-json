# JSON to HTML Report Generator

A Go script that converts Snyk Red Team CLI output into beautiful HTML reports with Snyk-themed styling.

## Features

- 🎨 Beautiful Snyk-themed design with gradient backgrounds
- 📊 Summary dashboard showing severity counts
- 🔍 Detailed view of each security finding
- 📱 Responsive design that works on all devices
- 🎯 Easy-to-read conversation turns with scrollable sections
- 🔗 Clickable links to target URLs

## Quick Start (macOS Apple Silicon)

**Just download the binary and use it directly:**

1. Download the pre-built binary: [json-to-html](https://github.com/lcrowther-snyk/readteam-html-json/raw/main/json-to-html)
2. Make it executable: `chmod +x json-to-html`
3. Run it: `./json-to-html results.json report.html`

```bash
# Download the binary
curl -LO https://github.com/lcrowther-snyk/readteam-html-json/raw/main/json-to-html
chmod +x json-to-html

# Use with Snyk Red Team
snyk redteam --experimental > results.json
./json-to-html results.json report.html
```

> **Note:** Replace `OWNER/REPO` with your actual GitHub username and repository name.

## Prerequisites

- [Snyk CLI](https://docs.snyk.io/snyk-cli/install-the-snyk-cli) installed and authenticated

## Installation

### Option 1: Use Pre-built Binary (macOS Apple Silicon)

**Method A: Direct Download**

Download just the binary file (no need to clone the entire repository):

```bash
curl -LO https://github.com/OWNER/REPO/raw/main/json-to-html
chmod +x json-to-html
```

**Method B: Clone Repository**

Or clone the full repository:

```bash
git clone https://github.com/OWNER/REPO.git
cd redteam-html-json
./json-to-html results.json report.html
```

### Option 2: Build from Source

If you're on a different platform or want to build from source:

**Requirements:** Go 1.16 or higher

```bash
# For your current platform
go build -o json-to-html json-to-html.go

# Or for specific platforms:
# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o json-to-html json-to-html.go

# Linux
GOOS=linux GOARCH=amd64 go build -o json-to-html json-to-html.go

# Windows
GOOS=windows GOARCH=amd64 go build -o json-to-html.exe json-to-html.go
```

## Usage

### Primary Usage: With Snyk Red Team CLI

The tool is designed to work directly with the Snyk Red Team experimental CLI:

```bash
# Run Snyk Red Team and generate HTML report in one command
snyk redteam --experimental > results.json
./json-to-html results.json report.html
```

Or pipe directly (if your shell supports it):

```bash
snyk redteam --experimental | tee results.json | ./json-to-html results.json report.html
```

### Standalone Usage

If you already have a JSON report file:

```bash
# Generate report with default name (report.html)
./json-to-html results.json

# Generate report with custom output name
./json-to-html results.json my-security-report.html
```

### Using with `go run`

If you prefer not to build the binary:

```bash
snyk redteam --experimental > results.json
go run json-to-html.go results.json report.html
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

## Workflow Tips

### Complete Red Team Testing Workflow

```bash
# 1. Clone the repository (if not already done)
git clone <repository-url>
cd redteam-html-json

# 2. Run Snyk Red Team assessment
snyk redteam --experimental > results.json

# 3. Generate HTML report (using pre-built binary)
./json-to-html results.json report.html

# 4. Open report in browser
open report.html  # macOS
# or
xdg-open report.html  # Linux
# or
start report.html  # Windows
```

## Security Notes

- The tool validates input files to prevent path traversal attacks
- Only .json files from the current directory are processed
- All file paths are sanitized before use
- Successfully scanned with Snyk Code with 0 vulnerabilities

## Requirements

- Go 1.16 or higher

## License

MIT

