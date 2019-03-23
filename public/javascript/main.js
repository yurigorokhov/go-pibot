// key codes
const FORWARD = 38;
const BACK = 40;
const LEFT = 37;
const RIGHT = 39;

// commands
const COMMAND_FORWARD = 0;
const COMMAND_BACK = 1;
const COMMAND_LEFT = 2;
const COMMAND_RIGHT = 3;
const COMMAND_STOP = 4;
const COMMAND_FORWARD_RIGHT = 5;
const COMMAND_FORWARD_LEFT = 6;
const COMMAND_BACK_RIGHT = 7;
const COMMAND_BACK_LEFT = 8;

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
function sendCommand(ws, keysPressed, speed) {
    var command;
    if(keysPressed.get(FORWARD) === true) {
        if(keysPressed.get(RIGHT) === true) {
            command = COMMAND_FORWARD_RIGHT;
        } else if(keysPressed.get(LEFT) === true) {
            command = COMMAND_FORWARD_LEFT;
        } else {
            command = COMMAND_FORWARD;
        }
    } else if(keysPressed.get(BACK) === true) {
        if(keysPressed.get(RIGHT) === true) {
            command = COMMAND_BACK_RIGHT;
        } else if(keysPressed.get(LEFT) === true) {
            command = COMMAND_BACK_LEFT;
        } else {
            command = COMMAND_BACK;
        }
    } else if(keysPressed.get(LEFT)) {
        command = COMMAND_LEFT;
    } else if(keysPressed.get(RIGHT)) {
        command = COMMAND_RIGHT;
    } else {
        command = COMMAND_STOP;
    }
    let commandJson = generateCommand(command, speed);
    ws.send(JSON.stringify(commandJson));
}
(() => {
    let currentOrientation = {
        alpha: 0,
        beta: 0,
        gamma: 0
    }
    let calibration = {
        gamma: 0,
        beta: 0,
        gamma: 0
    }

    var videoHref = "http://" + window.location.hostname + ":9090/stream/video.mjpeg";
    //document.getElementById('main-video-img').setAttribute("src", videoHref);
    const slider = document.getElementById('slider');
    let speed = 50;
    const keysPressed = new Map([
        [LEFT, false],
        [FORWARD, false],
        [RIGHT, false],
        [BACK, false]
    ]);
    document.addEventListener('DOMContentLoaded', (e) => {
        const ws = connectWS();
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
        if(window.DeviceOrientationEvent) {
            window.addEventListener('deviceorientation', function(e) {
                currentOrientation.alpha = e.alpha
                currentOrientation.beta = e.beta
                currentOrientation.gamma = e.gamma

                const relativeGamma = currentOrientation.gamma - calibration.gamma
                const relativeBeta = currentOrientation.beta - calibration.beta

                keysPressed.set(FORWARD, false)
                keysPressed.set(BACK, false)
                keysPressed.set(RIGHT, false)
                keysPressed.set(LEFT, false)
                if (relativeGamma < -10 && relativeGamma > -50) {
                    speed = (Math.abs(relativeGamma) - 10) * (100/40)
                    keysPressed.set(FORWARD, true)
                } else if (relativeGamma > 10 && relativeGamma < 50) {
                    speed = (Math.abs(relativeGamma) - 10) * (100/40)
                    keysPressed.set(BACK, true)
                }
                if(relativeBeta < -20) {
                    keysPressed.set(RIGHT, true)
                } else if (relativeBeta > 20) {
                    keysPressed.set(LEFT, true)
                }
            })
        }
        document.addEventListener('keydown', (e) => {

            // Speed controls
            if(e.keyCode === 187 || e.keyCode === 189) {
                if(e.shiftKey) {
                    let modifier = e.keyCode === 187 ? 1 : -1;
                    speed = speed + (modifier);
                    slider.noUiSlider.set(speed)
                }
                return;
            }
            if(keysPressed.has(e.keyCode) === false) {
                 return;
            }

            // Motion control
            keysPressed.set(e.keyCode, true);
        });
        document.addEventListener('keyup', (e) => {
            if(keysPressed.has(e.keyCode) === false) {
                return;
            }

            // Motion control
            keysPressed.set(e.keyCode, false);
        });
        document.getElementById('calibrate').addEventListener('click', () => {
            calibration.alpha = currentOrientation.alpha
            calibration.beta = currentOrientation.beta
            calibration.gamma = currentOrientation.gamma
        })
        setInterval(() => {
            sendCommand(ws, keysPressed, speed);
        }, 100);
    });
})();
