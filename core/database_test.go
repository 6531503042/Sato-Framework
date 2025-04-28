package core

import (
	"sync"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestMongoDBProvider(t *testing.T) {
	// Create a new MongoDB provider
	provider := NewMongoDBProvider("mongodb://localhost:27017", "test")

	// Test connection
	err := provider.Connect()
	if err != nil {
		t.Skip("MongoDB not available, skipping test:", err)
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
		t.Skip("MongoDB not available, skipping connection test:", err)
	}

	// Test disconnecting all providers
	err = manager.DisconnectAll()
	if err != nil {
		t.Errorf("Failed to disconnect all providers: %v", err)
	}
}

// MockMongoDBProvider implements a mock MongoDB provider for testing
type MockMongoDBProvider struct {
	client   *mongo.Client
	db       *mongo.Database
	mu       sync.RWMutex
}

func NewMockMongoDBProvider() *MockMongoDBProvider {
	return &MockMongoDBProvider{}
}

func (p *MockMongoDBProvider) Connect() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return nil
}

func (p *MockMongoDBProvider) Disconnect() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return nil
}

func (p *MockMongoDBProvider) GetDB() interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.db
}

func (p *MockMongoDBProvider) GetCollection(name string) *mongo.Collection {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.db == nil {
		return nil
	}
	return p.db.Collection(name)
}

func TestMockMongoDBProvider(t *testing.T) {
	provider := NewMockMongoDBProvider()

	// Test connection
	err := provider.Connect()
	if err != nil {
		t.Errorf("Failed to connect mock provider: %v", err)
	}

	// Test getting database
	db := provider.GetDB()
	if db == nil {
		t.Error("Failed to get mock database")
	}

	// Test getting collection
	collection := provider.GetCollection("test")
	if collection == nil {
		t.Error("Failed to get mock collection")
	}

	// Test disconnection
	err = provider.Disconnect()
	if err != nil {
		t.Errorf("Failed to disconnect mock provider: %v", err)
	}
} 