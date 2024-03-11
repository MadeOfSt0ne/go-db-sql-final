package main

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := sq.Insert("parcel").
		Columns("client", "status", "address", "created_at").
		Values(p.Client, p.Status, p.Address, p.CreatedAt).
		RunWith(s.db).Exec()
	if err != nil {
		return 0, fmt.Errorf("error inserting parcel into db: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last inserted id: %v", err)
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	row := sq.Select("number", "client", "status", "address", "created_at").
		From("parcel").
		Where(sq.Eq{"number": number}).
		RunWith(s.db).QueryRow()
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	//return p, fmt.Errorf("error scanning row: %v", err)
	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := sq.Select("number", "client", "status", "address", "created_at").
		From("parcel").
		Where(sq.Eq{"client": client}).
		RunWith(s.db).Query()
	if err != nil {
		//return nil, fmt.Errorf("error getting parcels by client: %v", err)
		return nil, err
	}
	defer rows.Close()
	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			//return nil, fmt.Errorf("error scanning row: %v", err)
			return nil, err
		}
		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		//return nil, fmt.Errorf("rows error: %v", err)
		return nil, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := sq.Update("parcel").
		Set("status", status).
		Where(sq.Eq{"number": number}).
		RunWith(s.db).Exec()
	//return fmt.Errorf("error updating status: %v", err)
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := sq.Update("parcel").
		Set("address", address).
		Where(sq.Eq{"number": number, "status": ParcelStatusRegistered}).
		RunWith(s.db).Exec()
	//return fmt.Errorf("error updating address: %v", err)
	return err
}

func (s ParcelStore) Delete(number int) error {
	_, err := sq.Delete("parcel").
		Where(sq.Eq{"number": number, "status": ParcelStatusRegistered}).
		RunWith(s.db).Exec()
	//return fmt.Errorf("error deleting parcel: %v", err)
	return err
}
