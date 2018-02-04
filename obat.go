package front_igdsanglah

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/search"
)

func getObatPage(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	w.WriteHeader(200)
	fmt.Fprint(w, GenTemplate(w, ctx, nil, "obat-page"))
	// SendBackSuccess(w, nil, GenTemplate(w, ctx, nil, "obat-page"), "", "")
}
func cariObat(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	// log.Infof(ctx, "Key word adalah: %v", js.Data1)
	index, err := search.Open("DataObat")
	if err != nil {
		DocumentError(w, ctx, "membuat index pencarian", err, 500)
		return
	}
	t := index.Search(ctx, js.Data1, nil)
	list := []IndexObat{}
	for {
		var us IndexObat
		_, err := t.Next(&us)
		if err == search.Done {
			break
		}
		if err != nil {
			DocumentError(w, ctx, "mencari obat", err, 500)
			break
		}
		list = append(list, us)
	}
	// log.Infof(ctx, "List obat adaalh: %v", GenTemplate(w, ctx, list, "list-obat"))
	SendBackSuccess(w, nil, GenTemplate(w, ctx, list, "list-obat"), "", "")
}

func tambahObat(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	// max, _ := strconv.ParseFloat(js.Data6, 64)
	// min, _ := strconv.Parse
	obat := &Obat{
		MerkDagang:    js.Data1,
		Kandungan:     js.Data2,
		Keterangan:    js.Data8,
		SediaanObat:   js.Data3,
		MinDose:       ConvertStringtoFloat(js.Data6),
		MaxDose:       ConvertStringtoFloat(js.Data7),
		Takaran:       js.Data4,
		JmlPerTakaran: js.Data5,
		Submitter:     js.Data9,
		TglEdit:       timeNowIndonesia().In(ZonaIndo()),
	}
	ind := &IndexObat{
		MerkDagang: obat.MerkDagang,
		Kandungan:  obat.Kandungan,
	}

	k := datastore.NewKey(ctx, "Obat", "", 0, datastore.NewKey(ctx, "IGD", "fasttrack", 0, nil))
	ke, err := datastore.Put(ctx, k, obat)
	if err != nil {
		DocumentError(w, ctx, "menyimpan obat", err, 500)
		return
	}
	ind.Link = ke.Encode()
	index, err := search.Open("DataObat")
	if err != nil {
		DocumentError(w, ctx, "membuat index", err, 500)
		return
	}

	_, err = index.Put(ctx, ke.Encode(), ind)
	if err != nil {
		DocumentError(w, ctx, "menyimpan index", err, 500)
		return
	}
	SendBackSuccess(w, nil, "", "Berhasil menyimpan obat", "")
}

func ConvertStringtoFloat(s string) float64 {
	var m float64
	if s == "" {
		m = 0
		return m
	} else {
		n, _ := strconv.ParseFloat(s, 64)
		m = n
		return m
	}
}

func getIsianObat(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	w.WriteHeader(200)

	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	fmt.Fprint(w, GenTemplate(w, ctx, js, "get-info-obat"))
	// k, err := datastore.DecodeKey(js.Data1)
	// if err != nil {
	// 	DocumentError(w, ctx, "mendecode key", err, 500)
	// 	return
	// }
	// obt := &Obat{}
	// err = datastore.Get(ctx, k, obt)
	// if err != nil {
	// 	DocumentError(w, ctx, "mengambil data", err, 500)
	// 	return
	// }

}

func getDataObat(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	k, err := datastore.DecodeKey(js.Data1)
	if err != nil {
		DocumentError(w, ctx, "mendecode key", err, 500)
		return
	}
	obt := &Obat{}
	err = datastore.Get(ctx, k, obt)
	if err != nil {
		DocumentError(w, ctx, "mengambil data", err, 500)
		return
	}
	obt.LinkID = k.Encode()
	fin := ObatFinal{}
	fin.Obat = *obt
	if js.Data2 == "2" && obt.MaxDose != 0 {
		fin.MinDoseFinal = ConvertStringtoFloat(js.Data3) * obt.MinDose
		fin.MaxDoseFinal = ConvertStringtoFloat(js.Data3) * obt.MaxDose
		fin.Dewasa = false
	} else if js.Data2 == "2" && obt.MaxDose == 0 {
		fin.MinDoseFinal = ConvertStringtoFloat(js.Data3) * obt.MinDose
		fin.Dewasa = false
	} else {
		fin.Dewasa = true
	}
	w.WriteHeader(200)
	fmt.Fprint(w, GenTemplate(w, ctx, fin, "view-detail-obat"))
}
