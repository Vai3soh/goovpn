var cfg = '';

function getSelect(){ 
    return window.go.gui.Gui.GetConfigsAndChangeCWD().then(result => {
        var arr = [];
        for (var i = 0; i < result.length; i++) {
            arr[i] = [result[i],result[i]];
        }
        const select_ = new ItcCustomSelect('#inputGroupSelect01', {
            name: 'config', 
            options: arr, 
            onSelected(select) {
                cfg = select.value; 
            },
        });
        return select_;
    }); 
}

let Promise__ = getSelect().then(
    function(value){ 
        return value;
});

function ClearLog(){
    document.getElementById('FormControlTextarea1').innerHTML = "";
}

function getElementByXpath(path) {
    return document.evaluate(path, document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue;
}

var connect = document.getElementById('btnConnect');
var disconnect = document.getElementById('btnDisconnect');
var box = document.getElementsByClassName('itc-select__toggle');
var settings = getElementByXpath("/html/body/div[1]/ul/li");

var isEvent = false;
runtime.EventsOn("rcv:read_log", (msg,count,countReconn) => {
    
    document.getElementById('FormControlTextarea1').innerHTML += msg + "\n";

    if (!isEvent) {
        Promise__.then(select => {
            select.dispose();
            connect.disabled = true;
            connect.style.cursor = "not-allowed"; 
            settings.style.visibility = "hidden";
        });
        isEvent = true;
    }

    if (msg.includes("event name: RECONNECTING")) {
        box.config.style.color = "black";
    }

    if (msg.includes("Server poll timeout, trying next remote entry...")) {
        box.config.style.color = "black";
        count++;
        if (count == countReconn) {
            Promise__.then(select => {
                select.enable();
                disconnect.disabled = true; 
                disconnect.style.cursor = "not-allowed"; 
            });
        }
    }

    if (msg.includes("event name: DISCONNECTED") || msg.includes("UNKNOWN/UNSUPPORTED OPTIONS")) {
        connect.disabled = false; 
        connect.style.cursor = "pointer";
        box.config.style.color = "black";
        settings.style.visibility = "visible";

        Promise__.then(select => {
            select.enable();
            disconnect.disabled = true; 
            disconnect.style.cursor = "not-allowed"; 
        });
        
    }

    if (msg.includes("event name: CONNECTING")) {
        disconnect.disabled = false;
        disconnect.style.cursor = "pointer";
    }

    if (msg.includes("event name: CONNECTED")) {
        box.config.style.color = "green";
    }
})

function Runners() {   
    if (cfg != "") {
        ClearLog();
        window.go.openvpn.TransportOvpnClient.Connect(cfg);
        isEvent = false;
    }
}

function Destroy() {                  
   window.go.openvpn.TransportOvpnClient.Disconnect();
}

function Exist() {
    Getresponse('exist')
}