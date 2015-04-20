package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/juju/errgo"
	"github.com/martini-contrib/render"
	"github.com/nt3rp/convos/db"
)

func returnEnvelope(r render.Render, obj interface{}, err error) {
	// We are able to distinguish between multiple error types,
	// but for all the errors we have, they indicate an internal server error

	switch errgo.Cause(err) {
	case nil:
		// We could issue more specific http status codes for 'ok' (especially when creating objects)
		// but `ok` should be good enough for now.
		r.JSON(http.StatusOK, NewJsonEnvelopeFromObj(obj))
	default:
		r.JSON(http.StatusInternalServerError, NewJsonEnvelopeFromError(err))
	}
}

func getConvoFromRequest(req *http.Request) (*db.Convo, error) {
	decoder := json.NewDecoder(req.Body)

	var convo *db.Convo
	err := decoder.Decode(&convo)

	return convo, err
}

func getJsonFromRequest(req *http.Request) (map[string]string, error) {
	decoder := json.NewDecoder(req.Body)

	var jsonObj map[string]string
	err := decoder.Decode(&jsonObj)

	return jsonObj, err
}

func GetConvos(r render.Render) {
	convos, err := db.GetConvos()
	returnEnvelope(r, convos, err)
}

func GetConvo(params martini.Params, r render.Render) {
	id := params["id"]
	convo, err := db.GetConvo(id)
	returnEnvelope(r, convo, err)
}

func DeleteConvo(params martini.Params, r render.Render) {
	id := params["id"]
	err := db.DeleteConvo(id)
	returnEnvelope(r, "success", err)
}

func UpdateConvo(req *http.Request, params martini.Params, r render.Render) {
	patch, err := getJsonFromRequest(req)

	if err != nil {
		r.JSON(500, NewJsonEnvelopeFromError(err))
		return
	}

	id := params["id"]
	convo, err := db.UpdateConvo(id, patch["body"])
	returnEnvelope(r, convo, err)
}

func CreateConvo(req *http.Request, params martini.Params, r render.Render) {
	convo, err := getConvoFromRequest(req)

	if err != nil {
		r.JSON(500, NewJsonEnvelopeFromError(err))
		return
	}

	// For now, just suppress the conversion error
	id, _ := strconv.Atoi(params["id"])
	if id > 0 {
		convo.Parent = id

		// This incurs an extra DB call, but it seems like the simplest course of action to maintain the subject
		parent, err := db.GetConvo(params["id"])
		if err != nil {
			r.JSON(500, NewJsonEnvelopeFromError(err))
			return
		}

		convo.Subject = parent.Subject
	}

	// TODO: Need to return the saved object from the DB...
	newConvo, err := db.CreateConvo(convo)

	returnEnvelope(r, newConvo, err)
}
