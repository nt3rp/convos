package handlers

import (
	"bytes"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/go-martini/martini"
	"github.com/juju/errgo"
	"github.com/martini-contrib/render"
	"github.com/nt3rp/convos/db"
	"github.com/nt3rp/convos/handlers/mocks"
)

var (
	firstPost *db.Convo = &db.Convo{
		Sender: 1, Recipient: 2, Subject: "First Post", Body: "Message Body", Read: true, Children: nil,
	}
)

/* Utilities */

func setupConvoHandlerTest(t *testing.T) {
	db.Initialize("test_convos")

	// In case we had paniced previously
	tearDownConvoHandlerTest(t)

	if err := db.AddUser("1", "Alice"); err != nil {
		t.Fatal(err)
	}

	if err := db.AddUser("2", "Bob"); err != nil {
		t.Fatal(err)
	}
}

func tearDownConvoHandlerTest(t *testing.T) {
	tables := []string{"read_status", "convos", "users"}

	for _, table := range tables {
		if err := db.TruncateTable(table); err != nil {
			t.Fatalf("Truncate table (%s): %s\n", table, err)
		}
	}
}

func generateTestRequest(authKey string, bodyStr string) *http.Request {
	body := &ClosingBuffer{bytes.NewBufferString(bodyStr)}
	request := &http.Request{
		Body:   body,
		Header: http.Header{},
	}

	if authKey != "" {
		request.Header.Set("X-USER-API-KEY", authKey)
	}

	return request
}

type HandlerPrerequisites struct {
	Req    *http.Request
	Params martini.Params
	Render render.Render
}

func generateHandlerPrerequisites(authorized bool, body string) HandlerPrerequisites {
	var userId string

	if authorized {
		userId = "1"
	} else {
		userId = "0"
	}

	request := generateTestRequest(userId, body)
	UserAuthorizationMiddleware(request)

	return HandlerPrerequisites{
		request,
		martini.Params{},
		&mocks.Render{},
	}
}

/* Tests */

func Test_CreateConvo_Authorized(t *testing.T) {
	setupConvoHandlerTest(t)
	defer tearDownConvoHandlerTest(t)

	post := &db.Convo{
		Recipient: 2, Subject: "First Post", Body: "Message Body",
	}

	p := generateHandlerPrerequisites(true, post.ToJson())
	savedPost := post
	savedPost.Id = 1
	savedPost.Parent = 1  // Parent set to the same as id for top-level conversation
	savedPost.Sender = 1  // User should be set to logged in user
	savedPost.Read = true // Automatically mark as read on post
	expected := NewJsonEnvelopeFromObj(savedPost)

	CreateConvo(p.Req, p.Params, p.Render)

	// Verify the result
	renderer, _ := p.Render.(*mocks.Render)
	if renderer.StatusCode != http.StatusOK {
		t.Errorf("Wrong Status Code set. Expected: %v. Actual: %v", http.StatusOK, renderer.StatusCode)
	}

	if !reflect.DeepEqual(renderer.Response, expected) {
		t.Errorf("JSON Envelopes do not match.\nExpected: %#v\nActual  : %#v", expected, renderer.Response)
	}
}

func Test_CreateConvo_Unauthorized(t *testing.T) {
	setupConvoHandlerTest(t)
	defer tearDownConvoHandlerTest(t)

	p := generateHandlerPrerequisites(false, "")

	// TODO: Should this 500, or should this return 403 / 401?
	// TODO: Why EOF?
	expected := NewJsonEnvelopeFromError(errgo.New("EOF"))

	CreateConvo(p.Req, p.Params, p.Render)

	// Verify the result
	renderer, _ := p.Render.(*mocks.Render)
	if renderer.StatusCode != http.StatusInternalServerError {
		t.Errorf("Wrong Status Code set. Expected: %v. Actual: %v", http.StatusInternalServerError, renderer.StatusCode)
	}

	if !reflect.DeepEqual(renderer.Response, expected) {
		t.Errorf("JSON Envelopes do not match.\nExpected: %#v\nActual  : %#v", expected, renderer.Response)
	}
}

func Test_GetConvos_Authorized(t *testing.T) {
	setupConvoHandlerTest(t)
	defer tearDownConvoHandlerTest(t)

	p := generateHandlerPrerequisites(true, "")
	convo, err := db.CreateConvo("1", firstPost)
	if err != nil {
		t.Error(err)
	}

	// Set Expectations
	expected := NewJsonEnvelopeFromObj([]*db.Convo{convo})

	GetConvos(p.Render)

	// Verify the result
	renderer, _ := p.Render.(*mocks.Render)
	if renderer.StatusCode != http.StatusOK {
		t.Errorf("Wrong Status Code set. Expected: %v. Actual: %v", http.StatusOK, renderer.StatusCode)
	}

	if !reflect.DeepEqual(renderer.Response, expected) {
		t.Errorf("JSON Envelopes do not match.\nExpected: %#v\nActual  : %#v", expected, renderer.Response)
	}
}

func Test_GetConvos_Unauthorized(t *testing.T) {
	setupConvoHandlerTest(t)
	defer tearDownConvoHandlerTest(t)

	p := generateHandlerPrerequisites(false, "")

	_, err := db.CreateConvo("1", firstPost)
	if err != nil {
		t.Error(err)
	}

	// Set Expectations
	var emptyList []*db.Convo
	expected := NewJsonEnvelopeFromObj(emptyList)

	GetConvos(p.Render)

	// Verify the result
	renderer, _ := p.Render.(*mocks.Render)
	if renderer.StatusCode != http.StatusOK {
		t.Errorf("Wrong Status Code set. Expected: %v. Actual: %v", http.StatusOK, renderer.StatusCode)
	}

	if !reflect.DeepEqual(renderer.Response, expected) {
		t.Errorf("JSON Envelopes do not match.\nExpected: %#v\nActual  : %#v", expected, renderer.Response)
	}
}

func Test_GetConvo_Authorized(t *testing.T) {
	setupConvoHandlerTest(t)
	defer tearDownConvoHandlerTest(t)

	p := generateHandlerPrerequisites(true, "")

	convo, err := db.CreateConvo("1", firstPost)
	if err != nil {
		t.Error(err)
	}
	p.Params["id"] = strconv.Itoa(convo.Id)

	// Set Expectations
	expected := NewJsonEnvelopeFromObj(convo)

	GetConvo(p.Params, p.Render)

	// Verify the result
	renderer, _ := p.Render.(*mocks.Render)
	if renderer.StatusCode != http.StatusOK {
		t.Errorf("Wrong Status Code set. Expected: %v. Actual: %v", http.StatusOK, renderer.StatusCode)
	}

	if !reflect.DeepEqual(renderer.Response, expected) {
		t.Errorf("JSON Envelopes do not match.\nExpected: %#v\nActual  : %#v", expected, renderer.Response)
	}
}

func Test_GetConvo_Unauthorized(t *testing.T) {
	setupConvoHandlerTest(t)
	defer tearDownConvoHandlerTest(t)

	p := generateHandlerPrerequisites(false, "")

	convo, err := db.CreateConvo("1", firstPost)
	if err != nil {
		t.Error(err)
	}
	p.Params["id"] = strconv.Itoa(convo.Id)

	// Set Expectations
	expected := NewJsonEnvelopeFromError(errgo.Newf("Unable to find convo with id '%d'.: sql: no rows in result set", convo.Id))

	GetConvo(p.Params, p.Render)

	// Verify the result
	renderer, _ := p.Render.(*mocks.Render)
	if renderer.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong Status Code set. Expected: %v. Actual: %v", http.StatusNotFound, renderer.StatusCode)
	}

	if !reflect.DeepEqual(renderer.Response, expected) {
		t.Errorf("JSON Envelopes do not match.\nExpected: %#v\nActual  : %#v", expected, renderer.Response)
	}
}

func Test_DeleteConvo_Authorized(t *testing.T) {
	setupConvoHandlerTest(t)
	defer tearDownConvoHandlerTest(t)

	p := generateHandlerPrerequisites(true, "")

	convo, err := db.CreateConvo("1", firstPost)
	if err != nil {
		t.Error(err)
	}
	p.Params["id"] = strconv.Itoa(convo.Id)

	// Set Expectations
	expected := NewJsonEnvelopeFromObj("success")

	DeleteConvo(p.Params, p.Render)

	// Verify the result
	renderer, _ := p.Render.(*mocks.Render)
	if renderer.StatusCode != http.StatusOK {
		t.Errorf("Wrong Status Code set. Expected: %v. Actual: %v", http.StatusOK, renderer.StatusCode)
	}

	if !reflect.DeepEqual(renderer.Response, expected) {
		t.Errorf("JSON Envelopes do not match.\nExpected: %#v\nActual  : %#v", expected, renderer.Response)
	}
}

func Test_DeleteConvo_Unauthorized(t *testing.T) {
	setupConvoHandlerTest(t)
	defer tearDownConvoHandlerTest(t)

	p := generateHandlerPrerequisites(false, "")

	convo, err := db.CreateConvo("1", firstPost)
	if err != nil {
		t.Error(err)
	}
	p.Params["id"] = strconv.Itoa(convo.Id)

	// Set Expectations
	expected := NewJsonEnvelopeFromError(errgo.Newf("Unable to find convo with id '%d'.", convo.Id))

	// Do what we need to do
	DeleteConvo(p.Params, p.Render)

	// Verify the result
	renderer, _ := p.Render.(*mocks.Render)
	if renderer.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong Status Code set. Expected: %v. Actual: %v", http.StatusNotFound, renderer.StatusCode)
	}

	if !reflect.DeepEqual(renderer.Response, expected) {
		t.Errorf("JSON Envelopes do not match.\nExpected: %#v\nActual  : %#v", expected, renderer.Response)
	}
}

func Test_ReplyConvo_Authorized(t *testing.T) {
	setupConvoHandlerTest(t)
	defer tearDownConvoHandlerTest(t)

	post := &db.Convo{
		Recipient: 2, Subject: "First Post", Body: "Message Body",
	}
	p := generateHandlerPrerequisites(true, post.ToJson())

	convo, err := db.CreateConvo("1", firstPost)
	if err != nil {
		t.Error(err)
	}
	p.Params["id"] = strconv.Itoa(convo.Id)

	// Set Expectations
	savedPost := post
	savedPost.Id = convo.Id + 1
	savedPost.Parent = convo.Id // Parent set to the same as id for top-level conversation
	savedPost.Sender = 1        // User should be set to logged in user
	savedPost.Read = true       // Automatically mark as read on post
	expected := NewJsonEnvelopeFromObj(savedPost)

	CreateConvo(p.Req, p.Params, p.Render)

	// Verify the result
	renderer, _ := p.Render.(*mocks.Render)
	if renderer.StatusCode != http.StatusOK {
		t.Errorf("Wrong Status Code set. Expected: %v. Actual: %v", http.StatusOK, renderer.StatusCode)
	}

	if !reflect.DeepEqual(renderer.Response, expected) {
		t.Errorf("JSON Envelopes do not match.\nExpected: %#v\nActual  : %#v", expected, renderer.Response)
	}
}

func Test_ReplyConvo_Unauthorized(t *testing.T) {
	setupConvoHandlerTest(t)
	defer tearDownConvoHandlerTest(t)

	post := &db.Convo{
		Sender: 1, Recipient: 2, Subject: "First Post", Body: "Message Body", Read: true, Children: nil,
	}
	p := generateHandlerPrerequisites(false, post.ToJson())

	convo, err := db.CreateConvo("1", firstPost)
	if err != nil {
		t.Error(err)
	}
	p.Params["id"] = strconv.Itoa(convo.Id)

	// Set Expectations
	expected := NewJsonEnvelopeFromError(errgo.Newf("Unable to find convo with id '%d'.: sql: no rows in result set", convo.Id))

	CreateConvo(p.Req, p.Params, p.Render)

	// Verify the result
	renderer, _ := p.Render.(*mocks.Render)
	if renderer.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong Status Code set. Expected: %v. Actual: %v", http.StatusNotFound, renderer.StatusCode)
	}

	if !reflect.DeepEqual(renderer.Response, expected) {
		t.Errorf("JSON Envelopes do not match.\nExpected: %#v\nActual  : %#v", expected, renderer.Response)
	}
}

func Test_UpdateConvo_Authorized(t *testing.T) {
	setupConvoHandlerTest(t)
	defer tearDownConvoHandlerTest(t)

	// BUG! This should change the read status...
	p := generateHandlerPrerequisites(true, "{\"read\": \"false\"}")

	convo, err := db.CreateConvo("1", firstPost)
	if err != nil {
		t.Error(err)
	}
	p.Params["id"] = strconv.Itoa(convo.Id)

	// Set Expectations
	patchedConvo := convo
	patchedConvo.Read = false
	expected := NewJsonEnvelopeFromObj(patchedConvo)

	UpdateConvo(p.Req, p.Params, p.Render)

	// Verify the result
	renderer, _ := p.Render.(*mocks.Render)
	if renderer.StatusCode != http.StatusOK {
		t.Errorf("Wrong Status Code set. Expected: %v. Actual: %v", http.StatusOK, renderer.StatusCode)
	}

	env := renderer.Response.(JsonEnvelope)
	t.Log(env.Response)
	if !reflect.DeepEqual(renderer.Response, expected) {
		t.Errorf("JSON Envelopes do not match.\nExpected: %#v\nActual  : %#v", expected, renderer.Response)
	}
}

func Test_UpdateConvo_Unauthorized(t *testing.T) {
	setupConvoHandlerTest(t)
	defer tearDownConvoHandlerTest(t)

	p := generateHandlerPrerequisites(false, "{\"read\": \"false\"}")

	convo, err := db.CreateConvo("1", firstPost)
	if err != nil {
		t.Error(err)
	}
	p.Params["id"] = strconv.Itoa(convo.Id)

	// Set Expectations
	expected := NewJsonEnvelopeFromError(errgo.Newf("Unable to find convo with id '%d'.: sql: no rows in result set", convo.Id))

	UpdateConvo(p.Req, p.Params, p.Render)

	// Verify the result
	renderer, _ := p.Render.(*mocks.Render)
	if renderer.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong Status Code set. Expected: %v. Actual: %v", http.StatusNotFound, renderer.StatusCode)
	}

	env := renderer.Response.(JsonEnvelope)
	t.Log(env.Response)
	if !reflect.DeepEqual(renderer.Response, expected) {
		t.Errorf("JSON Envelopes do not match.\nExpected: %#v\nActual  : %#v", expected, renderer.Response)
	}
}

func Test_Persist_Delete(t *testing.T) {
	setupConvoHandlerTest(t)
	defer tearDownConvoHandlerTest(t)

	p := generateHandlerPrerequisites(true, "")

	var parentId, childId string
	parent, err := db.CreateConvo("1", firstPost)
	if err != nil {
		t.Error(err)
	}
	parentId = strconv.Itoa(parent.Id)

	parent.Id = 0
	child, err := db.CreateConvo("1", parent)
	if err != nil {
		t.Error(err)
	}
	childId = strconv.Itoa(child.Id)

	// Delete
	p.Params["id"] = parentId
	DeleteConvo(p.Params, p.Render)

	// Get the parent
	p.Params["id"] = parentId
	GetConvo(p.Params, p.Render)

	// Verify it no longer exists
	renderer, _ := p.Render.(*mocks.Render)
	if renderer.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong Status Code set. Expected: %v. Actual: %v", http.StatusNotFound, renderer.StatusCode)
	}

	// Get the child
	p.Params["id"] = childId
	GetConvo(p.Params, p.Render)

	// Verify it no longer exists
	renderer, _ = p.Render.(*mocks.Render)
	if renderer.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong Status Code set. Expected: %v. Actual: %v", http.StatusNotFound, renderer.StatusCode)
	}
}
