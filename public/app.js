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

        switch(command) {

            case "d":

                let x = parseInt(fields[1]);
                let y = parseInt(fields[2]);
                let rgb = fields[3].split("-");

                ctx.fillStyle = "rgb("+rgb[0]+","+rgb[1]+","+rgb[2]+")";

                ctx.fillRect( x, y, 1, 1 );

                break;

            case "r":

                let height = fields[1];
                let width = fields[2];

                let screen = document.getElementById('screen')
                screen.setAttribute("height", height);
                screen.setAttribute("width", width);
                ctx = screen.getContext('2d');

                break;

        }
    };
    ws.onclose = function(){
        console.log('closed!');
        window.location.reload(true)
    };



});