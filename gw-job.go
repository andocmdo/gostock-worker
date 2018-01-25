package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	gostock "github.com/andocmdo/gostockd/common"
)

// Job is an alias for adding our own methods to the common gostock.Job struct
type Job gostock.Job

// NewJob is a constructor for Job structs (init Args map)
func NewJob() *Job {
	var j Job
	j.Args = make(map[string]string)
	return &j
}

func (job *Job) setRunning(master *Server, wrkr *Worker) error {
	job.Running = true
	job.WorkerID = wrkr.ID

	log.Printf("called setRunning")
	log.Printf("%+v", wrkr)
	log.Printf("%+v", job)

	jsonWorker, _ := json.Marshal(*job)
	resp, err := http.Post(master.URLjobs+"/"+strconv.Itoa(job.ID), jsonData, bytes.NewBuffer(jsonWorker))
	//resp, err := http.PostForm(requestURL, url.Values{"port": {sPort}})
	if err != nil {
		//log.Printf("worker %d: error setting READY with master server", wn)
		//log.Println(err)
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//log.Printf(err.Error())
		return err
	}
	resp.Body.Close()
	if err = json.Unmarshal(body, job); err != nil {
		//log.Printf(err.Error())
		return err
	}
	if job.Valid != true {
		//log.Printf("worker %d: master server returned worker object with false VALID flag when setting READY!", wn)
		return errors.New("master server response was returned as invalid")
	}
	master.Valid = true
	master.LastContact = time.Now()
	master.LastUpdate = time.Now()

	log.Printf("%+v", wrkr)
	log.Printf("%+v", job)
	log.Printf("exiting setRunning")

	return nil
}
