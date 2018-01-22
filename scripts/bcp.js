function getBCPBulan(){
    var payload = {
        "data01": this.dataset.link,
        "data02": document.getElementById("email").innerHTML,
        "data03": this.innerHTML,
    }

    sendPost("/get-bcp-bulan", JSON.stringify(payload), changeMain)
}


function getBCPBulanIni(){
    var payload = {
        "data01": document.getElementById("email").innerHTML
    }
    sendPost("/get-bcp-bulan-ini", JSON.stringify(payload), changeMain)
}


function changeMain(){
    var js = JSON.parse(document.getElementById("server-response").innerHTML)
    document.getElementById("main-content").innerHTML = js.script
    var x = document.getElementsByClassName("no-urut-tabel-iki-small")
    var y = document.getElementsByClassName("no-urut-tabel-iki-large")
    var i
    for (i = 0; i < x.length; i++){
        x[i].innerHTML = i+1;
        y[i].innerHTML = i+1;
    }
    // server-response berisi data dari type TableBCP
    document.getElementById("server-response").innerHTML = js.modal
    document.getElementById("bcp-pdf-data").value = js.modal
    document.getElementById("modal01-tombol-tambahan").dataset.kursor = js.tambahan
    document.getElementById("modal01-tombol-tambahan").removeEventListener("click", hapusDataKunjungan)
    document.getElementById("modal01-tombol-tambahan").removeEventListener("click", sendEditData)
    document.getElementById("modal01-tombol-tambahan").removeEventListener("click", ubahDataKunjunganPasien)
    // console.log(document.getElementById("server-response").innerHTML)
}

function modifyEntri(){
    // console.log("Link adalah: " + this.dataset.link)
    var payload = {
        "data01": this.dataset.link
    }
    // console.log("Data adalah: " + JSON.stringify(payload))
    switch (this.value){
        case "1":
            // console.log("payload adalah: " + JSON.stringify(payload))
            sendPost("/get-data-pasien", JSON.stringify(payload), getDataEditPasien)
            
            // modifyModal("Edit Entri Kunjungan", "bla,bla", "more bla, bla", "link", "Edit", null)
            
            // sendPost("/edit-data-kunjungan", JSON.stringify(payload), showEditModal)
            break;
        case "2":
            modifyModal("Hapus Entri Kunjungan", "Yakin ingin menghapus entri ini?", "", this.dataset.link, "Hapus", null)
            document.getElementById("modal01-tombol-tambahan").addEventListener("click", hapusDataKunjungan)
            // sendPost("/hapus-data-kunjungan", JSON.stringify(payload), showHapusModal)
            break;
        case "3":
            sendPost("/get-data-tanggal-kunjungan-pasien", JSON.stringify(payload), getDataTanggalKunjunganPasien)
            // modifyModal("Ubah Tanggal Kunjungan", "bla,bla", "more bla, bla", "link", "Edit", null)
            // sendPost("/ubah-tanggal-kunjungan", JSON.stringify(payload), showUbahTanggalKunjungan)
            break;
        case "4":
            // modifyModal("Buat Resep", "bla,bla", "more bla, bla", "link", "Buat Resep", null)
            sendPost("/get-detail-pasien", JSON.stringify(payload), viewDetailPasien)
            break;
        default:
    }
    this.selectedIndex = 0
}

function getDataEditPasien(){
    var js = JSON.parse(document.getElementById("server-response").innerHTML)
    modifyModal("Ubah Data Kunjungan", "", js.script, js.modal, "Ubah")
    document.getElementById("ats").selectedIndex = parseInt(js.data.ats)
    document.getElementById("bagian").selectedIndex = parseInt(js.data.bagian)
    document.getElementById("iki").selectedIndex = parseInt(js.data.iki)
    document.getElementById("shift-jaga").selectedIndex = parseInt(js.data.shift)
    document.getElementById("modal01-tombol-tambahan").addEventListener("click", sendEditData)
}

function sendEditData(){
    var payload = {
        "link" : document.getElementById("modal01-tombol-tambahan").dataset.link,
        "diag" : document.getElementById("diagnosis").value, 
        "bagian" : document.getElementById("bagian").value,
        "iki" : document.getElementById("iki").value,
        "ats" : document.getElementById("ats").value, 
        "shift" : document.getElementById("shift-jaga").value,
        "dokter": document.getElementById("modal01-tombol-tambahan").dataset.index
    }
    // console.log(JSON.stringify(payload))
    sendPost("/edit-data-kunjungan", JSON.stringify(payload), responseEditModal)
}

function responseEditModal(){
    var kur = document.getElementById("bcp-pdf-kursor").value
    if (kur == ""){
        var payload = {
            "data01": document.getElementById("email").innerHTML
        }
        sendPost("/get-bcp-bulan-ini", JSON.stringify(payload), changeMain)
    } else {
        var payload = {
            "data01": document.getElementById("bcp-pdf-kursor").value,
            // "data01": document.getElementById("modal01-tombol-tambahan").dataset.link,
            "data02": document.getElementById("email").innerHTML,
            // "data03": document.getElementById("tanggal-kunjungan").innerHTML
            "data03": document.getElementById("bcp-pdf-bulan").value
        }
        sendPost("/get-bcp-bulan", JSON.stringify(payload), changeMain)
    }
    document.getElementById("modal01-tombol-tambahan").removeEventListener("click", hapusDataKunjungan)
    document.getElementById("modal01-tombol-tambahan").removeEventListener("click", sendEditData)
    document.getElementById("modal01-tombol-tambahan").removeEventListener("click", ubahDataKunjunganPasien)
    document.getElementById("modal01-title-header").innerHTML = "Berhasil"
    document.getElementById("modal01-content").innerHTML = "Berhasil mengubah data kunjungan"
    document.getElementById("modal01-content-02").innerHTML = ""
    document.getElementById("modal01-tombol-tambahan").style.display = "none"
    document.getElementById("modal01").style.display = "block"
}

function getPDF(){
    document.getElementById("bcp-pdf-form").submit()
}

function hapusDataKunjungan(){
    var payload = {
        "data01": document.getElementById("modal01-tombol-tambahan").dataset.link,
        "data02": document.getElementById("server-response").innerHTML
    }
    // console.log(JSON.stringify(payload))
    sendPost("/hapus-data-kunjungan", JSON.stringify(payload), responseEditModal)
}

function getDataTanggalKunjunganPasien(){
    // console.log(document.getElementById("server-response").innerHTML)
    var js = JSON.parse(document.getElementById("server-response").innerHTML)
    document.getElementById("modal01-title-header").innerHTML = "Ubah Tanggal Kunjungan"
    document.getElementById("modal01-content-02").innerHTML = js.script
    document.getElementById("modal01-tombol-tambahan").dataset.link = js.modal
    document.getElementById("modal01-tombol-tambahan").innerHTML = "Ubah"
    document.getElementById("modal01-tombol-tambahan").style.display = "block"
    document.getElementById("modal01-tombol-tambahan").addEventListener("click", ubahDataKunjunganPasien)
    document.getElementById("modal01").style.display = "block"
}

function ubahDataKunjunganPasien() {
    var payload = {
        "data01": document.getElementById("modal01-tombol-tambahan").dataset.link,
        "data03": document.getElementById("bcp-pdf-data").value,
        "data02": document.getElementById("jam-datang-baru").value
    }
    // console.log(JSON.stringify(payload))
    sendPost("/ubah-data-tanggal-kunjungan", JSON.stringify(payload), responseEditModal)
}
