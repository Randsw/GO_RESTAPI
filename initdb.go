package main

import (
	"github.com/randsw/ha-postgres/records"
)
// Records
type initData [5]records.Record

var testPeople initData

func initialize() initData {
	testPeople = [5]records.Record{
		{
			Id:      1,
			Name:    "Michail",
			Surname: "Garibaldi",
			Gender:  "Male",
			Email:   "garibaldi@example.com",
		},
		{
			Id:      2,
			Name:    "Sasha",
			Surname: "Grey",
			Gender:  "Female",
			Email:   "grey@example.com",
		},
		{
			Id:      3,
			Name:    "Aaron",
			Surname: "Rodgers",
			Gender:  "Male",
			Email:   "rodgers@example.com",
		},
		{
			Id:      4,
			Name:    "Margo",
			Surname: "Robby",
			Gender:  "Female",
			Email:   "robby@example.com",
		},
		{
			Id:      5,
			Name:    "Cortney",
			Surname: "Cox",
			Gender:  "Female",
			Email:   "cox@example.com",
		},
	}
	return testPeople
}
