package cli

import (
	"fmt"
	"os"
	"strings"
)

func GenerateModule(name string) {
	basePath := fmt.Sprintf("module/%s", name)

	dirs := []string{
		basePath + "/handler",
		basePath + "/service",
		basePath + "/repository",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Println("Failed to create dir:", dir)
			panic(err)
		}
	}

	writeFile(basePath+"/handler/"+name+"_controller.go", controllerTemplate(name))
	writeFile(basePath+"/service/"+name+"_service.go", serviceTemplate(name))
	writeFile(basePath+"/repository/"+name+"_repository.go", repositoryTemplate(name))

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

func controllerTemplate(name string) string {
	className := strings.Title(name)
	return fmt.Sprintf(`package handler

import (
	"github.com/6531503042/sato-framework/core"
	"github.com/gofiber/fiber/v2"
)

@core.Controller(core.ControllerOptions{
	RoutePrefix: "/%s",
})
type %sController struct{}

func (c *%sController) Hello(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello from %s Controller")
}
`, name, className, className, className)
}

func serviceTemplate(name string) string {
	className := strings.Title(name)
	return fmt.Sprintf(`package service

	type %sService struct{}

	func New%sService() *%sService {
		return &%sService{}
	}
`, className, className, className, className)
}

func repositoryTemplate(name string) string {
	className := strings.Title(name)
	return fmt.Sprintf(`package repository

type %sRepository struct{}

func New%sRepository() *%sRepository {
	return &%sRepository{}
}
`, className, className, className, className)
}
