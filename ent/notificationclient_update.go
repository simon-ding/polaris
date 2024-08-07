// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"polaris/ent/notificationclient"
	"polaris/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// NotificationClientUpdate is the builder for updating NotificationClient entities.
type NotificationClientUpdate struct {
	config
	hooks    []Hook
	mutation *NotificationClientMutation
}

// Where appends a list predicates to the NotificationClientUpdate builder.
func (ncu *NotificationClientUpdate) Where(ps ...predicate.NotificationClient) *NotificationClientUpdate {
	ncu.mutation.Where(ps...)
	return ncu
}

// SetName sets the "name" field.
func (ncu *NotificationClientUpdate) SetName(s string) *NotificationClientUpdate {
	ncu.mutation.SetName(s)
	return ncu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ncu *NotificationClientUpdate) SetNillableName(s *string) *NotificationClientUpdate {
	if s != nil {
		ncu.SetName(*s)
	}
	return ncu
}

// SetService sets the "service" field.
func (ncu *NotificationClientUpdate) SetService(s string) *NotificationClientUpdate {
	ncu.mutation.SetService(s)
	return ncu
}

// SetNillableService sets the "service" field if the given value is not nil.
func (ncu *NotificationClientUpdate) SetNillableService(s *string) *NotificationClientUpdate {
	if s != nil {
		ncu.SetService(*s)
	}
	return ncu
}

// SetSettings sets the "settings" field.
func (ncu *NotificationClientUpdate) SetSettings(s string) *NotificationClientUpdate {
	ncu.mutation.SetSettings(s)
	return ncu
}

// SetNillableSettings sets the "settings" field if the given value is not nil.
func (ncu *NotificationClientUpdate) SetNillableSettings(s *string) *NotificationClientUpdate {
	if s != nil {
		ncu.SetSettings(*s)
	}
	return ncu
}

// SetEnabled sets the "enabled" field.
func (ncu *NotificationClientUpdate) SetEnabled(b bool) *NotificationClientUpdate {
	ncu.mutation.SetEnabled(b)
	return ncu
}

// SetNillableEnabled sets the "enabled" field if the given value is not nil.
func (ncu *NotificationClientUpdate) SetNillableEnabled(b *bool) *NotificationClientUpdate {
	if b != nil {
		ncu.SetEnabled(*b)
	}
	return ncu
}

// Mutation returns the NotificationClientMutation object of the builder.
func (ncu *NotificationClientUpdate) Mutation() *NotificationClientMutation {
	return ncu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ncu *NotificationClientUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, ncu.sqlSave, ncu.mutation, ncu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ncu *NotificationClientUpdate) SaveX(ctx context.Context) int {
	affected, err := ncu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ncu *NotificationClientUpdate) Exec(ctx context.Context) error {
	_, err := ncu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ncu *NotificationClientUpdate) ExecX(ctx context.Context) {
	if err := ncu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ncu *NotificationClientUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(notificationclient.Table, notificationclient.Columns, sqlgraph.NewFieldSpec(notificationclient.FieldID, field.TypeInt))
	if ps := ncu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ncu.mutation.Name(); ok {
		_spec.SetField(notificationclient.FieldName, field.TypeString, value)
	}
	if value, ok := ncu.mutation.Service(); ok {
		_spec.SetField(notificationclient.FieldService, field.TypeString, value)
	}
	if value, ok := ncu.mutation.Settings(); ok {
		_spec.SetField(notificationclient.FieldSettings, field.TypeString, value)
	}
	if value, ok := ncu.mutation.Enabled(); ok {
		_spec.SetField(notificationclient.FieldEnabled, field.TypeBool, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ncu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{notificationclient.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ncu.mutation.done = true
	return n, nil
}

// NotificationClientUpdateOne is the builder for updating a single NotificationClient entity.
type NotificationClientUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *NotificationClientMutation
}

// SetName sets the "name" field.
func (ncuo *NotificationClientUpdateOne) SetName(s string) *NotificationClientUpdateOne {
	ncuo.mutation.SetName(s)
	return ncuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ncuo *NotificationClientUpdateOne) SetNillableName(s *string) *NotificationClientUpdateOne {
	if s != nil {
		ncuo.SetName(*s)
	}
	return ncuo
}

// SetService sets the "service" field.
func (ncuo *NotificationClientUpdateOne) SetService(s string) *NotificationClientUpdateOne {
	ncuo.mutation.SetService(s)
	return ncuo
}

// SetNillableService sets the "service" field if the given value is not nil.
func (ncuo *NotificationClientUpdateOne) SetNillableService(s *string) *NotificationClientUpdateOne {
	if s != nil {
		ncuo.SetService(*s)
	}
	return ncuo
}

// SetSettings sets the "settings" field.
func (ncuo *NotificationClientUpdateOne) SetSettings(s string) *NotificationClientUpdateOne {
	ncuo.mutation.SetSettings(s)
	return ncuo
}

// SetNillableSettings sets the "settings" field if the given value is not nil.
func (ncuo *NotificationClientUpdateOne) SetNillableSettings(s *string) *NotificationClientUpdateOne {
	if s != nil {
		ncuo.SetSettings(*s)
	}
	return ncuo
}

// SetEnabled sets the "enabled" field.
func (ncuo *NotificationClientUpdateOne) SetEnabled(b bool) *NotificationClientUpdateOne {
	ncuo.mutation.SetEnabled(b)
	return ncuo
}

// SetNillableEnabled sets the "enabled" field if the given value is not nil.
func (ncuo *NotificationClientUpdateOne) SetNillableEnabled(b *bool) *NotificationClientUpdateOne {
	if b != nil {
		ncuo.SetEnabled(*b)
	}
	return ncuo
}

// Mutation returns the NotificationClientMutation object of the builder.
func (ncuo *NotificationClientUpdateOne) Mutation() *NotificationClientMutation {
	return ncuo.mutation
}

// Where appends a list predicates to the NotificationClientUpdate builder.
func (ncuo *NotificationClientUpdateOne) Where(ps ...predicate.NotificationClient) *NotificationClientUpdateOne {
	ncuo.mutation.Where(ps...)
	return ncuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ncuo *NotificationClientUpdateOne) Select(field string, fields ...string) *NotificationClientUpdateOne {
	ncuo.fields = append([]string{field}, fields...)
	return ncuo
}

// Save executes the query and returns the updated NotificationClient entity.
func (ncuo *NotificationClientUpdateOne) Save(ctx context.Context) (*NotificationClient, error) {
	return withHooks(ctx, ncuo.sqlSave, ncuo.mutation, ncuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ncuo *NotificationClientUpdateOne) SaveX(ctx context.Context) *NotificationClient {
	node, err := ncuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ncuo *NotificationClientUpdateOne) Exec(ctx context.Context) error {
	_, err := ncuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ncuo *NotificationClientUpdateOne) ExecX(ctx context.Context) {
	if err := ncuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ncuo *NotificationClientUpdateOne) sqlSave(ctx context.Context) (_node *NotificationClient, err error) {
	_spec := sqlgraph.NewUpdateSpec(notificationclient.Table, notificationclient.Columns, sqlgraph.NewFieldSpec(notificationclient.FieldID, field.TypeInt))
	id, ok := ncuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "NotificationClient.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ncuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, notificationclient.FieldID)
		for _, f := range fields {
			if !notificationclient.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != notificationclient.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ncuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ncuo.mutation.Name(); ok {
		_spec.SetField(notificationclient.FieldName, field.TypeString, value)
	}
	if value, ok := ncuo.mutation.Service(); ok {
		_spec.SetField(notificationclient.FieldService, field.TypeString, value)
	}
	if value, ok := ncuo.mutation.Settings(); ok {
		_spec.SetField(notificationclient.FieldSettings, field.TypeString, value)
	}
	if value, ok := ncuo.mutation.Enabled(); ok {
		_spec.SetField(notificationclient.FieldEnabled, field.TypeBool, value)
	}
	_node = &NotificationClient{config: ncuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ncuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{notificationclient.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ncuo.mutation.done = true
	return _node, nil
}
