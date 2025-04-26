package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Transaction holds the schema definition for the Transaction entity.
type Transaction struct {
	ent.Schema
}

// Fields of the Transaction.
func (Transaction) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Comment("ID of the transaction (string primary key)").
			NotEmpty().
			Immutable().
			StructTag(`json:"id,omitempty"`),

		field.Int("user_id").
			Comment("ID of the user, to which the transaction belongs"),

		field.Float("amount").
			Comment("Amount of the transaction"),

		field.String("currency").
			Default("USD").
			Comment("Currency of the transaction"),

		field.Enum("type").
			Values("deposit", "withdrawal").
			Comment("Type of the transaction: deposit, withdrawal"),

		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("Time of the transaction creation"),
	}
}

// Edges of the Transaction.
func (Transaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("transactions").
			Unique().
			Required().
			Field("user_id").
			Comment("User, to which the transaction belongs"),
	}
}

// Indexes of the Transaction.
func (Transaction) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("created_at"),
	}
}
