function startBarcodeScan(){
        console.log("Scanner: " + _scannerIsRunning)
        if (_scannerIsRunning) {
            Quagga.stop();
        } else {
            startScanner();
        }
}