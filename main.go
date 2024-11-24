package main

import (
	"bufio"
	"cli/github"
	"fmt"
	"gofr.dev/pkg/gofr"
	vertex_ai "gofr.dev/pkg/gofr/ai/vertex-ai"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const serviceAccountJSON = ``

func main() {
	app := gofr.NewCMD()

	creds := serviceAccountJSON

	vertexAIConfigs := &vertex_ai.Configs{
		ProjectID:         "endless-fire-437206-j7",
		LocationID:        "us-central1",
		APIEndpoint:       "us-central1-aiplatform.googleapis.com",
		ModelID:           "gemini-1.5-pro-002",
		Credentials:       creds,
		Datastore:         "projects/endless-fire-437206-j7/locations/global/collections/default_collection/dataStores/gofr-datastore_1732298621027",
		SystemInstruction: "The response has to be verified from the documentation. \n\nUpon giving the code in the response make sure that the code has to be given only from a single document, it should not be combined from various documents.-> ensure this to be done effectively.\n\nThe prompt should be responded only from the documents\n\nUpon giving the same prompt at the many time make sure to give the relevant response for the prompt without any changes in the response. \n\nDon't mention the gs:// bucket paths anywhere in the response\n\n\nGive the code only from a single document dont give it by combining with various documents.\n\nEnable langchain to verify it with the previous response and answer in the sessions.",
	}

	vertexAIClient, err := vertex_ai.NewVertexAIClientWithKey(vertexAIConfigs)
	if err != nil {
		app.Logger().Fatalf("failed to create vertex AI client: %v", err)
	}

	app.UseVertexAI(vertexAIClient)

	app.SubCommand("example hello-world", func(c *gofr.Context) (interface{}, error) {
		prompt := "hi can you generate boiler plate for hello world api code for gofr with directory structure"

		fmt.Println("generating response....")
		response, err := c.VertexAI.SendMessageUsingSystemInstruction([]map[string]string{{"text": prompt, "role": "user"}}, []string{"use documentation and example and dont mix multiple files just use one and please have file name on top of every code or file data"})
		if err != nil {
			return nil, err
		}

		_ = extractStructure(response)

		return "example go code generated", nil
	})

	app.SubCommand("example redis", func(c *gofr.Context) (interface{}, error) {
		prompt := "can you give sample api for redis connection in gofr"

		fmt.Println("generating response....")
		response, err := c.VertexAI.SendMessageUsingSystemInstruction([]map[string]string{{"text": prompt, "role": "user"}}, []string{"use example in http-server-using-redis.pdf and can you please start explanation lines with ** or *"})
		if err != nil {
			return nil, err
		}

		_ = extractStructure(response)

		return "example go code generated", nil
	})

	app.SubCommand("example sql", func(c *gofr.Context) (interface{}, error) {
		prompt := "can you give sample api for sql connection in gofr"

		fmt.Println("generating response....")
		response, err := c.VertexAI.SendMessageUsingSystemInstruction([]map[string]string{{"text": prompt, "role": "user"}}, []string{"use example in using-migrations.pdf fo go code generation and can you please start explanation lines with ** or *"})
		if err != nil {
			return nil, err
		}

		_ = extractStructure(response)

		return "example go code generated", nil
	})

	app.SubCommand("issue raise", github.RaiseIssue)
	app.SubCommand("issue list", github.GetIssue)
	app.SubCommand("article list", func(c *gofr.Context) (interface{}, error) {
		prompt := "can you list article written for gofr with their links"

		fmt.Println("generating response....")
		response, err := c.VertexAI.SendMessageUsingSystemInstruction([]map[string]string{{"text": prompt, "role": "user"}}, []string{"use gofr article file to fetch info and give response and name - description - link, dont write anything else, format response for printing as a cli output"})
		if err != nil {
			return nil, err
		}

		return response, nil
	})

	app.SubCommand("article describe", func(c *gofr.Context) (interface{}, error) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("enter article name/title you want to get info about")
		article, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error getting article input %v", err)
		}
		prompt := fmt.Sprintf("can you describe this article written for gofr with link %s", article)

		fmt.Println("generating response....")
		response, err := c.VertexAI.SendMessageUsingSystemInstruction([]map[string]string{{"text": prompt, "role": "user"}}, []string{"use gofr article file to fetch info if not present any article related to prompt print no article found, format response for printing as a cli output"})
		if err != nil {
			return nil, err
		}

		return response, nil
	})

	app.SubCommand("release list", github.GetRelease)

	app.SubCommand("release latest", github.GetReleaseLatest)

	app.SubCommand("check error", func(c *gofr.Context) (interface{}, error) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("enter the error you want to check")
		errorString, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error gertting error %v", err)
		}

		fmt.Println("generating response...")

		prompt := fmt.Sprintf("can you check this error related to gofr and provide solution %s also go through issues and give detail if found something similar to this error just provide detail of the issue and prs", errorString)

		response, err := c.VertexAI.SendMessageUsingSystemInstruction([]map[string]string{{"text": prompt, "role": "user"}}, []string{"check all_github_issues.pdf and all_github_prs.pdf in gofr-issues title and body or description if found the error in any of them just give detail about issue and pr it like link for the same and format ir correctly, also add a promt to raise and issue on gofr.dev repo if not able to solve it"})
		if err != nil {
			return nil, err
		}

		return response, nil
	})

	app.SubCommand("analyze", func(c *gofr.Context) (interface{}, error) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("enter go file path you want to analyze")
		filePath, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error gertting file input %v", err)
		}

		filePath = strings.TrimSuffix(filePath, "\n")
		// File path

		// Read the entire file content
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("error gertting file %v", err)
		}

		fmt.Println("analyzing code...")
		prompt := fmt.Sprintf("can you please analyze this go code and provide suggettions to improve this%s", string(data))

		response, err := c.VertexAI.SendMessageUsingSystemInstruction([]map[string]string{{"text": prompt, "role": "user"}}, []string{"analyze code based on go basic rules and gofr and describe logic dont give code for analyzing this"})
		if err != nil {
			return nil, err
		}

		return response, nil
	})

	// Run the application
	app.Run()
}

func extractStructure(response string) []string {
	scanner := bufio.NewScanner(os.Stdin)

	basePath := ""
	fmt.Println("enter project name")

	// Scan the next line from the input
	if scanner.Scan() {
		basePath = scanner.Text() // Get the input string
	}

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}

	basePath = "./" + basePath

	// Create the directories
	err := os.MkdirAll(basePath, os.ModePerm) // os.ModePerm gives permissions 0777
	if err != nil {
		log.Fatalf("Error creating directories: %v\n", err)
	}
	// Prepare containers
	files := []string{}
	docs := []string{}
	codes := map[string]string{}

	// Iterate over lines
	lines := strings.Split(response, "\n")
	var currentCodeFile string
	inCodeBlock := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Detect start of a code block
		if strings.HasPrefix(line, "package ") {
			inCodeBlock = true
			currentCodeFile = "main.go" // Assuming this based on context, adjust as needed
			codes[currentCodeFile] = line + "\n"
			continue
		}

		// Detect end of a code block
		if inCodeBlock {
			if strings.HasPrefix(line, "```") {
				inCodeBlock = false
			} else {
				codes[currentCodeFile] += line + "\n"
			}
			continue
		}

		if strings.Contains(line, "* ") || strings.Contains(line, "**") || strings.HasPrefix(line, "//") {
			docs = append(docs, line)
		}
	}

	codes["go.mod"] = `module project-name

go 1.22

require gofr.dev v1.27.1
`
	codes["configs/.env"] = `HTTP_PORT=9000

REDIS_HOST=
REDIS_PORT=

DB_HOST=
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_PORT=
DB_DIALECT=
DB_CHARSET=`

	// Create the directories
	err = os.MkdirAll(basePath+"/configs", os.ModePerm) // os.ModePerm gives permissions 0777
	if err != nil {
		log.Fatalf("Error creating directories: %v\n", err)
	}

	for _, file := range files {
		fullPath := filepath.Join(basePath, file)

		// Check if the path is for a directory or file
		if strings.HasSuffix(file, "/") {
			// Create a directory
			if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
				fmt.Printf("Error creating directory %s: %v\n", fullPath, err)
				continue
			}
			fmt.Printf("Created directory: %s\n", fullPath)
		} else {
			// Create a file
			dir := filepath.Dir(fullPath)
			// Ensure the parent directory exists
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				fmt.Printf("Error creating parent directory %s: %v\n", dir, err)
				continue
			}
			// Create the file
			if _, err := os.Create(fullPath); err != nil {
				fmt.Printf("Error creating file %s: %v\n", fullPath, err)
				continue
			}
			fmt.Printf("Created file: %s\n", fullPath)
		}
	}

	// Step 2: Write code into respective files
	for fileName, code := range codes {
		filePath := filepath.Join(basePath, fileName)

		// Write code to the file
		if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
			log.Fatalf("Error writing code to file %s: %v\n", filePath, err)
		}

		fmt.Printf("Added code to: %s\n", filePath)
	}

	// Step 3: Create a documentation file
	docFilePath := filepath.Join(basePath, "docs.txt")
	docContent := strings.Join(docs, "\n")
	if err := os.WriteFile(docFilePath, []byte(docContent), 0644); err != nil {
		fmt.Printf("Error creating documentation file %s: %v\n", docFilePath, err)
	} else {
		fmt.Printf("Created documentation file: %s\n", docFilePath)
	}

	chmodCmd := exec.Command("chmod", "+x", basePath)
	_, err = chmodCmd.CombinedOutput()
	if err != nil {
		fmt.Println("failed change permission for project", err)
	}

	cmd := exec.Command("go", "get", "gofr.dev")
	cmd.Path = basePath
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("")
	}

	tidycmd := exec.Command("go", "mod", "tidy")
	tidycmd.Path = basePath
	_, err = tidycmd.CombinedOutput()
	if err != nil {
		fmt.Println("")
	}

	return nil
}

func parseFilePath(line string) string {
	line = strings.TrimPrefix(line, "├──")
	line = strings.TrimPrefix(line, "└──")
	line = strings.TrimPrefix(line, "│   └──")
	line = strings.TrimSpace(line)
	return line
}
