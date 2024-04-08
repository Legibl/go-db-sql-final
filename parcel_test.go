package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// get

	res, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, res.Client, parcel.Client)
	assert.Equal(t, res.Status, parcel.Status)
	assert.Equal(t, res.Address, parcel.Address)
	assert.Equal(t, res.CreatedAt, parcel.CreatedAt)

	// delete

	err = store.Delete(id)
	require.NoError(t, err)

	res, err = store.Get(parcel.Number)
	require.Error(t, err, sql.ErrNoRows)
	require.Empty(t, res)
	assert.ErrorIs(t, err, sql.ErrNoRows)

}

func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()
	// add

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// set address

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// check

	res, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, res.Address, "new test address")
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// set status

	err = store.SetStatus(id, "delivered")
	require.NoError(t, err)
	// check

	res, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, res.Status, "delivered")
}

func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}

	client := randRange.Intn(10_000_000)
	for i, p := range parcels {
		p.Client = client
		p.Number, err = store.Add(p)
		require.NoError(t, err)
		parcels[i] = p
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.ElementsMatch(t, storedParcels, parcels)
}
