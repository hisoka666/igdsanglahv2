function getAdminPage(){
    // console.log("admin button fired")
    document.getElementById("modal01-tombol-tutup").removeEventListener("click", getAdminPage)
    document.getElementById("modal01-tombol-tutup").removeEventListener("click", sendHapusStaf)
    servePage("get-admin-page", "main-content")
}

function tambahStaf(){
    var payload = {
        "data01": document.getElementById("admin-add-email").value,
        "data02": document.getElementById("admin-add-nama").value, 
        "data03": document.getElementById("admin-add-peran").value
    }
    // console.log(JSON.stringify(payload))
    if (payload.data01 == "" || payload.data02 == "" || payload.data03 == "0") {
        document.getElementById("admin-warning").style.display = "block"
    } else {
        document.getElementById("admin-warning").style.display = "none"
        sendPost("/add-staf", JSON.stringify(payload), responseAdmin)
    }
}

function responseAdmin(){
    var js = JSON.parse(document.getElementById("server-response").innerHTML)
    // console.log("berhasil : " + js.modal)
    document.getElementById("modal01-tombol-tambahan").style.display = "none"
    document.getElementById("modal01-title-header").innerHTML = "Berhasil"
    document.getElementById("modal01-content").innerHTML = js.modal
    document.getElementById("modal01-tombol-tutup").addEventListener("click", getAdminPage)
    document.getElementById("modal01").style.display = "block"
}

// function modifyStaf(){
//     if (this.value == "1") {
//         hapusStaf(this.dataset.link)
//     } else {
//         ubahDataStaf(this.dataset.link)
//     }
// }

function hapusStaf(){
    var link = this.dataset.link
    document.getElementById("modal01-tombol-tambahan").addEventListener("click", sendHapusStaf)
    modifyModal("Hapus Staf", "Yakin ingin menghapus entri ini?", "", link, "Hapus")
}

function sendHapusStaf(){
    var payload = {
        "data01": document.getElementById("modal01-tombol-tambahan").dataset.link
    }
    sendPost("/hapus-staf", JSON.stringify(payload), responseAdmin)
}

