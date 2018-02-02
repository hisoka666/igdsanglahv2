package front_igdsanglah

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func getBCPPDF(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	data := &TableBCP{}
	json.Unmarshal([]byte(r.FormValue("data")), data)
	// log.Infof(ctx, "Data adalah: %v", data.Pasien)
	for i, j := 0, len(data.Pasien)-1; i < j; i, j = i+1, j-1 {
		data.Pasien[i], data.Pasien[j] = data.Pasien[j], data.Pasien[i]
	}
	pdf, err := CreateBCPPDF(ctx, data.Pasien, data.IKI, r.FormValue("email"))
	if err != nil {
		DocumentError(w, ctx, "membuat pdf", err, 500)
		return
	}
	w.Header().Set("Content-type", "application/pdf")
	if _, err := pdf.WriteTo(w); err != nil {
		fmt.Fprintf(w, "%s", err)
	}

}

func TglBCPPDF(tgl time.Time, shift string) string {
	if tgl.Hour() < 12 && shift == "3" {
		return tgl.AddDate(0, 0, -1).Format("02-01-2006")
	} else {
		return tgl.Format("02-01-2006")
	}
}
func getDocDetail(c context.Context, k *datastore.Key) (*DetailStaf, error) {
	doc := &Staff{}
	q := datastore.NewQuery("DetailStaf").Ancestor(k)
	t := q.Run(c)
	det := &DetailStaf{}
	for {
		_, err := t.Next(det)
		if err == datastore.Done {
			break
		}
		if err == datastore.ErrNoSuchEntity {
			err := datastore.Get(c, k, doc)
			if err != nil {
				return nil, err
			}
			detail := &DetailStaf{
				NamaLengkap: doc.NamaLengkap,
				LinkID:      datastore.NewKey(c, "DetailStaf", "", 0, k).Encode(),
			}
			_, err = datastore.Put(c, datastore.NewKey(c, "DetailStaf", "", 0, k), detail)
			if err != nil {
				return nil, err
			}
			return detail, nil
		}
	}
	return det, nil
}
func GetResumeIKI(l []IKI) (string, string, string, string, string, string, string) {
	var a, b, c, d int

	for k, v := range l {
		switch {
		case k < 16:
			a = a + v.IKI1
			b = b + v.IKI2
		case k >= 16:
			c = c + v.IKI1
			d = d + v.IKI2
		}
	}
	f := float32(a+c) * 0.0032
	g := float32(b+d) * 0.01
	e := f + g
	return strconv.Itoa(a), strconv.Itoa(b), strconv.Itoa(a + c), strconv.Itoa(b + d), fmt.Sprintf("%.4f", f), fmt.Sprintf("%.4f", g), fmt.Sprintf("%.4f", e)
}
func CreateBCPPDF(c context.Context, pts []Pasien, iki []IKI, email string) (*bytes.Buffer, error) {
	staf, err := CekStaff(c, email)
	if err != nil {
		return nil, err
	}
	k, err := datastore.DecodeKey(staf.LinkID)
	if err != nil {
		return nil, err
	}

	det, err := getDocDetail(c, k)
	if err != nil {
		return nil, err
	}
	ikia, ikib, ikic, ikid, ikie, ikif, ikig := GetResumeIKI(iki)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)
	// Tabel IKI \\\\\\\\\\\\\\\\///////////////////////////////////////////////
	pdf.AddPageFormat("L", gofpdf.SizeType{Wd: 210, Ht: 297})
	pdf.Cell(160, 6, "Bukti Kegiatan Harian")
	pdf.Cell(120, 6, ("Nama Pegawai: " + det.NamaLengkap))
	pdf.Ln(-1)
	pdf.Cell(160, 6, "Pegawai RSUP Sanglah Denpasar")
	pdf.Cell(120, 6, ("NIP/Gol: " + det.NIP + "/" + det.GolonganPNS))
	pdf.Ln(-1)
	pdf.Cell(160, 6, ("Bulan: " + pts[0].TglKunjungan.Format("Jan, 2006")))
	pdf.Cell(120, 6, "Tempat Tugas: IGD RSUP Sanglah")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(10, 20, "No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 20, "Uraian", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 20, "Poin", "1", 0, "C", false, 0, "")
	pdf.CellFormat(176, 10, "Jumlah Kegiatan Harian", "1", 2, "C", false, 0, "")
	// range list iki

	for i := 1; i < 17; i++ {
		pdf.CellFormat(11, 10, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	}
	pdf.SetXY(266, 28)
	pdf.CellFormat(25, 20, "Jumlah Poin", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(10, 24, "1", "1", 0, "C", false, 0, "")

	pdf.MultiCell(50, 6, "Melakukan pelayanan medik umum (per pasien : pemeriksaan rawat jalan, IGD, visite rawat inap, tim medis diskusi", "1", "L", false)
	pdf.SetXY(70, 48)
	pdf.CellFormat(20, 24, "0,0032", "1", 0, "C", false, 0, "")
	for k, v := range iki {
		if k < 16 {
			if v.IKI1 == 0 && v.IKI2 == 0 {
				pdf.CellFormat(11, 24, "", "1", 0, "C", false, 0, "")
			} else {
				pdf.CellFormat(11, 24, strconv.Itoa(v.IKI1), "1", 0, "C", false, 0, "")
			}
		}
	}
	// for i := 1; i < 17; i++ {
	// 	pdf.CellFormat(11, 24, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	// }
	pdf.CellFormat(25, 24, ikia, "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(10, 12, "2", "1", 0, "C", false, 0, "")
	pdf.MultiCell(50, 6, "Melakukan tindakan medik umum tingkat sederhana (per tindakan)", "1", "L", false)
	pdf.SetXY(70, 72)
	pdf.CellFormat(20, 12, "0,01", "1", 0, "C", false, 0, "")
	for k, v := range iki {
		if k < 16 {
			if v.IKI1 == 0 && v.IKI2 == 0 {
				pdf.CellFormat(11, 12, "", "1", 0, "C", false, 0, "")
			} else {
				pdf.CellFormat(11, 12, strconv.Itoa(v.IKI2), "1", 0, "C", false, 0, "")
			}
		}
	}
	// for i := 1; i < 17; i++ {
	// 	pdf.CellFormat(11, 12, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	// }
	pdf.CellFormat(25, 12, ikib, "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.Ln(-1)
	// Baris ke dua
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(10, 20, "No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 20, "Uraian", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 20, "Poin", "1", 0, "C", false, 0, "")
	pdf.CellFormat(176, 10, "Jumlah Kegiatan Harian", "1", 2, "C", false, 0, "")
	for i := 17; i < 32; i++ {
		pdf.CellFormat(11, 10, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	}
	pdf.SetFont("Arial", "B", 7)
	pdf.MultiCell(11, 5, "Jumlah Poin", "1", "C", false)
	pdf.SetFont("Arial", "B", 9)
	pdf.SetXY(266, 96)
	pdf.MultiCell(25, 20, "Jumlah X Poin", "1", "C", false)
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(10, 24, "1", "1", 0, "C", false, 0, "")

	pdf.MultiCell(50, 6, "Melakukan pelayanan medik umum (per pasien : pemeriksaan rawat jalan, IGD, visite rawat inap, tim medis diskusi", "1", "L", false)
	pdf.SetXY(70, 116)
	pdf.CellFormat(20, 24, "0,0032", "1", 0, "C", false, 0, "")
	for k, v := range iki {
		if k >= 16 {
			if v.IKI1 == 0 && v.IKI2 == 0 {
				pdf.CellFormat(11, 24, "", "1", 0, "C", false, 0, "")
			} else {
				pdf.CellFormat(11, 24, strconv.Itoa(v.IKI1), "1", 0, "C", false, 0, "")
			}
		}
	}
	// for i := 17; i <= 32; i++ {
	// 	pdf.CellFormat(11, 24, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	// }
	pdf.CellFormat(11, 24, ikic, "1", 0, "C", false, 0, "")

	pdf.CellFormat(25, 24, ikie, "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(10, 12, "2", "1", 0, "C", false, 0, "")
	pdf.MultiCell(50, 6, "Melakukan tindakan medik umum tingkat sederhana (per tindakan)", "1", "L", false)
	pdf.SetXY(70, 140)
	pdf.CellFormat(20, 12, "0,01", "1", 0, "C", false, 0, "")
	for k, v := range iki {
		if k >= 16 {
			if v.IKI1 == 0 && v.IKI2 == 0 {
				pdf.CellFormat(11, 12, "", "1", 0, "C", false, 0, "")
			} else {
				pdf.CellFormat(11, 12, strconv.Itoa(v.IKI2), "1", 0, "C", false, 0, "")
			}
		}
	}
	// for i := 17; i <= 32; i++ {
	// 	pdf.CellFormat(11, 12, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	// }
	pdf.CellFormat(11, 12, ikid, "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 12, ikif, "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(256, 6, "Jumlah Point X Volume kegiatan pelayanan", "1", 0, "R", false, 0, "")
	pdf.CellFormat(25, 6, ikig, "1", 1, "C", false, 0, "")
	pdf.CellFormat(256, 6, "Target Point kegiatan pelayanan", "1", 0, "R", false, 0, "")
	pdf.CellFormat(25, 6, "1,111", "1", 1, "C", false, 0, "")
	pdf.Ln(-1)

	kegiatan := getKegiatanDiLuarIGD(pts)
	if len(kegiatan) != 0 {
		pdf.Cell(40, 6, "Kegiatan di luar fungsional: ")
		for k, v := range kegiatan {
			pdf.Cell(30, 6, (strconv.Itoa(k+1) + ". " + v.NamaKegiatan + " (" + v.TanggalPelaksanaan + ") "))
			pdf.Ln(-1)
			pdf.Cell(40, 6, "")
		}
	}

	////////////////// Buku Catatan Pasien ///////////////////////////////
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	wd := pdf.GetStringWidth("Buku Catatan Pribadi")
	pdf.SetX((210 - wd) / 2)
	pdf.Cell(wd, 9, "Buku Catatan Pribadi")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(20, 5, "Nama")
	pdf.Cell(105, 5, (": " + det.NamaLengkap))
	pdf.Ln(-1)
	pdf.Cell(20, 5, "Bulan")
	pdf.Cell(105, 5, (": " + pts[0].TglKunjungan.Format("Jan, 2006")))
	pdf.Ln(-1)
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(9, 20, "No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(18, 20, "Tanggal", "1", 0, "C", false, 0, "")
	pdf.CellFormat(17, 20, "No. CM", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 20, "Nama", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 20, "Diagnosis", "1", 0, "C", false, 0, "")

	pdf.MultiCell(20, 5, "Melakukan pelayanan medik umum", "1", "C", false)

	pdf.SetXY(174, 35)
	pdf.MultiCell(25, 4, "Melakukan tindakan medik umum tingkat sederhana", "1", "C", false)

	for k, v := range pts {
		pdf.SetFont("Arial", "", 8)
		diag := ProperCapital(v.Diagnosis)
		if len(diag) > 20 {
			diag = diag[:21]
		}
		// 11/02/1987
		tang := v.TglKunjungan.Format("02-01-2006")
		if v.TglKunjungan.Hour() > 0 && v.TglKunjungan.Hour() < 12 && v.ShiftJaga == "3" {
			tang = v.TglKunjungan.AddDate(0, 0, -1).Format("02-01-2006")
		}
		num := strconv.Itoa(k + 1)
		nocm := v.NoCM
		nam := ProperCapital(v.NamaPasien)
		if len(nam) > 25 {
			nam = nam[:26]
		}
		pdf.CellFormat(9, 7, num, "1", 0, "C", false, 0, "")
		pdf.CellFormat(18, 7, tang, "1", 0, "C", false, 0, "")
		pdf.CellFormat(17, 7, nocm, "1", 0, "C", false, 0, "")
		pdf.CellFormat(60, 7, nam, "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 7, diag, "1", 0, "L", false, 0, "")
		pdf.SetFont("ZapfDingbats", "", 8)
		if v.IKI == "1" {
			pdf.CellFormat(20, 7, "4", "1", 0, "C", false, 0, "")
			pdf.CellFormat(25, 7, "", "1", 0, "C", false, 0, "")
			pdf.Ln(-1)
		} else {
			pdf.CellFormat(20, 7, "", "1", 0, "C", false, 0, "")
			pdf.CellFormat(25, 7, "4", "1", 0, "C", false, 0, "")
			pdf.Ln(-1)
		}
	}

	t := new(bytes.Buffer)
	err = pdf.Output(t)
	if err != nil {
		return nil, err
		// fmt.Fprint(w, "gagal membuat pdf")
		// log.Fatalf("Error reading pdf %v", err)
	}
	return t, nil
}

func getKegiatanDiLuarIGD(p []Pasien) []KegiatanDiLuarIGD {

	keg := []KegiatanDiLuarIGD{}
	for _, v := range p {
		if v.NoCM == "00000000" || v.NoCM == "00000001" || v.NoCM == "00000002" {
			g := KegiatanDiLuarIGD{
				NamaKegiatan:       v.Diagnosis,
				TanggalPelaksanaan: v.TglKunjungan.Format("02"),
			}
			keg = append(keg, g)
		}
	}
	return keg
}
