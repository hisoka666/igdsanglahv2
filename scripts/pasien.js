function sendNoCM(){
    // console.log("button pressed")
    var large = document.getElementById("get-no-cm-large").value
    var small = document.getElementById("get-no-cm-small").value
    var nocm = ""
    if (large === "") {
        nocm = small
    } else {
        nocm = large
    }
    // console.log(nocm)
    var reg = /^\d+$/
    if (reg.test(nocm) == false){
        document.getElementById("warning-msg").innerHTML = "Harap masukkan angka"
        document.getElementById("warning-no-cm").style.display = "block"
        nocm = ""
    }else if (nocm.length < 8) {
        document.getElementById("warning-msg").innerHTML = "Nomor CM kurang"
        document.getElementById("warning-no-cm").style.display = "block"
        document.getElementById("form-input-pasien").innerHTML = ""
    } else if (nocm === ""){
        document.getElementById("form-input-pasien").innerHTML = ""
    } else {
        // document.getElementById("warning-no-cm").style.display = "none"
        document.getElementById("warning-no-cm").style.display = "none"
        document.getElementById("loading-animation").style.display="block"
        var payload = {
            "data01": nocm
        }
        // alert("no cm adalah: " + nocm)
        sendPost("/get-info-nocm", JSON.stringify(payload), showNoCMInfo)
        // console.log("no cm adalah: " + nocm)
    }
}

function showNoCMInfo(){
    var js = JSON.parse(document.getElementById("server-response").innerHTML)
    document.getElementById("form-input-pasien").innerHTML = js.script
    document.getElementById("tambah-entri-kunjungan").dataset.link = js.modal
    // document.getElementById("modal01-title-header").innerHTML = "Berhasil"
    // document.getElementById("modal01-content").innerHTML = js.modal
    // document.getElementById("modal01").style.display = "block"
}

function tambahDataKunjungan(){
    var large = document.getElementById("get-no-cm-large").value
    var small = document.getElementById("get-no-cm-small").value
    var nocm = ""
    if (large === "") {
        nocm = small
    } else {
        nocm = large
    }
    var payload = {
        "data01": this.dataset.link,
        "data02": nocm,
        "data03": document.getElementById("nama-pasien").value,
        "data04": document.getElementById("tanggal-lahir-pasien").value,
        "data05": document.getElementById("diagnosis").value,
        "data06": document.getElementById("ats").value,
        "data07": document.getElementById("bagian").value,
        "data08": document.getElementById("iki").value,
        "data09": document.getElementById("shift-jaga").value,
        "data10": document.getElementById("email").innerHTML
    }
    // console.log("payload adalah: " +  JSON.stringify(payload))
    if (payload.data02 === "" || payload.data03 === "" || payload.data04 === "" || payload.data05 === "" || payload.data06 === "0" || payload.data07 === "0" || payload.data08 === "0" || payload.data09 === "0"){
        document.getElementById("warning-msg").innerHTML = "Form belum lengkap"
        document.getElementById("warning-no-cm").style.display = "block"
    } else {
        document.getElementById("warning-no-cm").style.display = "none"
        sendPost("/tambah-data-kunjungan", JSON.stringify(payload), eraseForm)
    }
    // console.log(JSON.stringify(payload))
    
}

function viewDetailPasien(){
    var js = JSON.parse(document.getElementById("server-response").innerHTML)
    document.getElementById("main-content").innerHTML = js.script
    document.getElementById("detail-edit-button").addEventListener("click", viewEditInput)
}

function viewEditInput(){
    document.getElementById("detail-edit-button").removeEventListener("click", viewEditInput)
    var nama = document.getElementById("detail-nama-pts").innerHTML
    var alamat = document.getElementById("detail-alamat").innerHTML
    var jenkel = document.getElementById("jenkel-hidden").innerHTML
    // console.log("jenis kelamin adalah: " + jenkel)
    var tgllahir = document.getElementById("tgl-lahir-hidden").innerHTML
    var jk = '<select name="" id="ubah-detail-jenis-kelamin" class="w3-border w3-input w3-round">' 
            + '<option value="0" selected disabled>Pilih salah satu...</option>'
            + '<option value="1">Laki-laki</option>' 
            + '<option value="2">Perempuan</option></select>'
    document.getElementById("detail-nama-pts").innerHTML = '<input type="text" name="" id="ubah-detail-nama-pts" class="w3-input w3-round w3-border">'
    document.getElementById("detail-tanggal-lahir").innerHTML = '<input type="date" name="" id="ubah-detail-tanggal-lahir" class="w3-round w3-border w3-input">'
    document.getElementById("detail-alamat").innerHTML = '<input type="text" name="" id="ubah-detail-alamat" class="w3-round w3-border w3-input">'
    document.getElementById("detail-jenis-kelamin").innerHTML = jk
    document.getElementById("ubah-detail-nama-pts").value = nama
    document.getElementById("ubah-detail-alamat").value = alamat
    document.getElementById("ubah-detail-tanggal-lahir").value = tgllahir
    document.getElementById("ubah-detail-jenis-kelamin").selectedIndex = parseInt(jenkel)
    document.getElementById("detail-edit-button").innerHTML = "Simpan"
    document.getElementById("detail-edit-button").addEventListener("click", simpanDataBaruPasien)
}

function simpanDataBaruPasien(){
    document.getElementById("detail-edit-button").removeEventListener("click", simpanDataBaruPasien)
    var payload = {
        "data01": document.getElementById("ubah-detail-nama-pts").value,
        "data02": document.getElementById("ubah-detail-tanggal-lahir").value,
        "data03": document.getElementById("ubah-detail-alamat").value,
        "data04": document.getElementById("ubah-detail-jenis-kelamin").value,
        "data05": document.getElementById("detail-edit-button").dataset.link
    }
    // console.log(JSON.stringify(payload))
    // document.getElementById("detail-edit-button").addEventListener("click", viewEditInput)
    sendPost("/ubah-detail-pasien", JSON.stringify(payload), refreshDetailPasien)
    
}

function refreshDetailPasien(){
    document.getElementById("detail-edit-button").addEventListener("click", viewEditInput)
    document.getElementById("detail-edit-button").innerHTML = "Ubah"
    var js = JSON.parse(document.getElementById("server-response").innerHTML)
    document.getElementById("modal01-title-header").innerHTML = "Berhasil"
    document.getElementById("modal01-content").innerHTML = js.modal
    document.getElementById("modal01-tombol-tambahan").style.display = "none"
    document.getElementById("modal01").style.display = "block"
    document.getElementById("modal01-tombol-tutup").addEventListener("click", function(){
        var payload = {
            "data01": document.getElementsByClassName("link-kunjungan")[0].dataset.link
        }
        console.log(JSON.stringify(payload))
        sendPost("/get-detail-pasien", JSON.stringify(payload), viewDetailPasien)
    })
}