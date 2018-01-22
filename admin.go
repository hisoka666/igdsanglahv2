package front_igdsanglah

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func getAdminPage(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	q := datastore.NewQuery("Staff")
	st := []Staff{}
	k, err := q.GetAll(ctx, &st)
	if err != nil && err != datastore.ErrNoSuchEntity {
		DocumentError(w, ctx, "mengambil data staf", err, 500)
		return
	}
	staf := []Staff{}
	for m, n := range st {
		n.LinkID = k[m].Encode()
		staf = append(staf, n)
	}
	fmt.Fprintf(w, GenTemplate(w, ctx, staf, "admin"))
}

func addStaf(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	staf := &Staff{
		NamaLengkap: js.Data2,
		Email:       js.Data1,
		Peran:       js.Data3,
	}
	k := datastore.NewKey(ctx, "Staff", "", 0, datastore.NewKey(ctx, "IGD", "fasttrack", 0, nil))
	_, err := datastore.Put(ctx, k, staf)
	if err != nil {
		DocumentError(w, ctx, "menyimpan data", err, 503)
		return
	}
	SendBackSuccess(w, nil, "", "Berhasil menambahkan data", "")
}

func hapusDataStaf(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	js := &CatchDataJson{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	k, err := datastore.DecodeKey(js.Data1)
	if err != nil {
		DocumentError(w, ctx, "mendecode key", err, 500)
		return
	}
	err = datastore.Delete(ctx, k)
	if err != nil {
		DocumentError(w, ctx, "menghapus data", err, 500)
		return
	}
	SendBackSuccess(w, nil, "", "Berhasil menghapus data staf", "")
}
