/*
     KEA 3rd Semester web-programming exercise
     Copyright (C) 2017  Emil Hummel Clausen <sequoiia>
     This program is free software: you can redistribute it and/or modify
     it under the terms of the GNU General Public License as published by
     the Free Software Foundation, either version 3 of the License, or
     (at your option) any later version.
     This program is distributed in the hope that it will be useful,
     but WITHOUT ANY WARRANTY; without even the implied warranty of
     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
     GNU General Public License for more details.
     You should have received a copy of the GNU General Public License
     along with this program(LICENSE in root directory).  If not, see <http://www.gnu.org/licenses/>.
     Email: emil@hummel.yt
*/

import styles from './css/main.scss';

var changeAdvertBody = function (text) {
    var advertSpot = document.getElementsByClassName("advert")[0];
    advertSpot.innerHTML = text;
};

var advertButton = document.getElementById("changeAdvert");
var timerStartStopButton = document.getElementById("timerStartStop");
var timerResetButton = document.getElementById("timerReset");

var timer = document.getElementById("timer");
var timerStatus = 0;

var timerDOMSeconds = timer.getElementsByClassName("seconds")[0];
var timerDOMMinutes = timer.getElementsByClassName("minutes")[0];


advertButton.addEventListener("click", function () {
    var text = window.prompt("Enter advert text");

    if (text !== null) {
        changeAdvertBody(text);
    }
});

timerResetButton.addEventListener("click", function () {
    timerDOMSeconds.innerHTML = "00";
    timerDOMMinutes.innerHTML = "00";
    timerStatus = 0;
});

timerStartStopButton.addEventListener("click", function () {
    if (timerStatus === 0) {
        timerStatus = 1;
    } else if (timerStatus === 1) {
        timerStatus = 0;
    }
});