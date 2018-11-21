package app

import "encoding/json"

type QueryRequest struct {
	Type  string `json:"type"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type QueryRequests struct {
	Queries []QueryRequest `json:"queries"`
}

func ParseQueryRequest(request []byte) (QueryRequests, error) {
	var qr QueryRequests

	if err := json.Unmarshal(request, &qr); err != nil {
		panic(err)
	}

	return qr, nil
}
