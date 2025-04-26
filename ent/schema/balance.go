package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Balance holds the schema definition for the Balance entity.
type Balance struct {
	ent.Schema
}

// Fields of the Balance.
func (Balance) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Comment("ID of the balance").
			StructTag(`json:"id,omitempty"`),

		field.Int("user_id").
			Comment("ID of the user, to which the balance belongs"),

		field.String("currency").
			NotEmpty().
			Comment("Currency code (e.g. USD, EUR, RUB)"),

		field.Float("amount").
			Default(0).
			Comment("Amount of the balance in the specified currency"),

		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("Time of the balance creation"),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Time of the last balance update"),
	}
}

// Edges of the Balance.
func (Balance) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("balances").
			Unique().
			Required().
			Field("user_id").
			Comment("User, to which the balance belongs"),
	}
}

// Indexes of the Balance.
func (Balance) Indexes() []ent.Index {
	return []ent.Index{
		// Index for fast search of balances by currency
		index.Fields("currency"),

		// Index for fast search of balances by user_id
		index.Fields("user_id"),

		// Combined index for uniqueness of the combination of user+currency
		index.Fields("user_id", "currency").
			Unique(),
	}
}

func (Balance) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Adds a named CHECK constraint
		entsql.Checks(map[string]string{
			"balance_amount_non_negative": "amount >= 0",
		}),
	}
}
