// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package parsers

import (
	"strings"

	"parsers/sqlparser"
)

// 编译获得 AST，Statement 相当于 Expr/Stmt。Statement 同时也是一个 SQLNode。
// "select * from t;" 是一个 Statement；
// SQLNode 相当于是 Expr，可以在 SQLNode 上调用 Walk(Visitor)。
// Expr 可以是树结构。
func Parse(sql string) (sqlparser.Statement, error) {
	node, err := sqlparser.ParseStrictDDL(sql)
	if err != nil && strings.HasPrefix(strings.ToLower(sql), "insert") {
		if strings.HasSuffix(strings.ToLower(sql), "values") {
			sql += "('fill up')"
		} else {
			sql += " values('fill up')"
		}
		return sqlparser.ParseStrictDDL(sql)
	}
	return node, err
}
