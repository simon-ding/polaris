// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"polaris/ent/history"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// HistoryCreate is the builder for creating a History entity.
type HistoryCreate struct {
	config
	mutation *HistoryMutation
	hooks    []Hook
}

// SetMediaID sets the "media_id" field.
func (hc *HistoryCreate) SetMediaID(i int) *HistoryCreate {
	hc.mutation.SetMediaID(i)
	return hc
}

// SetEpisodeNums sets the "episode_nums" field.
func (hc *HistoryCreate) SetEpisodeNums(i []int) *HistoryCreate {
	hc.mutation.SetEpisodeNums(i)
	return hc
}

// SetSeasonNum sets the "season_num" field.
func (hc *HistoryCreate) SetSeasonNum(i int) *HistoryCreate {
	hc.mutation.SetSeasonNum(i)
	return hc
}

// SetNillableSeasonNum sets the "season_num" field if the given value is not nil.
func (hc *HistoryCreate) SetNillableSeasonNum(i *int) *HistoryCreate {
	if i != nil {
		hc.SetSeasonNum(*i)
	}
	return hc
}

// SetSourceTitle sets the "source_title" field.
func (hc *HistoryCreate) SetSourceTitle(s string) *HistoryCreate {
	hc.mutation.SetSourceTitle(s)
	return hc
}

// SetDate sets the "date" field.
func (hc *HistoryCreate) SetDate(t time.Time) *HistoryCreate {
	hc.mutation.SetDate(t)
	return hc
}

// SetTargetDir sets the "target_dir" field.
func (hc *HistoryCreate) SetTargetDir(s string) *HistoryCreate {
	hc.mutation.SetTargetDir(s)
	return hc
}

// SetSize sets the "size" field.
func (hc *HistoryCreate) SetSize(i int) *HistoryCreate {
	hc.mutation.SetSize(i)
	return hc
}

// SetNillableSize sets the "size" field if the given value is not nil.
func (hc *HistoryCreate) SetNillableSize(i *int) *HistoryCreate {
	if i != nil {
		hc.SetSize(*i)
	}
	return hc
}

// SetDownloadClientID sets the "download_client_id" field.
func (hc *HistoryCreate) SetDownloadClientID(i int) *HistoryCreate {
	hc.mutation.SetDownloadClientID(i)
	return hc
}

// SetNillableDownloadClientID sets the "download_client_id" field if the given value is not nil.
func (hc *HistoryCreate) SetNillableDownloadClientID(i *int) *HistoryCreate {
	if i != nil {
		hc.SetDownloadClientID(*i)
	}
	return hc
}

// SetIndexerID sets the "indexer_id" field.
func (hc *HistoryCreate) SetIndexerID(i int) *HistoryCreate {
	hc.mutation.SetIndexerID(i)
	return hc
}

// SetNillableIndexerID sets the "indexer_id" field if the given value is not nil.
func (hc *HistoryCreate) SetNillableIndexerID(i *int) *HistoryCreate {
	if i != nil {
		hc.SetIndexerID(*i)
	}
	return hc
}

// SetLink sets the "link" field.
func (hc *HistoryCreate) SetLink(s string) *HistoryCreate {
	hc.mutation.SetLink(s)
	return hc
}

// SetNillableLink sets the "link" field if the given value is not nil.
func (hc *HistoryCreate) SetNillableLink(s *string) *HistoryCreate {
	if s != nil {
		hc.SetLink(*s)
	}
	return hc
}

// SetHash sets the "hash" field.
func (hc *HistoryCreate) SetHash(s string) *HistoryCreate {
	hc.mutation.SetHash(s)
	return hc
}

// SetNillableHash sets the "hash" field if the given value is not nil.
func (hc *HistoryCreate) SetNillableHash(s *string) *HistoryCreate {
	if s != nil {
		hc.SetHash(*s)
	}
	return hc
}

// SetStatus sets the "status" field.
func (hc *HistoryCreate) SetStatus(h history.Status) *HistoryCreate {
	hc.mutation.SetStatus(h)
	return hc
}

// SetCreateTime sets the "create_time" field.
func (hc *HistoryCreate) SetCreateTime(t time.Time) *HistoryCreate {
	hc.mutation.SetCreateTime(t)
	return hc
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (hc *HistoryCreate) SetNillableCreateTime(t *time.Time) *HistoryCreate {
	if t != nil {
		hc.SetCreateTime(*t)
	}
	return hc
}

// Mutation returns the HistoryMutation object of the builder.
func (hc *HistoryCreate) Mutation() *HistoryMutation {
	return hc.mutation
}

// Save creates the History in the database.
func (hc *HistoryCreate) Save(ctx context.Context) (*History, error) {
	hc.defaults()
	return withHooks(ctx, hc.sqlSave, hc.mutation, hc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (hc *HistoryCreate) SaveX(ctx context.Context) *History {
	v, err := hc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (hc *HistoryCreate) Exec(ctx context.Context) error {
	_, err := hc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (hc *HistoryCreate) ExecX(ctx context.Context) {
	if err := hc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (hc *HistoryCreate) defaults() {
	if _, ok := hc.mutation.Size(); !ok {
		v := history.DefaultSize
		hc.mutation.SetSize(v)
	}
	if _, ok := hc.mutation.CreateTime(); !ok {
		v := history.DefaultCreateTime()
		hc.mutation.SetCreateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (hc *HistoryCreate) check() error {
	if _, ok := hc.mutation.MediaID(); !ok {
		return &ValidationError{Name: "media_id", err: errors.New(`ent: missing required field "History.media_id"`)}
	}
	if _, ok := hc.mutation.SourceTitle(); !ok {
		return &ValidationError{Name: "source_title", err: errors.New(`ent: missing required field "History.source_title"`)}
	}
	if _, ok := hc.mutation.Date(); !ok {
		return &ValidationError{Name: "date", err: errors.New(`ent: missing required field "History.date"`)}
	}
	if _, ok := hc.mutation.TargetDir(); !ok {
		return &ValidationError{Name: "target_dir", err: errors.New(`ent: missing required field "History.target_dir"`)}
	}
	if _, ok := hc.mutation.Size(); !ok {
		return &ValidationError{Name: "size", err: errors.New(`ent: missing required field "History.size"`)}
	}
	if _, ok := hc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "History.status"`)}
	}
	if v, ok := hc.mutation.Status(); ok {
		if err := history.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "History.status": %w`, err)}
		}
	}
	return nil
}

func (hc *HistoryCreate) sqlSave(ctx context.Context) (*History, error) {
	if err := hc.check(); err != nil {
		return nil, err
	}
	_node, _spec := hc.createSpec()
	if err := sqlgraph.CreateNode(ctx, hc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	hc.mutation.id = &_node.ID
	hc.mutation.done = true
	return _node, nil
}

func (hc *HistoryCreate) createSpec() (*History, *sqlgraph.CreateSpec) {
	var (
		_node = &History{config: hc.config}
		_spec = sqlgraph.NewCreateSpec(history.Table, sqlgraph.NewFieldSpec(history.FieldID, field.TypeInt))
	)
	if value, ok := hc.mutation.MediaID(); ok {
		_spec.SetField(history.FieldMediaID, field.TypeInt, value)
		_node.MediaID = value
	}
	if value, ok := hc.mutation.EpisodeNums(); ok {
		_spec.SetField(history.FieldEpisodeNums, field.TypeJSON, value)
		_node.EpisodeNums = value
	}
	if value, ok := hc.mutation.SeasonNum(); ok {
		_spec.SetField(history.FieldSeasonNum, field.TypeInt, value)
		_node.SeasonNum = value
	}
	if value, ok := hc.mutation.SourceTitle(); ok {
		_spec.SetField(history.FieldSourceTitle, field.TypeString, value)
		_node.SourceTitle = value
	}
	if value, ok := hc.mutation.Date(); ok {
		_spec.SetField(history.FieldDate, field.TypeTime, value)
		_node.Date = value
	}
	if value, ok := hc.mutation.TargetDir(); ok {
		_spec.SetField(history.FieldTargetDir, field.TypeString, value)
		_node.TargetDir = value
	}
	if value, ok := hc.mutation.Size(); ok {
		_spec.SetField(history.FieldSize, field.TypeInt, value)
		_node.Size = value
	}
	if value, ok := hc.mutation.DownloadClientID(); ok {
		_spec.SetField(history.FieldDownloadClientID, field.TypeInt, value)
		_node.DownloadClientID = value
	}
	if value, ok := hc.mutation.IndexerID(); ok {
		_spec.SetField(history.FieldIndexerID, field.TypeInt, value)
		_node.IndexerID = value
	}
	if value, ok := hc.mutation.Link(); ok {
		_spec.SetField(history.FieldLink, field.TypeString, value)
		_node.Link = value
	}
	if value, ok := hc.mutation.Hash(); ok {
		_spec.SetField(history.FieldHash, field.TypeString, value)
		_node.Hash = value
	}
	if value, ok := hc.mutation.Status(); ok {
		_spec.SetField(history.FieldStatus, field.TypeEnum, value)
		_node.Status = value
	}
	if value, ok := hc.mutation.CreateTime(); ok {
		_spec.SetField(history.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = value
	}
	return _node, _spec
}

// HistoryCreateBulk is the builder for creating many History entities in bulk.
type HistoryCreateBulk struct {
	config
	err      error
	builders []*HistoryCreate
}

// Save creates the History entities in the database.
func (hcb *HistoryCreateBulk) Save(ctx context.Context) ([]*History, error) {
	if hcb.err != nil {
		return nil, hcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(hcb.builders))
	nodes := make([]*History, len(hcb.builders))
	mutators := make([]Mutator, len(hcb.builders))
	for i := range hcb.builders {
		func(i int, root context.Context) {
			builder := hcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*HistoryMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, hcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, hcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, hcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (hcb *HistoryCreateBulk) SaveX(ctx context.Context) []*History {
	v, err := hcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (hcb *HistoryCreateBulk) Exec(ctx context.Context) error {
	_, err := hcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (hcb *HistoryCreateBulk) ExecX(ctx context.Context) {
	if err := hcb.Exec(ctx); err != nil {
		panic(err)
	}
}
