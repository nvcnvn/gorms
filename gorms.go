package gorms

import (
	"github.com/nvcnvn/gorms/dbctx"
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

var Roles [28]string = [28]string{
	"Chủ hộ",
	"Ông",
	"Bà",
	"Ông cố",
	"Bà cố",
	"Ông nội",
	"Bà nội",
	"Ông ngoại",
	"Bà ngoại",
	"Cha",
	"Mẹ",
	"Cha chồng",
	"Mẹ chồng",
	"Cha vợ",
	"Mẹ vợ",
	"Anh",
	"Chị",
	"Em",
	"Anh rễ",
	"Chị dâu",
	"Em họ",
	"Anh họ",
	"Chị họ",
	"Em rễ",
	"Em dâu",
	"Con",
	"Cháu",
	"Khác",
}

var Quals [15]string = [15]string{
	"Không đi học",
	"Lớp 1",
	"Lớp 2",
	"Lớp 3",
	"Lớp 4",
	"Lớp 5",
	"Lớp 6",
	"Lớp 7",
	"Lớp 8",
	"Lớp 9",
	"Lớp 10",
	"Lớp 11",
	"Lớp 12",
	"Đại Học",
	"Sau Đại Học",
}

func Data(c *Controller) {
	p := c.Request().URL.Path

	var err error
	if "/data/submit.html" == p {
		h := dbctx.House{}
		h.Group = c.Post("Group", true)
		h.Block = c.Post("Block", true)
		h.Address = c.Post("Address", true)

		m := dbctx.Person{}
		m.HouseInfo = h

		m.FullName = c.Post("FullName", true)
		m.Area = c.Post("Area", true)
		m.Desire = c.Post("Desire", true)
		m.Note = c.Post("Note", true)
		m.Health = c.Post("Health", true)

		m.Birth, err = time.Parse("02/01/2006", c.Post("Birth", false))
		if err != nil {
			println("invalid birthday format")
			return
		}

		m.Quals, _ = strconv.Atoi(c.Post("Quals", false))
		m.HI, _ = strconv.Atoi(c.Post("HI", true))

		orgs := c.Request().Form["Orgs"]
		if len(orgs) > 0 {
			for _, v := range orgs {
				if bson.IsObjectIdHex(v) {
					m.Orgs = append(m.Orgs, bson.ObjectIdHex(v))
				}
			}
		}

		if c.Post("Gender", false) == "1" {
			m.Gender = true
		}

		if c.Post("AttendingSchool", false) == "1" {
			m.AttendingSchool = true
			m.Class.Title = c.Post("SchoolTitle", true)
			m.Class.School = c.Post("School", true)
		}

		if c.Post("Working", false) == "1" {
			m.Working = true
			m.Work.Title = c.Post("WorkTitle", true)
			m.Work.Office = c.Post("Office", true)
		}

		if c.Post("Incomes", false) == "1" {
			amount, err := strconv.Atoi(c.Post("Amount", true))
			if err == nil {
				m.Incomes = []dbctx.Income{dbctx.Income{amount, c.Post("Form", true)}}
			}
		}
		if person := c.Post("person", false); bson.IsObjectIdHex(person) {
			//Edit person info
			m.PersonId = bson.ObjectIdHex(person)
			hid := c.Post("house", false)
			if bson.IsObjectIdHex(hid) {
				println("here......")
				m.HouseId = bson.ObjectIdHex(hid)
				m.Role, _ = strconv.Atoi(c.Post("Roles", false))
				err := c.db.EditPerson(&m)
				if err != nil {
					println(err.Error())
				}
				c.Redirect("/data/edit.html?h="+m.HouseId.Hex(), 303)
			}
		} else if house := c.Post("house", false); bson.IsObjectIdHex(house) {
			//Add member to a house
			m.Role, _ = strconv.Atoi(c.Post("Roles", false))
			c.db.AddMember(&m, bson.ObjectIdHex(house))
			c.Redirect("/data/edit.html?h="+m.HouseId.Hex(), 303)
		} else {
			//Add new host member - add new house
			err = c.db.SaveHouse(&m)
			if err != nil {
				data := c.NewViewData("Error")
				data["Error"] = err.Error()
				data["Quals"] = Quals
				data["Roles"] = Roles
				c.View("/data/edit.html", data)
			} else {
				c.Redirect("/data/edit.html?h="+m.PersonId.Hex(), 303)
			}
		}
	} else if "/data/edit.html" == p {
		if house := c.Get("h", false); bson.IsObjectIdHex(house) {
			//House detail
			data := c.NewViewData("Thong tin ho dan")
			data["Quals"] = Quals
			data["Roles"] = Roles
			h, _ := c.db.GetHouse(bson.ObjectIdHex(house))
			data["House"] = h
			o, _ := c.db.AllOrgs()
			orgs := make(map[string]string)
			for _, v := range o {
				orgs[v.OrgId.Hex()] = v.Name
			}
			data["Orgs"] = orgs
			c.View("housedetail.tmpl", data)
		} else if person := c.Get("p", false); bson.IsObjectIdHex(person) {
			pData, err := c.db.GetPerson(bson.ObjectIdHex(person))
			if err != nil {
				c.Print("Thong tin khong ton tai.")
				return
			}

			data := c.NewViewData("Chinh sua thogn tin")
			o, _ := c.db.AllOrgs()
			orgs := make(map[string]string)
			for _, v := range o {
				orgs[v.OrgId.Hex()] = v.Name
			}
			data["Orgs"] = orgs
			data["Quals"] = Quals
			data["Roles"] = Roles
			data["Person"] = pData
			c.View("personedit.tmpl", data)
		} else {
			//Add Host member
			data := c.NewViewData("Add")
			data["Quals"] = Quals
			data["Roles"] = Roles
			o, _ := c.db.AllOrgs()
			orgs := make(map[string]string)
			for _, v := range o {
				orgs[v.OrgId.Hex()] = v.Name
			}
			data["Orgs"] = orgs
			c.View("houseadd.tmpl", data)
		}
	} else if "/data/print.html" == p {
		c.View("housedetail.tmpl", nil)
	} else {
		//List all person
		data := c.NewViewData("Danh sach ho dan")
		//filter parameter
		flags := struct {
			House, Frole, Fqual, Fgender, Fattending, Attending       bool
			Fworking, Working, Fincome, Gender, Forg, Fhealth, Health bool
			Fdesire, Desire, Fnote, Note, Fhi                         bool
			Squal, Squaldesc, Sage, Sagedesc, Sincome, Sincomedesc    bool
			Role, Qual, Incomefrom, Incometo, Hi, Offset, Limit       int
			Orgs                                                      []bson.ObjectId
			Offsetid                                                  bson.ObjectId
		}{}
		flags.Limit = 50
		if c.Get("type", false) == "1" {
			flags.House = true
		}

		if qual := c.Get("sort_qual", false); len(qual) == 1 {
			flags.Squal = true
			if qual == "1" {
				flags.Squaldesc = true
			}
		}

		if age := c.Get("sort_age", false); len(age) == 1 {
			flags.Sage = true
			if age == "1" {
				flags.Sagedesc = true
			}
		}

		if income := c.Get("sort_income", false); len(income) == 1 {
			flags.Sincome = true
			if income == "1" {
				flags.Sincomedesc = true
			}
		}

		if r, err := strconv.Atoi(c.Get("role", false)); err == nil {
			flags.Role = r
			flags.Frole = true
		}

		if q, err := strconv.Atoi(c.Get("qual", false)); err == nil {
			flags.Qual = q
			flags.Fqual = true
		}

		if gen := c.Get("gender", false); len(gen) == 1 {
			flags.Fgender = true
			if gen == "1" {
				flags.Gender = true
			}
		}

		if n := len(c.Request().Form["orgs"]); n > 0 {
			flags.Orgs = make([]bson.ObjectId, 0, n)
			for _, v := range c.Request().Form["orgs"] {
				if bson.IsObjectIdHex(v) {
					flags.Forg = true
					flags.Orgs = append(flags.Orgs, bson.ObjectIdHex(v))
				}
			}
		}

		if a := c.Get("attending", false); len(a) == 1 {
			flags.Fattending = true
			if a == "1" {
				flags.Attending = true
			}
		}

		if w := c.Get("working", false); len(w) == 1 {
			flags.Fworking = true
			if w == "1" {
				flags.Working = true
			}
		}

		if from, err := strconv.Atoi(c.Get("incomefrom", false)); err == nil {
			flags.Fincome = true
			flags.Incomefrom = from
		}

		if to, err := strconv.Atoi(c.Get("incometo", false)); err == nil {
			flags.Fincome = true
			flags.Incometo = to
		}

		if h, err := strconv.Atoi(c.Get("hi", false)); err == nil {
			flags.Fhi = true
			flags.Hi = h
		}

		if h := c.Get("health", false); len(h) == 1 {
			flags.Fhealth = true
			if h == "1" {
				flags.Health = true
			}
		}

		if d := c.Get("desire", false); len(d) == 1 {
			flags.Fdesire = true
			if d == "1" {
				flags.Desire = true
			}
		}

		if n := c.Get("note", false); len(n) == 1 {
			flags.Fnote = true
			if n == "1" {
				flags.Note = true
			}
		}

		persons, _ := c.db.Filter(flags.House, flags.Frole, flags.Role,
			flags.Fqual, flags.Qual, flags.Fgender, flags.Gender, flags.Forg,
			flags.Orgs, flags.Fattending, flags.Attending, flags.Fworking,
			flags.Working, flags.Fincome, flags.Incomefrom, flags.Incometo,
			flags.Fhi, flags.Hi, flags.Fhealth, flags.Health, flags.Fdesire,
			flags.Desire, flags.Fnote, flags.Note, flags.Offsetid, flags.Squal,
			flags.Squaldesc, flags.Sage, flags.Sagedesc, flags.Sincome,
			flags.Sincomedesc, flags.Offset, flags.Limit)
		o, _ := c.db.AllOrgs()

		orgsMap := make(map[string]string)
		for _, v := range o {
			orgsMap[v.OrgId.Hex()] = v.Name
		}
		data["Orgs"] = orgsMap
		data["Quals"] = Quals
		data["Roles"] = Roles
		data["Persons"] = persons
		data["Flags"] = flags
		c.View("houselst.tmpl", data)
	}
}
