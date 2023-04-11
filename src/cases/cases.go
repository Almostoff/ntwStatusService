package cases

import (
	"fProject/src/countryAlpha2"
	"fProject/src/entity"
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Usecase interface {
	MakeSMSNote(s []string) [][]*entity.SMSData
	MakeMMSNote(b []byte) [][]*entity.MMSData
	MakeVoiceNote(s []string) []*entity.VoiceData
	MakeEmailNote(s []string) map[string][][]*entity.EmailData
	MakeBillingNote(b []byte) *entity.BillingData
	MakeSupportDataNote(b []byte) []int
	MakeIncidentDataNote(b []byte) []*entity.IncidentData
	MakeResultSetT([][]*entity.SMSData, [][]*entity.MMSData, []*entity.VoiceData, map[string][][]*entity.EmailData, *entity.BillingData, []int, []*entity.IncidentData) (*entity.ResultSetT, bool)
}

type Repository interface {
	MakeSMSNote(s []string) []*entity.SMSData
	MakeMMSNote(b []byte) []*entity.MMSData
	MakeVoiceNote(s []string, f float32, i1 int, i2 int, i3 int) []*entity.VoiceData
	MakeEmailNote(s []string, i int) []*entity.EmailData
	MakeBillingNote(b []bool) *entity.BillingData
	MakeSupportDataNote(b []byte) []*entity.SupportData
	MakeIncidentDataNote(b []byte) []*entity.IncidentData
}

func NewUseCase(repository Repository) *usecase {
	return &usecase{
		repository: repository,
	}
}

type usecase struct {
	repository Repository
}

func (u *usecase) MakeResultSetT(sms [][]*entity.SMSData, mms [][]*entity.MMSData, v []*entity.VoiceData, m map[string][][]*entity.EmailData, b *entity.BillingData, s []int, i []*entity.IncidentData) (*entity.ResultSetT, bool) {
	result := &entity.ResultSetT{sms, mms, v, m, b, s, i}
	bo := true
	if sms == nil || mms == nil || v == nil || m == nil || b == nil || s == nil || i == nil {
		bo = false
	}
	return result, bo
}

func (u *usecase) MakeIncidentDataNote(b []byte) []*entity.IncidentData {
	data := u.repository.MakeIncidentDataNote(b)
	result := sortIncidentSlice(data)
	return result
}

func (u *usecase) MakeSupportDataNote(b []byte) []int {
	result := make([]int, 0)
	data := u.repository.MakeSupportDataNote(b)
	load := 0
	averageMinPerTikcket := 60 / 18
	for _, v := range data {
		load += v.ActiveTickets
	}
	if load < 9 {
		result = append(result, 1)
	} else if load >= 9 && load <= 16 {
		result = append(result, 2)
	} else {
		result = append(result, 3)
	}
	wT := float64(load * averageMinPerTikcket)
	result = append(result, int(wT))
	return result
}

func (u *usecase) MakeBillingNote(b []byte) *entity.BillingData {
	boolSlice := make([]bool, 0)
	slice := make([]string, 0)
	var bResult float64
	for _, v := range b {
		slice = append(slice, string(v))
	}
	for i, _ := range slice {
		if slice[len(slice)-i-1] == "1" {
			bResult += math.Pow(2, float64(i))
			boolSlice = append(boolSlice, true)
		} else {
			boolSlice = append(boolSlice, false)
		}
	}
	data := u.repository.MakeBillingNote(boolSlice)
	return data
}

func (u *usecase) MakeEmailNote(s []string) map[string][][]*entity.EmailData {
	res := make([]*entity.EmailData, 0)
	for _, v := range s {
		r := strings.Split(v, ";")
		if len(r) == 3 && filter(r[1], "Email") {
			dT, err := strconv.Atoi(r[2])
			if err != nil {
				log.Fatal(err)
			}
			for code, _ := range countryAlpha2.CA2() {
				if r[0] == code {
					for _, v := range u.repository.MakeEmailNote(r, dT) {
						res = append(res, v)
					}
				}
			}
		}
	}
	result := sortEmailSlice(res)
	return result
}

func (u *usecase) MakeVoiceNote(s []string) []*entity.VoiceData {
	res := make([]*entity.VoiceData, 0)
	for _, v := range s {
		r := strings.Split(v, ";")
		if len(r) == 8 && filter(r[3], "Voice") {
			f, t, v, m := parseData(r)
			for code, _ := range countryAlpha2.CA2() {
				if r[0] == code {
					for _, v := range u.repository.MakeVoiceNote(r, f, t, v, m) {
						res = append(res, v)
					}
				}
			}
		}
	}
	return res
}

func (u *usecase) MakeMMSNote(b []byte) [][]*entity.MMSData {
	res := make([]*entity.MMSData, 0)
	data := u.repository.MakeMMSNote(b)
	for _, v := range data {
		if filter(v.Provider, "SMS") {
			for code, c := range countryAlpha2.CA2() {
				if v.Country == code {
					v.Country = c
					res = append(res, v)
				}
			}
		}
	}
	s1, s2 := sortMMSSlice(res)
	result := [][]*entity.MMSData{s1, s2}
	return result
}

func (u *usecase) MakeSMSNote(sliceString []string) [][]*entity.SMSData {
	res := make([]*entity.SMSData, 0)
	for _, v := range sliceString {
		r := strings.Split(v, ";")
		if len(r) == 4 && filter(r[3], "SMS") {
			for code, v := range countryAlpha2.CA2() {
				if r[0] == code {
					r[0] = v
					for _, v := range u.repository.MakeSMSNote(r) {
						res = append(res, v)
					}
				}
			}
		}
	}
	s1, s2 := sortSMSSlice(res)
	result := [][]*entity.SMSData{s1, s2}
	return result
}

func sortIncidentSlice(slice []*entity.IncidentData) []*entity.IncidentData {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Status < slice[j].Status
	})
	return slice
}

func sortEmailSlice(slice []*entity.EmailData) map[string][][]*entity.EmailData {
	m := make(map[string][][]*entity.EmailData, 0)
	ch := ""
	sort.SliceStable(slice, func(i, j int) bool {
		return slice[i].DeliveryTime < slice[j].DeliveryTime
	})
	sort.SliceStable(slice, func(i, j int) bool {
		return slice[i].Country < slice[j].Country
	})
	for i, v := range slice {
		if v.Country != ch {
			emailFastest := make([]*entity.EmailData, 0)
			emailFastest = append(emailFastest, slice[i:i+3]...)
			prom := make([][]*entity.EmailData, 0)
			prom = append(prom, emailFastest)
			ch = v.Country
			m[ch] = append(m[ch], emailFastest)
		}
	}
	sort.SliceStable(slice, func(i, j int) bool {
		return slice[i].DeliveryTime > slice[j].DeliveryTime
	})
	sort.SliceStable(slice, func(i, j int) bool {
		return slice[i].Country < slice[j].Country
	})
	for i, v := range slice {
		if v.Country != ch {
			emailSlowest := make([]*entity.EmailData, 0)
			emailSlowest = append(emailSlowest, slice[i:i+3]...)
			prom := make([][]*entity.EmailData, 0)
			prom = append(prom, emailSlowest)
			ch = v.Country
			m[ch] = append(m[ch], emailSlowest)
		}
	}
	return m
}

func sortMMSSlice(slice []*entity.MMSData) ([]*entity.MMSData, []*entity.MMSData) {
	sort.Slice(slice, func(i, j int) bool {
		a, b := slice[i].Country, slice[j].Country
		for utf8.RuneCountInString(a) < utf8.RuneCountInString(b) {
			a += " "
		}
		for utf8.RuneCountInString(b) < utf8.RuneCountInString(a) {
			b += " "
		}
		return a < b
	})
	mmsByCountry := make([]*entity.MMSData, len(slice))
	copy(mmsByCountry, slice)
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Provider < slice[j].Provider
	})
	mmsByProvider := make([]*entity.MMSData, len(slice))
	copy(mmsByProvider, slice)
	return mmsByProvider, mmsByCountry
}

func sortSMSSlice(slice []*entity.SMSData) ([]*entity.SMSData, []*entity.SMSData) {
	sort.Slice(slice, func(i, j int) bool {
		a, b := slice[i].Country, slice[j].Country
		for utf8.RuneCountInString(a) < utf8.RuneCountInString(b) {
			a += " "
		}
		for utf8.RuneCountInString(b) < utf8.RuneCountInString(a) {
			b += " "
		}
		return a < b
	})
	smsByCountry := make([]*entity.SMSData, len(slice))
	copy(smsByCountry, slice)
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Provider < slice[j].Provider
	})
	smsByProvider := make([]*entity.SMSData, len(slice))
	copy(smsByProvider, slice)
	return smsByProvider, smsByCountry
}

func parseData(s []string) (float32, int, int, int) {
	fmt.Println()
	f, err := strconv.ParseFloat(s[4], 32)
	t, err := strconv.Atoi(s[5])
	v, err := strconv.Atoi(s[6])
	m, err := strconv.Atoi(s[7])
	if err != nil {
		log.Fatal(err)
	}
	return float32(f), t, v, m
}

func filter(provider string, code string) bool {
	validSMSProviders := []string{"Rond", "Kildy", "Topolo"}
	validVoiceProviders := []string{"TransparentCalls", "E-Voice", "JustPhone"}
	validEmailProviders := []string{"Gmail", "Yahoo", "Hotmail", "MSN", "Orange", "Comcast", "AOL", "Live", "RediffMail", "GMX", "Protonmail", "Yandex", "Mail.ru"}
	switch code {
	case "SMS":
		for _, v := range validSMSProviders {
			if provider == v {
				return true
			}
		}
	case "Voice":
		for _, v := range validVoiceProviders {
			if provider == v {
				return true
			}
		}
	case "Email":
		for _, v := range validEmailProviders {
			if provider == v {
				return true
			}
		}
	}
	return false
}
