package front_igdsanglah

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func hapusDataKunjungan(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	bcp := &TableBCP{}
	err := json.Unmarshal([]byte(js.Data2), bcp)
	if err != nil {
		DocumentError(w, ctx, "membaca json", err, 500)
		return
	}
	k, err := datastore.DecodeKey(js.Data1)
	if err != nil {
		DocumentError(w, ctx, "mendecode key dari string", err, 500)
		return
	}
	kun := &KunjunganPasien{}
	err = datastore.Get(ctx, k, kun)
	if err != nil {
		DocumentError(w, ctx, "mengambil data", err, 500)
		return
	}
	kun.Hide = true
	_, err = datastore.Put(ctx, k, kun)
	if err != nil {
		DocumentError(w, ctx, "menghapus data", err, 500)
		return
	}
	tab, err := getTableBCPbyCursor(ctx, bcp.Kursor, bcp.Email, bcp.StringTgl)
	if err != nil {
		DocumentError(w, ctx, "mengambil list bcp", err, 500)
		return
	}
	jss, _ := json.Marshal(tab)
	SendBackSuccess(w, nil, GenTemplate(w, ctx, tab, "bcp-content"), string(jss), bcp.Kursor)
}
func getBCPBulan(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	data, err := getTableBCPbyCursor(ctx, js.Data1, js.Data2, js.Data3)
	if err != nil {
		DocumentError(w, ctx, "mengambil data table bcp", err, 500)
		return
	}

	jss, _ := json.Marshal(data)
	SendBackSuccess(w, nil, GenTemplate(w, ctx, data, "bcp-content"), string(jss), js.Data1)
}
func getBCPBulanIni(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	// email := js.Data1
	tab, err := getTableBCPbyCursor(ctx, "", js.Data1, "")
	if err != nil {
		DocumentError(w, ctx, "mengambil data tabel bcp", err, 500)
		return
	}

	jss, _ := json.Marshal(tab)
	SendBackSuccess(w, nil, GenTemplate(w, ctx, tab, "bcp-content"), string(jss), "")
}

func getListPasien(c context.Context, kur, email, tgl string) ([]Pasien, error) {
	q := datastore.NewQuery("KunjunganPasien").Filter("Dokter =", email).Order("-JamDatang")
	list := []KunjunganPasien{}
	if kur != "" {
		tglkur, err := time.ParseInLocation("2006/01/02 15:04", tgl+"/01 08:00", ZonaIndo())
		if err != nil {
			return nil, err
		}
		sor, err := datastore.DecodeCursor(kur)
		if err != nil {
			return nil, err
		}
		q = q.Start(sor)
		m, err := iterateList(c, q, tglkur)
		if err != nil {
			return nil, err
		}
		list = append(list, m...)
	} else if tgl != "" && kur == "" {
		now := time.Date(timeNowIndonesia().Year(), timeNowIndonesia().Month(), 1, 8, 0, 0, 0, ZonaIndo())
		m, err := iterateList(c, q, now)
		if err != nil {
			return nil, err
		}
		list = append(list, m...)
	} else {
		now := time.Date(timeNowIndonesia().Year(), timeNowIndonesia().Month(), 1, 8, 0, 0, 0, ZonaIndo())
		m, err := iterateList(c, q, now)
		if err != nil {
			return nil, err
		}
		list = append(list, m...)
	}
	pts, err := convertToListPasien(c, list)
	if err != nil {
		return nil, err
	}
	return pts, nil

}
func CountIKI(g context.Context, n []Pasien, tgl time.Time) (string, []IKI) {
	wkt := time.Date(tgl.Year(), tgl.Month(), 1, 8, 0, 0, 0, ZonaIndo())
	jmlhari := wkt.AddDate(0, 1, -1).Day()
	list := []IKI{}
	for i := 0; i < jmlhari; i++ {
		timeToCount := wkt.AddDate(0, 0, +i)
		timeAfter := timeToCount.AddDate(0, 0, 1)
		// timeAfter := time.Date(timeToCount.Year(), timeToCount.Month(), (timeToCount.Day() + 1), 12, 0, 0, 0, ZonaIndo())
		var a int
		b := &a
		var c int
		d := &c
		for _, v := range n {
			if v.TglKunjungan.After(timeAfter) {
				continue
			}
			if v.TglKunjungan.Before(timeToCount) {
				break
			}
			if v.NoCM == "00000000" || v.NoCM == "00000001" || v.NoCM == "00000002" {
				continue
			}
			if v.IKI == "1" {
				*b = *b + 1
			} else {
				*d = *d + 1
			}
		}
		e := IKI{
			Tanggal: timeToCount.Format("02/01/2006"),
			IKI1:    a,
			IKI2:    c,
		}
		list = append(list, e)
	}
	var iki1 int
	var iki2 int
	iki01 := &iki1
	iki02 := &iki2
	for _, v := range list {
		*iki02 = *iki02 + v.IKI2
		*iki01 = *iki01 + v.IKI1
	}
	total := float32(iki1)*0.0032 + float32(iki2)*0.01
	return fmt.Sprintf("%.4f", total), list
}
func iterateList(c context.Context, q *datastore.Query, tgl time.Time) ([]KunjunganPasien, error) {
	m := []KunjunganPasien{}
	t := q.Run(c)
	for {
		d := &KunjunganPasien{}
		k, err := t.Next(d)
		if err != nil && err == datastore.Done {
			break
		}
		if err != nil && err != datastore.Done {
			return nil, err
		}
		if d.Hide == true {
			continue
		}
		jam := d.JamDatang.In(ZonaIndo())
		if jam.Before(tgl) {
			break
		}
		// d.JamDatang = d.JamDatang.In(ZonaIndo())
		// d.JamDatang = jam
		d.LinkID = k.Encode()
		m = append(m, *d)
	}
	return m, nil
}

func convertToListPasien(c context.Context, kun []KunjunganPasien) ([]Pasien, error) {
	list := []Pasien{}
	for _, v := range kun {
		ke, err := datastore.DecodeKey(v.LinkID)
		if err != nil {
			return nil, err
		}
		par := ke.Parent()
		n := &DataPasien{}
		err = datastore.Get(c, par, n)
		if err != nil {
			return nil, err
		}

		m := Pasien{
			NamaPasien:   n.NamaPasien,
			TglKunjungan: v.JamDatang.In(ZonaIndo()),
			ShiftJaga:    v.ShiftJaga,
			ATS:          v.ATS,
			Dept:         v.Bagian,
			NoCM:         par.StringID(),
			Diagnosis:    v.Diagnosis,
			IKI:          v.GolIKI,
			LinkID:       v.LinkID,
			TglAsli:      v.JamDatangRiil,
			TglLahir:     n.TglLahir,
		}
		if v.ShiftJaga == "3" && v.JamDatang.Hour() > 5 && v.JamDatang.Hour() < 12 {
			m.TglKunjungan = time.Date(v.JamDatang.Year(), v.JamDatang.Month(), v.JamDatang.Day(), 6, v.JamDatang.Minute(), v.JamDatang.Second(), v.JamDatang.Nanosecond(), ZonaIndo())
		}
		list = append(list, m)
	}
	return list, nil
}

func editDataKunjungan(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	kun := &KunjunganPasien{}
	json.NewDecoder(r.Body).Decode(kun)
	defer r.Body.Close()
	k, err := datastore.DecodeKey(kun.LinkID)
	if err != nil {
		DocumentError(w, ctx, "mendecode key", err, 500)
		return
	}
	kunlama := &KunjunganPasien{}
	err = datastore.Get(ctx, k, kunlama)
	if err != nil {
		DocumentError(w, ctx, "mengambil data", err, 500)
		return
	}
	kunlama.ATS = kun.ATS
	kunlama.Bagian = kun.Bagian
	kunlama.Diagnosis = kun.Diagnosis
	kunlama.GolIKI = kun.GolIKI
	kunlama.ShiftJaga = kun.ShiftJaga
	kunlama.LinkID = kun.LinkID
	_, err = datastore.Put(ctx, k, kunlama)
	if err != nil {
		DocumentError(w, ctx, "menyimpan data", err, 500)
		return
	}
	tab, err := getTableBCPbyCursor(ctx, kun.Dokter, kunlama.Dokter, TanggalBCP(kunlama.JamDatang, kun.ShiftJaga))
	jss, _ := json.Marshal(tab)
	SendBackSuccess(w, nil, GenTemplate(w, ctx, tab, "bcp-content"), string(jss), kun.Dokter)
}

func TanggalBCP(tgl time.Time, shift string) string {
	if tgl.Hour() < 12 && shift == "3" {
		return tgl.AddDate(0, 0, -1).Format("2006/01")
	} else {
		return tgl.Format("2006/01")
	}
}

func getDataTanggalKunjunganPasien(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	k, err := datastore.DecodeKey(js.Data1)
	if err != nil {
		DocumentError(w, ctx, "mendecode key", err, 500)
		return
	}
	kun := &KunjunganPasien{}
	err = datastore.Get(ctx, k, kun)
	if err != nil {
		DocumentError(w, ctx, "mengambil data", err, 500)
		return
	}
	kun.LinkID = k.Encode()
	// log.Infof(ctx, "html adalah: %v", GenTemplate(w, ctx, kun, "edit-tanggal-kunjungan"))
	SendBackSuccess(w, nil, GenTemplate(w, ctx, kun, "edit-tanggal-kunjungan"), k.Encode(), "")
}

func getTableBCPbyCursor(c context.Context, kur, email, tgl string) (*TableBCP, error) {
	tab := &TableBCP{}
	if kur == "" {
		list, err := getListPasien(c, "", email, "")
		if err != nil {
			return nil, err
		}
		a, b := CountIKI(c, list, timeNowIndonesia())
		data := &TableBCP{
			Pasien:    list,
			Title:     timeNowIndonesia().Format("Jan, 2006"),
			StringTgl: timeNowIndonesia().Format("2006/01"),
			Total:     a,
			Kursor:    "",
			Email:     email,
			IKI:       b,
		}
		tab = data
	} else {
		list, err := getListPasien(c, kur, email, tgl)
		if err != nil {
			return nil, err
		}
		tang, err := time.ParseInLocation("2006/01/02", tgl+"/01", ZonaIndo())
		if err != nil {
			return nil, err
		}
		a, b := CountIKI(c, list, tang)
		data := &TableBCP{
			Pasien:    list,
			Title:     tang.Format("Jan, 2006"),
			StringTgl: tang.Format("2006/01"),
			Total:     a,
			Kursor:    kur,
			Email:     email,
			IKI:       b,
		}
		tab = data
	}
	// log.Infof(c, "IKI adalah: %v", tab.IKI)
	// log.Infof(c, "List adalah: %v", tab.Pasien)
	return tab, nil
}

func ubahDataTanggalKunjungan(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	k, err := datastore.DecodeKey(js.Data1)
	if err != nil {
		DocumentError(w, ctx, "mendecode key", err, 500)
		return
	}

	kun := &KunjunganPasien{}
	err = datastore.Get(ctx, k, kun)
	if err != nil {
		DocumentError(w, ctx, "mengambil data", err, 500)
		return
	}
	tglbaru := ChangeStringtoTime(js.Data2)
	jamubah := kun.JamDatangRiil.In(ZonaIndo())
	kun.JamDatang = time.Date(tglbaru.Year(), tglbaru.Month(), tglbaru.Day(), jamubah.Hour(), jamubah.Minute(), jamubah.Second(), jamubah.Nanosecond(), ZonaIndo())
	_, err = datastore.Put(ctx, k, kun)
	if err != nil {
		DocumentError(w, ctx, "gagal menyimpan data", err, 500)
		return
	}
	tb := &TableBCP{}
	json.Unmarshal([]byte(js.Data3), tb)
	tab, err := getTableBCPbyCursor(ctx, tb.Kursor, tb.Email, tb.StringTgl)
	if err != nil {
		DocumentError(w, ctx, "mengambil list", err, 500)
		return
	}
	jss, _ := json.Marshal(tab)
	SendBackSuccess(w, nil, "", string(jss), "")
}
