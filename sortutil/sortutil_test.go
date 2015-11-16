package sortutil

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2015 Essential Kaos                         //
//      Essential Kaos Open Source License <http://essentialkaos.com/ekol?en>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	. "gopkg.in/check.v1"
	"testing"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type SortSuite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&SortSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (ss *SortSuite) TestSorting(c *C) {
	v := []string{"1", "2.1", "2", "2.3.4", "1.3", "1.6.5", "2.3.3", "14.0", "6"}

	Versions(v)

	c.Assert(v, DeepEquals, []string{"1", "1.3", "1.6.5", "2", "2.1", "2.3.3", "2.3.4", "6", "14.0"})
}