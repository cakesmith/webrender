document.addEventListener("DOMContentLoaded", function() {

    'use strict';

    var ws = null;
    var ctx = document.getElementById('screen').getContext('2d');

    ws = new WebSocket("ws:/" + location.host + "/ws");

    ws.onopen = function(){
        console.log('connected!');
    };

    ws.onmessage = function(e){
        let fields = e.data.split(" ");
        let command = fields[0];

        if (command === "d") {

            // console.log("draw command " + command);

            let x = parseInt(fields[1]);
            let y = parseInt(fields[2]);
            let rgb = fields[3].split("-");

            ctx.fillStyle = "rgb("+rgb[0]+","+rgb[1]+","+rgb[2]+")";

            ctx.fillRect( x, y, 1, 1 );

        }
    };
    ws.onclose = function(){
        console.log('closed!');
        window.location.reload(true)
    };



});