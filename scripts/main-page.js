function startBarcodeScan(){
        console.log("Scanner: " + _scannerIsRunning)
        if (_scannerIsRunning) {
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