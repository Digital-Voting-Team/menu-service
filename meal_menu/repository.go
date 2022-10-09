package meal_menu

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	queryCreateTable = `CREATE TABLE IF NOT EXISTS public.meal_menu
	(
		id integer NOT NULL GENERATED BY DEFAULT AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
		meal integer NOT NULL,
		menu integer NOT NULL,
		CONSTRAINT meal_menu_pkey PRIMARY KEY (id),
		CONSTRAINT meal_menu_meal_fkey FOREIGN KEY (meal)
			REFERENCES public.meals (id) MATCH SIMPLE
			ON UPDATE NO ACTION
			ON DELETE NO ACTION,
		CONSTRAINT meal_menu_menu_fkey FOREIGN KEY (menu)
			REFERENCES public.menu (id) MATCH SIMPLE
			ON UPDATE NO ACTION
			ON DELETE NO ACTION
	)
	
	TABLESPACE pg_default;
	
	ALTER TABLE IF EXISTS public.meal_menu
    OWNER to postgres;`

	queryDeleteTable = `DROP TABLE public.meal_menu`

	queryInsert = `INSERT INTO public.meal_menu(meal, menu)
	VALUES ($1, $2) RETURNING id;`

	querySelect = `SELECT * FROM public.meal_menu;`

	queryUpdate = `UPDATE public.meal_menu
	SET meal=$2, menu=$3
	WHERE id=$1;`

	queryDelete = `DELETE FROM public.meal_menu
	WHERE id=$1;`

	queryCleanDb = `DELETE FROM public.meal_menu;`

	queryResetCounter = `alter sequence meal_menu_id_seq restart with 1`
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) Insert(mealMenu *MealMenu) (int, error) {
	rows, err := repo.db.Queryx(queryInsert, mealMenu.MealId, mealMenu.MenuId)
	defer rows.Close()
	id := -1
	if err != nil {
		return id, err
	}

	rows.Next()
	err = rows.Scan(&id)
	return id, err
}

func (repo *Repository) CreateTable() error {
	_, err := repo.db.Exec(queryCreateTable)
	return err
}

func (repo *Repository) DeleteTable() error {
	_, err := repo.db.Exec(queryDeleteTable)
	return err
}

func (repo *Repository) Select() ([]MealMenu, error) {
	rows, err := repo.db.Queryx(querySelect)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	mealMenu := MealMenu{}
	var mealMenuArray []MealMenu
	for rows.Next() {
		err = rows.StructScan(&mealMenu)
		if err != nil {
			return nil, err
		}
		mealMenuArray = append(mealMenuArray, mealMenu)
	}
	return mealMenuArray, err
}

func (repo *Repository) Delete(id int) error {
	_, err := repo.db.Exec(queryDelete, id)
	return err
}

func (repo *Repository) Update(id int, mealMenu *MealMenu) error {
	_, err := repo.db.Queryx(queryUpdate, id, mealMenu.MealId, mealMenu.MenuId)
	return err
}

func (repo *Repository) Clean() error {
	_, err := repo.db.Exec(queryCleanDb)
	return err
}

func (repo *Repository) ResetCounter() error {
	_, err := repo.db.Exec(queryResetCounter)
	return err
}
