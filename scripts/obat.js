function openObatPage(){
    servePage("/get-obat-page", "main-content")
}

function buatObat(){
    document.getElementById("obat-main-content").style.display = "none"
    document.getElementById("obat-add-obat").style.display = "block"
}

function viewObatMain(){
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
    var js = document.getElementById("server-response").innerHTML
    document.getElementById("obat-div-hasil").innerHTML = js.script
}

function optionSediaan(){
    // console.log(this.value)
    if (this.value == "2" || this.value == "3" || this.value == "4" || this.value == "5") {
        document.getElementById("khusus-cairan").style.display="block"
    } else {
        document.getElementById("khusus-cairan").style.display="none"
    }
}