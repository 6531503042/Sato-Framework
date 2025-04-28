package cli

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// CRUDTemplateData holds the data for CRUD template generation
type CRUDTemplateData struct {
	Name       string
	TitleName  string
	ModulePath string
}

// GenerateCRUD generates CRUD operations for a given entity
func GenerateCRUD(name string) error {
	data := CRUDTemplateData{
		Name:       name,
		TitleName:  strings.Title(name),
		ModulePath: fmt.Sprintf("module/%s", name),
	}

	// Create module directory structure
	dirs := []string{
		data.ModulePath + "/controller",
		data.ModulePath + "/service",
		data.ModulePath + "/repository",
		data.ModulePath + "/dto",
		data.ModulePath + "/entity",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	// Generate files
	templates := map[string]string{
		fmt.Sprintf("%s/controller/%s_controller.go", data.ModulePath, name): `package controller

import (
	"github.com/gofiber/fiber/v2"
	"{{.ModulePath}}/service"
)

type {{.TitleName}}Controller struct {
	service *service.{{.TitleName}}Service
}

func New{{.TitleName}}Controller(service *service.{{.TitleName}}Service) *{{.TitleName}}Controller {
	return &{{.TitleName}}Controller{service: service}
}

func (c *{{.TitleName}}Controller) Create(ctx *fiber.Ctx) error {
	return ctx.SendString("Create {{.Name}}")
}

func (c *{{.TitleName}}Controller) FindAll(ctx *fiber.Ctx) error {
	return ctx.SendString("Find all {{.Name}}")
}

func (c *{{.TitleName}}Controller) FindOne(ctx *fiber.Ctx) error {
	return ctx.SendString("Find one {{.Name}}")
}

func (c *{{.TitleName}}Controller) Update(ctx *fiber.Ctx) error {
	return ctx.SendString("Update {{.Name}}")
}

func (c *{{.TitleName}}Controller) Delete(ctx *fiber.Ctx) error {
	return ctx.SendString("Delete {{.Name}}")
}`,
		fmt.Sprintf("%s/service/%s_service.go", data.ModulePath, name): `package service

type {{.TitleName}}Service struct {
}

func New{{.TitleName}}Service() *{{.TitleName}}Service {
	return &{{.TitleName}}Service{}
}`,
		fmt.Sprintf("%s/repository/%s_repository.go", data.ModulePath, name): `package repository

type {{.TitleName}}Repository struct {
}

func New{{.TitleName}}Repository() *{{.TitleName}}Repository {
	return &{{.TitleName}}Repository{}
}`,
		fmt.Sprintf("%s/dto/%s_dto.go", data.ModulePath, name): `package dto

type Create{{.TitleName}}Dto struct {
}

type Update{{.TitleName}}Dto struct {
}`,
		fmt.Sprintf("%s/entity/%s_entity.go", data.ModulePath, name): `package entity

type {{.TitleName}} struct {
	ID        uint   ` + "`" + `json:"id"` + "`" + `
	CreatedAt int64  ` + "`" + `json:"createdAt"` + "`" + `
	UpdatedAt int64  ` + "`" + `json:"updatedAt"` + "`" + `
}`,
	}

	for file, tmpl := range templates {
		if err := generateFile(file, tmpl, data); err != nil {
			return err
		}
	}

	return nil
}

func generateFile(path, tmpl string, data CRUDTemplateData) error {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer f.Close()

	if err := t.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
} 