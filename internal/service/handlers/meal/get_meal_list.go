package handlers

import (
	"github.com/Digital-Voting-Team/menu-service/internal/data"
	"github.com/Digital-Voting-Team/menu-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/menu-service/internal/service/requests/meal"
	"github.com/Digital-Voting-Team/menu-service/resources"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetMealList(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetMealListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	mealsQ := helpers.MealsQ(r)
	applyFilters(mealsQ, request)
	meals, err := mealsQ.Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get meals")
		ape.Render(w, problems.InternalError())
		return
	}
	categories, err := helpers.CategoriesQ(r).FilterById(getCategoryIds(meals)...).Select()

	response := resources.MealListResponse{
		Data:     newMealsList(meals),
		Links:    helpers.GetOffsetLinks(r, request.OffsetPageParams),
		Included: newMealIncluded(categories),
	}
	ape.Render(w, response)
}

func applyFilters(q data.MealsQ, request requests.GetMealListRequest) {
	q.Page(request.OffsetPageParams)

	if len(request.FilterMealName) > 0 {
		q.FilterByNames(request.FilterMealName...)
	}

	if len(request.FilterPriceFrom) > 0 {
		q.FilterByPriceFrom(request.FilterPriceFrom...)
	}

	if len(request.FilterPriceTo) > 0 {
		q.FilterByPriceTo(request.FilterPriceTo...)
	}

	if len(request.FilterAmount) > 0 {
		q.FilterByAmount(request.FilterAmount...)
	}

	if len(request.FilterCategoryId) > 0 {
		q.FilterByCategoryId(request.FilterCategoryId...)
	}
}

func newMealsList(meals []data.Meal) []resources.Meal {
	result := make([]resources.Meal, len(meals))
	for i, meal := range meals {
		result[i] = resources.Meal{
			Key: resources.NewKeyInt64(meal.Id, resources.MEAL),
			Attributes: resources.MealAttributes{
				MealName: meal.MealName,
				Price:    meal.Price,
				Amount:   meal.Amount,
			},
			Relationships: resources.MealRelationships{
				Category: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(meal.CategoryId, 10),
						Type: resources.CATEGORY,
					},
				},
			},
		}
	}
	return result
}

func getCategoryIds(meals []data.Meal) []int64 {
	categoryIDs := make([]int64, len(meals))
	for i := 0; i < len(meals); i++ {
		categoryIDs[i] = meals[i].CategoryId
	}
	return categoryIDs
}

func newMealIncluded(categories []data.Category) resources.Included {
	result := resources.Included{}
	for _, item := range categories {
		resource := newCategoryModel(item)
		result.Add(&resource)
	}
	return result
}

func newCategoryModel(category data.Category) resources.Category {
	return resources.Category{
		Key: resources.NewKeyInt64(category.Id, resources.CATEGORY),
		Attributes: resources.CategoryAttributes{
			CategoryName: category.CategoryName,
			Unit:         category.Unit,
		},
	}
}
