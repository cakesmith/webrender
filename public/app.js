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

    document.body.addEventListener('touchmove', function(e) {
        e.preventDefault()
    });

    let screen = document.getElementById('screen');

    let ctx = screen.getContext('2d');

    let proto = location.protocol === "http:" ? "ws:" : "wss:";

    let ws = new WebSocket(proto + "//" + location.host + "/ws");

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

    screen.addEventListener('click', function(e) {
        e.preventDefault();
        ws.send("mc " + e.button + " " + e.offsetX + " " + e.offsetY)
    });

        // var mylatesttap;

    screen.addEventListener('touchmove', function(e) {
       e.preventDefault();
    });

    screen.addEventListener('touchstart', function(e) {

        // var now = new Date().getTime();
        // var timesince = now - mylatesttap;
        //
        // if((timesince < 600) && (timesince > 0)){
        //     return
        // }else{
        //     e.preventDefault()
        // }
        //
        // mylatesttap = new Date().getTime();

        let rect = screen.getBoundingClientRect();
        let x = e.clientX - rect.left;
        let y = e.clientY - rect.top;
        ws.send("tc " + x + " " + y)
    });


});