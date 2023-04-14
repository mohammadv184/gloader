package driver

import (
	"strings"
)

// Sort is a sort for a query.
type Sort struct {
	Key       string
	Direction Direction
}

// GetDirection returns the direction of the sort.
func (s *Sort) GetDirection() Direction {
	return s.Direction
}

// GetKey returns the key of the sort.
func (s *Sort) GetKey() string {
	return s.Key
}

// SortableConnection is a connection that can be sorted.
type SortableConnection interface {
	OrderBy(dataCollection, key string, direction ...Direction) SortableConnection
	OrderByRoot(key string, direction ...Direction) SortableConnection
	GetSorts(dataCollection string) []*Sort
	GetAllSorts() map[string][]*Sort
	GetRootSorts() []*Sort
	ResetSorts(dataCollection string)
	ResetAllSorts()
}

// Direction is direction of a sort.
type Direction uint8

const (
	Asc  Direction = iota // Asc is the ascending direction.
	Desc                  // Desc is the descending direction.
)

var directionStringMap = map[Direction]string{
	Asc:  "ASC",
	Desc: "DESC",
}

// String returns the string representation of the direction.
func (d Direction) String() string {
	return directionStringMap[d]
}

// GetDirectionFromString returns the direction from a string.
func GetDirectionFromString(direction string) Direction {
	for k, v := range directionStringMap {
		if strings.EqualFold(v, direction) {
			return k
		}
	}
	return 0
}

// DefaultSortBuilder is the default implementation of SortableConnection.
// That can be used as embedded struct in a driver to implement simple filterable connection.
// And also this builder can build sql query from the sorts.
type DefaultSortBuilder struct {
	sorts     map[string][]*Sort
	rootSorts []*Sort // root sorts is a general sort that will apply to all data collections.
}

// OrderBy adds a sort to the builder.
func (sb *DefaultSortBuilder) OrderBy(dataCollection, key string, direction ...Direction) SortableConnection {
	sb.initSortsIfIsNil() // allocate memory for the sorts map if it is nil, for preventing panic.

	dir := Asc
	if len(direction) > 0 {
		dir = direction[0]
	}
	sb.sorts[dataCollection] = append(sb.sorts[dataCollection], &Sort{
		Direction: dir,
		Key:       key,
	})
	return sb
}

func (sb *DefaultSortBuilder) OrderByRoot(key string, direction ...Direction) SortableConnection {
	sb.initSortsIfIsNil() // allocate memory for the sorts map if it is nil, for preventing panic.

	dir := Asc
	if len(direction) > 0 {
		dir = direction[0]
	}
	sb.rootSorts = append(sb.rootSorts, &Sort{
		Direction: dir,
		Key:       key,
	})
	return sb
}

// GetSorts returns the sorts of the builder.
func (sb *DefaultSortBuilder) GetSorts(dataCollection string) []*Sort {
	sb.initSortsIfIsNil() // allocate memory for the sorts map if it is nil, for preventing panic.
	return append(sb.rootSorts, sb.sorts[dataCollection]...)
}

// GetAllSorts returns all sorts of the builder.
func (sb *DefaultSortBuilder) GetAllSorts() map[string][]*Sort {
	sb.initSortsIfIsNil() // allocate memory for the sorts map if it is nil, for preventing panic.
	return sb.sorts
}

func (sb *DefaultSortBuilder) GetRootSorts() []*Sort {
	sb.initSortsIfIsNil() // allocate memory for the sorts map if it is nil, for preventing panic.
	return sb.rootSorts
}

// ResetSorts resets the specified sorts of a data collection.
func (sb *DefaultSortBuilder) ResetSorts(dataCollection string) {
	sb.initSortsIfIsNil() // allocate memory for the sorts map if it is nil, for preventing panic.
	sb.sorts[dataCollection] = []*Sort{}
}

// ResetAllSorts resets all data collection sorts.
func (sb *DefaultSortBuilder) ResetAllSorts() {
	sb.sorts = make(map[string][]*Sort)
	sb.rootSorts = []*Sort{}
}

// BuildSortSQL builds the sort SQL.
func (sb *DefaultSortBuilder) BuildSortSQL(dataCollection string) string {
	sb.initSortsIfIsNil() // allocate memory for the sorts map if it is nil, for preventing panic.

	if len(sb.rootSorts) == 0 && len(sb.sorts[dataCollection]) == 0 {
		return ""
	}
	var sql strings.Builder
	sql.WriteString(" ORDER BY ")
	for i, sort := range append(sb.rootSorts, sb.sorts[dataCollection]...) {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString(sort.Key)
		sql.WriteString(" ")
		sql.WriteString(sort.Direction.String())
	}
	return sql.String()
}

func (sb *DefaultSortBuilder) initSortsIfIsNil() {
	if sb.sorts == nil {
		sb.sorts = make(map[string][]*Sort)
	}
}
