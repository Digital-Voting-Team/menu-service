package handlers

import (
	"fmt"
	"menu-service/internal/service/helpers"
	requests "menu-service/internal/service/requests/menu"
	"menu-service/resources"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetMenu(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetMenuRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	resultMenu, err := helpers.MenusQ(r).FilterById(request.MenuId).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get menu from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	if resultMenu == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	resultCafe, err := helpers.ValidateCafe(r.Header.Get("Authorization"), r.Context().Value("cafeEndpoint").(string), resultMenu.CafeId)
	fmt.Println(resultCafe)

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
