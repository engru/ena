package main

import "time"

// User is orm for user
// +genTable: user
type User struct {
	Id      int64
	Name    string `xorm:"unique notnull"`
	Salt    string
	Age     int
	Passwd  string    `xorm:"varchar(200)"`
	Created time.Time `xorm:"created"`
	Update  time.Time `xorm:"updated"`
}

/*
+genTable:
*/
// User2 and User3
type (
	// User2
	// +genTable: user2
	User2 struct { // user2 line
		Uid int `xorm:"pk 'uid'"`
	}

	// comment group

	// User3
	// +genTable:
	User3 struct { // user3 line
	}
)

type I1 interface {
}
