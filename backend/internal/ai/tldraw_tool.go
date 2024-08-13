package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/liushuangls/go-anthropic/v2"
)

type TldrawToolOutput struct {
	Id   string `json:"id"`
	Icon string `json:"icon"`
	Tool string `json:"tool"`
	Util string `json:"util"`
}

func GenTldrawTool(query string) (TldrawToolOutput, error) {
	api_key := os.Getenv("ANTHROPIC_API_KEY")
	if api_key == "" {
		return TldrawToolOutput{}, errors.New("anthropic API Key env var not found")
	}

	anthropic_client := anthropic.NewClient(api_key)

	resp, err := anthropic_client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model:     anthropic.ModelClaude3Dot5Sonnet20240620,
		MaxTokens: 4096,
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(query),
		},
		System: SystemPromptGenTldrawTool,
	})
	if err != nil {
		return TldrawToolOutput{}, err
	}

	log.Println("API Resp", resp.Content[0].GetText())

	tool, err := parseTldrawToolXML(resp.Content[0].GetText())
	if err != nil {
		return TldrawToolOutput{}, err
	}

	writeToolFiles(tool, "/Users/13point5/projects/tlcrazy/frontend")

	return tool, nil
}

type TldrawXML struct {
	Id    string
	Files []TldrawXMLFile
}

type TldrawXMLFile struct {
	Name    string
	Content string
}

func customXMLParser(xmlString string) (TldrawXML, error) {
	var tldraw TldrawXML

	// Find the <tool> tag and extract the id attribute
	toolStart := strings.Index(xmlString, "<tool")
	toolEnd := strings.Index(xmlString, ">")
	if toolStart == -1 || toolEnd == -1 {
		return tldraw, fmt.Errorf("invalid XML: missing <tool> tag")
	}

	// Extract the id attribute from <tool>
	toolTag := xmlString[toolStart:toolEnd]
	idIndex := strings.Index(toolTag, `id="`)
	if idIndex != -1 {
		idStart := idIndex + len(`id="`)
		idEnd := strings.Index(toolTag[idStart:], `"`)
		tldraw.Id = toolTag[idStart : idStart+idEnd]
	}

	// Process each <file> tag
	for {
		fileStart := strings.Index(xmlString, "<file")
		if fileStart == -1 {
			break
		}

		fileEnd := strings.Index(xmlString[fileStart:], ">")
		if fileEnd == -1 {
			return tldraw, fmt.Errorf("invalid XML: malformed <file> tag")
		}
		fileEnd += fileStart

		// Extract the name attribute from <file>
		fileTag := xmlString[fileStart:fileEnd]
		nameIndex := strings.Index(fileTag, `name="`)
		var file TldrawXMLFile
		if nameIndex != -1 {
			nameStart := nameIndex + len(`name="`)
			nameEnd := strings.Index(fileTag[nameStart:], `"`)
			file.Name = fileTag[nameStart : nameStart+nameEnd]
		}

		// Find the closing </file> tag
		fileCloseStart := strings.Index(xmlString[fileEnd:], "</file>")
		if fileCloseStart == -1 {
			return tldraw, fmt.Errorf("invalid XML: missing </file> tag")
		}
		fileCloseEnd := fileCloseStart + len("</file>")
		fileContent := xmlString[fileEnd+1 : fileEnd+fileCloseStart]

		// Assign content and append to list of files
		file.Content = fileContent
		tldraw.Files = append(tldraw.Files, file)

		// Move the cursor forward
		xmlString = xmlString[fileEnd+fileCloseEnd:]
	}

	return tldraw, nil
}

func parseTldrawToolXML(xmlString string) (TldrawToolOutput, error) {
	parsedXML, err := customXMLParser(xmlString)
	if err != nil {
		return TldrawToolOutput{}, err
	}

	out := TldrawToolOutput{Id: parsedXML.Id}

	for _, file := range parsedXML.Files {
		if file.Name == "tool.ts" {
			out.Tool = file.Content
		}

		if file.Name == "util.tsx" {
			out.Util = file.Content
		}

		if file.Name == "icon.svg" {
			out.Icon = file.Content
		}
	}

	return out, nil
}

type WriteFileResult struct {
	path string
	err  error
}

func writeToolFiles(tool TldrawToolOutput, appPath string) (bool, []error) {
	errors := []error{}

	// Check if appPath is valid
	if _, err := os.Stat(appPath); err != nil {
		errors = append(errors, err)
		return false, errors
	}

	toolsJSONPath := filepath.Join(appPath, "components/tldraw-custom-tools/tools.json")
	iconPath := filepath.Join(appPath, "public/custom-tool-icons", fmt.Sprintf("%s.svg", tool.Id))

	toolFolderPath := filepath.Join(appPath, "components/tldraw-custom-tools", tool.Id)
	if err := ensureDirectoryExists(toolFolderPath); err != nil {
		errors = append(errors, err)
		return false, errors
	}

	toolPath := filepath.Join(toolFolderPath, "tool.ts")
	utilPath := filepath.Join(toolFolderPath, "util.tsx")

	files := []string{toolsJSONPath, toolPath, utilPath, iconPath}

	// Write to files concurrently and store errors in a channel
	wg := sync.WaitGroup{}
	resChan := make(chan WriteFileResult, len(files))

	wg.Add(len(files))
	go appendToolId(toolsJSONPath, tool.Id, resChan, &wg)
	go writeToolFile(toolPath, tool.Tool, resChan, &wg)
	go writeToolFile(utilPath, tool.Util, resChan, &wg)
	go writeToolFile(iconPath, tool.Icon, resChan, &wg)
	wg.Wait()

	for range len(files) {
		writeRes := <-resChan
		if writeRes.err != nil {
			errors = append(errors, writeRes.err)
			log.Println("ERROR:", writeRes.err.Error())
		}
	}

	// TODO: undo operations if errors exist

	close(resChan)

	return true, errors
}

func writeToolFile(path, content string, resChan chan WriteFileResult, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println("Writing to", path)

	err := os.WriteFile(path, []byte(content), 0644)
	resChan <- WriteFileResult{path, err}
	if err == nil {
		log.Println("Finished writing to", path)
	}
}

type ToolsFileContent struct {
	Ids []string `json:"ids"`
}

func appendToolId(path, toolId string, resChan chan WriteFileResult, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("Appending Tool ID: %s", toolId)

	// Read file content as string
	fileContent, err := os.ReadFile(path)
	if err != nil {
		resChan <- WriteFileResult{path, err}
		return
	}

	// Parse string as struct
	var data ToolsFileContent
	if err := json.Unmarshal(fileContent, &data); err != nil {
		resChan <- WriteFileResult{path, err}
		return
	}

	// Prevent duplicates
	toolIdIndex := sort.SearchStrings(data.Ids, toolId)
	fmt.Println("data before", data.Ids)
	fmt.Println(toolIdIndex, len(data.Ids))
	if toolIdIndex == len(data.Ids) {
		data.Ids = append(data.Ids, toolId)
	}
	fmt.Println("data after", data.Ids)

	// Convert struct to string
	dataStr, err := json.Marshal(data)
	if err != nil {
		resChan <- WriteFileResult{path, err}
		return
	}

	// Write new data
	err = os.WriteFile(path, []byte(dataStr), 0644)
	resChan <- WriteFileResult{path, err}
	if err == nil {
		log.Printf("Finished appending Tool ID: %s", toolId)
	}
}

func ensureDirectoryExists(toolFolderPath string) error {
	// Check if the directory exists
	if _, err := os.Stat(toolFolderPath); os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.MkdirAll(toolFolderPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
		fmt.Println("Directory created:", toolFolderPath)
	} else if err != nil {
		// Some other error occurred
		return fmt.Errorf("failed to check directory: %v", err)
	} else {
		fmt.Println("Directory already exists:", toolFolderPath)
	}
	return nil
}
