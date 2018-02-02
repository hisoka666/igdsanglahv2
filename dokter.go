package front_igdsanglah

import (
	"encoding/json"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func getDocProfile(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	q := datastore.NewQuery("Staff").Filter("Email=", js.Data1)
	st := []Staff{}
	k, err := q.GetAll(ctx, &st)
	if err != nil {
		DocumentError(w, ctx, "mengambil data dokter", err, 500)
		return
	}
	// log.Infof(ctx, "Key adalah: %v", k)
	// log.Infof(ctx, "Staff adalah: %v", st)
	// if len(st) == 0 {
	// det := []DetailStaf{}
	// 	SendBackSuccess(w, nil, GenTemplate(w, ctx, det, "doc-profile"), "", "")
	// } else {
	qu := datastore.NewQuery("DetailStaf").Ancestor(k[0])
	det := []DetailStaf{}
	ke, err := qu.GetAll(ctx, &det)
	if err != nil {
		DocumentError(w, ctx, "mengambil detail dokter", err, 500)
		return
	}
	// log.Infof(ctx, "length det: %v", len(det))
	if len(det) == 0 {
		// log.Infof(ctx, "keypar adalH: %v", k[0])
		de := DetailStaf{
			Umur: k[0].Encode(),
		}
		det = append(det, de)
	} else {
		det[0].LinkID = ke[0].Encode()
		det[0].Umur = k[0].Encode()
		det[0].TanggalLahir = det[0].TanggalLahir.In(ZonaIndo())
	}
	// log.Infof(ctx, "det adalah: %v", det[0].LinkID)
	SendBackSuccess(w, nil, GenTemplate(w, ctx, det[0], "doc-profile"), "", "")

}

func ubahDetailDokter(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	// log.Infof(ctx, "json adalah : %v", js.Data7)
	var k *datastore.Key
	if js.Data7 == "" {
		m, err := datastore.DecodeKey(js.Data8)
		if err != nil {
			DocumentError(w, ctx, "mendecode keyparent", err, 500)
			return
		}
		j := datastore.NewKey(ctx, "DetailStaf", "", 0, m)
		k = j
	} else {
		m, err := datastore.DecodeKey(js.Data7)
		if err != nil {
			DocumentError(w, ctx, "mendecode key", err, 500)
			return
		}
		log.Infof(ctx, "Key ultimate adalH: %v", m)
		k = m
	}
	// log.Infof(ctx, "Key adalah: %v", k)
	det := &DetailStaf{
		NamaLengkap:  js.Data1,
		Bagian:       js.Data4,
		TanggalLahir: ChangeStringtoTime(js.Data2),
		GolonganPNS:  js.Data3,
	}
	if det.GolonganPNS == "" {
		det.NPP = js.Data5
	} else {
		det.NIP = js.Data5
	}
	_, err := datastore.Put(ctx, k, det)
	if err != nil {
		DocumentError(w, ctx, "menyimpan data dokter", err, 500)
		return
	}
	SendBackSuccess(w, nil, "", "Berhasil menyimpan data dokter", "")
}
