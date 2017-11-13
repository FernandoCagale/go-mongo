package main

import (
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	Database   = "golang-project"
	Collection = "projects"
)

type Project struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson:"description"`
	Tasks       []Task        `json:"tasks" bson:"tasks"`
}

type Task struct {
	Description string   `json:"description" bson:"description"`
	Tags        []string `json:"tags" bson:"tags"`
}

func postProjects(c echo.Context) error {
	m := c.Get("mongo").(*mgo.Collection)

	project := new(Project)
	if err := c.Bind(project); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "bad request",
		})
	}

	project.ID = bson.NewObjectId()

	err := m.Insert(&project)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "bad request",
		})
	}

	return c.JSON(http.StatusOK, project)
}

func putProjects(c echo.Context) error {
	m := c.Get("mongo").(*mgo.Collection)
	id := bson.ObjectIdHex(c.Param("id"))

	project := new(Project)
	if err := c.Bind(project); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "bad request",
		})
	}

	err := m.UpdateId(id, &project)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	project.ID = id
	return c.JSON(http.StatusOK, project)
}

func getProjects(c echo.Context) error {
	m := c.Get("mongo").(*mgo.Collection)

	var projects []Project
	err := m.Find(nil).Sort("-start").All(&projects)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "bad request",
		})
	}

	return c.JSON(http.StatusOK, projects)
}

func getIDProjects(c echo.Context) error {
	m := c.Get("mongo").(*mgo.Collection)
	id := bson.ObjectIdHex(c.Param("id"))

	project := Project{}
	err := m.FindId(id).One(&project)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, project)
}

func deleteProjects(c echo.Context) error {
	m := c.Get("mongo").(*mgo.Collection)
	id := bson.ObjectIdHex(c.Param("id"))

	err := m.Remove(bson.M{"_id": id})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "ok")
}

func main() {
	session, err := getSession()
	if err != nil {
		panic(err)
	}

	defer session.Close()

	mongo := session.DB(Database).C(Collection)

	e := echo.New()
	e.Use(bindMongo(mongo))

	e.GET("/projects", getProjects)
	e.POST("/projects", postProjects)
	e.GET("/projects/:id", getIDProjects)
	e.PUT("/projects/:id", putProjects)
	e.DELETE("/projects/:id", deleteProjects)

	e.Logger.Fatal(e.Start(":8000"))
}

func getSession() (*mgo.Session, error) {
	Host := []string{
		"127.0.0.1:27017",
	}
	return mgo.DialWithInfo(&mgo.DialInfo{
		Addrs: Host,
	})
}

func bindMongo(mongo *mgo.Collection) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("mongo", mongo)
			return next(c)
		}
	}
}
