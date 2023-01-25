package driver

import (
	"strings"
)

type SortableConnection interface {
	OrderBy(key string, direction ...Direction) SortableConnection
	GetSorts() []*Sort
	ResetSorts()
}

type Sort struct {
	Direction Direction
	Key       string
}

func (s *Sort) GetDirection() Direction {
	return s.Direction
}

func (s *Sort) GetKey() string {
	return s.Key
}

type Direction uint8

const (
	Asc Direction = iota
	Desc
)

var directionStringMap = map[Direction]string{
	Asc:  "ASC",
	Desc: "DESC",
}

func (d Direction) String() string {
	return directionStringMap[d]
}

func GetDirectionFromString(direction string) Direction {
	for k, v := range directionStringMap {
		if strings.ToLower(v) == strings.ToLower(direction) {
			return k
		}
	}
	return 0
}

type DefaultSortBuilder struct {
	sorts []*Sort
}

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
func (d *DefaultSortBuilder) GetSorts() []*Sort {
	return d.sorts
}
func (d *DefaultSortBuilder) ResetSorts() {
	d.sorts = []*Sort{}
}

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
