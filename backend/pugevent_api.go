package backend

import (
	"context"
	"net/http"

	"github.com/favclip/ucon/swagger"
	"github.com/favclip/ucon"
	"time"
)

func setupPugEvent(swPlugin *swagger.Plugin) {
	api := &PugEventAPI{}
	tag := swPlugin.AddTag(&swagger.Tag{Name: "PugEvent", Description: "PugEvent list"})
	var hInfo *swagger.HandlerInfo

	hInfo = swagger.NewHandlerInfo(api.Get)
	ucon.Handle(http.MethodGet, "/api/1/event", hInfo)
	hInfo.Description, hInfo.Tags = "get to events", []string{tag.Name}
}

type PugEventAPI struct{}

type PugEventAPIGetResponse struct {
	List []*PugEvent `json:"list"`
}

func (api *PugEventAPI) Get(c context.Context) (*PugEventAPIGetResponse, error) {
	e := PugEvent{
		Title: "GCPUG sinmetal Day 2018 May",
		StartAt: time.Now(),
	}

	list := []*PugEvent{}
	list = append(list, &e)

	return &PugEventAPIGetResponse{
		List : list,
	}, nil
}