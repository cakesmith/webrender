function styleFrom(rgb) {
 return "rgb("+rgb[0]+","+rgb[1]+","+rgb[2]+")";
}


function processCommand(ctx, data) {

    let fields = data.split(" ");
    let command = fields[0];

    switch(command) {

        case "r":

            let x1 = fields[1];
            let y1 = fields[2];
            let w = fields[3];
            let h = fields[4];

            ctx.fillStyle = styleFrom(fields[5].split('-'));

            ctx.fillRect(x1, y1, w, h);

            break;

    }
}

document.addEventListener("DOMContentLoaded", function() {

    'use strict';

    var ws = null;
    var ctx = document.getElementById('screen').getContext('2d');

    var proto = location.protocol === "http:" ? "ws:" : "wss:";

    ws = new WebSocket(proto + "//" + location.host + "/ws");

    ws.onopen = function(){
        console.log('connected!');
    };

    ws.onmessage = function(e){

        let commands = e.data.split('\n');

        commands.forEach(function(data) {
            processCommand(ctx, data)
        })

    };
    ws.onclose = function(){
        console.log('closed!');
        window.location.reload(true)
    };



});