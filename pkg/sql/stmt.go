/*
Copyright 2021 CodeNotary, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sql

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/codenotary/immudb/embedded/store"
)

const (
	catalogDatabasePrefix = "CATALOG.DATABASE." // (key=CATALOG.DATABASE.{dbID}, value={dbNAME})
	catalogTablePrefix    = "CATALOG.TABLE."    // (key=CATALOG.TABLE.{dbID}{tableID}{pkID}, value={tableNAME})
	catalogColumnPrefix   = "CATALOG.COLUMN."   // (key=CATALOG.COLUMN.{dbID}{tableID}{colID}{colTYPE}, value={colNAME})
	catalogIndexPrefix    = "CATALOG.INDEX."    // (key=CATALOG.INDEX.{dbID}{tableID}{colID}, value={})
	rowPrefix             = "ROW."              // (key=ROW.{dbID}{tableID}{colID}({valLen}{val})?{pkVal}, value={})
)

type SQLValueType = string

const (
	IntegerType   SQLValueType = "INTEGER"
	BooleanType                = "BOOLEAN"
	StringType                 = "STRING"
	BLOBType                   = "BLOB"
	TimestampType              = "TIMESTAMP"
)

type AggregateFn = int

const (
	COUNT AggregateFn = iota
	SUM
	MAX
	MIN
	AVG
)

type CmpOperator = int

const (
	EQ CmpOperator = iota
	NE
	LT
	LE
	GT
	GE
)

type LogicOperator = int

const (
	AND LogicOperator = iota
	OR
)

type JoinType = int

const (
	InnerJoin JoinType = iota
	LeftJoin
	RightJoin
)

type SQLStmt interface {
	isDDL() bool
	CompileUsing(e *Engine) (ces []*store.KV, des []*store.KV, err error)
}

type TxStmt struct {
	stmts []SQLStmt
}

func (stmt *TxStmt) isDDL() bool {
	for _, stmt := range stmt.stmts {
		if stmt.isDDL() {
			return true
		}
	}
	return false
}

func (stmt *TxStmt) CompileUsing(e *Engine) (ces []*store.KV, des []*store.KV, err error) {
	for _, stmt := range stmt.stmts {
		cs, ds, err := stmt.CompileUsing(e)
		if err != nil {
			return nil, nil, err
		}

		ces = append(ces, cs...)
		ds = append(ds, ds...)
	}
	return
}

type CreateDatabaseStmt struct {
	db string
}

// for writes, always needs to be up the date, doesn't matter the snapshot...
// for reading, a snapshot is created. It will wait until such tx is indexed.
// still writing to the catalog will wait the index to be up to date and locked
// conditional lock on writeLocked
func (stmt *CreateDatabaseStmt) isDDL() bool {
	return true
}

func (stmt *CreateDatabaseStmt) CompileUsing(e *Engine) (ces []*store.KV, des []*store.KV, err error) {
	db, err := e.catalog.newDatabase(stmt.db)
	if err != nil {
		return nil, nil, err
	}

	kv := &store.KV{
		Key:   e.mapKey(catalogDatabasePrefix, encodeID(db.id)),
		Value: []byte(stmt.db),
	}

	ces = append(ces, kv)

	return
}

type UseDatabaseStmt struct {
	db string
}

func (stmt *UseDatabaseStmt) isDDL() bool {
	return false
}

func (stmt *UseDatabaseStmt) CompileUsing(e *Engine) (ces []*store.KV, des []*store.KV, err error) {
	exists := e.catalog.ExistDatabase(stmt.db)
	if !exists {
		return nil, nil, ErrDatabaseDoesNotExist
	}

	e.implicitDatabase = stmt.db

	return
}

type UseSnapshotStmt struct {
	since, upTo string
}

func (stmt *UseSnapshotStmt) isDDL() bool {
	return false
}

func (stmt *UseSnapshotStmt) CompileUsing(e *Engine) (ces []*store.KV, des []*store.KV, err error) {
	return nil, nil, errors.New("not yet supported")
}

type CreateTableStmt struct {
	table    string
	colsSpec []*ColSpec
	pk       string
}

func (stmt *CreateTableStmt) isDDL() bool {
	return true
}

func (stmt *CreateTableStmt) CompileUsing(e *Engine) (ces []*store.KV, des []*store.KV, err error) {
	if e.implicitDatabase == "" {
		return nil, nil, ErrNoDatabaseSelected
	}

	db := e.catalog.dbsByName[e.implicitDatabase]

	table, err := db.newTable(stmt.table, stmt.colsSpec, stmt.pk)
	if err != nil {
		return nil, nil, err
	}

	for colID, col := range table.colsByID {
		ce := &store.KV{
			Key:   e.mapKey(catalogColumnPrefix, encodeID(db.id), encodeID(table.id), encodeID(colID), []byte(col.colType)),
			Value: nil,
		}
		ces = append(ces, ce)
	}

	te := &store.KV{
		Key:   e.mapKey(catalogTablePrefix, encodeID(db.id), encodeID(table.id), encodeID(table.pk.id)),
		Value: []byte(table.name),
	}
	ces = append(ces, te)

	return
}

type ColSpec struct {
	colName string
	colType SQLValueType
}

type CreateIndexStmt struct {
	table string
	col   string
}

func (stmt *CreateIndexStmt) isDDL() bool {
	return true
}

func (stmt *CreateIndexStmt) CompileUsing(e *Engine) (ces []*store.KV, des []*store.KV, err error) {
	if e.implicitDatabase == "" {
		return nil, nil, ErrNoDatabaseSelected
	}

	table, exists := e.catalog.dbsByName[e.implicitDatabase].tablesByName[stmt.table]
	if !exists {
		return nil, nil, ErrTableDoesNotExist
	}

	if table.pk.colName == stmt.col {
		return nil, nil, ErrIndexAlreadyExists
	}

	col, exists := table.colsByName[stmt.col]
	if !exists {
		return nil, nil, ErrColumnDoesNotExist
	}

	_, exists = table.indexes[col.id]
	if exists {
		return nil, nil, ErrIndexAlreadyExists
	}

	table.indexes[col.id] = struct{}{}

	te := &store.KV{
		Key:   e.mapKey(catalogIndexPrefix, encodeID(table.db.id), encodeID(table.id), encodeID(col.id)),
		Value: []byte(table.name),
	}
	ces = append(ces, te)

	return
}

type AddColumnStmt struct {
	table   string
	colSpec *ColSpec
}

func (stmt *AddColumnStmt) isDDL() bool {
	return true
}

func (stmt *AddColumnStmt) CompileUsing(e *Engine) (ces []*store.KV, des []*store.KV, err error) {
	return nil, nil, errors.New("not yet supported")
}

type UpsertIntoStmt struct {
	tableRef *TableRef
	cols     []string
	rows     []*RowSpec
}

type RowSpec struct {
	Values []Value
}

func (r *RowSpec) Bytes(t *Table, cols []string) ([]byte, error) {
	valbuf := bytes.Buffer{}

	// len(stmt.cols)
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(len(cols)))
	_, err := valbuf.Write(b[:])
	if err != nil {
		return nil, err
	}

	for i, val := range r.Values {
		col, _ := t.colsByName[cols[i]]

		// len(colName) + colName
		b := make([]byte, 4+len(col.colName))
		binary.BigEndian.PutUint32(b, uint32(len(col.colName)))
		copy(b[4:], []byte(col.colName))

		_, err = valbuf.Write(b)
		if err != nil {
			return nil, err
		}

		valb, err := encodeValue(val, col.colType, !asKey)
		if err != nil {
			return nil, err
		}

		_, err = valbuf.Write(valb)
		if err != nil {
			return nil, err
		}
	}

	return valbuf.Bytes(), nil
}

func (stmt *UpsertIntoStmt) isDDL() bool {
	return false
}

func (stmt *UpsertIntoStmt) Validate(table *Table) (map[uint64]int, error) {
	pkIncluded := false
	selByColID := make(map[uint64]int, len(stmt.cols))

	for i, c := range stmt.cols {
		col, exists := table.colsByName[c]
		if !exists {
			return nil, ErrInvalidColumn
		}

		if table.pk.colName == c {
			pkIncluded = true
		}

		_, duplicated := selByColID[col.id]
		if duplicated {
			return nil, ErrDuplicatedColumn
		}

		selByColID[col.id] = i
	}

	if !pkIncluded {
		return nil, ErrPKCanNotBeNull
	}

	return selByColID, nil
}

func (stmt *UpsertIntoStmt) CompileUsing(e *Engine) (ces []*store.KV, des []*store.KV, err error) {
	table, err := stmt.tableRef.referencedTable(e)
	if err != nil {
		return nil, nil, err
	}

	cs, err := stmt.Validate(table)
	if err != nil {
		return nil, nil, err
	}

	for _, row := range stmt.rows {
		if len(row.Values) != len(stmt.cols) {
			return nil, nil, ErrInvalidNumberOfValues
		}

		pkVal := row.Values[cs[table.pk.id]]
		pkEncVal, err := encodeValue(pkVal, table.pk.colType, asKey)
		if err != nil {
			return nil, nil, err
		}

		bs, err := row.Bytes(table, stmt.cols)
		if err != nil {
			return nil, nil, err
		}

		// create entry for the column which is the pk
		pke := &store.KV{
			Key:   e.mapKey(rowPrefix, encodeID(table.db.id), encodeID(table.id), encodeID(table.pk.id), pkEncVal),
			Value: bs,
		}
		des = append(des, pke)

		// create entries for each indexed column, with value as value for pk column
		for colID := range table.indexes {
			cVal := row.Values[cs[colID]]
			encVal, err := encodeValue(cVal, table.colsByID[colID].colType, asKey)
			if err != nil {
				return nil, nil, err
			}

			ie := &store.KV{
				Key:   e.mapKey(rowPrefix, encodeID(table.db.id), encodeID(table.id), encodeID(colID), encVal, pkEncVal),
				Value: nil,
			}
			des = append(des, ie)
		}
	}

	return
}

type Value interface {
	Value() interface{}
	jointColumnTo(col *Column) (*ColSelector, error)
}

type Number struct {
	val uint64
}

func (v *Number) Value() interface{} {
	return v.val
}

func (v *Number) jointColumnTo(col *Column) (*ColSelector, error) {
	return nil, ErrJointColumnNotFound
}

type String struct {
	val string
}

func (v *String) Value() interface{} {
	return v.val
}

func (v *String) jointColumnTo(col *Column) (*ColSelector, error) {
	return nil, ErrJointColumnNotFound
}

type Bool struct {
	val bool
}

func (v *Bool) Value() interface{} {
	return v.val
}

func (v *Bool) jointColumnTo(col *Column) (*ColSelector, error) {
	return nil, ErrJointColumnNotFound
}

type Blob struct {
	val []byte
}

func (v *Blob) Value() interface{} {
	return v.val
}

func (v *Blob) jointColumnTo(col *Column) (*ColSelector, error) {
	return nil, ErrJointColumnNotFound
}

type SysFn struct {
	fn string
}

func (v *SysFn) Value() interface{} {
	return nil
}

func (v *SysFn) jointColumnTo(col *Column) (*ColSelector, error) {
	return nil, ErrJointColumnNotFound
}

type Param struct {
	id string
}

func (v *Param) Value() interface{} {
	return nil
}

func (v *Param) jointColumnTo(col *Column) (*ColSelector, error) {
	return nil, ErrJointColumnNotFound
}

type Comparison int

const (
	EqualTo Comparison = iota
	LowerThan
	LowerOrEqualTo
	GreaterThan
	GreaterOrEqualTo
)

/*
func (e *Engine) tableFrom(colSel *ColSelector) (*Table, error) {
	if colSel == nil {
		return nil, ErrIllegalArguments
	}

	if e.implicitDatabase == "" && colSel.db == "" {
		return nil, ErrNoDatabaseSelected
	}

	dbName := e.implicitDatabase
	if colSel.db != "" {
		dbName = colSel.db
	}

	db, exist := e.catalog.dbsByName[dbName]
	if !exist {
		return nil, ErrDatabaseDoesNotExist
	}

	table, exist := db.tablesByName[colSel.table]
	if !exist {
		return nil, ErrTableDoesNotExist
	}

	return table, nil
}
*/

type DataSource interface {
	Resolve(e *Engine, snap *store.Snapshot, ordCol *OrdCol) (RowReader, error)
}

type SelectStmt struct {
	distinct  bool
	selectors []Selector
	ds        DataSource
	joins     []*JoinSpec
	where     BoolExp
	groupBy   []*ColSelector
	having    BoolExp
	limit     uint64
	orderBy   []*OrdCol
	as        string
}

func (stmt *SelectStmt) isDDL() bool {
	return false
}

func (stmt *SelectStmt) CompileUsing(e *Engine) (ces []*store.KV, des []*store.KV, err error) {
	if len(stmt.orderBy) > 1 {
		return nil, nil, ErrLimitedOrderBy
	}

	if len(stmt.orderBy) > 0 {
		tableRef, ok := stmt.ds.(*TableRef)
		if !ok {
			return nil, nil, ErrLimitedOrderBy
		}

		table, err := tableRef.referencedTable(e)
		if err != nil {
			return nil, nil, err
		}

		col, colExists := table.colsByName[stmt.orderBy[0].sel.col]
		if !colExists {
			return nil, nil, ErrLimitedOrderBy
		}

		if table.pk.id == col.id {
			return nil, nil, nil
		}

		_, indexed := table.indexes[col.id]
		if !indexed {
			return nil, nil, ErrLimitedOrderBy
		}
	}

	return nil, nil, nil
}

func (stmt *SelectStmt) Resolve(e *Engine, snap *store.Snapshot, ordCol *OrdCol) (RowReader, error) {
	// Ordering is only supported at TableRef level
	if ordCol != nil {
		return nil, ErrLimitedOrderBy
	}

	_, _, err := stmt.CompileUsing(e)
	if err != nil {
		return nil, err
	}

	var orderByCol *OrdCol

	if len(stmt.orderBy) > 0 {
		orderByCol = stmt.orderBy[0]
	}

	rowReader, err := stmt.ds.Resolve(e, snap, orderByCol)
	if err != nil {
		return nil, err
	}

	if stmt.joins != nil {
		rowReader, err = e.newJointRowReader(snap, rowReader, stmt.joins)
		if err != nil {
			return nil, err
		}
	}

	if stmt.where != nil {
		// filteredRowReader
	}

	//	rowBuilder := newRowBuilder(stmt.selectors)
	// another to filter selected rows
	//	&RowReaderWithrowReader

	if stmt.groupBy != nil {
		// groupedRowReader
	}

	if stmt.having != nil {
		// filteredRowReader
	}

	return rowReader, err

	/*cols := make(map[string]*Column, len(stmt.selectors))

	for _, s := range stmt.selectors {
		colSel := s.(*ColSelector)
		cols[colSel.col] = table.cols[colSel.col]
	}

	return &RowReader{reader: r, cols: cols}, nil
	*/
}

type TableRef struct {
	db    string
	table string
	as    string
}

func (stmt *TableRef) referencedTable(e *Engine) (*Table, error) {
	if e == nil {
		return nil, ErrIllegalArguments
	}

	var db string

	if db != "" {
		exists := e.catalog.ExistDatabase(stmt.db)
		if !exists {
			return nil, ErrDatabaseDoesNotExist
		}

		db = stmt.db
	}

	if db == "" {
		if e.implicitDatabase == "" {
			return nil, ErrNoDatabaseSelected
		}

		db = e.implicitDatabase
	}

	table, exists := e.catalog.dbsByName[db].tablesByName[stmt.table]
	if !exists {
		return nil, ErrTableDoesNotExist
	}

	return table, nil
}

func (stmt *TableRef) Resolve(e *Engine, snap *store.Snapshot, ordCol *OrdCol) (RowReader, error) {
	if e == nil || snap == nil || (ordCol != nil && ordCol.sel == nil) {
		return nil, ErrIllegalArguments
	}

	table, err := stmt.referencedTable(e)
	if err != nil {
		return nil, err
	}

	colName := table.pk.colName
	cmp := GreaterOrEqualTo
	var initKeyVal []byte

	if ordCol != nil {
		if ordCol.sel.db != "" && ordCol.sel.db != table.db.name {
			return nil, ErrInvalidColumn
		}

		if ordCol.sel.table != "" && ordCol.sel.table != table.name {
			return nil, ErrInvalidColumn
		}

		col, exist := table.colsByName[ordCol.sel.col]
		if !exist {
			return nil, ErrColumnDoesNotExist
		}

		// if it's not PK then it must be an indexed column
		if table.pk.colName != ordCol.sel.col {
			_, indexed := table.indexes[col.id]
			if !indexed {
				return nil, ErrColumnNotIndexed
			}
		}

		colName = col.colName
		cmp = ordCol.cmp

		if ordCol.useInitKeyVal {
			if len(initKeyVal) > len(maxKeyVal(col.colType)) {
				return nil, ErrIllegalArguments
			}
			initKeyVal = ordCol.initKeyVal
		}

		if !ordCol.useInitKeyVal && (cmp == LowerThan || cmp == LowerOrEqualTo) {
			initKeyVal = maxKeyVal(col.colType)
		}
	}

	return e.newRawRowReader(snap, table, colName, cmp, initKeyVal)
}

type JoinSpec struct {
	joinType JoinType
	ds       DataSource
	cond     BoolExp
}

type GroupBySpec struct {
	cols []string
}

type OrdCol struct {
	sel           *ColSelector
	cmp           Comparison
	initKeyVal    []byte
	useInitKeyVal bool
}

type Selector interface {
}

type ColSelector struct {
	db    string
	table string
	col   string
	as    string
}

func (sel *ColSelector) resolve(implicitDatabase string) string {
	if sel.as != "" {
		return sel.as
	}

	if sel.db == "" {
		return implicitDatabase + "." + sel.table + "." + sel.col
	}

	return sel.db + "." + sel.table + "." + sel.col
}

type AggSelector struct {
	aggFn AggregateFn
	as    string
}

type AggColSelector struct {
	aggFn AggregateFn
	db    string
	table string
	col   string
	as    string
}

type BoolExp interface {
	jointColumnTo(col *Column) (*ColSelector, error)
}

func (bexp *ColSelector) jointColumnTo(col *Column) (*ColSelector, error) {
	if bexp.db != "" && bexp.db != col.table.db.name {
		return nil, ErrJointColumnNotFound
	}

	if bexp.table != "" && bexp.table != col.table.name {
		return nil, ErrJointColumnNotFound
	}

	if bexp.col != col.colName {
		return nil, ErrJointColumnNotFound
	}

	return bexp, nil
}

type NotBoolExp struct {
	exp BoolExp
}

func (bexp *NotBoolExp) jointColumnTo(col *Column) (*ColSelector, error) {
	return bexp.exp.jointColumnTo(col)
}

type LikeBoolExp struct {
	col     *ColSelector
	pattern string
}

func (bexp *LikeBoolExp) jointColumnTo(col *Column) (*ColSelector, error) {
	return nil, ErrJointColumnNotFound
}

type CmpBoolExp struct {
	op          CmpOperator
	left, right BoolExp
}

func (bexp *CmpBoolExp) jointColumnTo(col *Column) (*ColSelector, error) {
	if bexp.op != EQ {
		return nil, ErrJointColumnNotFound
	}

	selLeft, okLeft := bexp.left.(*ColSelector)
	selRight, okRight := bexp.right.(*ColSelector)

	if !okLeft || !okRight {
		return nil, ErrJointColumnNotFound
	}

	_, errLeft := selLeft.jointColumnTo(col)
	_, errRight := selRight.jointColumnTo(col)

	if errLeft != nil && errLeft != ErrJointColumnNotFound {
		return nil, errLeft
	}

	if errRight != nil && errRight != ErrJointColumnNotFound {
		return nil, errRight
	}

	if errLeft == nil && errRight == nil {
		return nil, ErrInvalidJointColumn
	}

	if errLeft == nil {
		return selRight, nil
	}

	return selLeft, nil
}

type BinBoolExp struct {
	op          LogicOperator
	left, right BoolExp
}

func (bexp *BinBoolExp) jointColumnTo(col *Column) (*ColSelector, error) {
	jcolLeft, errLeft := bexp.left.jointColumnTo(col)
	jcolRight, errRight := bexp.left.jointColumnTo(col)

	if errLeft != nil && errLeft != ErrJointColumnNotFound {
		return nil, errLeft
	}

	if errRight != nil && errRight != ErrJointColumnNotFound {
		return nil, errRight
	}

	if errLeft == ErrJointColumnNotFound && errRight == ErrJointColumnNotFound {
		return nil, ErrJointColumnNotFound
	}

	if errLeft == nil && errRight == nil && jcolLeft != jcolRight {
		return nil, ErrInvalidJointColumn
	}

	if errLeft == nil {
		return jcolLeft, nil
	}

	return jcolRight, nil
}

type ExistsBoolExp struct {
	q *SelectStmt
}

func (bexp *ExistsBoolExp) jointColumnTo(col *Column) (*ColSelector, error) {
	return nil, ErrJointColumnNotFound
}