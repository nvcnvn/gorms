package dbctx

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type Class struct {
	Title  string
	School string
}

type Work struct {
	Title  string
	Office string
}

type Income struct {
	Amount int
	From   string
}

type Organization struct {
	OrgId bson.ObjectId `bson:"_id"`
	Name  string
}

type Person struct {
	PersonId        bson.ObjectId `bson:"_id"`
	HouseId         bson.ObjectId
	Role            int
	HouseInfo       House
	FullName        string
	Gender          bool
	Birth           time.Time
	Quals           int
	Area            string          `bson:",omitempty"`
	Orgs            []bson.ObjectId `bson:",omitempty"`
	AttendingSchool bool
	Class           Class
	Working         bool
	Work            Work
	AvgIncome       int      `bson:",omitempty"`
	Incomes         []Income `bson:",omitempty"`
	HI              int      `bson:",omitempty"`
	Health          string   `bson:",omitempty"`
	Desire          string   `bson:",omitempty"`
	Note            string   `bson:",omitempty"`
}

type House struct {
	Group    string
	Block    string
	Ward     string `bson:",omitempty"`
	District string `bson:",omitempty"`
	City     string `bson:",omitempty"`
	Street   string `bson:",omitempty"`
	Address  string
}
