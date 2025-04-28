package core

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatabaseProvider defines the interface for database providers
type DatabaseProvider interface {
	Connect() error
	Disconnect() error
	GetDB() interface{}
}

// MySQLProvider implements MySQL database provider
type MySQLProvider struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	db       *sql.DB
}

// NewMySQLProvider creates a new MySQL provider
func NewMySQLProvider(host string, port int, user, password, database string) *MySQLProvider {
	return &MySQLProvider{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
	}
}

// Connect implements MySQL connection
func (p *MySQLProvider) Connect() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", p.User, p.Password, p.Host, p.Port, p.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	p.db = db
	return nil
}

// Disconnect implements MySQL disconnection
func (p *MySQLProvider) Disconnect() error {
	return p.db.Close()
}

// GetDB returns the database connection
func (p *MySQLProvider) GetDB() *sql.DB {
	return p.db
}

// PostgreSQLProvider implements PostgreSQL database provider
type PostgreSQLProvider struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	db       *sql.DB
}

// NewPostgreSQLProvider creates a new PostgreSQL provider
func NewPostgreSQLProvider(host string, port int, user, password, database string) *PostgreSQLProvider {
	return &PostgreSQLProvider{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
	}
}

// Connect implements PostgreSQL connection
func (p *PostgreSQLProvider) Connect() error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.Database)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	p.db = db
	return nil
}

// Disconnect implements PostgreSQL disconnection
func (p *PostgreSQLProvider) Disconnect() error {
	return p.db.Close()
}

// GetDB returns the database connection
func (p *PostgreSQLProvider) GetDB() *sql.DB {
	return p.db
}

// MongoDBProvider implements MongoDB database provider
type MongoDBProvider struct {
	URI      string
	Database string
	client   *mongo.Client
	db       *mongo.Database
	mu       sync.RWMutex
}

// NewMongoDBProvider creates a new MongoDB provider
func NewMongoDBProvider(uri, database string) *MongoDBProvider {
	return &MongoDBProvider{
		URI:      uri,
		Database: database,
	}
}

// Connect implements MongoDB connection
func (p *MongoDBProvider) Connect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(p.URI))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	p.client = client
	p.db = client.Database(p.Database)
	return nil
}

// Disconnect implements MongoDB disconnection
func (p *MongoDBProvider) Disconnect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return p.client.Disconnect(ctx)
	}
	return nil
}

// GetDB returns the MongoDB database instance
func (p *MongoDBProvider) GetDB() interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.db
}

// GetCollection returns a MongoDB collection
func (p *MongoDBProvider) GetCollection(name string) *mongo.Collection {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.db.Collection(name)
}

// DatabaseManager manages database connections
type DatabaseManager struct {
	providers map[string]DatabaseProvider
}

// NewDatabaseManager creates a new database manager
func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{
		providers: make(map[string]DatabaseProvider),
	}
}

// RegisterProvider registers a database provider
func (m *DatabaseManager) RegisterProvider(name string, provider DatabaseProvider) error {
	if _, exists := m.providers[name]; exists {
		return fmt.Errorf("provider %s already exists", name)
	}

	m.providers[name] = provider
	return nil
}

// GetProvider returns a database provider
func (m *DatabaseManager) GetProvider(name string) (DatabaseProvider, error) {
	provider, exists := m.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}

	return provider, nil
}

// ConnectAll connects all registered providers
func (m *DatabaseManager) ConnectAll() error {
	for _, provider := range m.providers {
		if err := provider.Connect(); err != nil {
			return err
		}
	}

	return nil
}

// DisconnectAll disconnects all registered providers
func (m *DatabaseManager) DisconnectAll() error {
	for _, provider := range m.providers {
		if err := provider.Disconnect(); err != nil {
			return err
		}
	}

	return nil
} 