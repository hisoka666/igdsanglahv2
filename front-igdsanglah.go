package front_igdsanglah

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/leekchan/accounting"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func init() {
	http.HandleFunc("/", indexPage)
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	if u == nil {
		url, _ := user.LoginURL(ctx, "/")
		fmt.Fprintf(w, "<a href='%s'>Sign in or Register</a>", url)
		return
	}
	logout, _ := user.LogoutURL(ctx, "/")
	if u.Email != "suryasedana@gmail.com" {
		http.Redirect(w, r, logout, 403)
	} else {
		me := FrontPage{
			LogOut:   logout,
			UserName: "I Wayan Surya Sedana",
		}
		front := GenTemplate(w, ctx, me, "index")
		fmt.Fprint(w, front)
	}
}

func GenTemplate(w http.ResponseWriter, c context.Context, n interface{}, temp ...string) string {
	b := new(bytes.Buffer)
	funcs := template.FuncMap{
		"jam": func(t time.Time) string {
			zone, _ := time.LoadLocation("Asia/Makassar")
			return t.In(zone).Format("15:04")
		},
		"inc": func(i int) int {
			return i + 1
		},
		"rp": func(i int) string {
			ac := accounting.Accounting{
				Symbol:    "Rp ",
				Precision: 2,
				Thousand:  ".",
				Decimal:   ",",
			}
			m := fmt.Sprint(ac.FormatMoney(i))
			return m
		},
		"rpi": func(s string) string {
			rp, _ := strconv.Atoi(s)
			ac := accounting.Accounting{
				Symbol:    "Rp ",
				Precision: 2,
				Thousand:  ".",
				Decimal:   ",",
			}
			m := fmt.Sprint(ac.FormatMoney(rp))
			return m
		},
		"strtgl": func(t time.Time) string {
			return t.Format("Mon, 02/01/2006")
		},
		"istimezero": func(t time.Time) bool {
			return t.IsZero()
		},
		"convstrjaga": func(j string) string {
			var m string
			switch j {
			case "1":
				m = "Pagi"
			case "2":
				m = "Sore"
			case "3":
				m = "Malam"
			}
			return m
		},
	}

	tmpl := template.New("")
	for k, v := range temp {
		if k == 0 {
			tmp := template.Must(template.New(v + ".html").Funcs(funcs).ParseFiles("templates/" + v + ".html"))
			tmpl = tmp
		}
	}

	for k, v := range temp {
		if k != 0 {
			temp, err := template.Must(tmpl.Clone()).ParseFiles("templates/" + v + ".html")
			if err != nil {
				DocumentError(w, c, "parse template multiple", err, 500)
				return ""
			}
			tmpl = temp
		}
	}
	err := tmpl.Execute(b, n)
	if err != nil {
		DocumentError(w, c, "eksekusi template", err, 500)
		return ""
	}

	return b.String()
}

func DocumentError(w http.ResponseWriter, c context.Context, topik string, err error, kode int) {
	msg := "Telah terjadi kesalahan dalam " + topik + " : %v"
	log.Errorf(c, msg, err)
	w.WriteHeader(kode)
	fmt.Fprint(w, "Gagal "+topik)
}
