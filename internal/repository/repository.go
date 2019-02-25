package repository

import (
	"errors"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/plusspeed/payments-api/internal/model"
	"time"
)

//Repository has a postgres sql connection and implements the PaymentTransaction interface
type Repository struct {
	Database pg.DB
}

//PaymentTransaction contains all the DB operations for a Payment
type PaymentTransaction interface {
	Get(id string) (*model.Payment, error)
	Create(*model.Payment) error
	Update(*model.Payment) error
	Delete(id string) error
	List(limit, start int) ([]model.Payment, error)
}

//ErrNotFound is returned when no payment is returned
var ErrNotFound = errors.New("payment not found")

//New connects to a postgres sql db and creates Payment if does not exist.
func New(pgAddress, dbName, pgUsername, pgPassword string) *Repository {
	db := pg.Connect(&pg.Options{
		Addr:     pgAddress,
		User:     pgUsername,
		Password: pgPassword,
		Database: dbName,
	}).WithTimeout(5 * time.Second)

	err := createSchema(db)
	if err != nil {
		panic(err)
	}
	return &Repository{Database: *db}
}

func createSchema(db *pg.DB) error {
	for _, m := range []interface{}{(*model.Payment)(nil)} {
		err := db.CreateTable(m, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

//Get returns a model.Payment
//ErrNoRows if not found
func (d *Repository) Get(id string) (*model.Payment, error) {
	payment := &model.Payment{ID: id}
	err := d.Database.Select(payment)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return payment, nil
}

//Create inserts a model.Payment
func (d *Repository) Create(payment *model.Payment) error {
	err := d.Database.Insert(payment)
	if err != nil {
		return err
	}
	return nil
}

//Update modify an existing model.Payment
//ErrNoRows if not found
func (d *Repository) Update(m *model.Payment) error {
	err := d.Database.Update(m)
	if err != nil {
		if err == pg.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	return nil
}

//Delete deletes an existing model.Payment
//ErrNoRows if not found
func (d *Repository) Delete(id string) error {
	payment := &model.Payment{ID: id}
	err := d.Database.Delete(payment)
	if err != nil {
		if err == pg.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	return nil
}

//List returns a list of model.Payment for a offsset and limit orderly by ID Desc
//ErrNoRows if not found
func (d *Repository) List(offset, limit int) ([]model.Payment, error) {
	var ts []model.Payment
	err := d.Database.Model(&ts).
		Offset(offset).Limit(limit).Order("id DESC").
		Select()
	if err != nil {
		return nil, err
	}

	if len(ts) == 0 {
		return nil, ErrNotFound
	}
	return ts, nil
}
