package repository

import (
	"github.com/go-pg/pg/orm"
	"github.com/pborman/uuid"
	"github.com/plusspeed/payments-api/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	// values from the docker compose
	pgUsername = "test"
	pgPassword = "example"
	pgAddress  = "127.0.0.1:5432"
)

func TestDatabase_Create_UPDATE_GET_DELETE(t *testing.T) {
	dbTest := New(pgAddress, "", pgUsername, pgPassword)
	defer dbTest.Database.Close()
	clearDB(*dbTest)

	var paymentID = uuid.NewRandom().String()
	const org1 = "1"
	const org2 = "2"

	p1 := &model.Payment{
		ID:             paymentID,
		OrganisationID: org1,
	}
	err := dbTest.Create(p1)
	assert.Nil(t, err)

	p2, err := dbTest.Get(paymentID)
	assert.Nil(t, err)

	assert.Equal(t, p1, p2, "should be equal %+v %+v", p1, p2)

	p3 := &model.Payment{
		ID:             paymentID,
		OrganisationID: org2,
	}

	err = dbTest.Update(p3)
	assert.Nil(t, err)

	p4, err := dbTest.Get(paymentID)
	assert.Nil(t, err)

	assert.Equal(t, p3, p4, "should be equal %+v %+v", p3, p4)

	err = dbTest.Delete(paymentID)
	assert.Nil(t, err)

	p5, err := dbTest.Get(paymentID)
	assert.Nil(t, p5)
	assert.NotNil(t, err)
}

func TestDatabase_NotFound(t *testing.T) {
	dbTest := New(pgAddress, "", pgUsername, pgPassword)
	defer dbTest.Database.Close()
	clearDB(*dbTest)

	var paymentID = uuid.NewRandom().String()

	_, err := dbTest.Get(paymentID)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNotFound, err, "should be equal %+v %+v", ErrNotFound, err)

	err = dbTest.Update(&model.Payment{
		ID: paymentID,
	})
	assert.NotNil(t, err)
	assert.Equal(t, ErrNotFound, err, "should be equal %+v %+v", ErrNotFound, err)

	err = dbTest.Delete(paymentID)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNotFound, err, "should be equal %+v %+v", ErrNotFound, err)
}

func TestDatabase_List(t *testing.T) {
	dbTest := New(pgAddress, "", pgUsername, pgPassword)
	defer dbTest.Database.Close()
	clearDB(*dbTest)

	// when there are no results, return ErrNotFound
	_, err := dbTest.List(0, 1)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNotFound, err, "should be equal %+v %+v", ErrNotFound, err)

	//Add data
	var paymentID1 = uuid.NewRandom().String()
	var paymentID2 = uuid.NewRandom().String()
	var paymentID3 = uuid.NewRandom().String()

	p1 := &model.Payment{
		ID: paymentID1,
	}
	err = dbTest.Create(p1)
	assert.Nil(t, err)

	p2 := &model.Payment{
		ID: paymentID2,
	}
	err = dbTest.Create(p2)
	assert.Nil(t, err)

	p3 := &model.Payment{
		ID: paymentID3,
	}
	err = dbTest.Create(p3)
	assert.Nil(t, err)

	//Test List
	ts1, err := dbTest.List(0, 5)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(ts1), "the length should be 3 instead of", len(ts1))

	ts2, err := dbTest.List(0, 3)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(ts2), "the length should be 3 instead of", len(ts2))

	_, err = dbTest.List(3, 5)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNotFound, err, "should be equal %+v %+v", ErrNotFound, err)

}

func clearDB(dbTest Repository) {
	err := dbTest.Database.DropTable(&model.Payment{}, &orm.DropTableOptions{
		IfExists: true,
		Cascade:  true,
	})

	err = dbTest.Database.CreateTable(&model.Payment{}, &orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		panic(err.Error())
	}
}
