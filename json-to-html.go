package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Data structures to match the JSON
type Report struct {
	ID      string   `json:"id"`
	Results []Result `json:"results"`
}

type Result struct {
	ID         string     `json:"id"`
	Definition Definition `json:"definition"`
	Severity   string     `json:"severity"`
	URL        string     `json:"url"`
	Turns      []Turn     `json:"turns"`
	Evidence   Evidence   `json:"evidence"`
}

type Definition struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Turn struct {
	Request  string `json:"request"`
	Response string `json:"response"`
}

type Evidence struct {
	Type    string          `json:"type"`
	Content EvidenceContent `json:"content"`
}

type EvidenceContent struct {
	Reason string `json:"reason"`
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Snyk AI Red Teaming Report</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 40px 20px;
            color: #333;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        
        .header {
            background: white;
            border-radius: 12px;
            padding: 40px;
            margin-bottom: 30px;
            box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
            display: flex;
            align-items: center;
            gap: 30px;
        }
        
        .logo {
            width: 120px;
            height: 120px;
            object-fit: contain;
        }
        
        .header-content {
            flex: 1;
        }
        
        .header h1 {
            color: #1a1a2e;
            font-size: 36px;
            margin-bottom: 10px;
            font-weight: 700;
        }
        
        .report-id {
            color: #666;
            font-size: 14px;
            font-family: 'Monaco', 'Courier New', monospace;
            background: #f5f5f5;
            padding: 8px 12px;
            border-radius: 6px;
            display: inline-block;
        }
        
        .summary {
            background: white;
            border-radius: 12px;
            padding: 30px;
            margin-bottom: 30px;
            box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
        }
        
        .summary-item {
            text-align: center;
            padding: 20px;
            border-radius: 8px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }
        
        .summary-item h3 {
            font-size: 14px;
            text-transform: uppercase;
            letter-spacing: 1px;
            margin-bottom: 10px;
            opacity: 0.9;
        }
        
        .summary-item .number {
            font-size: 48px;
            font-weight: 700;
        }
        
        .result-card {
            background: white;
            border-radius: 12px;
            padding: 30px;
            margin-bottom: 20px;
            box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
            transition: transform 0.2s, box-shadow 0.2s;
        }
        
        .result-card:hover {
            transform: translateY(-4px);
            box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
        }
        
        .result-header {
            display: flex;
            align-items: start;
            gap: 20px;
            margin-bottom: 20px;
            padding-bottom: 20px;
            border-bottom: 2px solid #f0f0f0;
        }
        
        .severity-badge {
            padding: 8px 16px;
            border-radius: 6px;
            font-weight: 600;
            font-size: 12px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            white-space: nowrap;
        }
        
        .severity-high {
            background: #ff4757;
            color: white;
        }
        
        .severity-medium {
            background: #ffa502;
            color: white;
        }
        
        .severity-low {
            background: #1e90ff;
            color: white;
        }
        
        .result-title {
            flex: 1;
        }
        
        .result-title h2 {
            color: #1a1a2e;
            font-size: 24px;
            margin-bottom: 8px;
        }
        
        .result-description {
            color: #666;
            font-size: 14px;
            line-height: 1.6;
        }
        
        .result-id {
            color: #999;
            font-size: 12px;
            font-family: 'Monaco', 'Courier New', monospace;
            margin-top: 8px;
        }
        
        .section {
            margin-bottom: 25px;
        }
        
        .section-title {
            color: #667eea;
            font-size: 16px;
            font-weight: 700;
            margin-bottom: 15px;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        .section-title::before {
            content: '';
            width: 4px;
            height: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border-radius: 2px;
        }
        
        .turn {
            margin-bottom: 20px;
            border-left: 3px solid #667eea;
            padding-left: 20px;
        }
        
        .turn-label {
            color: #667eea;
            font-weight: 600;
            font-size: 12px;
            text-transform: uppercase;
            letter-spacing: 1px;
            margin-bottom: 8px;
        }
        
        .turn-content {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 8px;
            font-size: 14px;
            line-height: 1.6;
            color: #333;
            white-space: pre-wrap;
            word-wrap: break-word;
            max-height: 300px;
            overflow-y: auto;
        }
        
        .evidence {
            background: #fff8e1;
            border: 2px solid #ffd54f;
            border-radius: 8px;
            padding: 20px;
        }
        
        .evidence-title {
            color: #f57c00;
            font-weight: 700;
            font-size: 14px;
            margin-bottom: 10px;
        }
        
        .evidence-reason {
            color: #555;
            font-size: 14px;
            line-height: 1.6;
        }
        
        .url-link {
            color: #667eea;
            text-decoration: none;
            font-size: 14px;
            word-break: break-all;
        }
        
        .url-link:hover {
            text-decoration: underline;
        }
        
        .footer {
            text-align: center;
            color: white;
            padding: 30px;
            font-size: 14px;
            opacity: 0.9;
        }
        
        .turns-container {
            max-height: 600px;
            overflow-y: auto;
            padding-right: 10px;
        }
        
        .turns-container::-webkit-scrollbar {
            width: 8px;
        }
        
        .turns-container::-webkit-scrollbar-track {
            background: #f1f1f1;
            border-radius: 4px;
        }
        
        .turns-container::-webkit-scrollbar-thumb {
            background: #667eea;
            border-radius: 4px;
        }
        
        .turns-container::-webkit-scrollbar-thumb:hover {
            background: #764ba2;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <img src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTTo9xa5PaelKqiwn19sGlbme64pSprOrzzgTt3OvUc-sLafDhiDLmC6c2Far3PmOxhrN0&usqp=CAU" alt="Snyk Logo" class="logo">
            <div class="header-content">
                <h1>Security Testing Report</h1>
                <div class="report-id">Report ID: {{.ID}}</div>
            </div>
        </div>
        
        <div class="summary">
            <div class="summary-item">
                <h3>Total Findings</h3>
                <div class="number">{{len .Results}}</div>
            </div>
            <div class="summary-item">
                <h3>High Severity</h3>
                <div class="number">{{.HighCount}}</div>
            </div>
            <div class="summary-item">
                <h3>Medium Severity</h3>
                <div class="number">{{.MediumCount}}</div>
            </div>
            <div class="summary-item">
                <h3>Low Severity</h3>
                <div class="number">{{.LowCount}}</div>
            </div>
        </div>
        
        {{range .Results}}
        <div class="result-card">
            <div class="result-header">
                <span class="severity-badge severity-{{.Severity}}">{{.Severity}}</span>
                <div class="result-title">
                    <h2>{{.Definition.Name}}</h2>
                    <p class="result-description">{{.Definition.Description}}</p>
                    <p class="result-id">Definition ID: {{.Definition.ID}}</p>
                    <p class="result-id">Result ID: {{.ID}}</p>
                </div>
            </div>
            
            <div class="section">
                <div class="section-title">Target URL</div>
                <a href="{{.URL}}" class="url-link" target="_blank">{{.URL}}</a>
            </div>
            
            {{if .Turns}}
            <div class="section">
                <div class="section-title">Conversation Turns ({{len .Turns}})</div>
                <div class="turns-container">
                    {{range $index, $turn := .Turns}}
                    <div class="turn">
                        <div class="turn-label">Request {{add $index 1}}</div>
                        <div class="turn-content">{{$turn.Request}}</div>
                    </div>
                    <div class="turn">
                        <div class="turn-label">Response {{add $index 1}}</div>
                        <div class="turn-content">{{$turn.Response}}</div>
                    </div>
                    {{end}}
                </div>
            </div>
            {{end}}
            
            {{if .Evidence.Content.Reason}}
            <div class="section">
                <div class="section-title">Evidence</div>
                <div class="evidence">
                    <div class="evidence-title">Analysis Reason</div>
                    <div class="evidence-reason">{{.Evidence.Content.Reason}}</div>
                </div>
            </div>
            {{end}}
        </div>
        {{end}}
        
        <div class="footer">
            <p>Generated by Snyk Security Testing Tool</p>
            <p>Report generated on {{.GeneratedAt}}</p>
        </div>
    </div>
</body>
</html>`

func main() {
	var jsonData []byte
	var err error
	var outputFile string

	// Check if we're reading from stdin (pipe) or from a file
	if len(os.Args) == 1 {
		// No arguments - read from stdin
		fmt.Fprintln(os.Stderr, "Reading from stdin...")
		jsonData, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("Error reading from stdin: %v", err)
		}
		outputFile = "report.html"
	} else if os.Args[1] == "-" {
		// Explicit stdin with "-"
		fmt.Fprintln(os.Stderr, "Reading from stdin...")
		jsonData, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("Error reading from stdin: %v", err)
		}
		if len(os.Args) >= 3 {
			outputFile = os.Args[2]
		} else {
			outputFile = "report.html"
		}
	} else {
		// Read from file
		inputFile := os.Args[1]
		if len(os.Args) >= 3 {
			outputFile = os.Args[2]
		} else {
			outputFile = "report.html"
		}

		// Sanitize and validate input file path to prevent path traversal
		// 1. Clean the path
		inputFile = filepath.Clean(inputFile)

		// 2. Ensure the file has a .json extension
		if !strings.HasSuffix(strings.ToLower(inputFile), ".json") {
			log.Fatalf("Error: Input file must have a .json extension")
		}

		// 3. Get basename to prevent directory traversal
		baseFile := filepath.Base(inputFile)

		// 4. Ensure no directory separators in basename (additional safety check)
		if strings.Contains(baseFile, "..") || strings.ContainsAny(baseFile, string(filepath.Separator)) {
			log.Fatalf("Error: Invalid file name")
		}

		// 5. Get working directory
		workDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting working directory: %v", err)
		}

		// 6. Construct safe path in current directory
		safeInputPath := filepath.Join(workDir, baseFile)

		// 7. Verify the file exists
		if _, err := os.Stat(safeInputPath); os.IsNotExist(err) {
			log.Fatalf("Error: File does not exist: %s", baseFile)
		}

		// Read JSON file using the sanitized path
		jsonData, err = ioutil.ReadFile(safeInputPath)
		if err != nil {
			log.Fatalf("Error reading JSON file: %v", err)
		}
	}

	// Sanitize output file path
	outputFile = filepath.Clean(outputFile)
	baseOutputFile := filepath.Base(outputFile)
	if !strings.HasSuffix(strings.ToLower(baseOutputFile), ".html") {
		baseOutputFile = baseOutputFile + ".html"
	}
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}
	safeOutputPath := filepath.Join(workDir, baseOutputFile)

	// Parse JSON
	var report Report
	err = json.Unmarshal(jsonData, &report)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Count severities
	highCount := 0
	mediumCount := 0
	lowCount := 0

	for _, result := range report.Results {
		switch strings.ToLower(result.Severity) {
		case "high":
			highCount++
		case "medium":
			mediumCount++
		case "low":
			lowCount++
		}
	}

	// Create template data
	data := map[string]interface{}{
		"ID":          report.ID,
		"Results":     report.Results,
		"HighCount":   highCount,
		"MediumCount": mediumCount,
		"LowCount":    lowCount,
		"GeneratedAt": "2025-11-10",
	}

	// Parse template with custom function
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	tmpl, err := template.New("report").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	// Create output file
	file, err := os.Create(safeOutputPath)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer file.Close()

	// Execute template
	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	fmt.Printf("✓ HTML report generated successfully: %s\n", baseOutputFile)
	fmt.Printf("  Total findings: %d\n", len(report.Results))
	fmt.Printf("  High: %d | Medium: %d | Low: %d\n", highCount, mediumCount, lowCount)
}
