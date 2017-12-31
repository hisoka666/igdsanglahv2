var _scannerIsRunning = false;

function startScanner() {
    document.getElementById("barcode-reader").style.display = "block"
    Quagga.init({
        inputStream: {
            name: "Live",
            type: "LiveStream",
            target: document.querySelector('#scanner-container'),
            constraints: {
                width: 100,
                height: 320,
                facingMode: "environment"
            }, 
        },
        decoder: {
            readers: ["code_128_reader"]
        },
        debug: false,

    }, function (err) {
        if (err) {
            console.log(err);
            return
        }

        console.log("Initialization finished. Ready to start");
        Quagga.start();
        var child = document.getElementsByTagName("canvas")[0]
        child.parentNode.removeChild(child)
        // Set flag to is running
        _scannerIsRunning = true;
    });

    Quagga.onDetected(function (result) {
        document.getElementById("barcode-reader").style.display = "none"
        _scannerIsRunning = false;
        document.getElementById("get-no-cm").value = result.codeResult.code
        Quagga.stop()
    });
}