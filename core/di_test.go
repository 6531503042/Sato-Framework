package core

import (
	"testing"
)

type TestService struct {
	Name string
}

func TestContainer(t *testing.T) {
	// Create a new container
	container := NewContainer()

	// Create a test service
	service := &TestService{Name: "test"}

	// Test registering service
	err := container.Register("test", service)
	if err != nil {
		t.Errorf("Failed to register service: %v", err)
	}

	// Test getting service
	result, err := container.Get("test")
	if err != nil {
		t.Errorf("Failed to get service: %v", err)
	}
	if result == nil {
		t.Error("Service is nil")
	}

	// Test getting non-existent service
	_, err = container.Get("non-existent")
	if err == nil {
		t.Error("Expected error when getting non-existent service")
	}

	// Test registering duplicate service
	err = container.Register("test", service)
	if err == nil {
		t.Error("Expected error when registering duplicate service")
	}
}

func TestContainerInject(t *testing.T) {
	// Create a new container
	container := NewContainer()

	// Create a test service
	service := &TestService{Name: "test"}

	// Register service
	err := container.Register("test", service)
	if err != nil {
		t.Errorf("Failed to register service: %v", err)
	}

	// Create a struct with inject tag
	type TestStruct struct {
		Service *TestService `inject:"test"`
	}

	// Create instance
	instance := &TestStruct{}

	// Test injection
	err = container.Inject(instance)
	if err != nil {
		t.Errorf("Failed to inject service: %v", err)
	}

	// Verify injection
	if instance.Service == nil {
		t.Error("Service was not injected")
	}
	if instance.Service.Name != "test" {
		t.Errorf("Expected service name to be 'test', got '%s'", instance.Service.Name)
	}
} 