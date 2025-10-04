package database

import (
	"log"

	documentEntities "testcase/internal/modules/document/entities"
	userEntities "testcase/internal/modules/user/entities"
)

type EntityRegistry struct {
	entities []interface{}
}

func NewEntityRegistry() *EntityRegistry {
	return &EntityRegistry{
		entities: make([]interface{}, 0),
	}
}

func (er *EntityRegistry) RegisterEntities() {
	er.addEntity(&userEntities.User{})
	er.addEntity(&documentEntities.Document{})
}

func (er *EntityRegistry) addEntity(entity interface{}) {
	er.entities = append(er.entities, entity)
	log.Printf("📝 Registered entity: %T", entity)
}

func (er *EntityRegistry) GetEntities() []interface{} {
	return er.entities
}

func (er *EntityRegistry) GetEntityCount() int {
	return len(er.entities)
}

func (er *EntityRegistry) RunMigrations(db *Database) error {
	if len(er.entities) == 0 {
		log.Println("⚠️  No entities registered for migration")
		return nil
	}

	log.Printf("🔄 Starting migration for %d entities...", len(er.entities))

	if err := db.AutoMigrate(er.entities...); err != nil {
		return err
	}

	log.Printf("✅ Successfully migrated %d entities", len(er.entities))
	return nil
}

func (er *EntityRegistry) RegisterAndMigrate(db *Database) error {
	log.Println("📋 Registering entities...")
	er.RegisterEntities()

	log.Println("🚀 Running database migrations...")
	return er.RunMigrations(db)
}

func (er *EntityRegistry) ListRegisteredEntities() {
	log.Println("📋 Registered Entities:")
	if len(er.entities) == 0 {
		log.Println("   - No entities registered")
		return
	}

	for i, entity := range er.entities {
		log.Printf("   %d. %T", i+1, entity)
	}
}
