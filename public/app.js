'use strict';

function styleFrom(rgb) {
 return "rgb("+rgb[0]+","+rgb[1]+","+rgb[2]+")";
}

function processCommand(ctx, data, ws) {

    let fields = data.split(" ");
    let command = fields[0];

    switch(command) {

        case "r":

            let x = fields[1];
            let y = fields[2];
            let w = fields[3];
            let h = fields[4];

            ctx.fillStyle = styleFrom(fields[5].split(':'));

            ctx.fillRect(x, y, w, h);

            break;


        case "req":

            let id = fields[1];

            switch(fields[2]) {

                case "h":
                    let height = getHeight();

                    ws.send("res " + id + " " + height);
                    
                    break;

                case "w":
                    let width = getWidth();

                    ws.send("res " + id + " " + width);

                    break;

            }


    }
}


function getHeight() {
    let height = Math.max(document.documentElement.clientHeight, window.innerHeight || 0);
    return height - 2
}

function getWidth() {
    let width = Math.max(document.documentElement.clientWidth, window.innerWidth || 0);
    return width - 2
}

document.addEventListener("DOMContentLoaded", function() {


    let screen = document.getElementById('screen');

    screen.height = getHeight();
    screen.width = getWidth();

    let ctx = screen.getContext('2d');

    let proto = location.protocol === "http:" ? "ws:" : "wss:";

    let ws = new WebSocket(proto + "//" + location.host + "/ws");

    ws.onopen = function(){
        console.log('connected!');
    };

    ws.onmessage = function(e){

        let commands = e.data.split('\n');

        commands.forEach(function(data) {
            processCommand(ctx, data, ws)
        })

    };

    ws.onclose = function(){
        console.log('closed!');
        window.location.reload(true)
    };

    screen.addEventListener('click', function(e) {
        e.preventDefault();
        ws.send("mc " + e.button + " " + e.offsetX + " " + e.offsetY)
    });

    document.addEventListener('keydown', function(e) {
        e.preventDefault();
        if (e.key.length !== 1) {
            return
        }
        ws.send("k " + e.key.charCodeAt(0))
    });


    let full = document.getElementById("full");

    full.addEventListener('touchmove', function(e) {
        e.preventDefault();
    });

    screen.addEventListener('touchmove', function(e) {
       e.preventDefault();
    });

    document.body.addEventListener('touchmove', function(e) {
        e.preventDefault()
    });


});