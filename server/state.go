package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	log "github.com/denysvitali/gc_log"
	"net/http"
)

type State struct {
	Db *sql.DB
}


type Params struct {
	Alpha float64 `json:"alpha"`
	R int32 `json:"r"`
	Start_Temp int32 `json:"start_temperature"`
}

type Algorithm struct {
	Name string `json:"name"`
	Seed int64 `json:"seed"`
	Mode string `json:"mode"`
	Params Params `json:"params"`
}

type TSPSol struct {
	From string `json:"from"`
	Problem string `json:"problem"`
	Time_Elapsed int64 `json:"time_elapsed"`
	Algorithms []Algorithm `json:"algorithms"`
	Rl int `json:"rl"`
	Bk int `json:"bk"`
	Perf float64 `json:"perf"`
}

type GenericResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
}

 func (s State) UploadJson(w http.ResponseWriter, r *http.Request) {
	var tspSol TSPSol
	err := json.NewDecoder(r.Body).Decode(&tspSol)

	if err != nil {
		errorResult := GenericResponse{
			Success: false,
			Message: "Invalid JSON!",
		}
		fmt.Printf("Error is %v\n", err)
		_ = json.NewEncoder(w).Encode(errorResult)
		return
	}

	_ = json.NewEncoder(w).Encode(GenericResponse{
		true, "OK",
	})


	mjson, _ := json.Marshal(tspSol)

	fmt.Printf("TSP Sol: %v\n", tspSol)
	statement, err:= s.Db.Prepare("INSERT INTO tsp_results (problem, \"from\", time_elapsed, rl, \"json\") VALUES ($1," +
		"$2, $3, $4, $5);")

	 if err != nil {
		 log.Error("Unable to create statement: ", err)
		 return
	 }

	_, err = statement.Exec(tspSol.Problem, tspSol.From, tspSol.Time_Elapsed, tspSol.Rl, string(mjson))

	 if err != nil {
		 log.Error("Error is: ", err)
	 }
}