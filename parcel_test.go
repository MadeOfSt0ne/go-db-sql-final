package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	newStatus = "new test status"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
	db         *sql.DB
	err        error
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func TestMain(m *testing.M) {
	db, err = sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	exitCode := m.Run()
	os.Exit(exitCode)
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	store := NewParcelStore(db)
	parcel := getTestParcel()

	parcel.Number, err = store.Add(parcel)
	require.NoError(t, err)
	require.Positive(t, parcel.Number)

	p, err := store.Get(parcel.Number)
	require.NoError(t, err)
	require.Equal(t, parcel, p)

	err = store.Delete(parcel.Number)
	require.NoError(t, err)
	_, err = store.Get(parcel.Number)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	store := NewParcelStore(db)
	parcel := getTestParcel()

	parcel.Number, err = store.Add(parcel)
	require.NoError(t, err)
	require.Positive(t, parcel.Number)

	newAddress := "new test address"
	err = store.SetAddress(parcel.Number, newAddress)
	require.NoError(t, err)

	updated, err := store.Get(parcel.Number)
	require.NoError(t, err)
	require.Equal(t, newAddress, updated.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	store := NewParcelStore(db)
	parcel := getTestParcel()

	parcel.Number, err = store.Add(parcel)
	require.NoError(t, err)
	require.Positive(t, parcel.Number)

	err = store.SetStatus(parcel.Number, newStatus)
	require.NoError(t, err)

	updated, err := store.Get(parcel.Number)
	require.NoError(t, err)
	require.Equal(t, newStatus, updated.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.Positive(t, id)

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, len(parcels), len(storedParcels))

	for _, parcel := range storedParcels {
		storedParcel, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		require.Equal(t, storedParcel, parcel)
	}
}
