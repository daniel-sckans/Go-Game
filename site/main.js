let c;
let ctx;

const ws = new WebSocket("wss://race-game-go.herokuapp.com/ws")
let version = '0.1.1'

class Box {
    constructor(x, y, w, h) {
        this.x = x
        this.y = y
        this.width = w
        this.height = h
    }

    update() {
        ctx.fillRect(this.x, this.y, this.width, this.height)
    }

    checkClick(x, y){
        if(x <= this.x+this.width && x >= this.x){
            if(y <= this.y+this.height && y >= this.y){
                return true
            }
        }
    }
}

window.onload = init

let box1 = new Box(0,0,0,0)
let box2 = new Box(0,0,0,0)
let color

let started
document.onclick = (e) => {
    if(!started){
        started = true
    } else {
        const data = {}
        if(box1.checkClick(e.clientX, e.clientY)){
            if(color>=.5){
                ws.send("green")
            } else {
                ws.send("red")
            }
        }
        if(box2.checkClick(e.clientX, e.clientY)){
            if(color<.5){
                ws.send("green")
            } else {
                ws.send("red")
            }
        }
    }
}

function resize(){
    ctx.canvas.width  = window.innerWidth;
    ctx.canvas.height = window.innerHeight;
}

let score = 0

function init(){
    c = document.querySelector('canvas')
    ctx = c.getContext('2d')
    resize()

    
    ws.onopen = () => {console.log("Websocket Opened");ws.send('loaded')}
    ws.onclose = () => {console.log("Websocket Closed")}
    ws.onerror = (error) => {console.log("Websocket Error: ", error)}
    
    ws.onmessage = (msg) => {
        message_data = JSON.parse(msg.data)
        box1.x = message_data.x1 * c.width
        box1.y = message_data.y1 * c.height
        box1.width = message_data.w1 * c.width
        box1.height = message_data.h1 * c.height
        box2.x = message_data.x2 * c.width
        box2.y = message_data.y2 * c.height
        box2.width = message_data.w2 * c.width
        box2.height = message_data.h2 * c.height
        color = message_data.c
        if(message_data.message){
            score = eval(score+message_data.message)
        }
    }

    window.requestAnimationFrame(gameLoop)
}

function gameLoop(timeStamp){
    window.onresize = resize()
    draw()
    window.requestAnimationFrame(gameLoop)
}

function draw(){
    ctx.fillRect(0,0,c.width,c.height)

    if(color>=.5){
        ctx.fillStyle = 'green'
        box1.update()
        ctx.fillStyle = 'red'
        box2.update()
    }else{
        ctx.fillStyle = 'green'
        box2.update()
        ctx.fillStyle = 'red'
        box1.update()
    }

    ctx.fillStyle = 'white'
    ctx.font = 'bold 40px Verdana'
    ctx.textAlign = "left"
    ctx.fillText(''+score, 5, 45)

    if(!started){
        ctx.fillStyle = 'black'
        ctx.font = 'bold 20px Verdana'
        ctx.fillRect(0,0,c.width,c.height)
        ctx.textAlign = "center"
        ctx.translate(c.width/2,c.height/2)
        ctx.strokeStyle = '#ffffff'
        ctx.fillStyle = '#ffffff'
        ctx.lineWidth = 5
        ctx.strokeRect(0-90, -25, 180, 50)
        ctx.fillText('CLICK TO BEGIN', 0, 7, 150)
        ctx.textAlign = "right"
        ctx.fillText(version,c.width/2-5, c.height/2-10)
    }
}