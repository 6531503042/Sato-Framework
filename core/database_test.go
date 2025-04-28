package core

import (
	"testing"
)

func TestMongoDBProvider(t *testing.T) {
	// Create a new MongoDB provider
	provider := NewMongoDBProvider("mongodb://localhost:27017", "test")

	// Test connection
	err := provider.Connect()
	if err != nil {
		t.Errorf("Failed to connect to MongoDB: %v", err)
	}

	// Test getting database
	db := provider.GetDB()
	if db == nil {
		t.Error("Failed to get database")
	}

	// Test getting collection
	collection := provider.GetCollection("test")
	if collection == nil {
		t.Error("Failed to get collection")
	}

	// Test disconnection
	err = provider.Disconnect()
	if err != nil {
		t.Errorf("Failed to disconnect from MongoDB: %v", err)
	}
}

func TestDatabaseManager(t *testing.T) {
	// Create a new database manager
	manager := NewDatabaseManager()

	// Create a MongoDB provider
	mongoProvider := NewMongoDBProvider("mongodb://localhost:27017", "test")

	// Test registering provider
	err := manager.RegisterProvider("mongo", mongoProvider)
	if err != nil {
		t.Errorf("Failed to register provider: %v", err)
	}

	// Test getting provider
	provider, err := manager.GetProvider("mongo")
	if err != nil {
		t.Errorf("Failed to get provider: %v", err)
	}
	if provider == nil {
		t.Error("Provider is nil")
	}

	// Test connecting all providers
	err = manager.ConnectAll()
	if err != nil {
		t.Errorf("Failed to connect all providers: %v", err)
	}

	// Test disconnecting all providers
	err = manager.DisconnectAll()
	if err != nil {
		t.Errorf("Failed to disconnect all providers: %v", err)
	}
} 