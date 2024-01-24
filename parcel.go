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
		Suffix("RETURNING \"id\""). // можно написать RETURNING number
		RunWith(s.db).Exec()
	if err != nil {
		fmt.Printf("Error inserting parcel into db: %v", err)
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("Error getting last inserted id: %v", err)
		return 0, err
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
	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := sq.Select("number", "client", "status", "address", "created_at").
		From("parcel").
		Where(sq.Eq{"client": client}).
		RunWith(s.db).Query()
	if err != nil {
		fmt.Printf("Error getting parcels by client: %v", err)
		return nil, err
	}
	defer rows.Close()
	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := sq.Update("parcel").
		Set("status", status).
		Where(sq.Eq{"number": number}).
		RunWith(s.db).Exec()
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := sq.Update("parcel").
		Set("address", address).
		Where(sq.Eq{"number": number, "status": ParcelStatusRegistered}).
		RunWith(s.db).Exec()
	return err
}

func (s ParcelStore) Delete(number int) error {
	_, err := sq.Delete("parcel").
		Where(sq.Eq{"number": number, "status": ParcelStatusRegistered}).
		RunWith(s.db).Exec()
	return err
}
