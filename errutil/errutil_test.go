package errutil

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2015 Essential Kaos                         //
//      Essential Kaos Open Source License <http://essentialkaos.com/ekol?en>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"errors"
	. "gopkg.in/check.v1"
	"testing"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

type ErrSuite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&ErrSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *ErrSuite) TestPositive(c *C) {
	errs := NewErrors()

	errs.Add(errors.New("1"))
	errs.Add(errors.New("2"))
	errs.Add(errors.New("3"))
	errs.Add(errors.New("4"))
	errs.Add(errors.New("5"))

	c.Assert(errs.All(), HasLen, 5)
	c.Assert(errs.HasErrors(), Equals, true)
	c.Assert(errs.Last(), DeepEquals, errors.New("5"))
	c.Assert(errs.All(), DeepEquals,
		[]error{
			errors.New("1"),
			errors.New("2"),
			errors.New("3"),
			errors.New("4"),
			errors.New("5"),
		},
	)
}

func (s *ErrSuite) TestNegative(c *C) {
	errs := NewErrors()

	c.Assert(errs.All(), HasLen, 0)
	c.Assert(errs.HasErrors(), Equals, false)
	c.Assert(errs.Last(), IsNil)
}