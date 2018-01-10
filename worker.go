package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Worker contains state data for workers
type Worker struct {
	ID         int       `json:"id"`
	JobID      int       `json:"jobID"` // updateable
	Valid      bool      `json:"valid"`
	Created    time.Time `json:"created"`
	IPAddr     string    `json:"ipaddr"`
	Port       string    `json:"port"`
	Ready      bool      `json:"ready"`      // updateable
	Working    bool      `json:"working"`    // updateable
	LastUpdate time.Time `json:"lastUpdate"` // updateable
}

func (wrkr *Worker) run(wn int, master Server) {
	log.Printf("worker thread %d: started", wn)

	// register worker with gostockd api server
	if err := wrkr.register(&master); err != nil {
		// handle error
		log.Printf("worker %d: Error registering with master server", wn)
		log.Printf(err.Error())
	}
	log.Printf("worker %d: Successfully registered with master server", wn)

	// open listening port for jobs
	job := NewJob() // initialize an empty job to place incoming JSON job
	// TODO 0mq stuff here

	// loop here
	for {
		// update server showing this worker as ready to accept jobs
		if err := wrkr.setReady(&master); err != nil {
			//handle error
			log.Printf("worker %d: Error setting READY with master server", wn)
			log.Printf(err.Error())
		}
		log.Printf("worker %d: Successfully notified master server, READY to accept jobs", wn)

		// wait/listen to port for incoming jobs
		time.Sleep(time.Second * 15)

		// decode incoming job
		log.Printf("worker %d: recieved job", wn)
		time.Sleep(time.Second * 5)

		// start job
		log.Printf("worker %d: started job", wn)
		time.Sleep(time.Second * 5)

		// update server that we are working, and that job is running on this worker
		if err := wrkr.setWorking(&master, job); err != nil {
			//handle error
			log.Printf("worker %d: Error setting READY with master server", wn)
			log.Printf(err.Error())
		}
		log.Printf("worker %d: updated master for running job NUMBER WHAT?", wn)
		time.Sleep(time.Second * 5)

		// wait for job to finish
		log.Printf("worker %d: completed job", wn)
		time.Sleep(time.Second * 5)

		// update server of job completion status
		log.Printf("worker %d: updated master of successful job completion", wn)
	}
}

func (wrkr *Worker) register(master *Server) error {
	jsonWorker, _ := json.Marshal(wrkr)
	resp, err := http.Post(master.URLworkers, jsonData, bytes.NewBuffer(jsonWorker))
	//resp, err := http.PostForm(requestURL, url.Values{"port": {wrkr.Port}})
	if err != nil {
		//log.Printf("worker %d: error registering with master server", wn)
		log.Print("error registering worker with master server: ", err.Error())
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("error registering worker with master server: ", err.Error())
		return err
	}
	resp.Body.Close()
	if err = json.Unmarshal(body, &wrkr); err != nil {
		log.Print("error registering worker with master server: ", err.Error())
		return err
	}
	if wrkr.Valid != true {
		//log.Printf("worker %d: master server returned worker object with valid flag set false!", wn)
		return errors.New("master server response was returned as invalid")
	}
	//log.Printf("worker %d: Successfully registered with master server", wn)
	master.Valid = true
	master.LastContact = time.Now()
	master.LastUpdate = time.Now()

	return nil
}

func (wrkr *Worker) setReady(master *Server) error {
	// set this worker as ready to accept jobs
	wrkr.Ready = true
	jsonWorker, _ := json.Marshal(wrkr)
	resp, err := http.Post(master.URLworkers, jsonData, bytes.NewBuffer(jsonWorker))
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
	if err = json.Unmarshal(body, &wrkr); err != nil {
		//log.Printf(err.Error())
		return err
	}
	if wrkr.Valid != true {
		//log.Printf("worker %d: master server returned worker object with false VALID flag when setting READY!", wn)
		return errors.New("master server response was returned as invalid")
	}
	master.Valid = true
	master.LastContact = time.Now()
	master.LastUpdate = time.Now()

	return nil
}

func (wrkr *Worker) setWorking(master *Server, job *Job) error {
	// set this worker as ready to accept jobs
	wrkr.Ready = false
	wrkr.Working = true
	wrkr.JobID = job.ID
	jsonWorker, _ := json.Marshal(wrkr)
	resp, err := http.Post(master.URLworkers, jsonData, bytes.NewBuffer(jsonWorker))
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
	if err = json.Unmarshal(body, &wrkr); err != nil {
		//log.Printf(err.Error())
		return err
	}
	if wrkr.Valid != true {
		//log.Printf("worker %d: master server returned worker object with false VALID flag when setting READY!", wn)
		return errors.New("master server response was returned as invalid")
	}
	master.Valid = true
	master.LastContact = time.Now()
	master.LastUpdate = time.Now()

	return nil
}

func (job *Job) setRunning(master *Server, wrkr *Worker) error {
	return nil
}