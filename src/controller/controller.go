package controller

import (
	"encoding/json"
	"fProject/src/cases"
	"fProject/src/entity"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type ResultT struct {
	Status bool               `json:"status"`
	Data   *entity.ResultSetT `json:"data"`
	Error  string             `json:"error"`
}

type Controller struct {
	usecase cases.Usecase
}

func NewController(usecase cases.Usecase) *Controller {
	return &Controller{
		usecase: usecase,
	}
}

func Build(r *mux.Router, usecase cases.Usecase) {
	ctr := NewController(usecase)
	ctr.getResultT()
	fmt.Println("--work done--")
	r.HandleFunc("/", ctr.handleConnection)
}

func (c *Controller) handleConnection(w http.ResponseWriter, r *http.Request) {
	result := ResultT{}
	w.WriteHeader(200)
	gR, bo := c.getResultT()
	if bo == true {
		result = ResultT{bo, gR, ""}
	} else {
		result = ResultT{bo, gR, "Error on collect data"}
	}
	jResult, _ := json.MarshalIndent(result, "", " ")
	fmt.Fprint(w, string(jResult))
}

func (c *Controller) getResultT() (*entity.ResultSetT, bool) {
	done := make(chan struct{})
	sms := c.sms()
	mms := c.mmsGetReq()
	voice := c.voiceCall()
	email := c.email()
	billing := c.billing()
	sD := c.supportData()
	inc := c.incidentData()
	go func(d chan struct{}) {
		sms = c.sms()
		mms = c.mmsGetReq()
		voice = c.voiceCall()
		email = c.email()
		billing = c.billing()
		sD = c.supportData()
		inc = c.incidentData()
		close(d)
	}(done)
	<-done
	result, bo := c.usecase.MakeResultSetT(sms, mms, voice, email, billing, sD, inc)
	return result, bo
}

func (c *Controller) sms() [][]*entity.SMSData {
	data, err := ioutil.ReadFile("simulator/sms.data")
	if err != nil {
		log.Fatal(err)
	}
	sliceString := strings.Split(string(data), "\n")
	r := c.usecase.MakeSMSNote(sliceString)
	return r
}

func (c *Controller) mmsGetReq() [][]*entity.MMSData {
	resp, err := http.Get("http://127.0.0.1:8383/mms")
	if err != nil {
		log.Fatal(err)
	}
	if resp.Status == "200 OK" {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		r := c.usecase.MakeMMSNote(body)
		return r
	} else {
		return c.usecase.MakeMMSNote(nil)
	}
}

func (c *Controller) voiceCall() []*entity.VoiceData {
	data, err := ioutil.ReadFile("simulator/voice.data")
	if err != nil {
		log.Fatal(err)
	}
	sliceString := strings.Split(string(data), "\n")
	r := c.usecase.MakeVoiceNote(sliceString)
	return r
}

func (c *Controller) email() map[string][][]*entity.EmailData {
	data, err := ioutil.ReadFile("simulator/email.data")
	if err != nil {
		log.Fatal(err)
	}
	sliceString := strings.Split(string(data), "\n")
	r := c.usecase.MakeEmailNote(sliceString)
	return r
}

func (c *Controller) billing() *entity.BillingData {
	data, err := ioutil.ReadFile("simulator/billing.data")
	if err != nil {
		log.Fatal(err)
	}
	r := c.usecase.MakeBillingNote(data)
	return r
}

func (c *Controller) supportData() []int {
	resp, err := http.Get("http://127.0.0.1:8383/support")
	if err != nil {
		log.Fatal(err)
	}
	if resp.Status == "200 OK" {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		r := c.usecase.MakeSupportDataNote(body)
		return r
	} else {
		return c.usecase.MakeSupportDataNote(nil)
	}
}

func (c *Controller) incidentData() []*entity.IncidentData {
	resp, err := http.Get("http://127.0.0.1:8383/accendent")
	if err != nil {
		log.Fatal(err)
	}
	if resp.Status == "200 OK" {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		r := c.usecase.MakeIncidentDataNote(body)
		return r
	} else {
		return c.usecase.MakeIncidentDataNote(nil)
	}
}
