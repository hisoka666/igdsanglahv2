var _scannerIsRunning = false;
fixOrientation = function(w, h){
    var md = new MobileDetect(window.navigator.userAgent), d = {
        w: w,
        h: h
    };

    if (md.phone() || md.tablet()){
        if (window.matchMedia('(orientation:portrait)').matches){
            if (md.userAgent() !== 'Safari'){
                d.w = h;
                d.h = w;
            }
        }
    }
    return d;
}

function startScanner() {
    // var dim = fixOrientation(320, 100)
    document.getElementById("interactive").style.display = "block"
    Quagga.init({
        inputStream: {
            name: "Live",
            type: "LiveStream",
            // target: document.querySelector('#scanner-container'),
            constraints: {
                // width: dim.w,
                // height: dim.h,
                facingMode: "environment"
            }, 
        },
        decoder: {
            readers: ["code_128_reader"]
        },
        // debug: false,
        multiple: false,

    }, function (err) {
        if (err) {
            console.log(err);
            return
        }

        console.log("Initialization finished. Ready to start");
        Quagga.start();
        // var child = document.getElementsByTagName("canvas")[0]
        // child.parentNode.removeChild(child)
        // Set flag to is running
        _scannerIsRunning = true;
    });

    Quagga.onDetected(function (result) {
        document.getElementById("get-no-cm-small").value = result.codeResult.code
        var payload = {
            "data01": result.codeResult.code
        }
        // alert("hasil adalah: " + result.codeResult.code)
        sendPost("/get-info-nocm", JSON.stringify(payload), showNoCMInfo)
        sendPost("/get-info-nocm", JSON.stringify(payload), showNoCMInfoForKegiatan)
        document.getElementById("interactive").style.display = "none"
        _scannerIsRunning = false;
        Quagga.stop()
    });
}