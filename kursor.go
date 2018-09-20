package front_igdsanglah

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func createKursor(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var staf []Staff
	q := datastore.NewQuery("Staff")
	_, err := q.GetAll(ctx, &staf)
	if err != nil {
		DocumentError(w, ctx, "mengambil staf", err, 500)
		return
	}

	for _, v := range staf {
		par := datastore.NewKey(ctx, "Dokter", v.Email, 0, datastore.NewKey(ctx, "IGD", "fasttrack", 0, nil))
		CreateEndKursor(ctx, par, v.Email)
	}
	CreateKursorIGD(ctx)
}

func CreateKursorIGD(c context.Context) {
	q := datastore.NewQuery("KunjunganPasien").Filter("Hide=", false).Order("-JamDatang")
	t := q.Run(c)
	kur := KursorIGD{}
	kun := KunjunganPasien{}
	hariini := time.Date(timeNowIndonesia().Year(), timeNowIndonesia().Month(), 1, 8, 0, 0, 0, ZonaIndo())
	tgl := hariini.AddDate(0, -1, 0).Format("2006/01")
	k := datastore.NewKey(c, "KursorIGD", tgl, 0, datastore.NewKey(c, "IGD", "fasttrack", 0, nil))
	for {
		_, err := t.Next(&kun)
		if err == datastore.Done {
			break
		}
		if err != nil {
			log.Infof(c, "Kesalahan membaca database: %v", err)
			break
		}

		if kun.JamDatangRiil.Before(hariini) == true {
			cursor, _ := t.Cursor()
			kur.Point = cursor.String()
			kur.Bulan = tgl
			k, err := datastore.Put(c, k, &kur)
			if err != nil {
				log.Errorf(c, "Kesalahan menulis database: %v", err)
				break
			}
			log.Infof(c, "Berhasil menambahkan kursor %v", k)
			break
		}
	}
}
func getKursor(c context.Context, email string) ([]Kursor, error) {
	kur := []Kursor{}
	q := datastore.NewQuery("Kursor").Ancestor(datastore.NewKey(c, "Dokter", email, 0, datastore.NewKey(c, "IGD", "fasttrack", 0, nil)))
	n, err := q.GetAll(c, &kur)
	if err != nil {
		return nil, err
	}
	for k, v := range n {
		kur[k].Link = v.StringID()
	}
	if timeNowIndonesia().Day() == 1 && timeNowIndonesia().Hour() < 8 {
		ku := Kursor{
			Link: timeNowIndonesia().AddDate(0, -1, 0).Format("2006/01"),
		}
		kur = append(kur, ku)
	}
	for i, j := 0, len(kur)-1; i < j; i, j = i+1, j-1 {
		kur[i], kur[j] = kur[j], kur[i]
	}
	return kur, nil
}

func CreateEndKursor(c context.Context, par *datastore.Key, email string) {

	q := datastore.NewQuery("KunjunganPasien").Filter("Dokter=", email).Filter("Hide=", false).Order("-JamDatang")
	t := q.Run(c)
	kur := Kursor{}
	kun := KunjunganPasien{}
	tgl := timeNowIndonesia().AddDate(0, -1, 0).Format("2006/01")
	tglnow := timeNowIndonesia().Format("2006/01")
	tglend, err := time.ParseInLocation("2006/01/02 15:04", tglnow+"/01 07:30", ZonaIndo())
	if err != nil {
		log.Errorf(c, "Gagal memparse tglend : %v", err)
		return
	}
	k := datastore.NewKey(c, "Kursor", tgl, 0, par)
	for {
		_, err := t.Next(&kun)
		if err == datastore.Done {
			break
		}
		if err != nil {
			log.Errorf(c, "Gagal membaca data %v", err)
		}
		if IsThisCursor(kun.JamDatang.In(ZonaIndo()), tglend, kun.ShiftJaga) {
			cursor, _ := t.Cursor()
			kur.Point = cursor.String()
			if _, err := datastore.Put(c, k, &kur); err != nil {
				log.Errorf(c, "gagal menyimpan kursor : %v", err)
			}
			log.Infof(c, "Created Kursor for %v", email)
			break
		}
	}
}

func IsThisCursor(tglkun, tglawal time.Time, shift string) bool {
	if tglkun.Before(tglawal) && shift != "3" {
		return true
	} else if shift == "3" && tglkun.Hour() < 12 {
		return true
	} else {
		return false
	}
}
