package utils

import (
	"gorm.io/gorm"
)

/**
 * @Author: zze
 * @Date: 2023/3/29 11:46
 * @Desc: 分页插件
 */

type WhereClause struct {
	Logic   string         `json:"logic"`
	Columns []ColumnClause `json:"columns"`
}

type ColumnClause struct {
	Column string `json:"column"`
	Op     string `json:"op"`
	Value  any    `json:"value"`
}

type Pager struct {
	Page          int           `json:"page" valid:"required~分页页码不能为空"`
	PageSize      int           `json:"page_size" valid:"required~分页大小不能为空"`
	Order         string        `json:"order,optional"`
	Wheres        []WhereClause `json:"wheres"`
	SelectColumns []any         `json:"select_columns"`
	Search        string        `json:"search"`
}

func (p Pager) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p Pager) Limit() int {
	if p.PageSize == 0 {
		return 10
	}

	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

func Paginate[T any](db *gorm.DB, pager *Pager, results *[]*T) (int, error) {
	if pager.PageSize == 0 {
		pager.PageSize = 10
	} else if pager.PageSize > 100 {
		pager.PageSize = 100
	}

	query := db.Model(new(T))
	// 构建 where 子句
	for _, where := range pager.Wheres {
		subQuery := db
		for _, column := range where.Columns {
			var (
				conditionStr   string
				conditionValue any
			)
			// 判断是否为模糊查询
			if val, ok := column.Value.(string); ok && column.Op == "~" && !IsEmpty(val) {
				conditionStr = column.Column + " LIKE ?"
				conditionValue = "%" + column.Value.(string) + "%"
			} else {
				switch column.Op {
				case "!=":
					conditionStr = column.Column + " != ?"
				case "=":
					conditionStr = column.Column + " = ?"
				case ">":
					conditionStr = column.Column + " > ?"
				case "<":
					conditionStr = column.Column + " < ?"
				case ">=":
					conditionStr = column.Column + " >= ?"
				case "<=":
					conditionStr = column.Column + " <= ?"
				case "->":
					conditionStr = column.Column + " IN (?)"
				case "!->":
					conditionStr = column.Column + " NOT IN (?)"
				}
				conditionValue = column.Value
			}
			if !IsEmpty(conditionStr) {
				if where.Logic == "and" {
					subQuery = subQuery.Where(conditionStr, conditionValue)
				} else {
					subQuery = subQuery.Or(conditionStr, conditionValue)
				}
			}
		}
		query = query.Where(subQuery)
	}
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return 0, err
	}
	total := int(totalCount)
	// 计算 offset
	offset := (pager.Page - 1) * pager.PageSize
	// 构建查询语句
	query = query.Limit(pager.PageSize).Offset(offset)

	if !IsEmpty(pager.Order) {
		query = query.Order(pager.Order)
	}
	// 执行查询并返回结果
	if len(pager.SelectColumns) > 0 {
		query = query.Select(pager.SelectColumns[0], pager.SelectColumns[1:]...)
	}
	return total, query.Find(&results).Error
}
