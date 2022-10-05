package menu

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	queryCreateTable = `CREATE TABLE IF NOT EXISTS public.menu
	(
		id integer NOT NULL GENERATED BY DEFAULT AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
		cafe integer NOT NULL,
		CONSTRAINT menu_pkey PRIMARY KEY (id)
	)
	
	TABLESPACE pg_default;
	
	ALTER TABLE IF EXISTS public.menu
    OWNER to postgres;`

	queryDeleteTable = `DROP TABLE public.menu`

	queryInsert = `INSERT INTO public.menu(cafe)
	VALUES ($1) RETURNING id;`

	querySelect = `SELECT * FROM public.menu;`

	queryUpdate = `UPDATE public.menu
	SET cafe=$2
	WHERE id=$1;`

	queryDelete = `DELETE FROM public.menu
	WHERE id=$1;`

	queryCleanDb = `DELETE FROM public.menu;`
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) Insert(menu *Menu) (int, error) {
	rows, err := repo.db.Queryx(queryInsert, menu.CafeId)
	defer rows.Close()
	id := -1
	if err != nil {
		return id, err
	}

	rows.Next()
	err = rows.Scan(&id)
	return id, nil
}

func (repo *Repository) CreateTable() error {
	_, err := repo.db.Exec(queryCreateTable)
	return err
}

func (repo *Repository) DeleteTable() error {
	_, err := repo.db.Exec(queryDeleteTable)
	return err
}

func (repo *Repository) Select() ([]Menu, error) {
	rows, err := repo.db.Queryx(querySelect)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	menu := Menu{}
	var menuArray []Menu
	for rows.Next() {
		err = rows.StructScan(&menu)
		if err != nil {
			return nil, err
		}
		menuArray = append(menuArray, menu)
	}
	return menuArray, err
}

func (repo *Repository) Delete(id int) error {
	_, err := repo.db.Exec(queryDelete, id)
	return err
}

func (repo *Repository) Update(id int, menu *Menu) error {
	_, err := repo.db.Queryx(queryUpdate, id, menu.CafeId)
	return err
}

func (repo *Repository) Clean() error {
	_, err := repo.db.Exec(queryCleanDb)
	return err
}
