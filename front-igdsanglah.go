package front_igdsanglah

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/leekchan/accounting"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func init() {
	http.HandleFunc("/", indexPage)
	http.HandleFunc("/home", homePage)
	http.HandleFunc("/get-info-nocm", getInfoNoCM)
	http.HandleFunc("/tambah-data-kunjungan", tambahDataKunjungan)
	http.HandleFunc("/get-bcp-bulan-ini", getBCPBulanIni)
	http.HandleFunc("/get-bcp-bulan", getBCPBulan)
	http.HandleFunc("/test-db/", testDatabase)
	http.HandleFunc("/get-data-pasien", getDataPasienByLink)
	http.HandleFunc("/edit-data-kunjungan", editDataKunjungan)
	http.HandleFunc("/get-bcp-pdf", getBCPPDF)
	http.HandleFunc("/createkursor", createKursor)
	http.HandleFunc("/hapus-data-kunjungan", hapusDataKunjungan)
	http.HandleFunc("/get-data-tanggal-kunjungan-pasien", getDataTanggalKunjunganPasien)
	http.HandleFunc("/ubah-data-tanggal-kunjungan", ubahDataTanggalKunjungan)
	http.HandleFunc("/get-admin-page", getAdminPage)
	http.HandleFunc("/add-staf", addStaf)
	http.HandleFunc("/hapus-staf", hapusDataStaf)
	http.HandleFunc("/get-detail-pasien", getDetailPasien)
	http.HandleFunc("/ubah-detail-pasien", ubahDetailPasien)
	http.HandleFunc("/get-doc-profile", getDocProfile)
	http.HandleFunc("/ubah-detail-dokter", ubahDetailDokter)
	http.HandleFunc("/get-obat-page", getObatPage)
	http.HandleFunc("/cari-obat", cariObat)
	http.HandleFunc("/tambah-obat", tambahObat)
	http.HandleFunc("/get-isian-obat", getIsianObat)
	http.HandleFunc("/get-data-obat", getDataObat)
	http.HandleFunc("/edit-data-obat", editDataObat)
	http.HandleFunc("/kegiatan-dokter", getKegiatanDokter)
	http.HandleFunc("/tambah-kegiatan-dokter", addKegiatanDokter)
	http.HandleFunc("/hapus-kegiatan", hapusKegiatanDokter)
	http.HandleFunc("/get-kegiatan-bulanan", getContentKegiatanBulanan)
}

// homePage digunakan untuk menampilkan template halaman utama
func homePage(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	w.WriteHeader(200)
	fmt.Fprint(w, GenTemplate(w, ctx, nil, "front-content"))
}

// indexPage digunakan untuk menampilkan isi dari halaman Home.
// Fungsi ini menggunakan gmail sebagai alat login. Kemudian
// menggunakan appengine untuk mengecek apakah user adalah ang-
// gota dari Staff
func indexPage(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	if u == nil {
		url, _ := user.LoginURL(ctx, "/")
		fmt.Fprintf(w, "<a href='%s'>Sign in or Register</a>", url)
		return
	}
	// Membuat link untuk logout.
	logout, _ := user.LogoutURL(ctx, "/")
	if u.Email == "suryasedana@gmail.com" {
		kur, err := getKursor(ctx, u.Email)
		if err != nil {
			// log.Infof(ctx, "Gagal mengambil kursor: %v", err)
			DocumentError(w, ctx, "Gagal mengambil kursor", err, 500)
			return
		}
		// Menyiapkan struct untuk diexecute di template untuk Home
		me := FrontPage{
			LogOut:   logout,
			UserName: "I Wayan Surya Sedana",
			Email:    u.Email,
			Kursor:   kur,
			Peran:    "admin",
		}
		// Menyiapkan script untuk Home
		front := GenTemplate(w, ctx, me, "index", "front-content")
		// Response server
		fmt.Fprint(w, front)
		return
	} else {
		// Mengecek apakah email adalah anggota dari Staff
		staf, err := CekStaff(ctx, u.Email)
		// log.Infof
		if err != nil {
			log.Errorf(ctx, "Email tidak terdaftar, %v", err)
			// Jika bukan, secara otomatis user akan logout
			fmt.Fprintf(w, `<p>Maaf email anda tidak terdaftar dalam sistem. Hubungi admin. </p><a href="%s">Logout</a>`, logout)
			// http.Redirect(w, r, logout, 403)
		} else {
			kur, err := getKursor(ctx, u.Email)
			// log.Infof(ctx, "Email adalah : %v", staf.Email)
			if err != nil {
				DocumentError(w, ctx, "Gagal mengambil kursor", err, 500)
			}
			// Menyiapkan struct untuk diexecute di template untuk Home
			me := FrontPage{
				LogOut:   logout,
				UserName: staf.NamaLengkap,
				Email:    u.Email,
				Kursor:   kur,
				Peran:    staf.Peran,
			}
			// Menyiapkan script untuk Home
			front := GenTemplate(w, ctx, me, "index", "front-content")
			// Response server
			fmt.Fprint(w, front)
		}
	}
}

// GenTemplate menggenerate string HTML untuk ditampilkan di halaman web.
// Fungsi ini membutuhkan parameter w berupa http.ResponseWriter, c
// context.Context, dan n interface, serta slice string dari nama file
// HTML yang akan dijadikan template
func GenTemplate(w http.ResponseWriter, c context.Context, n interface{}, temp ...string) string {
	b := new(bytes.Buffer)
	// funcs adalah list fungsi yang akan diterapkan di template
	funcs := template.FuncMap{
		// jam digunakan untuk mengubah format waktu menjadi hh:mm
		"jam": func(t time.Time) string {
			zone, _ := time.LoadLocation("Asia/Makassar")
			return t.In(zone).Format("15:04")
		},
		// inc digunakan untuk mengubah index dari range naik 1 nomor
		// berguna untuk membuat nomor urut dari slice pada template
		"inc": func(i int) int {
			return i + 1
		},
		// rp digunakan untuk mengubah bilangan int menjadi string
		// dengan mata Uang Rupiah
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
		// rpi digunakan untuk mengubah string bilangan menjadi string
		// mata uang Rupiah
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
		// strtglhari digunakan untuk membuat string tanggal disertai
		// nama hari dengan format Mon, 02/01/2006
		"strtglhari": func(t time.Time) string {
			return t.Format("Mon, 02/01/2006")
		},
		// strtgl digunakan untuk membuat string tanggal dari
		// sebuah type Time dengan format 02/01/2006
		"strtgl": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		// istimezero digunakan untuk mencari tahu apakah type Time nol
		// digunakan untuk memberikan nilai false pada alur if dalam template
		"istimezero": func(t time.Time) bool {
			return t.IsZero()
		},
		// convstrjaga digunakan untuk mengubah string ShiftJaga yang berupa
		// angka menjadi Pagi, Sore dan Malam
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
		// propercaps digunakan untuk memperbaiki string dengan kapitalisasi
		// yang salah
		"propercaps": func(input string) string {
			words := strings.Fields(input)
			smallwords := " dan atau dr. "

			for index, word := range words {
				if strings.Contains(smallwords, " "+word+" ") {
					words[index] = word
				} else {
					words[index] = strings.Title(word)
				}
			}
			return strings.Join(words, " ")
		},
		// tglbcp digunakan untuk membuat nama link bcp tiap bulan
		"tglbcp": func(tgl time.Time, shift string) string {
			if tgl.Hour() < 12 && shift == "3" {
				return tgl.AddDate(0, 0, -1).Format("2006/01")
			} else {
				return tgl.Format("2006/01")
			}
		},
		"umur": func(lahir time.Time) string {
			skrng := timeNowIndonesia()
			yr := skrng.Year() - lahir.Year()
			mn := skrng.Month() - lahir.Month()
			dy := skrng.Day() - lahir.Day()
			if skrng.YearDay() < lahir.YearDay() {
				yr--
				mn = skrng.Month() + 12 - lahir.Month()
			}
			if skrng.Day() < lahir.Day() {
				dy = lahir.AddDate(0, 0, -(lahir.Day())).Day() + lahir.Day() - skrng.Day()
			}
			return fmt.Sprintf("%d Tahun %d Bulan %d Hari", yr, mn, dy)
		},
		"jenkel": func(jk string) string {
			if jk == "1" {
				return "Laki-laki"
			} else {
				return "Perempuan"
			}
		},
		"htmltgl": func(tgl time.Time) string {
			return tgl.Format("2006-01-02")
		},
		"sediaan": func(sed string) bool {
			if sed == "2" || sed == "3" || sed == "4" || sed == "5" {
				return true
			} else {
				return false
			}
		},
		"ubahsediaan": func(sed string) string {
			var m string
			switch sed {
			case "1":
				m = "Tablet"
			case "2":
				m = "Sirup"
			case "3":
				m = "Drop"
			case "4":
				m = "Supositori"
			case "5":
				m = "Vial"
			case "6":
				m = "Cream/Ointment"
			case "7":
				m = "Kapsul"
			case "8":
				m = "Ampul"
			}
			return m
		},
		"convertfloat": func(dosis float64) string {
			return fmt.Sprintf("%.4f", dosis)
		},
	}
	// Membuat template baru
	tmpl := template.New("")
	for k, v := range temp {
		if k == 0 {
			tmp := template.Must(template.New(v + ".html").Funcs(funcs).ParseFiles("templates/" + v + ".html"))
			tmpl = tmp
		}
	}
	// Tambahan template jika ada
	// if k != 0 untuk memastikan bahwa template yang akan diparse bukan yang pertama
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

// DocumentError digunakan untuk merekam error yang terjadi. Fungsi ini akan menambahkan error
// ke log di appengine dan mengirimkan error ke client.
func DocumentError(w http.ResponseWriter, c context.Context, topik string, err error, kode int) {
	msg := "Telah terjadi kesalahan dalam " + topik + " : %v"
	log.Errorf(c, msg, err)
	w.WriteHeader(kode)
	fmt.Fprint(w, "Gagal "+topik)
}

func checkError(w http.ResponseWriter, c context.Context, topik string, err error, kode int) {
	if err != nil {
		msg := "Telah terjadi kesalahan dalam " + topik + " : %v"
		log.Errorf(c, msg, err)
		w.WriteHeader(kode)
		fmt.Fprint(w, "Gagal "+topik)
		return
	}
}

// SendBackSuccess digunakan untuk mengirim response ke Client. Data digunakan untuk mengirim
// JSONArray sedangkan Script, ModalScript, ScriptTambahan masing-masing digunakan untuk
// mengirim script utama, script untuk di modal, dan sript untuk tombol tambahan di modal.
// Field Data sebaiknya berupa Struct yang bisa diubah ke bentuk JSON. Script, ModalScript,
// dan ScriptTambahan adalah data dengan type String
func SendBackSuccess(w http.ResponseWriter, dat interface{}, script, modal, tambahan string) {
	w.WriteHeader(200)
	res := &ResponseJson{
		Data:           dat,
		Script:         script,
		ModalScript:    modal,
		ScriptTambahan: tambahan,
	}
	json.NewEncoder(w).Encode(res)
}

// ChangeStringtoTime digunakan untuk
func ChangeStringtoTime(tgl string) time.Time {
	str, _ := time.ParseInLocation("2006-1-02", tgl, ZonaIndo())
	return str
}
func ZonaIndo() *time.Location {
	zone, _ := time.LoadLocation("Asia/Makassar")
	return zone
}
func timeNowIndonesia() time.Time {
	zone, _ := time.LoadLocation("Asia/Makassar")
	now := time.Now()
	return now.In(zone)
}

func convertJamJaga(tgl time.Time, shift string) time.Time {
	if shift == "3" && tgl.Hour() < 12 {
		tgladjust := time.Date(tgl.Year(), tgl.Month(), tgl.Day(), tgl.Hour()-3, 0, 0, 0, ZonaIndo())
		return tgladjust
	} else {
		return tgl
	}
}

func ProperCapital(input string) string {
	words := strings.Fields(input)
	smallwords := " dan atau dr. "

	for index, word := range words {
		if strings.Contains(smallwords, " "+word+" ") {
			words[index] = word
		} else {
			words[index] = strings.Title(word)
		}
	}
	return strings.Join(words, " ")
}

func CekStaff(ctx context.Context, email string) (*Staff, error) {
	var staf []Staff
	q := datastore.NewQuery("Staff").Filter("Email=", email)

	k, err := q.GetAll(ctx, &staf)
	if err != nil {
		return nil, err
	}
	// user = ""
	// token = ""
	// peran = ""
	// link = ""
	if len(staf) == 0 {
		// user = "no-access"
		return nil, datastore.ErrNoSuchEntity
	}
	doc := &Staff{}
	for _, v := range staf {
		// token = CreateToken(ctx, v.Email)
		// user = v.NamaLengkap
		// peran = v.Peran
		// link = k[0].Encode()
		doc.NamaLengkap = v.NamaLengkap
		doc.LinkID = k[0].Encode()
		doc.Peran = v.Peran
		doc.Email = email
	}
	// var kunci *datastore.Key
	// for _, n := range k {
	// 	kunci = n
	// }
	// log.Infof(ctx, "Link adalah: %v", link)
	if doc == nil {
		return nil, datastore.ErrNoSuchEntity
	} else {
		return doc, nil
	}
}
