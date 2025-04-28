package cli

import (
	"fmt"
	"os"
	"strings"
)

func GenerateModule(name string) {
	basePath := fmt.Sprintf("module/%s", name)

	// Create directory structure
	dirs := []string{
		basePath,
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Println("Failed to create dir:", dir)
			panic(err)
		}
	}

	// Generate files
	writeFile(basePath+"/"+name+".schema.go", schemaTemplate(name))
	writeFile(basePath+"/"+name+".controller.go", controllerTemplate(name))
	writeFile(basePath+"/"+name+".service.go", serviceTemplate(name))
	writeFile(basePath+"/"+name+".module.go", moduleTemplate(name))
	writeFile(basePath+"/"+name+".permissions.go", permissionsTemplate(name))

	// Update main.go to include the new module
	updateMainGo(name)

	fmt.Println("✅ Module", name, "generated successfully.")
}

func writeFile(filePath, content string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(content)
}

func schemaTemplate(name string) string {
	className := strings.Title(name)
	return fmt.Sprintf(`package %s

import (
	"time"
)

// %sSchema defines the schema for %s
type %sSchema struct {
	ID        string    `+"`json:\"id\" bson:\"_id,omitempty\"`"+`
	Name      string    `+"`json:\"name\" bson:\"name\"`"+`
	CreatedAt time.Time `+"`json:\"createdAt\" bson:\"createdAt\"`"+`
	UpdatedAt time.Time `+"`json:\"updatedAt\" bson:\"updatedAt\"`"+`
}

// ToEntity converts schema to entity
func (s *%sSchema) ToEntity() *%sEntity {
	return &%sEntity{
		ID:        s.ID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

// FromEntity converts entity to schema
func (s *%sSchema) FromEntity(e *%sEntity) {
	s.ID = e.ID
	s.Name = e.Name
	s.CreatedAt = e.CreatedAt
	s.UpdatedAt = e.UpdatedAt
}
`, name, className, name, className, className, className, className, className, className)
}

func controllerTemplate(name string) string {
	className := strings.Title(name)
	return fmt.Sprintf(`package %s

import (
	"github.com/6531503042/sato-framework/core"
	"github.com/gofiber/fiber/v2"
)

// %sController handles %s-related requests
type %sController struct {
	service *%sService
}

// New%sController creates a new %s controller
func New%sController(service *%sService) *%sController {
	return &%sController{service: service}
}

// Get%s handles GET /%s request
func (c *%sController) Get%s(ctx *fiber.Ctx) error {
	result, err := c.service.Get%s()
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(result)
}

// Create%s handles POST /%s request
func (c *%sController) Create%s(ctx *fiber.Ctx) error {
	var schema %sSchema
	if err := ctx.BodyParser(&schema); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	result, err := c.service.Create%s(schema)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(201).JSON(result)
}
`, name, className, name, className, className, className, name, className, className, className, className, name, className, className, className, className, className)
}

func serviceTemplate(name string) string {
	className := strings.Title(name)
	return fmt.Sprintf(`package %s

import (
	"time"
)

// %sService handles %s business logic
type %sService struct {
	// Add your dependencies here
}

// New%sService creates a new %s service
func New%sService() *%sService {
	return &%sService{}
}

// Get%s retrieves all %s
func (s *%sService) Get%s() ([]%sSchema, error) {
	// TODO: Implement actual database query
	return []%sSchema{
		{
			ID:        "1",
			Name:      "Example 1",
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
		{
			ID:        "2",
			Name:      "Example 2",
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
	}, nil
}

// Create%s creates a new %s
func (s *%sService) Create%s(schema %sSchema) (%sSchema, error) {
	// TODO: Implement actual database save
	schema.ID = "3"
	schema.CreatedAt = time.Now().Unix()
	schema.UpdatedAt = time.Now().Unix()
	return schema, nil
}
`, name, className, name, className, className, name, className, className, className, name, className, className, className, className, name, className, className, className)
}

func moduleTemplate(name string) string {
	className := strings.Title(name)
	return fmt.Sprintf(`package %s

import (
	"github.com/6531503042/sato-framework/core"
)

// %sModule represents the %s module
type %sModule struct {
	Controller *%sController
	Service    *%sService
}

// New%sModule creates a new %s module
func New%sModule() *%sModule {
	service := New%sService()
	controller := New%sController(service)
	return &%sModule{
		Controller: controller,
		Service:    service,
	}
}

// Register registers the module
func (m *%sModule) Register(app *core.App) {
	// Register routes
	app.GetFiber().Get("/%s", m.Controller.Get%s)
	app.GetFiber().Post("/%s", m.Controller.Create%s)
}
`, name, className, name, className, className, className, className, name, className, className, className, className, className, className, name, className, className, className)
}

func permissionsTemplate(name string) string {
	className := strings.Title(name)
	return fmt.Sprintf(`package %s

import "github.com/gofiber/fiber/v2"

// %sPermissions defines permissions for %s
type %sPermissions struct {
	Create string
	Read   string
	Update string
	Delete string
}

// New%sPermissions creates new %s permissions
func New%sPermissions() *%sPermissions {
	return &%sPermissions{
		Create: "%s:create",
		Read:   "%s:read",
		Update: "%s:update",
		Delete: "%s:delete",
	}
}

// CanCreate checks if user can create %s
func (p *%sPermissions) CanCreate(ctx *fiber.Ctx) error {
	// TODO: Implement permission check
	return ctx.Next()
}

// CanRead checks if user can read %s
func (p *%sPermissions) CanRead(ctx *fiber.Ctx) error {
	// TODO: Implement permission check
	return ctx.Next()
}

// CanUpdate checks if user can update %s
func (p *%sPermissions) CanUpdate(ctx *fiber.Ctx) error {
	// TODO: Implement permission check
	return ctx.Next()
}

// CanDelete checks if user can delete %s
func (p *%sPermissions) CanDelete(ctx *fiber.Ctx) error {
	// TODO: Implement permission check
	return ctx.Next()
}
`, name, className, name, className, className, name, className, className, className, name, name, name, name, className, name, className, name, className, name, className, name)
}

func updateMainGo(moduleName string) {
	mainPath := "main.go"
	content, err := os.ReadFile(mainPath)
	if err != nil {
		fmt.Printf("Warning: Could not read main.go: %v\n", err)
		return
	}

	// Add import
	importPath := fmt.Sprintf(`"module/%s"`, moduleName)
	importLine := fmt.Sprintf("\t%s\n", importPath)
	
	// Add module registration
	moduleReg := fmt.Sprintf(`
	// Register %s module
	%sModule := %s.New%sModule()
	%sModule.Register(framework)
`, moduleName, moduleName, moduleName, moduleName, moduleName)

	// Update main.go
	newContent := string(content)
	
	// Add import if not exists
	if !strings.Contains(newContent, importPath) {
		importIndex := strings.Index(newContent, "import (")
		if importIndex != -1 {
			newContent = newContent[:importIndex+8] + importLine + newContent[importIndex+8:]
		}
	}

	// Add module registration if not exists
	if !strings.Contains(newContent, fmt.Sprintf("Register %s module", moduleName)) {
		frameworkIndex := strings.Index(newContent, "framework := core.NewApp()")
		if frameworkIndex != -1 {
			newContent = newContent[:frameworkIndex+len("framework := core.NewApp()")] + moduleReg + newContent[frameworkIndex+len("framework := core.NewApp()"):]
		}
	}

	// Write back to main.go
	err = os.WriteFile(mainPath, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("Warning: Could not update main.go: %v\n", err)
	}
}
