package requests

import (
	"Menu-Service/internal/service/helpers"
	"Menu-Service/resources"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/urlval"
)

type UpdateMealRequest struct {
	MealId int64 `url:"-" json:"-"`
	Data   resources.Meal
}

func NewUpdateMealRequest(r *http.Request) (UpdateMealRequest, error) {
	request := UpdateMealRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.MealId = cast.ToInt64(chi.URLParam(r, "id"))

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *UpdateMealRequest) validate() error {
	return helpers.MergeErrors(validation.Errors{
		"/data/attributes/meal_name": validation.Validate(&r.Data.Attributes.MealName, validation.Required,
			validation.Length(3, 45)),
		"/data/attributes/price": validation.Validate(&r.Data.Attributes.Price, validation.Required,
			validation.By(helpers.IsFloat)),
		"/data/attributes/amount": validation.Validate(&r.Data.Attributes.Amount, validation.Required,
			validation.By(helpers.IsFloat)),
		"/data/relationships/category/data/id": validation.Validate(&r.Data.Relationships.Category.Data.ID,
			validation.Required, validation.By(helpers.IsInteger)),
	}).Filter()
}