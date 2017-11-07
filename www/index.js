var lastX, lastY, date, lastT, downT

function startup() {
	console.log("whatup")
	lastX = 0
	lastY = 0
	date = new Date()
	el = document.getElementById("main")
	el.addEventListener("touchstart", handleDown, false)
	el.addEventListener("touchend", handleUp, false)
	el.addEventListener("touchmove", handleMove, false)
}

function handleDown(e) {
	lastX = Math.floor(e.touches[0].screenX)
	lastY = Math.floor(e.touches[0].screenY)
	lastT = Date.now()
	downT = lastT
}

function handleUp(e) {
	if(Date.now() - downT < 100) {
		// construct an HTTP request
		var xhr = new XMLHttpRequest();
		xhr.open("post", "/clickmouse.go", true);
		// send the request
		xhr.send();
	}
}

function handleMove(e) {
	t = Date.now()
	dt = t - lastT

	x = e.touches[0].screenX
	y = e.touches[0].screenY
	dx = 10 * (x - lastX) / dt
	dy = 10 * (y - lastY) / dt
	dx = Math.ceil(dx * Math.abs(dx))
	dy = Math.ceil(dy * Math.abs(dy))

	var delta = {dx: dx, dy: dy}

	// construct an HTTP request
	var xhr = new XMLHttpRequest();
	xhr.open("post", "/movemouse.go", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');

	// send the collected data as JSON
	xhr.send(JSON.stringify(delta));

	lastX = x
	lastY = y
	lastT = t
}
