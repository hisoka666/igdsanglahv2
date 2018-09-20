function openKegiatanPage(){
    // alert(document.getElementById("email").innerHTML)
    var payload = {
        "data01": document.getElementById("email").innerHTML
    }
    sendPost("/kegiatan-dokter", JSON.stringify(payload), viewKegiatanDokter)
}

function viewKegiatanDokter(){
    var js = JSON.parse(document.getElementById("server-response").innerHTML)
    document.getElementById("main-content").innerHTML = js.script
}

function sendNoCMForKegiatan(){
    var large = document.getElementById("get-no-cm-large").value
    var small = document.getElementById("get-no-cm-small").value
    var nocm = ""
    if (large === "") {
        nocm = small
    } else {
        nocm = large
    }
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
        document.getElementById("warning-no-cm").style.display = "none"
        document.getElementById("loading-animation").style.display="block"
        var payload = {
            "data01": nocm
        }
        sendPost("/get-info-nocm", JSON.stringify(payload), showNoCMInfoForKegiatan)
    }
}

function showNoCMInfoForKegiatan(){
    var js = JSON.parse(document.getElementById("server-response").innerHTML)
    document.getElementById("namaptskegiatan").innerHTML = js.data.namapts
    document.getElementById("keynocm").value = js.modal
    document.getElementById("form-input-kegiatan").style.display = "block"
}

function tambahKegiatan(){
    var payload = {
        "data01": document.getElementById("keynocm").value,
        "data02": document.getElementById("kegiatan-tindakan").value,
        "data03": document.getElementById("email").innerHTML,
        "data04": document.getElementById("namaptskegiatan").innerHTML
    }
    sendPost("/tambah-kegiatan-dokter", JSON.stringify(payload), viewResponseKegiatan)
}

function viewResponseKegiatan(){
    var js = JSON.parse(document.getElementById("server-response").innerHTML)
    document.getElementById("main-content").innerHTML = js.script
    infoModal("Berhasil", js.modal)
}