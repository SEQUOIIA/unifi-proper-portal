import styles from './css/main.scss';

let voucherInput;
let voucherButton;

window.onload = function (e) {
    if (authstate === -1) {
        voucherButton = document.getElementById("voucherbutton");
        voucherInput = document.getElementById("vouchercode");
    }
    init();
};

function init() {

    if (authstate === -1) {
        voucherButton.addEventListener("click", function (e) {
            consumeVoucher();
        });
    }

    setInterval(checkAuth, 4000)
}

function checkAuth() {
    let request = new XMLHttpRequest();
    request.onload = function (e) {
        let data = JSON.parse(request.response);

        if (data.error !== 1) {
            if (data.data.authorisation !== authstate) {
                location.reload(true);
            }
        } else {
            window.location.replace("https://google.com");
        }
        data = null;
        request = null;
        e = null;
    };
    request.ontimeout = function (e) {
        e = null;
        request = null;
    };

    request.timeout = 2000;
    request.open("GET", "/api/status");
    request.send();

}

function consumeVoucher() {
    let code = voucherInput.value;
    let request = new XMLHttpRequest();
    let requestPayload = {"code": ""};
    requestPayload.code = code;

    request.onload = function (e) {
        let response = JSON.parse(request.response);

        if (response.success === 1) {
            location.reload(true);
        } else {
            switch (response.error) {
                case 1:
                    window.alert(response.error_message);
                    break;
                case 2:
                    window.location.replace("https://google.com");
                    break;
                default:
                    break;
            }
        }

        requestPayload = null;
        request = null;
        response = null;
        code = null;
    };

    request.open("POST", "/api/voucher");
    request.send(JSON.stringify(requestPayload));
}