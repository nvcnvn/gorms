package gorms

import (
	"github.com/nvcnvn/gorms/dbctx"
	"github.com/openvn/toys"
	"github.com/openvn/toys/secure/membership"
	"github.com/openvn/toys/secure/membership/session"
	"github.com/openvn/toys/view"
	"labix.org/v2/mgo"
	"net/http"
	"path"
	"time"
)

const (
	dbname string = "test"
)

type Controller struct {
	toys.Controller
	sess session.Provider
	auth membership.Authenticater
	tmpl *view.View
	db   *dbctx.DBCtx
}

func (c *Controller) NewViewData(title string) map[string]interface{} {
	m := make(map[string]interface{})
	m["Title"] = title
	m["DBCtx"] = c.db
	return m
}

func (c *Controller) View(page string, data interface{}) {
	c.tmpl.Load(c, page, data)
}

type Handler struct {
	fn     func(c *Controller)
	dbsess *mgo.Session
	tmpl   *view.View
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := Controller{}
	c.Init(w, r)

	dbsess := h.dbsess.Clone()
	defer dbsess.Close()

	//database collection (table)
	database := dbsess.DB(dbname)
	sessColl := database.C("toysSession")
	userColl := database.C("toysUser")

	rememberColl := database.C("toysUserRemember")

	houseColl := database.C("toysOrgs")
	personColl := database.C("toysPerson")

	//web session
	c.sess = session.NewMgoProvider(w, r, sessColl)

	//web authenthicator
	c.auth = membership.NewAuthDBCtx(w, r, c.sess, userColl, rememberColl)

	//database context
	c.db = dbctx.NewDBCtx(houseColl, personColl)

	//view template
	c.tmpl = h.tmpl

	//process
	h.fn(&c)
}

// NewHandler receive a controll function and a mongodb session
func NewHandler(f func(c *Controller), dbsess *mgo.Session, tmpl *view.View) *Handler {
	h := &Handler{}
	h.dbsess = dbsess
	h.tmpl = tmpl
	h.fn = f

	dbsess.DB(dbname).C("toysSession").EnsureIndex(mgo.Index{
		Key:         []string{"lastactivity"},
		ExpireAfter: 7200 * time.Second,
	})

	dbsess.DB(dbname).C("toysUser").EnsureIndex(mgo.Index{
		Key:    []string{"email"},
		Unique: true,
	})

	return h
}

// Match is a wrapper function for path.Math
func Match(pattern, name string) bool {
	ok, err := path.Match(pattern, name)
	if err != nil {
		return false
	}
	return ok
}
