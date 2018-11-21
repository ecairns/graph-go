package app

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type result struct {
	Type  string      `json:"type"`
	From  string      `json:"from"`
	To    string      `json:"to"`
	Path  interface{} `json:"path"`
	Score float32     `json:"score,omitempty"`
}

type resultSet struct {
	Answers []result `json:"answers"`
}

func Query(qrs QueryRequests) (string, error) {
	var rs resultSet

	for _, v := range qrs.Queries {
		//q := qr.Queries[k];
		r, err := execQuery(v)
		if err != nil {
			continue
		}

		rs.Answers = append(rs.Answers, r)
	}
	jsonBytes, _ := json.MarshalIndent(rs, "", "  ")
	return string(jsonBytes), nil
}

func execQuery(qr QueryRequest) (result, error) {
	var result result
	var query string
	var params []interface{}

	initQuery := `
WITH RECURSIVE paths(to_node, path, total, already_visited, cycle_detected) AS(
   SELECT e.to_node, concat(e.from_node, ',', e.to_node), cost, ARRAY[e.from_node], false 	
   FROM graph g, graph_edge ge, edge e
   WHERE ge.graph_id=g.id AND e.id=ge.edge_id AND e.from_node='a'
 UNION ALL
   SELECT e.to_node, concat(path, ',', e.to_node), total+e.cost, already_visited, e.from_node = ANY(already_visited)
   FROM graph g, graph_edge ge, edge e, paths p
   WHERE ge.graph_id=g.id AND e.id=ge.edge_id AND p.to_node = e.from_node  AND NOT cycle_detected
) SELECT path, total FROM paths p WHERE NOT cycle_detected`
	query = initQuery
	if qr.Start != "" && qr.End != "" {
		query += " AND (p.path LIKE $1 OR p.path=$2)"
		params = append(params, fmt.Sprintf("%v,%%,%v", qr.Start, qr.End))
		params = append(params, fmt.Sprintf("%v,%v", qr.Start, qr.End))
	} else if qr.Start != "" {
		query += " AND p.path LIKE $1"
		params = append(params, fmt.Sprintf("%v,%%", qr.Start))
	} else if qr.End != "" {
		query += " AND p.path LIKE $1"
		params = append(params, fmt.Sprintf("%%,%v", qr.End))
	}

	if qr.Type == "cheapest" {
		query += " ORDER BY total LIMIT 1"
	}

	queryStmt, err := DB.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := queryStmt.Query(params...)
	defer rows.Close()

	result.From = qr.Start
	result.To = qr.End
	result.Type = qr.Type
	result.Path = false

	var Paths [][]string
	for rows.Next() {
		var path string
		var total float32

		if err := rows.Scan(&path, &total); err != nil {
			return result, err
		}
		if qr.Type == "cheapest" {
			result.Path = strings.Split(path, ",")
			result.Score = total

		} else {
			Paths = append(Paths, strings.Split(path, ","))

		}
	}
	if len(Paths) > 0 {
		result.Path = Paths
	}

	return result, nil
}
