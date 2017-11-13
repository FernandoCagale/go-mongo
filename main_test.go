package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	ID           bson.ObjectId
	DatabaseTest = Database + "-test"
	projectJSON  = `{
		"name":"CRM",
		"description": "Description Project",
		"tasks": [{
			"description":"Golang",
			"tags": ["BACKEND"]
		 },
		 {
			"description":"Mongodb",
			"tags": ["BACKEND"]
		 },
		 {
			"description":"ReactJS",
			"tags": ["FRONTEND"]
		 }]
	  }`
	projectJSONPut = `{
		"name":"BLOG",
		"description": "Description Project alter",
		"tasks": [{
			"description":"GO",
			"tags": ["BACK"]
		 },
		 {
			"description":"Postgres",
			"tags": ["BACK"]
		 },
		 {
			"description":"AngularJS",
			"tags": ["FRONT"]
		 }]
	  }`
)

func init() {
	session, err := getSession()
	if err != nil {
		panic(err)
	}

	err = session.DB(DatabaseTest).DropDatabase()
	if err != nil {
		panic(err)
	}
}

func TestPostProjects(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/projects", strings.NewReader(projectJSON))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	session, err := getSession()
	if err != nil {
		panic(err)
	}

	defer session.Close()

	c.Set("mongo", getCollectionTest(session))

	project := Project{}

	// Assertions
	if assert.NoError(t, postProjects(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		err := json.Unmarshal([]byte(rec.Body.String()), &project)

		assert.Nil(t, err)
		assert.NotNil(t, project.ID)
		assert.Equal(t, project.Name, "CRM")
		assert.Equal(t, project.Description, "Description Project")
		assert.Equal(t, project.Tasks[0].Description, "Golang")
		assert.Equal(t, project.Tasks[0].Tags[0], "BACKEND")
		assert.Equal(t, project.Tasks[1].Description, "Mongodb")
		assert.Equal(t, project.Tasks[1].Tags[0], "BACKEND")
		assert.Equal(t, project.Tasks[2].Description, "ReactJS")
		assert.Equal(t, project.Tasks[2].Tags[0], "FRONTEND")

		ID = project.ID
	}
}

func TestPostProjectsBadRequest(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/projects", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	session, err := getSession()
	if err != nil {
		panic(err)
	}

	defer session.Close()

	c.Set("mongo", getCollectionTest(session))

	// Assertions
	if assert.NoError(t, postProjects(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, rec.Body.String(), `{"message":"bad request"}`)
	}
}

func TestGetProjects(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/projects", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	session, err := getSession()
	if err != nil {
		panic(err)
	}

	defer session.Close()

	c.Set("mongo", getCollectionTest(session))

	projects := []Project{}

	// Assertions
	if assert.NoError(t, getProjects(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		err := json.Unmarshal([]byte(rec.Body.String()), &projects)

		assert.Nil(t, err)
		assert.NotNil(t, projects[0].ID)
		assert.Equal(t, projects[0].Name, "CRM")
		assert.Equal(t, projects[0].Description, "Description Project")
		assert.Equal(t, projects[0].Tasks[0].Description, "Golang")
		assert.Equal(t, projects[0].Tasks[0].Tags[0], "BACKEND")
		assert.Equal(t, projects[0].Tasks[1].Description, "Mongodb")
		assert.Equal(t, projects[0].Tasks[1].Tags[0], "BACKEND")
		assert.Equal(t, projects[0].Tasks[2].Description, "ReactJS")
		assert.Equal(t, projects[0].Tasks[2].Tags[0], "FRONTEND")
	}
}

func TestGetIdProjects(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/projects/:id", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(ID.Hex())

	session, err := getSession()
	if err != nil {
		panic(err)
	}

	defer session.Close()

	c.Set("mongo", getCollectionTest(session))

	project := Project{}

	// Assertions
	if assert.NoError(t, getIDProjects(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		err := json.Unmarshal([]byte(rec.Body.String()), &project)

		assert.Nil(t, err)
		assert.NotNil(t, project.ID)
		assert.Equal(t, project.Name, "CRM")
		assert.Equal(t, project.Description, "Description Project")
		assert.Equal(t, project.Tasks[0].Description, "Golang")
		assert.Equal(t, project.Tasks[0].Tags[0], "BACKEND")
		assert.Equal(t, project.Tasks[1].Description, "Mongodb")
		assert.Equal(t, project.Tasks[1].Tags[0], "BACKEND")
		assert.Equal(t, project.Tasks[2].Description, "ReactJS")
		assert.Equal(t, project.Tasks[2].Tags[0], "FRONTEND")
	}
}

func TestGetIdProjectsNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/projects/:id", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("59d3db513405112c58d87480")

	session, err := getSession()
	if err != nil {
		panic(err)
	}

	defer session.Close()

	c.Set("mongo", getCollectionTest(session))

	// Assertions
	if assert.NoError(t, getIDProjects(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, rec.Body.String(), `{"message":"not found"}`)
	}
}

func TestPutProjects(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.PUT, "/projects/:id", strings.NewReader(projectJSONPut))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(ID.Hex())

	session, err := getSession()
	if err != nil {
		panic(err)
	}

	defer session.Close()

	c.Set("mongo", getCollectionTest(session))

	project := Project{}

	// Assertions
	if assert.NoError(t, putProjects(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		err := json.Unmarshal([]byte(rec.Body.String()), &project)

		assert.Nil(t, err)
		assert.NotNil(t, project.ID)
		assert.Equal(t, project.Name, "BLOG")
		assert.Equal(t, project.Description, "Description Project alter")
		assert.Equal(t, project.Tasks[0].Description, "GO")
		assert.Equal(t, project.Tasks[0].Tags[0], "BACK")
		assert.Equal(t, project.Tasks[1].Description, "Postgres")
		assert.Equal(t, project.Tasks[1].Tags[0], "BACK")
		assert.Equal(t, project.Tasks[2].Description, "AngularJS")
		assert.Equal(t, project.Tasks[2].Tags[0], "FRONT")
	}
}

func TestPutProjectsBadRequest(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.PUT, "/projects/:id", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(ID.Hex())

	session, err := getSession()
	if err != nil {
		panic(err)
	}

	defer session.Close()

	c.Set("mongo", getCollectionTest(session))

	// Assertions
	if assert.NoError(t, putProjects(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, rec.Body.String(), `{"message":"bad request"}`)
	}
}

func TestPutProjectsNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.PUT, "/projects/:id", strings.NewReader(projectJSONPut))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("59d3db513405112c58d87480")

	session, err := getSession()
	if err != nil {
		panic(err)
	}

	defer session.Close()

	c.Set("mongo", getCollectionTest(session))

	// Assertions
	if assert.NoError(t, putProjects(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, rec.Body.String(), `{"message":"not found"}`)
	}
}

func TestDeleteProjects(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, "/projects/:id", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(ID.Hex())

	session, err := getSession()
	if err != nil {
		panic(err)
	}

	defer session.Close()

	c.Set("mongo", getCollectionTest(session))

	// Assertions
	if assert.NoError(t, deleteProjects(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteProjectsNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, "/projects/:id", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("59d3db513405112c58d87480")

	session, err := getSession()
	if err != nil {
		panic(err)
	}

	defer session.Close()

	c.Set("mongo", getCollectionTest(session))

	// Assertions
	if assert.NoError(t, deleteProjects(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, rec.Body.String(), `{"message":"not found"}`)
	}
}

func getCollectionTest(session *mgo.Session) *mgo.Collection {
	return session.DB(DatabaseTest).C(Collection)
}
