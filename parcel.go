package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	query := "INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)"
	result, err := s.db.Exec(query, p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {

	query := "SELECT * FROM parcel WHERE number = ?"
	res := s.db.QueryRow(query, number)

	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	err := res.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	query := "SELECT * FROM parcel WHERE Client = ?"
	rows, err := s.db.Query(query, client)

	if err != nil {
		return nil, fmt.Errorf("error getting data about client with id %d: %v", client, err)
	}
	defer rows.Close()

	var parcels []Parcel

	for rows.Next() {
		var p Parcel

		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		parcels = append(parcels, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return parcels, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {

	query := "UPDATE parcel SET status = ? WHERE number = ?"
	_, err := s.db.Exec(query, status, number)
	if err != nil {
		return fmt.Errorf("error updating status of parcel: %v", err)
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {

	query := "UPDATE parcel SET address = ? WHERE number = ? AND status = ?"

	res, err := s.db.Exec(query, address, number, ParcelStatusRegistered)
	if err != nil {
		return fmt.Errorf("error updating address of parcel: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getiing rows affected")
	}

	if rowsAffected == 0 {
		return fmt.Errorf("unable to change addres for parcel, incorrect status or parcel not found")
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {

	query := "DELETE FROM parcel WHERE number = ? and status = ?"

	res, err := s.db.Exec(query, number, ParcelStatusRegistered)
	if err != nil {
		return fmt.Errorf("failed to delete parcel: %s", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %s", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("unable to delete parcel, incorrect status or parcel not found")
	}

	return nil
}
