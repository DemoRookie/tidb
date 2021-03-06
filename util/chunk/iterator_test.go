// Copyright 2017 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package chunk

import (
	"github.com/pingcap/check"
	"github.com/pingcap/tidb/mysql"
	"github.com/pingcap/tidb/types"
)

func (s *testChunkSuite) TestIterator(c *check.C) {
	fields := []*types.FieldType{types.NewFieldType(mysql.TypeLonglong)}
	chk := NewChunk(fields)
	n := 10
	var expected []int64
	for i := 0; i < n; i++ {
		chk.AppendInt64(0, int64(i))
		expected = append(expected, int64(i))
	}
	var rows []Row
	li := NewList(fields, 1)
	li2 := NewList(fields, 5)
	var ptrs []RowPtr
	var ptrs2 []RowPtr
	for i := 0; i < n; i++ {
		rows = append(rows, chk.GetRow(i))
		ptr := li.AppendRow(chk.GetRow(i))
		ptrs = append(ptrs, ptr)
		ptr2 := li2.AppendRow(chk.GetRow(i))
		ptrs2 = append(ptrs2, ptr2)
	}

	it := NewSliceIterator(rows)
	checkIterator(c, it, expected)
	it = NewChunkIterator(chk)
	checkIterator(c, it, expected)
	it = NewListIterator(li)
	checkIterator(c, it, expected)
	it = NewRowPtrIterator(li, ptrs)
	checkIterator(c, it, expected)
	it = NewListIterator(li2)
	checkIterator(c, it, expected)
	it = NewRowPtrIterator(li2, ptrs2)
	checkIterator(c, it, expected)

	it = NewSliceIterator(nil)
	c.Assert(it.Begin(), check.Equals, it.End())
	it = NewChunkIterator(new(Chunk))
	c.Assert(it.Begin(), check.Equals, it.End())
	it = NewListIterator(new(List))
	c.Assert(it.Begin(), check.Equals, it.End())
	it = NewRowPtrIterator(li, nil)
	c.Assert(it.Begin(), check.Equals, it.End())
}

func checkIterator(c *check.C, it Iterator, expected []int64) {
	var got []int64
	for row := it.Begin(); row != it.End(); row = it.Next() {
		got = append(got, row.GetInt64(0))
	}
	c.Assert(got, check.DeepEquals, expected)
}
