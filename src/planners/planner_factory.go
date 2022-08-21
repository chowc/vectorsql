// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package planners

import (
	"parsers"

	"base/errors"
	"parsers/sqlparser"
)

type planCreator func(ast sqlparser.Statement) IPlan

// 根据语句类型找到对应的执行计划生成类
var table = map[string]planCreator{
	sqlparser.NodeNameUse:            NewUsePlan,
	sqlparser.NodeNameSelect:         NewSelectPlan,
	sqlparser.NodeNameDatabaseCreate: NewCreateDatabasePlan,
	sqlparser.NodeNameDatabaseDrop:   NewDropDatabasePlan,
	sqlparser.NodeNameTableCreate:    NewCreateTablePlan,
	sqlparser.NodeNameTableDrop:      NewDropTablePlan,
	sqlparser.NodeNameShowDatabases:  NewShowDatabasesPlan,
	sqlparser.NodeNameShowTables:     NewShowTablesPlan,
	sqlparser.NodeNameInsert:         NewInsertPlan,
}

func PlanFactory(query string) (IPlan, error) {
	statement, err := parsers.Parse(query)
	if err != nil {
		return nil, err
	}

	creator, ok := table[statement.Name()]
	if !ok {
		return nil, errors.Errorf("Couldn't get the planner:%T", statement)
	}
	// creator 是一个函数，调用它会构造返回对应类型的 IPlan，调用 IPlan.Build 再根据 Statement 构造出执行计划。
	plan := creator(statement)
	return plan, plan.Build()
}
