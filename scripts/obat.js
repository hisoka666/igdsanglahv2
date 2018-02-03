function openObatPage(){
    servePage("/get-obat-page", "main-content")
}

function buatObat(){
    console.log("buat obat")
    // document.getElementById("obat-input-cari-obat-large").innerHTML = ""
    document.getElementById("obat-input-cari-obat").value = ""
    document.getElementById("obat-main-content").style.display = "none"
    document.getElementById("obat-add-obat").style.display = "block"
}

function viewObatMain(){
    document.getElementById("obat-input-cari-obat").value = ""
    document.getElementById("obat-add-obat").style.display="none"
    document.getElementById("obat-main-content").style.display = "block"
    document.getElementById("form-add-obat").reset()
}

function cariObatLive(){
    var payload = {
        "data01": this.value
    }
    sendPost("/cari-obat", JSON.stringify(payload), obatViewResult)
    // console.log(this.value)
}

function obatViewResult(){
    var js = JSON.parse(document.getElementById("server-response").innerHTML)
    // console.log(document.getElementById("server-response").innerHTML)
    if (js.script == "") {
        var obt = '<a href="#" class="w3-bar-item w3-button" style="width:100%" onclick="buatObat()">Data obat tidak ditemukan. Tambah data obat?</a>'
        document.getElementById("obat-div-hasil").innerHTML = obt
    } else {
        document.getElementById("obat-div-hasil").innerHTML = js.script
    }
}

function optionSediaan(){
    // console.log(this.value)
    if (this.value == "2" || this.value == "3" || this.value == "4" || this.value == "5") {
        document.getElementById("khusus-cairan").style.display="block"
    } else {
        document.getElementById("khusus-cairan").style.display="none"
    }
}

function simpanObat(){
    var payload = {
        "data01": document.getElementById("nama-obat").value,
        "data02": document.getElementById("kandungan-obat").value,
        "data03": document.getElementById("sediaan-obat").value,
        "data04": document.getElementById("takaran-sediaan-obat").value,
        "data05": document.getElementById("kandungan-sediaan-obat").value,
        "data06": document.getElementById("min-dose").value,
        "data07": document.getElementById("max-dose").value,
        "data08": document.getElementById("keterangan").value,
        "data09": document.getElementById("email").innerHTML
    }
    // console.log(JSON.stringify(payload))
    if (payload.data01 == "" || payload.data02 == "" || payload.data03 == "0") {
        console.log("this happend")
        // this.stopPropagation()
        document.getElementById("obat-warning").innerHTML = "Data tidak boleh kosong!"
    } else {
        sendPost("/tambah-obat", JSON.stringify(payload), obatSuccess)
    }
}

function obatSuccess(){
    var js = JSON.parse(document.getElementById("server-response").innerHTML)
    document.getElementById("modal01-title-header").innerHTML = "Berhasil"
    document.getElementById("modal01-content").innerHTML = js.modal
    document.getElementById("modal01-content-02").innerHTML = ""
    document.getElementById("modal01-tombol-tambahan").style.display = "none"
    document.getElementById("modal01").style.display = "block"
    viewObatMain()
}

function viewSearch(){
    var payload = {
        "data01": this.dataset.link
    }
    sendPost("/get-data-obat", JSON.stringify(payload), viewDataObat)
    // console.log(JSON.stringify(payload))
}

function viewDataObat(){
    var js = documen.getElementById("server-response")
}