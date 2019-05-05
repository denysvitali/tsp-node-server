package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	log "github.com/denysvitali/gc_log"
	"net/http"
	"time"
)

type State struct {
	Db *sql.DB
}


type Params struct {
	Alpha float64 `json:"alpha"`
	R int32 `json:"r"`
	Start_Temp float64 `json:"start_temperature"`
}

type Algorithm struct {
	Name string `json:"name"`
	Seed int64 `json:"seed"`
	Mode string `json:"mode"`
	Params Params `json:"params"`
}

type TSPSol struct {
	From string `json:"from"`
	Commit string `json:"commit"`
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
 	fmt.Printf("Received request!\n");
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
	statement, err:= s.Db.Prepare("INSERT INTO tsp_results (problem, \"commit\", \"from\", time_elapsed, rl, \"json\") VALUES ($1," +
		"$2, $3, $4, $5, $6);")

	 if err != nil {
		 log.Error("Unable to create statement: ", err)
		 return
	 }

	_, err = statement.Exec(tspSol.Problem, tspSol.Commit, tspSol.From, tspSol.Time_Elapsed, tspSol.Rl, string(mjson))

	 if err != nil {
		 log.Error("Error is: ", err)
	 }

	 _ = statement.Close()
}

func (s State) GetResults(w http.ResponseWriter, r *http.Request) {
	var err error
	var result *sql.Rows
	var statement *sql.Stmt

	if 	r.URL.Query().Get("problem") != "" &&
		r.URL.Query().Get("commit") != "" {
		statement, err = s.Db.Prepare("SELECT id, \"commit\", \"from\", problem, time_elapsed, rl, json, received_on FROM tsp_results WHERE problem=$1 AND \"commit\"=$2 ORDER BY rl ASC, time_elapsed ASC");
		if err != nil {
			log.Error("Invalid prepared statement: ", err)
			return
		}

		result, err = statement.Query(r.URL.Query().Get("problem"),
			r.URL.Query().Get("commit"))

		if err != nil {
			log.Error("Unable to get results...", err)
			return
		}
	} else {
		statement, err = s.Db.Prepare("SELECT id, \"commit\", \"from\", problem, time_elapsed, rl, json, received_on FROM tsp_results WHERE problem=$1 ORDER BY rl ASC, time_elapsed ASC");
		if err != nil {
			log.Error("Invalid prepared statement: ", err)
			return
		}

		result, err = statement.Query(r.URL.Query().Get("problem"))

		if err != nil {
			log.Error("Unable to get results...", err)
			return
		}
	}

	if statement == nil {
		log.Error("Statement cannot be nil");
		return
	}

	var mstruct struct{
		Id int32 `json:"id"`
		From string `json:"from"`
		Problem string `json:"problem"`
		Commit *string `json:"commit"`
		Time_Elapsed int64 `json:"time_elapsed"`
		Rl int64 `json:"rl"`
		Json string `json:"json"`
		ReceivedOn time.Time `json:"time"`
	}
	result.Next()
	err = result.Scan(
		&mstruct.Id,
		&mstruct.Commit,
		&mstruct.From,
		&mstruct.Problem,
		&mstruct.Time_Elapsed,
		&mstruct.Rl,
		&mstruct.Json,
		&mstruct.ReceivedOn)

	if err != nil {
		log.Error("Unable to scan result", err)
	}

	mar, err := json.Marshal(mstruct)

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(mar)

	_ = statement.Close()

}
