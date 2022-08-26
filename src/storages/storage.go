// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package storages

import (
	"columns"
	"datastreams"
	"sessions"
)

type IStorage interface {
	Name() string
	Columns() []*columns.Column
	// 获取输入流，用于从 Storage 读取数据。调用 IDataBlockInputStream.Read，数据可以在 Read 调用的时候再读取 -> 延迟读取。
	GetInputStream(*sessions.Session) (datastreams.IDataBlockInputStream, error)
	GetOutputStream(*sessions.Session) (datastreams.IDataBlockOutputStream, error)
	Close()
}
