function openDocProfile(){
    var payload = {
        "data01": document.getElementById("email").innerHTML
    }
    document.getElementById("modal01-tombol-tutup").removeEventListener("click", openDocProfile)
    sendPost("/get-doc-profile", JSON.stringify(payload), viewDocProfile)
}

function viewDocProfile(){
    var js = JSON.parse(document.getElementById("server-response").innerHTML)
    document.getElementById("main-content").innerHTML = js.script
    document.getElementById("ubah-data-dokter").addEventListener("click", viewUbahDataDokter)
}

function viewUbahDataDokter(){
    var tgllahir = ""
    if (document.getElementById("doc-tanggal-lahir").innerHTML == ""){
        tgllahir = ""
    }else {
        tgllahir = document.getElementById("html-doc-tgl-lahir").innerHTML
    }
    var nama = document.getElementById("doc-nama-lengkap").innerHTML
    var nipnpp = document.getElementById("doc-nip-npp").innerHTML
    var gol = document.getElementById("doc-golongan").innerHTML
    var bagian = document.getElementById("doc-tempat-tugas").innerHTML
    document.getElementById("doc-nama-lengkap").innerHTML = '<input type="text" name="" id="ubah-doc-nama-lengkap" class="w3-input w3-round w3-border">'
    document.getElementById("doc-tanggal-lahir").innerHTML = '<input type="date" name="" id="ubah-doc-tanggal-lahir" class="w3-input w3-round w3-border">'
    document.getElementById("doc-nip-npp").innerHTML = '<input type="text" name="" id="ubah-doc-nip-npp" class="w3-round w3-input w3-border">'
    document.getElementById("doc-golongan").innerHTML = '<input type="text" name="" id="ubah-doc-golongan" class="w3-round w3-border w3-input">'
    document.getElementById("doc-tempat-tugas").innerHTML = '<input type="text" name="" id="ubah-doc-tempat-tugas" class="w3-round w3-border w3-input">'
    document.getElementById("ubah-doc-nama-lengkap").value = nama
    document.getElementById("ubah-doc-tanggal-lahir").value = tgllahir
    document.getElementById("ubah-doc-golongan").value = gol
    document.getElementById("ubah-doc-tempat-tugas").value = bagian
    document.getElementById("ubah-doc-nip-npp").value = nipnpp
    document.getElementById("ubah-data-dokter").innerHTML = "Simpan"
    document.getElementById("ubah-data-dokter").addEventListener("click", simpanDataDokter)
    document.getElementById("ubah-data-dokter").removeEventListener("click", viewUbahDataDokter)
}

function simpanDataDokter(){
    var payload = {
        "data01": document.getElementById("ubah-doc-nama-lengkap").value,
        "data02": document.getElementById("ubah-doc-tanggal-lahir").value,
        "data03": document.getElementById("ubah-doc-golongan").value,
        "data04": document.getElementById("ubah-doc-tempat-tugas").value,
        "data05": document.getElementById("ubah-doc-nip-npp").value,
        "data06": document.getElementById("email").innerHTML,
        "data07": document.getElementById("ubah-data-dokter").dataset.link,
        "data08": document.getElementById("ubah-data-dokter").dataset.linkpar
    }
    // console.log(JSON.stringify(payload))
    document.getElementById("ubah-data-dokter").removeEventListener("click", simpanDataDokter)
    document.getElementById("ubah-data-dokter").addEventListener("click", viewUbahDataDokter)
    sendPost("/ubah-detail-dokter", JSON.stringify(payload), refreshViewDetailDokter)
}

function refreshViewDetailDokter(){
    var js = JSON.parse(document.getElementById("server-response").innerHTML)
    document.getElementById("modal01-title-header").innerHTML = "Berhasil"
    document.getElementById("modal01-content").innerHTML = js.modal
    document.getElementById("modal01-tombol-tutup").addEventListener("click", openDocProfile)
    document.getElementById("modal01").style.display = "block"
}