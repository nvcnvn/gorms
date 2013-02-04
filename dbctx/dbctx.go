package dbctx

import (
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type DBCtx struct {
	orgColl    *mgo.Collection
	personColl *mgo.Collection
}

// NewDBCtx receive 2 *mgo.Collection for house and person
func NewDBCtx(org, person *mgo.Collection) *DBCtx {
	ctx := &DBCtx{}
	ctx.orgColl = org
	ctx.personColl = person
	return ctx
}

func (ctx *DBCtx) SaveHouse(p *Person) error {
	if p.Role != 0 {
		errors.New("A host Person must in Role 0")
	}
	p.PersonId = bson.NewObjectId()
	p.HouseId = p.PersonId
	p.AvgIncome = 0
	if len(p.Incomes) > 0 {
		for _, v := range p.Incomes {
			p.AvgIncome += v.Amount
		}
	}
	return ctx.personColl.Insert(p)
}

func (ctx *DBCtx) GetHouse(id bson.ObjectId) ([]Person, error) {
	members := []Person{}
	err := ctx.personColl.Find(bson.M{"houseid": id}).
		Sort("role", "birth").All(&members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (ctx *DBCtx) AddMember(m *Person, id bson.ObjectId) error {
	p := &Person{}
	if ctx.personColl.FindId(id).One(p) != nil {
		return errors.New("dbctx: invalid host id")
	}

	m.PersonId = bson.NewObjectId()
	m.HouseId = id
	m.HouseInfo = p.HouseInfo
	return ctx.personColl.Insert(m)
}

func (ctx *DBCtx) AddOrg(o *Organization) error {
	o.OrgId = bson.NewObjectId()
	return ctx.orgColl.Insert(o)
}

func (ctx *DBCtx) AllOrgs() ([]Organization, error) {
	orgs := []Organization{}
	err := ctx.orgColl.Find(nil).All(&orgs)
	if err != nil {
		return nil, err
	}
	return orgs, nil
}
