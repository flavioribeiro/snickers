package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"
	"github.com/pivotal-golang/lager"
	"github.com/snickers/snickers/core"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers"
)

// CreateJob creates a job
func (sn *SnickersServer) CreateJob(w http.ResponseWriter, r *http.Request) {
	log := sn.logger.Session("create-job")
	log.Debug("started")
	defer log.Debug("finished")

	dbInstance, err := db.GetDatabase()
	if err != nil {
		log.Error("failed-getting-database", err)
		HTTPError(w, http.StatusBadRequest, "getting database", err)
		return
	}

	var jobInput snickers.JobInput
	if err := json.NewDecoder(r.Body).Decode(&jobInput); err != nil {
		log.Error("failed-unpacking-job", err)
		HTTPError(w, http.StatusBadRequest, "unpacking job", err)
		return
	}

	preset, err := dbInstance.RetrievePreset(jobInput.PresetName)
	if err != nil {
		log.Error("failed-retrieving-preset", err)
		HTTPError(w, http.StatusBadRequest, "retrieving preset", err)
		return
	}

	var job snickers.Job
	job.ID = uniuri.New()
	job.Source = jobInput.Source
	job.Destination = jobInput.Destination
	job.Preset = preset
	job.Status = snickers.JobCreated
	_, err = dbInstance.StoreJob(job)
	if err != nil {
		log.Error("failed-storing-job", err)
		HTTPError(w, http.StatusBadRequest, "storing job", err)
		return
	}

	result, err := json.Marshal(job)
	if err != nil {
		log.Error("failed-packaging-job-data", err)
		HTTPError(w, http.StatusBadRequest, "packing job data", err)
		return
	}
	fmt.Fprintf(w, "%s", result)
	log.Info("created", lager.Data{"id": job.ID})
}

// ListJobs lists all jobs
func (sn *SnickersServer) ListJobs(w http.ResponseWriter, r *http.Request) {
	log := sn.logger.Session("list-jobs")
	log.Debug("started")
	defer log.Debug("finished")

	dbInstance, err := db.GetDatabase()
	if err != nil {
		log.Error("failed-getting-database", err)
		HTTPError(w, http.StatusBadRequest, "getting database", err)
		return
	}

	jobs, _ := dbInstance.GetJobs()
	result, err := json.Marshal(jobs)
	if err != nil {
		log.Error("failed-getting-jobs", err)
		HTTPError(w, http.StatusBadRequest, "getting jobs", err)
		return
	}

	fmt.Fprintf(w, "%s", string(result))
	log.Info("got-jobs")
}

// GetJobDetails returns the details of a given job
func (sn *SnickersServer) GetJobDetails(w http.ResponseWriter, r *http.Request) {
	log := sn.logger.Session("get-job-details")
	log.Debug("started")
	defer log.Debug("finished")

	dbInstance, err := db.GetDatabase()
	if err != nil {
		log.Error("failed-getting-database", err)
		HTTPError(w, http.StatusBadRequest, "getting database", err)
		return
	}

	vars := mux.Vars(r)
	jobID := vars["jobID"]
	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		log.Error("failed-retrieving-job", err)
		HTTPError(w, http.StatusBadRequest, "retrieving job", err)
		return
	}

	result, err := json.Marshal(job)
	if err != nil {
		log.Error("failed-packaging-job-data", err)
		HTTPError(w, http.StatusBadRequest, "packing job data", err)
		return
	}

	fmt.Fprintf(w, "%s", result)
	log.Info("got-job-details", lager.Data{"id": job.ID})
}

// StartJob triggers an encoding process
func (sn *SnickersServer) StartJob(w http.ResponseWriter, r *http.Request) {
	log := sn.logger.Session("start-job")
	log.Debug("started")
	defer log.Debug("finished")

	dbInstance, err := db.GetDatabase()
	if err != nil {
		log.Error("failed-getting-database", err)
		HTTPError(w, http.StatusBadRequest, "getting database", err)
		return
	}

	vars := mux.Vars(r)
	jobID := vars["jobID"]
	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		log.Error("failed-retrieving-job", err)
		HTTPError(w, http.StatusBadRequest, "retrieving job", err)
		return
	}

	log.Debug("starting-job", lager.Data{"id": job.ID})
	w.WriteHeader(http.StatusOK)
	go core.StartJob(job)
}
