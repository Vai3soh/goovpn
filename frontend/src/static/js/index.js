
function removeElement(el){
    el.parentNode.removeChild(el);
}

function getElementByXpath(path) {
    return document.evaluate(path, document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue;
}

function removeElements() {
    window.go.gui.Gui.IsWindows().then(result => {
        if (result){
            var el = document.getElementById(`legacy_algo`);
            var el_ = document.getElementById(`use_systemd`);
            removeElement(el_);
            el_ = getElementByXpath("/html/body/fieldset[1]/div[1]");
            removeElement(el_);
            removeElement(el);
            el = getElementByXpath("/html/body/fieldset[2]/div[2]/label");
            removeElement(el);
        }
    });

}
removeElements()

function validate(evt,regexp) {
    var theEvent = evt || window.event;

    if (theEvent.type === 'paste') {
        key = event.clipboardData.getData('text/plain');
    } else {
        var key = theEvent.keyCode || theEvent.which;
        key = String.fromCharCode(key);
    }
    var regex = regexp;
    if( !regex.test(key) ) {
        theEvent.returnValue = false;
        if(theEvent.preventDefault) theEvent.preventDefault();
    }
}

function addSelect(id,opt,trg){
    const select_ = new ItcCustomSelect(id, {
        targetValue: trg,
        name: 'ssl', 
        options: opt,
        onSelected(select) {
            let val = select.value
            /*
            if (!isNaN(select.value)) {
                val = Math.trunc(select.value)
            } */
            let objJS = { atrId: id, value: val };
            //objJS = JSON.parse(JSON.stringify(message));
            window.go.gui.Gui.SaveData(objJS,"ssl_cmp");
        },
    });
    return select_;
}

const selectSsl = addSelect('#ssl',[[0,0],[1,1],[2,2]],0);

const selectCmp = addSelect('#cmp',[['yes','yes'],['no','no'],['asym','asym']],'yes');


function SetAndSaveValueOtherOpt(id,buckName) {
    var input = document.getElementById(id);
    if (input.value != "") {
        const record = {
            atrId: id,
            value: input.value,
        };
        window.go.gui.Gui.SaveData(record,buckName);
    } else {
        const record = {
            atrId: id,
            value: "",
        }; 
        window.go.gui.Gui.DeleteData(record,buckName);
    }
} 

function SetAndSaveValueOpenVpnLibrary(id,buckName) {
    var input = document.getElementById(id);
    input.setAttribute('type', input.type);
    let record = {
        atrId: id,
        value: input.type,
    }
    if (input.checked){
        window.go.gui.Gui.SaveData(record,buckName);
        return 
    }
    record = {
        atrId: id,
        value: "",
    }; 
    window.go.gui.Gui.DeleteData(record,buckName);
}

runtime.EventsOn("rcv:save_from_db_general_configure", (msg) => {

    for (var i = 0; i < msg.length; i++) {
        var input = document.getElementById(msg[i].AtrId);
        if (input.type == 'text') {
            input.setAttribute('value', (msg[i].Value));
        } else {
            input.checked = true;
        }
       
    }
});

runtime.EventsOn("rcv:save_from_db_checkbox", (msg) => {

    for (var i = 0; i < msg.length; i++) {
        var input = document.getElementById(msg[i].AtrId);
        if (input != null){
            input.checked = true;
        }
    }
});

runtime.EventsOn("rcv:save_from_db_select", (msg) => {

    if (msg.length != 0) {
        selectCmp.value = msg[0].Value; 
        selectSsl.value = msg[1].Value;
    }
});

runtime.EventsOn("rcv:save_from_db_input", (msg) => {

    for (var i = 0; i < msg.length; i++) {
        var input = document.getElementById(msg[i].AtrId);
        if (input.type == 'checkbox'){
            input.checked = true;
        } else {
            input.setAttribute('value', (msg[i].Value));
        }
    }
});

window.go.gui.Gui.SaveToFrontendParams();

