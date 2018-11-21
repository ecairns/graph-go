package app

import (
	"errors"
	"fmt"
)

type Node struct {
	Id   string `db:"id" xml:"id"`
	Name string `db:"name" xml:"name"`
}

func (this *Node) Save() bool {
	var query string
	var err error

	errs := this.Validate()
	if errs != nil {
		return false
	}
	query = "INSERT INTO node (id, name) SELECT $1, $2 WHERE NOT EXISTS (SELECT 1 FROM node WHERE id=$3);"
	_, err = DB.Exec(query,
		this.Id,
		this.Name,
		this.Id,
	)
	if err != nil {
		fmt.Printf("NODE SAVE ERROR: %v\n", err)
		return false
	}

	return true
}

func (this *Node) Validate() []error {
	var errs []error

	if this.Id == "" {
		errs = append(errs, errors.New("Node ID cannot be empty"))
	}

	if this.Name == "" {
		this.Name = this.Id
	}

	return errs
}
