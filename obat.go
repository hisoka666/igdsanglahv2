package front_igdsanglah

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "golang.org/x/net/context"
	"google.golang.org/appengine"
	_ "google.golang.org/appengine/datastore"
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
	SendBackSuccess(w, nil, GenTemplate(w, ctx, list, "list-obat"), "", "")
}
