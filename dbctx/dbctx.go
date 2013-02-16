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

func (ctx *DBCtx) EditPerson(p *Person) error {
	p2 := &Person{}
	if ctx.personColl.FindId(p.PersonId).One(p2) != nil {
		return errors.New("dbctx: non exist person data")
	}

	if p2.PersonId == p2.HouseId && (p2.Role != p.Role || p2.HouseId != p.HouseId) {
		//The main host change their role or change their house
		members, err := ctx.GetHouse(p2.HouseId)
		n := len(members)
		if err == nil && n > 1 {
			//the house have more than 1 members.
			//find if there is another host
			found := false
			for _, v := range members {
				if v.PersonId != v.HouseId {
					found = true
					selector := bson.M{"houseid": v.HouseId}
					if p2.HouseId != p.HouseId {
						selector["_id"] = bson.M{"$ne": p.PersonId}
					}

					ctx.personColl.Update(selector, bson.M{"$set": bson.M{"houseid": v.PersonId}})
					break
				}
			}
			if !found {
				return errors.New("dbctx: house must have another host")
			}
		}
	}

	return ctx.personColl.UpdateId(p.PersonId, p)
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

func (ctx *DBCtx) GetPerson(id bson.ObjectId) (*Person, error) {
	p := &Person{}
	if err := ctx.personColl.FindId(id).One(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (ctx *DBCtx) Filter(house, frole bool, role int, fqual bool, qual int,
	fgender, gender, forg bool, orgs []bson.ObjectId, fattending, attending,
	fworking, working, fincome bool, incomefrom, incometo int, fhi bool, hi int,
	fhealth, health, fdesire, desire, fnote, note bool, offsetid bson.ObjectId,
	squal, squaldesc, sage, sagedesc, sincome, sincomedesc bool, offset, limit int) ([]Person, error) {
	plst := []Person{}

	filter := bson.M{}

	if fnote {
		filter["note"] = bson.M{"$exists": note}
	}

	if fdesire {
		filter["desire"] = bson.M{"$exists": desire}
	}

	if fhealth {
		filter["health"] = bson.M{"$exists": health}
	}

	if fhi {
		filter["hi"] = hi
	}

	if fincome {
		if incometo >= incomefrom {
			filter["avgincome"] = bson.M{
				"$gte": incomefrom,
				"$lte": incometo,
			}
		} else {
			filter["avgincome"] = bson.M{"$gte": incomefrom}
		}
	}

	if fworking {
		filter["working"] = working
	}

	if fattending {
		filter["attendingschool"] = attending
	}

	if forg {
		filter["orgs"] = bson.M{"$in": orgs}
	}

	if frole {
		filter["role"] = role
	}

	if fqual {
		filter["quals"] = qual
	}

	if fgender {
		filter["gender"] = gender
	}

	query := ctx.personColl.Find(filter)

	if squal {
		if squaldesc {
			query.Sort("-quals")
		} else {
			query.Sort("quals")
		}
	}

	if sage {
		if sagedesc {
			query.Sort("-birth")
		} else {
			query.Sort("birth")
		}
	}

	if sincome {
		if sincomedesc {
			query.Sort("-avgincome")
		} else {
			query.Sort("avgincome")
		}
	}

	if house {
		query.Sort("houseid", "role", "birth")
	}

	err := query.Limit(limit).All(&plst)
	if err != nil {
		return nil, err
	}
	return plst, nil
}
