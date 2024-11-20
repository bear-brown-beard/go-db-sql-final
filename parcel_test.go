package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	// randSource - источник псевдослучайных чисел.
	// Для повышения уникальности в качестве seed используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
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

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	// создаём подключение к in-memory базе данных SQLite
	db, err := sql.Open("sqlite", ":memory:") // in-memory база данных
	require.NoError(t, err)

	// создаём таблицу
	_, err = db.Exec(`CREATE TABLE parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER,
		status TEXT,
		address TEXT,
		created_at TEXT
	)`)
	require.NoError(t, err)

	// создаём объект ParcelStore
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавляем посылку в БД и получаем её ID
	id, err := store.Add(parcel)
	require.NoError(t, err)
	parcel.Number = id

	// get
	// получаем только что добавленную посылку по номеру
	storedParcel, err := store.Get(parcel.Number)
	require.NoError(t, err)
	require.Equal(t, parcel, storedParcel)

	// delete
	// удаляем посылку из БД
	err = store.Delete(parcel.Number)
	require.NoError(t, err)

	// проверяем, что посылка удалена
	_, err = store.Get(parcel.Number)
	require.Error(t, err) // должна быть ошибка, так как посылка удалена
}

// TestSetAddress проверяет обновление адреса посылки
func TestSetAddress(t *testing.T) {
	// prepare
	// создаём подключение к in-memory базе данных SQLite
	db, err := sql.Open("sqlite", ":memory:") // in-memory база данных
	require.NoError(t, err)

	// создаём таблицу
	_, err = db.Exec(`CREATE TABLE parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER,
		status TEXT,
		address TEXT,
		created_at TEXT
	)`)
	require.NoError(t, err)

	// создаём объект ParcelStore и добавляем посылку
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// добавляем посылку в БД и получаем её ID
	id, err := store.Add(parcel)
	require.NoError(t, err)
	parcel.Number = id

	// set address
	// изменяем адрес посылки
	newAddress := "new test address"
	err = store.SetAddress(parcel.Number, newAddress)
	require.NoError(t, err)

	// check
	// проверяем, что адрес посылки обновился
	storedParcel, err := store.Get(parcel.Number)
	require.NoError(t, err)
	require.Equal(t, newAddress, storedParcel.Address)
}

// TestSetStatus проверяет обновление статуса посылки
func TestSetStatus(t *testing.T) {
	// prepare
	// создаём подключение к in-memory базе данных SQLite
	db, err := sql.Open("sqlite", ":memory:") // in-memory база данных
	require.NoError(t, err)

	// создаём таблицу
	_, err = db.Exec(`CREATE TABLE parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER,
		status TEXT,
		address TEXT,
		created_at TEXT
	)`)
	require.NoError(t, err)

	// создаём объект ParcelStore и добавляем посылку
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// добавляем посылку в БД и получаем её ID
	id, err := store.Add(parcel)
	require.NoError(t, err)
	parcel.Number = id

	// set status
	// изменяем статус посылки на "sent"
	err = store.SetStatus(parcel.Number, ParcelStatusSent)
	require.NoError(t, err)

	// check
	// проверяем, что статус посылки обновился
	storedParcel, err := store.Get(parcel.Number)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, storedParcel.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	// создаём подключение к in-memory базе данных SQLite
	db, err := sql.Open("sqlite", ":memory:") // in-memory база данных
	require.NoError(t, err)

	// создаём таблицу
	_, err = db.Exec(`CREATE TABLE parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER,
		status TEXT,
		address TEXT,
		created_at TEXT
	)`)
	require.NoError(t, err)

	// создаём объект ParcelStore и добавляем несколько посылок для одного клиента
	store := NewParcelStore(db)
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	client := randRange.Intn(10_000_000)
	for i := range parcels {
		parcels[i].Client = client
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		parcels[i].Number = id
	}

	// get by client
	// получаем все посылки для данного клиента
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)

	// check all parcels for the client
	// проверяем, что количество полученных посылок совпадает с количеством добавленных
	require.Len(t, storedParcels, len(parcels)) 

	// проверяем, что все посылки из storedParcels присутствуют в оригинальном списке
	for _, storedParcel := range storedParcels {
		require.Contains(t, parcels, storedParcel) // проверка наличия посылки в исходных данных
	}
}
