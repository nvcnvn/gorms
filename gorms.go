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

		if house := c.Post("house", false); bson.IsObjectIdHex(house) {
			//Add member to a house
			m.Role, _ = strconv.Atoi(c.Post("Roles", false))
			c.db.AddMember(&m, bson.ObjectIdHex(house))
			c.Redirect("/data/edit.html?h="+m.HouseId.Hex(), 303)
		} else if person := c.Post("person", false); bson.IsObjectIdHex(person) {

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
		id := c.Get("h", false)
		if bson.IsObjectIdHex(id) {
			//House detail
			data := c.NewViewData("View")
			data["Quals"] = Quals
			data["Roles"] = Roles
			h, _ := c.db.GetHouse(bson.ObjectIdHex(id))
			data["House"] = h
			o, _ := c.db.AllOrgs()
			orgs := make(map[string]string)
			for _, v := range o {
				orgs[v.OrgId.Hex()] = v.Name
			}
			data["Orgs"] = orgs
			c.View("housedetail.tmpl", data)
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
	}
}
