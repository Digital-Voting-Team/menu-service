package handlers

import (
	"Menu-Service/internal/data"
	"Menu-Service/internal/service/helpers"
	requests "Menu-Service/internal/service/requests/menu"
	"Menu-Service/resources"
	"github.com/spf13/cast"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func CreateMenu(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateMenuRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	Menu := data.Menu{
		CafeId: cast.ToInt64(request.Data.Relationships.Cafe.Data.ID),
	}

	var resultMenu data.Menu

	resultMenu, err = helpers.MenusQ(r).Insert(Menu)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create menu")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.MenuResponse{
		Data: resources.Menu{
			Key: resources.NewKeyInt64(resultMenu.Id, resources.MENU),
			Relationships: resources.MenuRelationships{
				Cafe: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(resultMenu.CafeId, 10),
						Type: resources.CAFE_REF,
					},
				},
			},
		},
	}
	ape.Render(w, result)
}