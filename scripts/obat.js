function openObatPage(){
    servePage("/get-obat-page", "main-content")
}

function buatObat(){
    // console.log("buat obat")
    // document.getElementById("obat-input-cari-obat-large").innerHTML = ""
    document.getElementById("obat-input-cari-obat").value = ""
    document.getElementById("obat-main-content").style.display = "none"
    document.getElementById("obat-add-obat").style.display = "block"
    document.getElementById("link-obat-edit").innerHTML = ""
}

function viewObatMain(){
    document.getElementById("obat-input-cari-obat").value = ""
    document.getElementById("obat-add-obat").style.display="none"
    document.getElementById("obat-main-content").style.display = "block"
    document.getElementById("form-add-obat").reset()
    document.getElementById("obat-div-hasil").innerHTML = ""
    document.getElementById("link-obat-edit").innerHTML = ""
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
        "data09": document.getElementById("email").innerHTML,
        "data10": document.getElementById("link-obat-edit").innerHTML
    }
    // console.log(JSON.stringify(payload))
    if (payload.data01 == "" || payload.data02 == "" || payload.data03 == "0") {
        // console.log("this happend")
        // this.stopPropagation()
        document.getElementById("obat-warning").style.display = "block"
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
        "data01": this.dataset.link,
        "data02": this.innerHTML
    }
    // console.log(JSON.stringify(payload))
    sendPost("/get-isian-obat", JSON.stringify(payload), viewDataObat)
    
}

function viewDataObat(){
    var js = JSON.parse(document.getElementById("server-response").innerHTML)
    document.getElementById("obat-div-hasil").innerHTML = js.script
    // console.log(js.script)
    // console.log("nama obat : " + js.modal)
    document.getElementById("obat-input-cari-obat").value = js.modal
    // document.getElementById("berat-badan").style.display = "none"
}

function funcDewasaAnak(){
    if (this.value == "2"){
        document.getElementById("berat-badan").style.display = "block"
    } else {
        document.getElementById("berat-badan").style.display = "none"
    }
}

function submitDataObat() {
    var payload = {
        "data01": document.getElementById("link-obat").innerHTML,
        "data02": document.getElementById("pilih-dws-anak").value,
        "data03": document.getElementById("berat-badan").value
    }
    // console.log(JSON.stringify(payload))
    if (payload.data02 == "0") {
        document.getElementById("warning-msg").innerHTML = "Pilih dewasa atau anak-anak"
        document.getElementById("obat-warning").style.display = "block"
    } else if (payload.data02 == "2" && payload.data03 == "") {
        document.getElementById("warning-msg").innerHTML = "Berat badan belum diisi"
        document.getElementById("obat-warning").style.display = "block"
    } else {
        // console.log(JSON.stringify(payload))
        sendPost("/get-data-obat", JSON.stringify(payload), viewObatFinal)
    }
}

function viewObatFinal(){
    document.getElementById("obat-div-hasil").innerHTML = document.getElementById("server-response").innerHTML
}

function editObat(){
    var payload = {
        "data01": this.dataset.link
    }
    sendPost("/edit-data-obat", JSON.stringify(payload), viewEditObat)
    // console.log(JSON.stringify(payload))
}

function viewEditObat(){
    var js = JSON.parse(document.getElementById("server-response").innerHTML)
    // console.log(document.getElementById("server-response").innerHTML)
    // console.log(JSON.stringify(js))
    buatObat()
    document.getElementById("nama-obat").value = js.data.merk
    document.getElementById("kandungan-obat").value = js.data.kandungan
    document.getElementById("sediaan-obat").selectedIndex = parseInt(js.data.sediaan)
    document.getElementById("takaran-sediaan-obat").value = js.data.takaran
    document.getElementById("kandungan-sediaan-obat").value = js.data.jmlpertakaran
    document.getElementById("min-dose").value = js.data.mindose
    document.getElementById("max-dose").value = js.data.maxdose
    document.getElementById("keterangan").value = js.data.keterangan
    document.getElementById("link-obat-edit").innerHTML = js.data.link
}