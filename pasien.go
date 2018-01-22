package front_igdsanglah

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/search"
)

func getInfoNoCM(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	pts, k, err := getDataByNoCM(ctx, js.Data1)
	if err != nil {
		DocumentError(w, ctx, "Gagal mengambil data pasien", err, 500)
		return
	} else if k == nil {
		SendBackSuccess(w, nil, GenTemplate(w, ctx, pts, "input-pts"), "", "")
	} else {
		SendBackSuccess(w, nil, GenTemplate(w, ctx, pts, "input-pts"), k.Encode(), "")
	}
}
func getDataPasienByLink(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	// log.Infof(ctx, "Link adalah: %v", js)
	dat, err := getPasienIDDanKunjunganFromLink(ctx, js.Data1)
	if err != nil {
		DocumentError(w, ctx, "gagal mengambil data pasien", err, 500)
		return
	}
	SendBackSuccess(w, dat.Kunjungan, GenTemplate(w, ctx, dat, "modal-edit"), js.Data1, "")
}
func getPasienIDDanKunjunganFromLink(c context.Context, link string) (*IDDanKunjungan, error) {
	k, err := datastore.DecodeKey(link)
	if err != nil {
		return nil, err
	}
	kun := &KunjunganPasien{}
	err = datastore.Get(c, k, kun)
	if err != nil {
		return nil, err
	}
	kun.LinkID = k.Encode()
	par := k.Parent()
	dat := &DataPasien{}
	err = datastore.Get(c, par, dat)
	if err != nil {
		return nil, err
	}
	dat.NomorCM = par.StringID()
	m := &IDDanKunjungan{
		ID:        *dat,
		Kunjungan: *kun,
	}
	return m, nil
}
func getDataByNoCM(c context.Context, nocm string) (*DataPasien, *datastore.Key, error) {
	pts := &DataPasien{}
	k := datastore.NewKey(c, "DataPasien", nocm, 0, datastore.NewKey(c, "IGD", "fasttrack", 0, nil))
	err := datastore.Get(c, k, pts)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return nil, nil, err
	} else if err == datastore.ErrNoSuchEntity {
		return pts, nil, nil
	} else {
		return pts, k, nil
	}
	// checkError(w, c, "mengambil data pasien", err, 500)
}

func tambahDataKunjungan(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	// var k *datastore.Key
	if js.Data1 == "" {
		k := datastore.NewKey(ctx, "DataPasien", js.Data2, 0, datastore.NewKey(ctx, "IGD", "fasttrack", 0, nil))
		ind := &IndexDataPasien{
			Nama: js.Data3,
			NoCM: js.Data2,
		}
		index, err := search.Open("DataPasienTest02")
		if err != nil {
			DocumentError(w, ctx, "Gagal membuat index", err, 500)
			return
		}
		_, err = index.Put(ctx, k.Encode(), ind)
		if err != nil {
			DocumentError(w, ctx, "gagal menyimpan index", err, 500)
			return
		}
		dat := &DataPasien{
			NamaPasien: js.Data3,
			TglLahir:   ChangeStringtoTime(js.Data4),
			TglDaftar:  timeNowIndonesia(),
		}
		_, err = datastore.Put(ctx, k, dat)
		if err != nil {
			DocumentError(w, ctx, "Gagal menyimpan data pasien", err, 500)
		}

		kun := &KunjunganPasien{
			Diagnosis:     js.Data5,
			GolIKI:        js.Data8,
			ATS:           js.Data6,
			Bagian:        js.Data7,
			JamDatang:     timeNowIndonesia(),
			JamDatangRiil: timeNowIndonesia(),
			Dokter:        js.Data10,
			Hide:          false,
			ShiftJaga:     js.Data9,
		}

		_, err = datastore.Put(ctx, datastore.NewKey(ctx, "KunjunganPasien", "", 0, k), kun)
		if err != nil {
			DocumentError(w, ctx, "Gagal menyimpan data kunjungan", err, 500)
		} else {
			SendBackSuccess(w, nil, "", "Berhasil menyimpan data kunjungan", "")
		}
	} else {
		kun := &KunjunganPasien{
			Diagnosis:     js.Data5,
			GolIKI:        js.Data8,
			ATS:           js.Data6,
			Bagian:        js.Data7,
			JamDatang:     timeNowIndonesia(),
			JamDatangRiil: timeNowIndonesia(),
			Dokter:        js.Data10,
			Hide:          false,
			ShiftJaga:     js.Data9,
		}
		k, _ := datastore.DecodeKey(js.Data1)
		_, err := datastore.Put(ctx, datastore.NewKey(ctx, "KunjunganPasien", "", 0, k), kun)
		if err != nil {
			DocumentError(w, ctx, "Gagal menyimpan data kunjungan", err, 500)
		} else {
			SendBackSuccess(w, nil, "", "Berhasil menyimpan data kunjungan", "")
		}
	}
}

func getDetailPasien(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	k, err := datastore.DecodeKey(js.Data1)
	if err != nil {
		DocumentError(w, ctx, "mendecode key", err, 500)
		return
	}
	q := datastore.NewQuery("KunjunganPasien").Ancestor(k.Parent()).Order("-JamDatang").Filter("Hide=", false)
	j := []KunjunganPasien{}
	m, err := q.GetAll(ctx, &j)
	if err != nil {
		DocumentError(w, ctx, "mengambil list kunjungan", err, 500)
		return
	}
	list := []KunjunganPasien{}
	for a, b := range j {
		if b.Dokter == "sunia.raharja@gmail.com" {
			continue
		}
		b.LinkID = m[a].Encode()
		list = append(list, b)
	}
	dat := &DataPasien{}
	err = datastore.Get(ctx, k.Parent(), dat)
	if err != nil {
		DocumentError(w, ctx, "mengambil data pasien", err, 500)
		return
	}
	pts := DetailPasienPage{
		Pasien:    *dat,
		Kunjungan: list,
		LinkID:    k.Parent().Encode(),
	}
	log.Infof(ctx, "List kunjungan adalah : %v", pts.Kunjungan)
	SendBackSuccess(w, nil, GenTemplate(w, ctx, pts, "detail-pasien-page"), "", "")
}
