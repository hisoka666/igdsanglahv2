function startBarcodeScan(){
        // console.log("Scanner: " + _scannerIsRunning)
        if (_scannerIsRunning) {
            document.getElementById("interactive").style.display = "none"
            // var child = document.getElementsByTagName("video")[0]
            // child.parentNode.removeChild(child)
            _scannerIsRunning = false;
            Quagga.stop();
        } else {
            startScanner();
        }
}

function openSideBar(){
    document.getElementById("main").style.marginLeft = "25%"
    document.getElementById("my-sidebar").style.width = "25%"
    document.getElementById("my-sidebar").style.display = "block"
    document.getElementById("open-sidebar").style.display = "none"
}

function closeSideBar(){
    document.getElementById("main").style.marginLeft = "0%"
    document.getElementById("my-sidebar").style.display = "none"
    document.getElementById("open-sidebar").style.display = "inline-block"
}

function openMobileSideBar(){
    document.getElementById("my-sidebar").style.width = "100%"
    document.getElementById("my-sidebar").style.display = "block"
}

function servePage(url, target){
    document.getElementById("loading-animation").style.display="block"
    var xhttp = createServConn(target)
    xhttp.open("GET", url, true)
    xhttp.send()
}

function createServConn (targetId){
    var xhttp = new XMLHttpRequest()
    xhttp.onreadystatechange = function (){
        if (this.readyState == 4 && this.status == 200) {
            document.getElementById(targetId).innerHTML = this.responseText
            document.getElementById("loading-animation").style.display="none"
        }
    }

    return xhttp
}

function sendPost(url, content, somefunc){
    document.getElementById("loading-animation").style.display="block"
    var xhttp = new XMLHttpRequest()
    xhttp.onreadystatechange = function (){
        if (this.readyState == 4 && this.status == 200) {
            document.getElementById("server-response").innerHTML = this.responseText
            document.getElementById("loading-animation").style.display="none"
            somefunc()
        } else if (this.readyState == 4 && this.status != 200) {
            // console.log(this.responseText)
            document.getElementById("loading-animation").style.display="none"
            document.getElementById("modal01-title-header").innerHTML = "Peringatan!"
            document.getElementById("modal01-content").innerHTML = this.responseText
            document.getElementById("modal01").style.display="block"
        }
    }
    xhttp.open("POST", url, true)
    xhttp.send(content)
}

function openBCP(id){
    var x = document.getElementById(id)
    if (x.className.indexOf("w3-show") == -1) {
        x.className += " w3-show";
    }else {
        x.className = x.className.replace(" w3-show", "")
    }
}

function pageHome(){
    servePage("/home", "main-content")
}

function eraseForm(){
    document.getElementById("get-no-cm-large").value = ""
    document.getElementById("get-no-cm-small").value = ""
    var js = JSON.parse(document.getElementById("server-response").innerHTML)
    document.getElementById("form-input-pasien").innerHTML = ""
    document.getElementById("modal01-title-header").innerHTML = "Berhasil"
    document.getElementById("modal01-content").innerHTML = js.modal
    document.getElementById("modal01").style.display = "block"
}


function openContent(){
    var id = this.dataset.iden
    var x = document.getElementById(id);
    if (x.className.indexOf("w3-show") == -1) {
        x.className += " w3-show";
    } else { 
        x.className = x.className.replace(" w3-show", "");
    }
}

function modifyModal(title, content01, content02, link, additionalBut){
    document.getElementById("modal01-title-header").innerHTML = title
    document.getElementById("modal01-content").innerHTML = content01
    document.getElementById("modal01-content-02").innerHTML = content02
    document.getElementById("modal01-tombol-tambahan").dataset.link = link
    if (additionalBut !== ""){
        document.getElementById("modal01-tombol-tambahan").innerHTML = additionalBut
        // document.getElementById("modal01-tombol-tambahan").addEventListener("click", someFunc())
        document.getElementById("modal01-tombol-tambahan").style.display = "block"
        document.getElementById("modal01").style.display = "block"
    } else {
        document.getElementById("modal01").style.display = "block"
    }
}

// window.addEventListener("DOMContentLoaded", function(){
//     var myDatepicker = document.querySelector("input[]")
// })