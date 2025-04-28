package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// EntityGenerator watches schema files and generates entities
type EntityGenerator struct {
	watcher *fsnotify.Watcher
	done    chan bool
}

// NewEntityGenerator creates a new entity generator
func NewEntityGenerator() (*EntityGenerator, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %v", err)
	}

	return &EntityGenerator{
		watcher: watcher,
		done:    make(chan bool),
	}, nil
}

// Start starts watching for schema changes
func (g *EntityGenerator) Start() error {
	// Watch module directory
	err := g.watcher.Add("module")
	if err != nil {
		return fmt.Errorf("failed to watch module directory: %v", err)
	}

	go g.watch()
	return nil
}

// Stop stops watching for schema changes
func (g *EntityGenerator) Stop() {
	g.done <- true
	g.watcher.Close()
}

// watch watches for schema file changes
func (g *EntityGenerator) watch() {
	for {
		select {
		case event := <-g.watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				if strings.HasSuffix(event.Name, ".schema.go") {
					g.generateEntity(event.Name)
				}
			}
		case err := <-g.watcher.Errors:
			fmt.Printf("Error watching files: %v\n", err)
		case <-g.done:
			return
		}
	}
}

// generateEntity generates an entity file from a schema file
func (g *EntityGenerator) generateEntity(schemaPath string) {
	// Read schema file
	content, err := os.ReadFile(schemaPath)
	if err != nil {
		fmt.Printf("Error reading schema file: %v\n", err)
		return
	}

	// Extract package name and struct name
	packageName := filepath.Base(filepath.Dir(schemaPath))
	structName := strings.TrimSuffix(filepath.Base(schemaPath), ".schema.go")
	structName = strings.Title(structName)

	// Parse schema content to extract fields
	schemaContent := string(content)
	fields := extractFields(schemaContent)

	// Generate entity file
	entityPath := filepath.Join(filepath.Dir(schemaPath), structName+".entity.go")
	entityContent := fmt.Sprintf(`package %s

import (
	"time"
)

// %sEntity represents the %s entity
type %sEntity struct {
%s
}

// ToSchema converts entity to schema
func (e *%sEntity) ToSchema() *%sSchema {
	return &%sSchema{
%s
	}
}

// FromSchema converts schema to entity
func (e *%sEntity) FromSchema(s *%sSchema) {
%s
}
`, packageName, structName, packageName, structName, fields, structName, structName, generateToSchemaFields(fields), structName, structName, generateFromSchemaFields(fields))

	// Write entity file
	err = os.WriteFile(entityPath, []byte(entityContent), 0644)
	if err != nil {
		fmt.Printf("Error writing entity file: %v\n", err)
		return
	}

	fmt.Printf("Generated entity file: %s\n", entityPath)
}

// extractFields extracts field definitions from schema content
func extractFields(schemaContent string) string {
	// Find struct definition
	start := strings.Index(schemaContent, "type")
	if start == -1 {
		return ""
	}

	// Find opening brace
	start = strings.Index(schemaContent[start:], "{")
	if start == -1 {
		return ""
	}
	start += len("{")

	// Find closing brace
	end := strings.Index(schemaContent[start:], "}")
	if end == -1 {
		return ""
	}

	// Extract fields
	fields := schemaContent[start : start+end]
	return strings.TrimSpace(fields)
}

// generateToSchemaFields generates field assignments for ToSchema method
func generateToSchemaFields(fields string) string {
	var assignments []string
	for _, field := range strings.Split(fields, "\n") {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		// Extract field name
		parts := strings.Fields(field)
		if len(parts) < 2 {
			continue
		}
		fieldName := parts[0]

		assignments = append(assignments, fmt.Sprintf("\t\t%s: e.%s,", fieldName, fieldName))
	}
	return strings.Join(assignments, "\n")
}

// generateFromSchemaFields generates field assignments for FromSchema method
func generateFromSchemaFields(fields string) string {
	var assignments []string
	for _, field := range strings.Split(fields, "\n") {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		// Extract field name
		parts := strings.Fields(field)
		if len(parts) < 2 {
			continue
		}
		fieldName := parts[0]

		assignments = append(assignments, fmt.Sprintf("\te.%s = s.%s", fieldName, fieldName))
	}
	return strings.Join(assignments, "\n")
} 