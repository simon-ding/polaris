// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"
	"polaris/ent/notificationclient"
	"polaris/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// NotificationClientQuery is the builder for querying NotificationClient entities.
type NotificationClientQuery struct {
	config
	ctx        *QueryContext
	order      []notificationclient.OrderOption
	inters     []Interceptor
	predicates []predicate.NotificationClient
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the NotificationClientQuery builder.
func (ncq *NotificationClientQuery) Where(ps ...predicate.NotificationClient) *NotificationClientQuery {
	ncq.predicates = append(ncq.predicates, ps...)
	return ncq
}

// Limit the number of records to be returned by this query.
func (ncq *NotificationClientQuery) Limit(limit int) *NotificationClientQuery {
	ncq.ctx.Limit = &limit
	return ncq
}

// Offset to start from.
func (ncq *NotificationClientQuery) Offset(offset int) *NotificationClientQuery {
	ncq.ctx.Offset = &offset
	return ncq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (ncq *NotificationClientQuery) Unique(unique bool) *NotificationClientQuery {
	ncq.ctx.Unique = &unique
	return ncq
}

// Order specifies how the records should be ordered.
func (ncq *NotificationClientQuery) Order(o ...notificationclient.OrderOption) *NotificationClientQuery {
	ncq.order = append(ncq.order, o...)
	return ncq
}

// First returns the first NotificationClient entity from the query.
// Returns a *NotFoundError when no NotificationClient was found.
func (ncq *NotificationClientQuery) First(ctx context.Context) (*NotificationClient, error) {
	nodes, err := ncq.Limit(1).All(setContextOp(ctx, ncq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{notificationclient.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (ncq *NotificationClientQuery) FirstX(ctx context.Context) *NotificationClient {
	node, err := ncq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first NotificationClient ID from the query.
// Returns a *NotFoundError when no NotificationClient ID was found.
func (ncq *NotificationClientQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ncq.Limit(1).IDs(setContextOp(ctx, ncq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{notificationclient.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (ncq *NotificationClientQuery) FirstIDX(ctx context.Context) int {
	id, err := ncq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single NotificationClient entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one NotificationClient entity is found.
// Returns a *NotFoundError when no NotificationClient entities are found.
func (ncq *NotificationClientQuery) Only(ctx context.Context) (*NotificationClient, error) {
	nodes, err := ncq.Limit(2).All(setContextOp(ctx, ncq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{notificationclient.Label}
	default:
		return nil, &NotSingularError{notificationclient.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (ncq *NotificationClientQuery) OnlyX(ctx context.Context) *NotificationClient {
	node, err := ncq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only NotificationClient ID in the query.
// Returns a *NotSingularError when more than one NotificationClient ID is found.
// Returns a *NotFoundError when no entities are found.
func (ncq *NotificationClientQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ncq.Limit(2).IDs(setContextOp(ctx, ncq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{notificationclient.Label}
	default:
		err = &NotSingularError{notificationclient.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (ncq *NotificationClientQuery) OnlyIDX(ctx context.Context) int {
	id, err := ncq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of NotificationClients.
func (ncq *NotificationClientQuery) All(ctx context.Context) ([]*NotificationClient, error) {
	ctx = setContextOp(ctx, ncq.ctx, "All")
	if err := ncq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*NotificationClient, *NotificationClientQuery]()
	return withInterceptors[[]*NotificationClient](ctx, ncq, qr, ncq.inters)
}

// AllX is like All, but panics if an error occurs.
func (ncq *NotificationClientQuery) AllX(ctx context.Context) []*NotificationClient {
	nodes, err := ncq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of NotificationClient IDs.
func (ncq *NotificationClientQuery) IDs(ctx context.Context) (ids []int, err error) {
	if ncq.ctx.Unique == nil && ncq.path != nil {
		ncq.Unique(true)
	}
	ctx = setContextOp(ctx, ncq.ctx, "IDs")
	if err = ncq.Select(notificationclient.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (ncq *NotificationClientQuery) IDsX(ctx context.Context) []int {
	ids, err := ncq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (ncq *NotificationClientQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, ncq.ctx, "Count")
	if err := ncq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, ncq, querierCount[*NotificationClientQuery](), ncq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (ncq *NotificationClientQuery) CountX(ctx context.Context) int {
	count, err := ncq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (ncq *NotificationClientQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, ncq.ctx, "Exist")
	switch _, err := ncq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (ncq *NotificationClientQuery) ExistX(ctx context.Context) bool {
	exist, err := ncq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the NotificationClientQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (ncq *NotificationClientQuery) Clone() *NotificationClientQuery {
	if ncq == nil {
		return nil
	}
	return &NotificationClientQuery{
		config:     ncq.config,
		ctx:        ncq.ctx.Clone(),
		order:      append([]notificationclient.OrderOption{}, ncq.order...),
		inters:     append([]Interceptor{}, ncq.inters...),
		predicates: append([]predicate.NotificationClient{}, ncq.predicates...),
		// clone intermediate query.
		sql:  ncq.sql.Clone(),
		path: ncq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.NotificationClient.Query().
//		GroupBy(notificationclient.FieldName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (ncq *NotificationClientQuery) GroupBy(field string, fields ...string) *NotificationClientGroupBy {
	ncq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &NotificationClientGroupBy{build: ncq}
	grbuild.flds = &ncq.ctx.Fields
	grbuild.label = notificationclient.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//	}
//
//	client.NotificationClient.Query().
//		Select(notificationclient.FieldName).
//		Scan(ctx, &v)
func (ncq *NotificationClientQuery) Select(fields ...string) *NotificationClientSelect {
	ncq.ctx.Fields = append(ncq.ctx.Fields, fields...)
	sbuild := &NotificationClientSelect{NotificationClientQuery: ncq}
	sbuild.label = notificationclient.Label
	sbuild.flds, sbuild.scan = &ncq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a NotificationClientSelect configured with the given aggregations.
func (ncq *NotificationClientQuery) Aggregate(fns ...AggregateFunc) *NotificationClientSelect {
	return ncq.Select().Aggregate(fns...)
}

func (ncq *NotificationClientQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range ncq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, ncq); err != nil {
				return err
			}
		}
	}
	for _, f := range ncq.ctx.Fields {
		if !notificationclient.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if ncq.path != nil {
		prev, err := ncq.path(ctx)
		if err != nil {
			return err
		}
		ncq.sql = prev
	}
	return nil
}

func (ncq *NotificationClientQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*NotificationClient, error) {
	var (
		nodes = []*NotificationClient{}
		_spec = ncq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*NotificationClient).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &NotificationClient{config: ncq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, ncq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (ncq *NotificationClientQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := ncq.querySpec()
	_spec.Node.Columns = ncq.ctx.Fields
	if len(ncq.ctx.Fields) > 0 {
		_spec.Unique = ncq.ctx.Unique != nil && *ncq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, ncq.driver, _spec)
}

func (ncq *NotificationClientQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(notificationclient.Table, notificationclient.Columns, sqlgraph.NewFieldSpec(notificationclient.FieldID, field.TypeInt))
	_spec.From = ncq.sql
	if unique := ncq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if ncq.path != nil {
		_spec.Unique = true
	}
	if fields := ncq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, notificationclient.FieldID)
		for i := range fields {
			if fields[i] != notificationclient.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := ncq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := ncq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := ncq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := ncq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (ncq *NotificationClientQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(ncq.driver.Dialect())
	t1 := builder.Table(notificationclient.Table)
	columns := ncq.ctx.Fields
	if len(columns) == 0 {
		columns = notificationclient.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if ncq.sql != nil {
		selector = ncq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if ncq.ctx.Unique != nil && *ncq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range ncq.predicates {
		p(selector)
	}
	for _, p := range ncq.order {
		p(selector)
	}
	if offset := ncq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := ncq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// NotificationClientGroupBy is the group-by builder for NotificationClient entities.
type NotificationClientGroupBy struct {
	selector
	build *NotificationClientQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ncgb *NotificationClientGroupBy) Aggregate(fns ...AggregateFunc) *NotificationClientGroupBy {
	ncgb.fns = append(ncgb.fns, fns...)
	return ncgb
}

// Scan applies the selector query and scans the result into the given value.
func (ncgb *NotificationClientGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ncgb.build.ctx, "GroupBy")
	if err := ncgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*NotificationClientQuery, *NotificationClientGroupBy](ctx, ncgb.build, ncgb, ncgb.build.inters, v)
}

func (ncgb *NotificationClientGroupBy) sqlScan(ctx context.Context, root *NotificationClientQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(ncgb.fns))
	for _, fn := range ncgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*ncgb.flds)+len(ncgb.fns))
		for _, f := range *ncgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*ncgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ncgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// NotificationClientSelect is the builder for selecting fields of NotificationClient entities.
type NotificationClientSelect struct {
	*NotificationClientQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ncs *NotificationClientSelect) Aggregate(fns ...AggregateFunc) *NotificationClientSelect {
	ncs.fns = append(ncs.fns, fns...)
	return ncs
}

// Scan applies the selector query and scans the result into the given value.
func (ncs *NotificationClientSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ncs.ctx, "Select")
	if err := ncs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*NotificationClientQuery, *NotificationClientSelect](ctx, ncs.NotificationClientQuery, ncs, ncs.inters, v)
}

func (ncs *NotificationClientSelect) sqlScan(ctx context.Context, root *NotificationClientQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ncs.fns))
	for _, fn := range ncs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ncs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ncs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
