package store

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type defaultFileSystemStoreTestSuite struct {
	suite.Suite
	store *defaultFileSystemStore
}

func (s *defaultFileSystemStoreTestSuite) TestSetOk() {
	r, err := s.store.Set("xxx", false, "xxx")
	s.Suite.NoError(err)

	fmt.Printf("%v\n", r)
}

func TestStoreTestSuite(t *testing.T) {
	s := &defaultFileSystemStoreTestSuite{}
	s.store = newDefaultFileSystemStore()
	suite.Run(t, s)
}
