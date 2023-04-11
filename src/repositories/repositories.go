package repositories

import (
	"encoding/json"
	"fProject/src/entity"
	"fmt"
)

type repositories struct {
	SMSCollection          []*entity.SMSData
	MMSCollection          []*entity.MMSData
	VoiceCollection        []*entity.VoiceData
	EmailCollection        []*entity.EmailData
	SupportDataCollection  []*entity.SupportData
	IncidentDataCollection []*entity.IncidentData
	ResultSetT             *entity.ResultSetT
}

func NewRepository() *repositories {
	return &repositories{
		SMSCollection:          make([]*entity.SMSData, 0),
		MMSCollection:          make([]*entity.MMSData, 0),
		VoiceCollection:        make([]*entity.VoiceData, 0),
		EmailCollection:        make([]*entity.EmailData, 0),
		SupportDataCollection:  make([]*entity.SupportData, 0),
		IncidentDataCollection: make([]*entity.IncidentData, 0),
	}
}

func (r repositories) MakeIncidentDataNote(b []byte) []*entity.IncidentData {
	mms := &r.IncidentDataCollection
	err := json.Unmarshal(b, &mms)
	if err != nil {
		fmt.Print(r.IncidentDataCollection)
	}
	return r.IncidentDataCollection
}

func (r repositories) MakeSupportDataNote(b []byte) []*entity.SupportData {
	mms := &r.SupportDataCollection
	err := json.Unmarshal(b, &mms)
	if err != nil {
		fmt.Print(r.MMSCollection)
	}
	return r.SupportDataCollection
}

func (r repositories) MakeBillingNote(b []bool) *entity.BillingData {
	data := &entity.BillingData{b[0], b[1], b[2], b[3], b[4], b[5]}
	return data
}

func (r repositories) MakeEmailNote(s []string, dT int) []*entity.EmailData {
	data := &entity.EmailData{s[0], s[1], dT}
	r.EmailCollection = append(r.EmailCollection, data)
	return r.EmailCollection
}

func (r repositories) MakeVoiceNote(s []string, cS float32, t int, vP int, mOC int) []*entity.VoiceData {
	data := &entity.VoiceData{s[0], s[1], s[2], s[3], cS, t, vP, mOC}
	r.VoiceCollection = append(r.VoiceCollection, data)
	return r.VoiceCollection
}

func (r repositories) MakeMMSNote(b []byte) []*entity.MMSData {
	mms := &r.MMSCollection
	err := json.Unmarshal(b, &mms)
	if err != nil {
		fmt.Print(r.MMSCollection)
	}
	return r.MMSCollection
}

func (r repositories) MakeSMSNote(s []string) []*entity.SMSData {
	data := &entity.SMSData{s[0], s[1], s[2], s[3]}
	r.SMSCollection = append(r.SMSCollection, data)
	return r.SMSCollection
}
