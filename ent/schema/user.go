package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Unique().Comment("username"),
		field.Int("age").Positive().Comment("age"),
		field.UUID("uid", uuid.UUID{}).Unique().Comment("unique user id"),
		field.Time("created_at").Default(time.Now()).Comment("created at"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("updated at"),
	}
}
