package front_igdsanglah

import "time"

type FrontPage struct {
	UserName string   `json:"username"`
	LogOut   string   `json:"logout"`
	Token    string   `json:"token"`
	Email    string   `json:"email"`
	Kursor   []Kursor `json:"kursor"`
	Peran    string   `json:"peran"`
}
type Kursor struct {
	Point string `json:"point"`
	Link  string `json:"link"`
}
type CatchDataJson struct {
	Data1  string `json:"data01"`
	Data2  string `json:"data02"`
	Data3  string `json:"data03"`
	Data4  string `json:"data04"`
	Data5  string `json:"data05"`
	Data6  string `json:"data06"`
	Data7  string `json:"data07"`
	Data8  string `json:"data08"`
	Data9  string `json:"data09"`
	Data10 string `json:"data10"`
}

type ResponseJson struct {
	Data           interface{} `json:"data"`
	Script         string      `json:"script"`
	ModalScript    string      `json:"modal"`
	ScriptTambahan string      `json:"tambahan"`
	Kursor         string      `json:"kursor"`
}

type DataPasien struct {
	NamaPasien string    `json:"namapts"`
	NomorCM    string    `json:"nocm"`
	JenKel     string    `json:"jenkel"`
	Alamat     string    `json:"alamat"`
	TglDaftar  time.Time `json:"tgldaf"`
	TglLahir   time.Time `json:"tgllhr"`
	Umur       time.Time `json:"umur"`
}

type KunjunganPasien struct {
	Diagnosis     string    `json:"diag"`
	LinkID        string    `json:"link"`
	GolIKI        string    `json:"iki"`
	ATS           string    `json:"ats"`
	ShiftJaga     string    `json:"shift"`
	JamDatang     time.Time `json:"jam"`
	Dokter        string    `json:"dokter"`
	Hide          bool      `json:"hide"`
	JamDatangRiil time.Time `json:"jamriil"`
	Bagian        string    `json:"bagian"`
}

type Pasien struct {
	StatusServer string    `json:"stat"`
	TglKunjungan time.Time `json:"tgl"`
	ShiftJaga    string    `json:"shift"`
	ATS          string    `json:"ats"`
	Dept         string    `json:"dept"`
	NoCM         string    `json:"nocm"`
	NamaPasien   string    `json:"nama"`
	Diagnosis    string    `json:"diag"`
	IKI          string    `json:"iki"`
	LinkID       string    `json:"link"`
	TglAsli      time.Time `json:"tglasli"`
	TglLahir     time.Time `json:"tgllahir"`
}

type IndexDataPasien struct {
	Nama string
	NoCM string
}

type TableBCP struct {
	Pasien    []Pasien    `json:"pasien"`
	Title     string      `json:"title"`
	StringTgl string      `json:"strtgl"`
	Total     string      `json:"total"`
	Kursor    string      `json:"kursor"`
	Email     string      `json:"email"`
	IKI       []IKI       `json:"iki"`
	Data06    interface{} `json:"data06"`
}

type IKI struct {
	Tanggal string `json:"tanggal"`
	IKI1    int    `json:"iki1"`
	IKI2    int    `json:"iki2"`
}

type IDDanKunjungan struct {
	ID        DataPasien      `json:"id"`
	Kunjungan KunjunganPasien `json:"kunjungan"`
}

type Staff struct {
	Email       string `json:"email"`
	NamaLengkap string `json:"nama"`
	LinkID      string `json:"link"`
	Peran       string `json:"peran"`
}

type DetailStaf struct {
	NamaLengkap  string    `json:"nama"`
	NIP          string    `json:"nip"`
	NPP          string    `json:"npp"`
	GolonganPNS  string    `json:"golpns"`
	Alamat       string    `json:"alamat"`
	Bagian       string    `json:"bagian"`
	LinkID       string    `json:"link"`
	TanggalLahir time.Time `json:"tgl"`
	Umur         string    `json:"umur"`
}

type KursorIGD struct {
	Bulan string `json:"bulan"`
	Point string `json:"point"`
}

type KegiatanDiLuarIGD struct {
	NamaKegiatan       string `json:"nama"`
	TanggalPelaksanaan string `json:"tanggal"`
}

type DetailPasienPage struct {
	Pasien    DataPasien        `json:"datapts"`
	Kunjungan []KunjunganPasien `json:"kunjungan"`
	LinkID    string            `json:"link"`
}

type IndexObat struct {
	MerkDagang string `json:"merk"`
	Kandungan  string `json:"kandungan"`
	Link       string `json:"link"`
}

type Obat struct {
	MerkDagang    string  `json:"merk"`
	Kandungan     string  `json:"kandungan"`
	LinkID        string  `json:"link"`
	Keterangan    string  `json:"keterangan"`
	SediaanObat   string  `json:"sediaan"`
	MinDose       float64 `json:"mindose"`
	MaxDose       float64 `json:"maxdose"`
	Takaran       float64 `json:"takaran"`
	JmlPerTakaran float64 `json:"jmlpertakaran"`
	Submitter     string  `json:"submitter"`
}
