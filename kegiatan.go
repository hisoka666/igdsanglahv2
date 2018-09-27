package front_igdsanglah

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// getKegiatanDokter digunakan untuk mengambil database
// untuk kegiatan user.
func getKegiatanDokter(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	tgl := timeNowIndonesia().Format("2006/01")
	// CatchDataJson adalah struct untuk 'menangkap' bermacam jenis data
	// yang dikirim bersama metode POST dari front end
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()

	keg, err := getKegiatanBulanan(ctx, js.Data1, tgl)
	log.Infof(ctx, "Data list adalah: %v", keg)
	if err != nil {
		DocumentError(w, ctx, "Gagal mengambil daftar kegiatan", err, 500)
		return
	}
	SendBackSuccess(w, nil, GenTemplate(w, ctx, keg, "kegiatan-dokter"), "", "")
}

// getKegiatanBulanan digunakan untuk mengambil kegiatan di bulan tgl dengan
// format "2006/01", dengan email user. Fungsi ini menghasilkan slice KegiatanDokter
// dan error jika ada
func getKegiatanBulanan(ctx context.Context, email string, tgl string) ([]KegiatanDokter, error) {
	bln, err := time.Parse("2006/01/02", tgl+"/01")
	if err != nil {
		return nil, err
	}
	par, err := getDocKey(ctx, email)
	if err != nil {
		return nil, err
	}
	q := datastore.NewQuery("KegiatanDokter").Ancestor(par).Filter("TglTindakan >=", bln.In(ZonaIndo())).Filter("TglTindakan <", bln.In(ZonaIndo()).AddDate(0, 1, 0)).Filter("Hide =", false).Order("-TglTindakan")
	t := q.Run(ctx)
	keg := []KegiatanDokter{}
	for {
		var k KegiatanDokter
		kun, err := t.Next(&k)
		if err == datastore.Done {
			break
		}
		// if k.TglTindakan.After(bln.AddDate(0, 1, 0)) == true {
		// 	continue
		// }
		// if k.TglTindakan.Before(bln) == true {
		// 	break
		// }
		if err != nil {
			return nil, err
		}

		k.KeyDataTindakan = kun.Encode()
		keg = append(keg, k)
	}
	return keg, nil
}

// addKegiatanDokter digunakan untuk menambahkan data kegiatan
// ke struct KegiatanDokter
func addKegiatanDokter(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	// log.Infof(ctx, "email adalah %s", js.Data3)
	k, err := getDocKey(ctx, js.Data3)
	if err != nil {
		DocumentError(w, ctx, "mengambil key dokter", err, 500)
		return
	}
	keg := &KegiatanDokter{
		IDPasien:     js.Data1,
		NamaTindakan: js.Data2,
		NamaPasien:   js.Data4,
		TglTindakan:  timeNowIndonesia(),
		Hide:         false,
	}
	// log.Infof(ctx, "Nama Tindakan adalah %s", keg.NamaTindakan)
	_, err = datastore.Put(ctx, datastore.NewKey(ctx, "KegiatanDokter", "", 0, k), keg)
	if err != nil {
		DocumentError(w, ctx, "Gagal menyimpan data kegiatan", err, 500)
	} else {
		kd, err := getKegiatanBulanan(ctx, js.Data3, timeNowIndonesia().Format("2006/01"))
		if err != nil {
			DocumentError(w, ctx, "Gagal mengambil data kegiatan", err, 500)
			return
		}
		SendBackSuccess(w, nil, GenTemplate(w, ctx, kd, "kegiatan-dokter"), "Berhasil menyimpan kegiatan", "")
	}
}

func hapusKegiatanDokter(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	// log.Infof(ctx, "key adalah: %s", js.Data1)
	kun, keg, err := getKegiatanByKey(ctx, js.Data1)
	if err != nil {
		DocumentError(w, ctx, "Gagal mengambil data", err, 500)
		return
	}
	keg.Hide = true
	_, err = datastore.Put(ctx, kun, &keg)
	if err != nil {
		DocumentError(w, ctx, "Gagal menghapus data", err, 500)
		return
	}
	kd, err := getKegiatanBulanan(ctx, js.Data3, timeNowIndonesia().Format("2006/01"))
	if err != nil {
		DocumentError(w, ctx, "Gagal mengambil data kegiatan", err, 500)
		return
	}
	SendBackSuccess(w, nil, GenTemplate(w, ctx, kd, "kegiatan-dokter"), "Berhasil menghapus kegiatan", "")
	// log.Infof(ctx, "Data adalah: %v", keg)

}

func getKegiatanByKey(ctx context.Context, kun string) (*datastore.Key, KegiatanDokter, error) {
	k, err := datastore.DecodeKey(kun)
	var ke KegiatanDokter
	if err != nil {
		return nil, ke, err
	}
	err = datastore.Get(ctx, k, &ke)
	if err != nil {
		return nil, ke, err
	}
	ke.KeyDataTindakan = k.Encode()
	return k, ke, nil
}

func getContentKegiatanBulanan(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	tgl := js.Data1[0:4] + "/" + js.Data1[5:7]
	kd, err := getKegiatanBulanan(ctx, js.Data3, tgl)
	if err != nil {
		DocumentError(w, ctx, "Gagal mengambil data kegiatan", err, 500)
		return
	}
	// log.Infof(ctx, "List kegiatan adalah: %v", kd)
	if len(kd) == 0 {
		ke, err := getKegiatanBulanan(ctx, js.Data3, timeNowIndonesia().Format("2006/01"))
		if err != nil {
			DocumentError(w, ctx, "Gagal mengambil data kegiatan", err, 500)
			return
		}
		SendBackSuccess(w, nil, GenTemplate(w, ctx, ke, "kegiatan-dokter"), "Tidak ada data kegiatan!", "")
	} else {
		SendBackSuccess(w, nil, GenTemplate(w, ctx, kd, "kegiatan-dokter"), "Berhasil menyimpan kegiatan", "")
	}
}
