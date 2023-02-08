package driver

import (
	"strings"
)

// Sort is a sort for a query.
type Sort struct {
	Direction Direction
	Key       string
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
	OrderBy(key string, direction ...Direction) SortableConnection
	GetSorts() []*Sort
	ResetSorts()
}

// Direction is direction of a sort.
type Direction uint8

const (
	Asc Direction = iota
	Desc
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

// DefaultSortBuilder is a default sort builder.
type DefaultSortBuilder struct {
	sorts []*Sort
}

// OrderBy adds a sort to the builder.
func (d *DefaultSortBuilder) OrderBy(key string, direction ...Direction) SortableConnection {
	dir := Asc
	if len(direction) > 0 {
		dir = direction[0]
	}
	d.sorts = append(d.sorts, &Sort{
		Direction: dir,
		Key:       key,
	})
	return d
}

// GetSorts returns the sorts of the builder.
func (d *DefaultSortBuilder) GetSorts() []*Sort {
	return d.sorts
}

// ResetSorts resets the sorts of the builder.
func (d *DefaultSortBuilder) ResetSorts() {
	d.sorts = []*Sort{}
}

// BuildSortSQL builds the sort SQL.
func (d *DefaultSortBuilder) BuildSortSQL() string {
	if len(d.sorts) == 0 {
		return ""
	}
	var sql strings.Builder
	sql.WriteString(" ORDER BY ")
	for i, sort := range d.sorts {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString(sort.Key)
		sql.WriteString(" ")
		sql.WriteString(sort.Direction.String())
	}
	return sql.String()
}
