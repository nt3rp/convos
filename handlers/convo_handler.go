package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/nt3rp/convos/db"
)


func GetConvos(r render.Render) {
	convos, err := db.GetConvos()

	if err != nil {
		r.JSON(500, err)
	} else {
		r.JSON(200, convos)
	}
}

// TODO: Envelope object
func GetConvo(params martini.Params, r render.Render) {
	id := params["id"]
	convo, err := db.GetConvo(id)

	if err != nil {
		r.JSON(500, err)
	} else {
		r.JSON(200, convo)
	}
}

func DeleteConvo(params martini.Params, r render.Render) {
	id := params["id"]
	err := db.DeleteConvo(id)

	if err != nil {
		r.JSON(500, err)
	} else {
		r.JSON(200, "success")
	}
}

func UpdateConvo() {
	// Do we want users to be able to fully edit conversations?
}

func CreateConvo(req *http.Request, r render.Render) {
	decoder := json.NewDecoder(req.Body)

	var convo *db.Convo
	err := decoder.Decode(&convo)

	if err != nil {
		r.JSON(500, err)
	}

	err = db.CreateConvo(convo)

	if err != nil {
		r.JSON(500, err)
	} else {
		r.JSON(200, convo)
	}
}
