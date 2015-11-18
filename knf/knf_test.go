package knf

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2015 Essential Kaos                         //
//      Essential Kaos Open Source License <http://essentialkaos.com/ekol?en>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	check "gopkg.in/check.v1"
	"io/ioutil"
	"os"
	"testing"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const _CONFIG_DATA = `
    [formating]
test1:      1
            test2:2

      test3: 3

[string]
  test1: test
  test2: true
  test3: 4500
  test4: !$%^&
  test5: long long long long text for test
  test6: 

[boolean]
  test1: true
  test2: false
  test3: 0
  test4: 1
  test5:
  test6: example for test

[integer]
  test1: 1
  test2: -5
  test3: 10000000
  test4: A
  test5: 0xFF
  test6: 123.4
  test7: 123.456789
  test8: 0xZZYY
  test9: ABCD

[file-mode]
  test1: 644
  test2: 0644
  test3: 0
  test4: ABC
  test5: true

[comment]
  test1: 100
  # test2: 100

[macro]
  test1: 100
  test2: {macro:test1}.50
  test3: Value is {macro:test2}
  test4: "{macro:test3}"
  test5: {ABC}
  test6: {}

[k]
  t: 1
`

const _CONFIG_MALF_DATA = `
  test1: 123
  test2: 111
`

const (
	_CONFIG_FILE_NAME           = "knf-config-test.conf"
	_CONFIG_EMTPY_FILE_NAME     = "knf-config-test-empty.conf"
	_CONFIG_MALFORMED_FILE_NAME = "knf-config-test-malf.conf"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type KNFSuite struct {
	ConfigPath          string
	EmptyConfigPath     string
	MalformedConfigPath string
}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = check.Suite(&KNFSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) {
	check.TestingT(t)
}

func (s *KNFSuite) SetUpSuite(c *check.C) {
	tmpdir := c.MkDir()

	s.ConfigPath = tmpdir + "/" + _CONFIG_FILE_NAME
	s.EmptyConfigPath = tmpdir + "/" + _CONFIG_EMTPY_FILE_NAME
	s.MalformedConfigPath = tmpdir + "/" + _CONFIG_MALFORMED_FILE_NAME

	err := ioutil.WriteFile(s.ConfigPath, []byte(_CONFIG_DATA), 0644)

	if err != nil {
		c.Fatal(err.Error())
	}

	err = ioutil.WriteFile(s.EmptyConfigPath, []byte(""), 0644)

	if err != nil {
		c.Fatal(err.Error())
	}

	err = ioutil.WriteFile(s.MalformedConfigPath, []byte(_CONFIG_MALF_DATA), 0644)

	if err != nil {
		c.Fatal(err.Error())
	}
}

func (s *KNFSuite) TestErrors(c *check.C) {
	global = nil

	err := Global("/_not_exists_")

	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "File /_not_exists_ is not exist")

	err = Global(s.EmptyConfigPath)

	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, fmt.Sprintf("File %s is empty", s.EmptyConfigPath))

	err = Global(s.MalformedConfigPath)

	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, fmt.Sprintf("Configuration file %s is malformed.", s.MalformedConfigPath))

	updated, err := Reload()

	c.Assert(updated, check.IsNil)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "Global config is not loaded")

	c.Assert(GetS("test"), check.Equals, "")
	c.Assert(GetI("test"), check.Equals, 0)
	c.Assert(GetF("test"), check.Equals, 0.0)
	c.Assert(GetB("test"), check.Equals, false)
	c.Assert(GetM("test"), check.Equals, os.FileMode(0))
	c.Assert(HasSection("test"), check.Equals, false)
	c.Assert(HasProp("test"), check.Equals, false)
	c.Assert(Sections(), check.HasLen, 0)
	c.Assert(Props("test"), check.HasLen, 0)
	c.Assert(Validate([]*Validator{}), check.Not(check.HasLen), 0)
	c.Assert(Validate([]*Validator{})[0].Error(), check.Equals, "Global config struct is nil")

	config := &Config{}

	c.Assert(config.GetS("test"), check.Equals, "")
	c.Assert(config.GetI("test"), check.Equals, 0)
	c.Assert(config.GetF("test"), check.Equals, 0.0)
	c.Assert(config.GetB("test"), check.Equals, false)
	c.Assert(config.GetM("test"), check.Equals, os.FileMode(0))
	c.Assert(config.HasSection("test"), check.Equals, false)
	c.Assert(config.HasProp("test"), check.Equals, false)
	c.Assert(config.Sections(), check.HasLen, 0)
	c.Assert(config.Props("test"), check.HasLen, 0)
	c.Assert(config.Validate([]*Validator{}), check.HasLen, 0)

	updated, err = config.Reload()

	c.Assert(updated, check.IsNil)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "Path to config file is empty (non initialized struct?)")

	config = &Config{file: "/_not_exists_"}

	updated, err = config.Reload()

	c.Assert(updated, check.IsNil)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "File /_not_exists_ is not exist")
}

func (s *KNFSuite) TestParsing(c *check.C) {
	err := Global(s.ConfigPath)

	c.Assert(err, check.IsNil)

	_, err = Reload()

	c.Assert(err, check.IsNil)
}

func (s *KNFSuite) TestSections(c *check.C) {
	var err error

	err = Global(s.ConfigPath)

	c.Assert(global, check.NotNil)
	c.Assert(err, check.IsNil)

	sections := Sections()

	c.Assert(sections, check.HasLen, 8)
	c.Assert(
		sections,
		check.DeepEquals,
		[]string{"formating", "string", "boolean", "integer", "file-mode", "comment", "macro", "k"},
	)
}

func (s *KNFSuite) TestProps(c *check.C) {
	var err error

	err = Global(s.ConfigPath)

	c.Assert(global, check.NotNil)
	c.Assert(err, check.IsNil)

	props := Props("file-mode")

	c.Assert(props, check.HasLen, 5)
	c.Assert(
		props,
		check.DeepEquals,
		[]string{"test1", "test2", "test3", "test4", "test5"},
	)
}

func (s *KNFSuite) TestCheckers(c *check.C) {
	var err error

	err = Global(s.ConfigPath)

	c.Assert(global, check.NotNil)
	c.Assert(err, check.IsNil)

	c.Assert(HasSection("string"), check.Equals, true)
	c.Assert(HasSection("strings"), check.Equals, false)

	c.Assert(HasProp("string:test1"), check.Equals, true)
	c.Assert(HasProp("string:test9"), check.Equals, false)
	c.Assert(HasProp("strings:test9"), check.Equals, false)
}

func (s *KNFSuite) TestFormating(c *check.C) {
	var err error

	err = Global(s.ConfigPath)

	c.Assert(global, check.NotNil)
	c.Assert(err, check.IsNil)

	c.Assert(GetI("formating:test1"), check.Equals, 1)
	c.Assert(GetI("formating:test2"), check.Equals, 2)
	c.Assert(GetI("formating:test3"), check.Equals, 3)
}

func (s *KNFSuite) TestStrings(c *check.C) {
	var err error

	err = Global(s.ConfigPath)

	c.Assert(global, check.NotNil)
	c.Assert(err, check.IsNil)

	c.Assert(GetS("string:test1"), check.Equals, "test")
	c.Assert(GetS("string:test2"), check.Equals, "true")
	c.Assert(GetS("string:test3"), check.Equals, "4500")
	c.Assert(GetS("string:test4"), check.Equals, "!$%^&")
	c.Assert(GetS("string:test5"), check.Equals, "long long long long text for test")
}

func (s *KNFSuite) TestBoolean(c *check.C) {
	var err error

	err = Global(s.ConfigPath)

	c.Assert(global, check.NotNil)
	c.Assert(err, check.IsNil)

	c.Assert(GetB("boolean:test1"), check.Equals, true)
	c.Assert(GetB("boolean:test2"), check.Equals, false)
	c.Assert(GetB("boolean:test3"), check.Equals, false)
	c.Assert(GetB("boolean:test4"), check.Equals, true)
	c.Assert(GetB("boolean:test5"), check.Equals, false)
	c.Assert(GetB("boolean:test6"), check.Equals, true)
}

func (s *KNFSuite) TestInteger(c *check.C) {
	var err error

	err = Global(s.ConfigPath)

	c.Assert(global, check.NotNil)
	c.Assert(err, check.IsNil)

	c.Assert(GetI("integer:test1"), check.Equals, 1)
	c.Assert(GetI("integer:test2"), check.Equals, -5)
	c.Assert(GetI("integer:test3"), check.Equals, 10000000)
	c.Assert(GetI("integer:test4"), check.Equals, 0)
	c.Assert(GetI("integer:test5"), check.Equals, 0xFF)
	c.Assert(GetF("integer:test6"), check.Equals, 123.4)
	c.Assert(GetF("integer:test7"), check.Equals, 123.456789)
	c.Assert(GetF("integer:test8"), check.Equals, 0.0)
	c.Assert(GetI("integer:test8"), check.Equals, 0)
	c.Assert(GetF("integer:test9"), check.Equals, 0.0)
}

func (s *KNFSuite) TestFileMode(c *check.C) {
	var err error

	err = Global(s.ConfigPath)

	c.Assert(global, check.NotNil)
	c.Assert(err, check.IsNil)

	c.Assert(GetM("file-mode:test1"), check.Equals, os.FileMode(0644))
	c.Assert(GetM("file-mode:test2"), check.Equals, os.FileMode(0644))
	c.Assert(GetM("file-mode:test3"), check.Equals, os.FileMode(0))
	c.Assert(GetM("file-mode:test4"), check.Equals, os.FileMode(0))
	c.Assert(GetM("file-mode:test5"), check.Equals, os.FileMode(0))
}

func (s *KNFSuite) TestComments(c *check.C) {
	var err error

	err = Global(s.ConfigPath)

	c.Assert(global, check.NotNil)
	c.Assert(err, check.IsNil)

	c.Assert(GetI("comment:test1"), check.Equals, 100)
	c.Assert(GetI("comment:test2"), check.Not(check.Equals), 100)
	c.Assert(HasProp("comment:test2"), check.Equals, false)
}

func (s *KNFSuite) TestMacro(c *check.C) {
	var err error

	err = Global(s.ConfigPath)

	c.Assert(global, check.NotNil)
	c.Assert(err, check.IsNil)

	c.Assert(GetS("macro:test1"), check.Equals, "100")
	c.Assert(GetI("macro:test1"), check.Equals, 100)
	c.Assert(GetS("macro:test2"), check.Equals, "100.50")
	c.Assert(GetS("macro:test3"), check.Equals, "Value is 100.50")
	c.Assert(GetS("macro:test4"), check.Equals, "\"Value is 100.50\"")
	c.Assert(GetS("macro:test5"), check.Equals, "{ABC}")
	c.Assert(GetS("macro:test6"), check.Equals, "{}")
}

func (s *KNFSuite) TestNil(c *check.C) {
	var nilConf *Config

	c.Assert(nilConf.GetS("formating:test1"), check.Equals, "")
	c.Assert(nilConf.GetI("formating:test1"), check.Equals, 0)
	c.Assert(nilConf.GetF("formating:test1"), check.Equals, 0.0)
	c.Assert(nilConf.GetB("formating:test1"), check.Equals, false)
	c.Assert(nilConf.GetM("formating:test1"), check.Equals, os.FileMode(0))
	c.Assert(nilConf.HasSection("formating"), check.Equals, false)
	c.Assert(nilConf.HasProp("formating:test1"), check.Equals, false)
	c.Assert(nilConf.Sections(), check.HasLen, 0)
	c.Assert(nilConf.Props("formating"), check.HasLen, 0)

	_, err := nilConf.Reload()

	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "Config is nil")

	errs := nilConf.Validate([]*Validator{})

	c.Assert(errs, check.Not(check.HasLen), 0)
	c.Assert(errs[0].Error(), check.Equals, "Config is nil")
}

func (s *KNFSuite) TestDefault(c *check.C) {
	var err error

	err = Global(s.ConfigPath)

	c.Assert(global, check.NotNil)
	c.Assert(err, check.IsNil)

	c.Assert(GetS("string:test100", "fail"), check.Equals, "fail")
	c.Assert(GetB("boolean:test100", true), check.Equals, true)
	c.Assert(GetI("integer:test100", 9999), check.Equals, 9999)
	c.Assert(GetF("integer:test100", 123.45), check.Equals, 123.45)
	c.Assert(GetM("file-mode:test100", 0755), check.Equals, os.FileMode(0755))
	c.Assert(GetS("string:test6", "fail"), check.Equals, "fail")
}

func (s *KNFSuite) TestValidation(c *check.C) {
	var err error

	err = Global(s.ConfigPath)

	c.Assert(global, check.NotNil)
	c.Assert(err, check.IsNil)

	var errs []error

	errs = Validate([]*Validator{
		&Validator{"integer:test1", Empty, nil},
		&Validator{"integer:test1", Less, 0},
		&Validator{"integer:test1", Less, 0.5},
		&Validator{"integer:test1", Greater, 10},
		&Validator{"integer:test1", Greater, 10.1},
		&Validator{"integer:test1", Equals, 10},
		&Validator{"integer:test1", Equals, 10.1},
		&Validator{"integer:test1", Equals, "123"},
	})

	c.Assert(errs, check.HasLen, 0)

	errs = Validate([]*Validator{
		&Validator{"boolean:test5", Empty, nil},
		&Validator{"integer:test1", Less, 10},
		&Validator{"integer:test1", Greater, 0},
		&Validator{"integer:test1", Equals, 1},
		&Validator{"integer:test1", Greater, "12345"},
	})

	c.Assert(errs, check.HasLen, 5)

	c.Assert(errs[0].Error(), check.Equals, "Property boolean:test5 can't be empty")
	c.Assert(errs[1].Error(), check.Equals, "Property integer:test1 can't be less than 10")
	c.Assert(errs[2].Error(), check.Equals, "Property integer:test1 can't be greater than 0")
	c.Assert(errs[3].Error(), check.Equals, "Property integer:test1 can't be equal 1")
	c.Assert(errs[4].Error(), check.Equals, "Wrong validator for property integer:test1")
}
