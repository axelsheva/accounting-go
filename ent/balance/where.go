// Code generated by ent, DO NOT EDIT.

package balance

import (
	"accounting/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Balance {
	return predicate.Balance(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Balance {
	return predicate.Balance(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Balance {
	return predicate.Balance(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Balance {
	return predicate.Balance(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Balance {
	return predicate.Balance(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Balance {
	return predicate.Balance(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Balance {
	return predicate.Balance(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Balance {
	return predicate.Balance(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Balance {
	return predicate.Balance(sql.FieldLTE(FieldID, id))
}

// UserID applies equality check predicate on the "user_id" field. It's identical to UserIDEQ.
func UserID(v int) predicate.Balance {
	return predicate.Balance(sql.FieldEQ(FieldUserID, v))
}

// Currency applies equality check predicate on the "currency" field. It's identical to CurrencyEQ.
func Currency(v string) predicate.Balance {
	return predicate.Balance(sql.FieldEQ(FieldCurrency, v))
}

// Amount applies equality check predicate on the "amount" field. It's identical to AmountEQ.
func Amount(v float64) predicate.Balance {
	return predicate.Balance(sql.FieldEQ(FieldAmount, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldEQ(FieldUpdatedAt, v))
}

// UserIDEQ applies the EQ predicate on the "user_id" field.
func UserIDEQ(v int) predicate.Balance {
	return predicate.Balance(sql.FieldEQ(FieldUserID, v))
}

// UserIDNEQ applies the NEQ predicate on the "user_id" field.
func UserIDNEQ(v int) predicate.Balance {
	return predicate.Balance(sql.FieldNEQ(FieldUserID, v))
}

// UserIDIn applies the In predicate on the "user_id" field.
func UserIDIn(vs ...int) predicate.Balance {
	return predicate.Balance(sql.FieldIn(FieldUserID, vs...))
}

// UserIDNotIn applies the NotIn predicate on the "user_id" field.
func UserIDNotIn(vs ...int) predicate.Balance {
	return predicate.Balance(sql.FieldNotIn(FieldUserID, vs...))
}

// CurrencyEQ applies the EQ predicate on the "currency" field.
func CurrencyEQ(v string) predicate.Balance {
	return predicate.Balance(sql.FieldEQ(FieldCurrency, v))
}

// CurrencyNEQ applies the NEQ predicate on the "currency" field.
func CurrencyNEQ(v string) predicate.Balance {
	return predicate.Balance(sql.FieldNEQ(FieldCurrency, v))
}

// CurrencyIn applies the In predicate on the "currency" field.
func CurrencyIn(vs ...string) predicate.Balance {
	return predicate.Balance(sql.FieldIn(FieldCurrency, vs...))
}

// CurrencyNotIn applies the NotIn predicate on the "currency" field.
func CurrencyNotIn(vs ...string) predicate.Balance {
	return predicate.Balance(sql.FieldNotIn(FieldCurrency, vs...))
}

// CurrencyGT applies the GT predicate on the "currency" field.
func CurrencyGT(v string) predicate.Balance {
	return predicate.Balance(sql.FieldGT(FieldCurrency, v))
}

// CurrencyGTE applies the GTE predicate on the "currency" field.
func CurrencyGTE(v string) predicate.Balance {
	return predicate.Balance(sql.FieldGTE(FieldCurrency, v))
}

// CurrencyLT applies the LT predicate on the "currency" field.
func CurrencyLT(v string) predicate.Balance {
	return predicate.Balance(sql.FieldLT(FieldCurrency, v))
}

// CurrencyLTE applies the LTE predicate on the "currency" field.
func CurrencyLTE(v string) predicate.Balance {
	return predicate.Balance(sql.FieldLTE(FieldCurrency, v))
}

// CurrencyContains applies the Contains predicate on the "currency" field.
func CurrencyContains(v string) predicate.Balance {
	return predicate.Balance(sql.FieldContains(FieldCurrency, v))
}

// CurrencyHasPrefix applies the HasPrefix predicate on the "currency" field.
func CurrencyHasPrefix(v string) predicate.Balance {
	return predicate.Balance(sql.FieldHasPrefix(FieldCurrency, v))
}

// CurrencyHasSuffix applies the HasSuffix predicate on the "currency" field.
func CurrencyHasSuffix(v string) predicate.Balance {
	return predicate.Balance(sql.FieldHasSuffix(FieldCurrency, v))
}

// CurrencyEqualFold applies the EqualFold predicate on the "currency" field.
func CurrencyEqualFold(v string) predicate.Balance {
	return predicate.Balance(sql.FieldEqualFold(FieldCurrency, v))
}

// CurrencyContainsFold applies the ContainsFold predicate on the "currency" field.
func CurrencyContainsFold(v string) predicate.Balance {
	return predicate.Balance(sql.FieldContainsFold(FieldCurrency, v))
}

// AmountEQ applies the EQ predicate on the "amount" field.
func AmountEQ(v float64) predicate.Balance {
	return predicate.Balance(sql.FieldEQ(FieldAmount, v))
}

// AmountNEQ applies the NEQ predicate on the "amount" field.
func AmountNEQ(v float64) predicate.Balance {
	return predicate.Balance(sql.FieldNEQ(FieldAmount, v))
}

// AmountIn applies the In predicate on the "amount" field.
func AmountIn(vs ...float64) predicate.Balance {
	return predicate.Balance(sql.FieldIn(FieldAmount, vs...))
}

// AmountNotIn applies the NotIn predicate on the "amount" field.
func AmountNotIn(vs ...float64) predicate.Balance {
	return predicate.Balance(sql.FieldNotIn(FieldAmount, vs...))
}

// AmountGT applies the GT predicate on the "amount" field.
func AmountGT(v float64) predicate.Balance {
	return predicate.Balance(sql.FieldGT(FieldAmount, v))
}

// AmountGTE applies the GTE predicate on the "amount" field.
func AmountGTE(v float64) predicate.Balance {
	return predicate.Balance(sql.FieldGTE(FieldAmount, v))
}

// AmountLT applies the LT predicate on the "amount" field.
func AmountLT(v float64) predicate.Balance {
	return predicate.Balance(sql.FieldLT(FieldAmount, v))
}

// AmountLTE applies the LTE predicate on the "amount" field.
func AmountLTE(v float64) predicate.Balance {
	return predicate.Balance(sql.FieldLTE(FieldAmount, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Balance {
	return predicate.Balance(sql.FieldLTE(FieldUpdatedAt, v))
}

// HasUser applies the HasEdge predicate on the "user" edge.
func HasUser() predicate.Balance {
	return predicate.Balance(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserWith applies the HasEdge predicate on the "user" edge with a given conditions (other predicates).
func HasUserWith(preds ...predicate.User) predicate.Balance {
	return predicate.Balance(func(s *sql.Selector) {
		step := newUserStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Balance) predicate.Balance {
	return predicate.Balance(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Balance) predicate.Balance {
	return predicate.Balance(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Balance) predicate.Balance {
	return predicate.Balance(sql.NotPredicates(p))
}
