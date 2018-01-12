package front_igdsanglah

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func testDatabase(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	tgl := r.URL.Path[9:]
	log.Infof(ctx, "Path adalah: %v", tgl)
	q := datastore.NewQuery("KunjunganPasien").Filter("Dokter =", "suryasedana@gmail.com").Order("-JamDatang").Limit(100)
	t := q.Run(ctx)
	list := []KunjunganPasien{}
	for {
		j := &KunjunganPasien{}
		_, err := t.Next(j)
		if err == datastore.Done {
			break
		}
		if err != nil {
			DocumentError(w, ctx, "Gagal mengambil database", err, 500)
			return
		}
		list = append(list, *j)
	}

	fmt.Fprint(w, fmt.Sprint(list))

}
