package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/firebase/genkit/go/ai"
)

type ToolDefinition struct {
	Name        string
	Description string
}

var ReadFileDefinition = ToolDefinition{
	Name:        "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
}

type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
}

func ReadFile(ctx *ai.ToolContext, input ReadFileInput) (string, error) {
	fmt.Printf("\u001b[92mtool\u001b[0m: %s(%s)\n", ReadFileDefinition.Name, input)

	content, err := os.ReadFile(input.Path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

var ListFilesDescription = ToolDefinition{
	Name:        "list_files",
	Description: "List files and directories at a given path. If no path is provided, lists files in the current directory.",
}

type ListFilesInput struct {
	Path string `json:"path,omitempty" jsonschema_description:"Optional relative path to list files from. Defaults to current directory if not provided."`
}

func ListFiles(ctx *ai.ToolContext, input ListFilesInput) (string, error) {
	fmt.Printf("\u001b[92mtool\u001b[0m: %s(%s)\n", ReadFileDefinition.Name, input)

	dir := "."
	if input.Path != "" {
		dir = input.Path
	}
	var files []string
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if relPath != "." {
			if info.IsDir() {
				files = append(files, relPath+"/")
			} else {
				files = append(files, relPath)
			}
		}
		return nil
	}); err != nil {
		return "", err
	}
	result, err := json.Marshal(files)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

var EditFileDescription = ToolDefinition{
	Name: "edit_file",
	Description: `Make edits to a text file.

Replaces 'old_str' with 'new_str' in the given file. 'old_str' and 'new_str' MUST be different from each other.

If the file specified with path doesn't exist, it will be created.
`,
}

type EditFileInput struct {
	Path   string `json:"path" jsonschema_description:"The path to the file"`
	OldStr string `json:"old_str" jsonschema_description:"Text to search for - must match exactly and must only have one match exactly"`
	NewStr string `json:"new_str" jsonschema_description:"Text to replace old_str with"`
}

func EditFile(ctx *ai.ToolContext, input EditFileInput) (string, error) {
	fmt.Printf("\u001b[92mtool\u001b[0m: %s(%s)\n", ReadFileDefinition.Name, input)

	if input.Path == "" || input.OldStr == input.NewStr {
		return "", fmt.Errorf("invalid input parameters")
	}

	content, err := os.ReadFile(input.Path)
	if err != nil {
		if os.IsNotExist(err) && input.OldStr == "" {
			return createNewFile(input.Path, input.NewStr)
		}
		return "", err
	}

	oldContent := string(content)
	newContent := strings.Replace(oldContent, input.OldStr, input.NewStr, -1)

	if oldContent == newContent && input.OldStr != "" {
		return "", fmt.Errorf("old_str not found in file")
	}

	err = os.WriteFile(input.Path, []byte(newContent), 0644)
	if err != nil {
		return "", err
	}

	return "OK", nil
}

func createNewFile(filePath, content string) (string, error) {
	dir := path.Dir(filePath)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create directory: %v", err)
		}
	}

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}

	return fmt.Sprintf(("Successfully created file %s"), filePath), nil
}
