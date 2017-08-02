import styles from './css/main.scss';

function init() {
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

init();