package app

import (
	"errors"
	"fmt"
)

type Edge struct {
	Id   string  `xml:"id"`
	From string  `xml:"from"`
	To   string  `xml:"to"`
	Cost float32 `xml:"cost"` //default will be 0
}

func (this *Edge) Save() bool {
	var query string
	var err error

	errs := this.Validate()
	if errs != nil {
		return false
	}

	query = "INSERT INTO edge (id, from_node, to_node, cost) SELECT $1, $2, $3, $4 WHERE NOT EXISTS (SELECT 1 FROM edge WHERE id=$5 or (from_node=$6 AND to_node=$7));"
	_, err = DB.Exec(query,
		this.Id,
		this.From,
		this.To,
		this.Cost,
		this.Id,
		this.From,
		this.To,
	)
	if err != nil {
		fmt.Printf("edge INSERT ERROR: %v\n", err)
		return false
	}

	return true
}

func (this *Edge) Validate() []error {
	var errs []error

	if this.Id == "" {
		errs = append(errs, errors.New("Edge ID cannot be empty"))
	}

	if this.From == "" {
		errs = append(errs, errors.New("Edge From cannot be empty"))
	}

	if this.To == "" {
		errs = append(errs, errors.New("Edge To cannot be empty"))
	}

	if this.Cost < 0 {
		errs = append(errs, errors.New("Edge Cost cannot be negative"))
	}

	return errs
}
