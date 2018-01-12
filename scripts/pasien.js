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