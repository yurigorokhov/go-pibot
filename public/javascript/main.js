function connectWS() {
    const uri = new URL(document.location.toString());
    let wsUri = `ws://${uri.host}/ws`;
    return new WebSocket(wsUri);
}
function generateCommand(direction, speed) {
    return {
        type: 0, 
        data: {
            speed: speed, 
            direction: direction
        }
    };
}
function filter500( value, type ){
    return value % 50 ? 2 : 1;
}
(() => {
    const slider = document.getElementById('slider');
    let speed = 50;
    const validEntry = {
        37: 3,
        38: 0,
        39: 2,
        40: 1
    };
    document.addEventListener('DOMContentLoaded', (e) => {
        const ws = connectWS();
        const buttons = document.querySelectorAll('.direction-controls i');
        noUiSlider.create(slider, {
            start: [ speed ],
            connect: 'lower',
            range: {
                'min': 0,
                '25%': 25,
                '50%': 50,
                '75%': 75,
                'max': 100
            },
            pips: {
                filter: filter500,
                mode: 'range',
                density: 3 
            }
        });
        slider.noUiSlider.on('update', (value, handle) => {
            speed = parseInt(value[0], 10);
        });
        document.addEventListener('keydown', (e) => {
            console.log(e.keyCode);
            if(validEntry.hasOwnProperty(e.keyCode)) {
                let command = generateCommand(validEntry[e.keyCode], speed);
                ws.send(JSON.stringify(command));
            }
        });
    });
})();