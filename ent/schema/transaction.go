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
			Comment("ID транзакции (строковый первичный ключ)").
			NotEmpty().
			Immutable().
			StructTag(`json:"id,omitempty"`),
		field.Int("user_id").
			Comment("ID пользователя, которому принадлежит транзакция"),
		field.Float("amount").
			Comment("Сумма транзакции"),

		field.String("currency").
			Default("USD").
			Comment("Валюта транзакции"),

		field.Enum("type").
			Values("deposit", "withdrawal", "transfer").
			Comment("Тип транзакции: пополнение, снятие, перевод"),

		field.String("description").
			Optional().
			Comment("Описание транзакции"),

		field.String("status").
			Default("pending").
			Comment("Статус транзакции: в обработке, выполнено, отклонено и т.д."),

		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("Время создания транзакции"),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Время последнего обновления транзакции"),

		field.Time("completed_at").
			Optional().
			Nillable().
			Comment("Время завершения транзакции"),
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
			Comment("Пользователь, которому принадлежит транзакция"),
	}
}

// Indexes of the Transaction.
func (Transaction) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("created_at"),
		index.Fields("status"),
		index.Fields("user_id", "status"),
	}
}
