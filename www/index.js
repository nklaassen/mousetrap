var lastX, lastY, date, lastT, downT, lastScrollY

function startup() {
	console.log("whatup")
	lastX = 0
	lastY = 0
	date = new Date()
	el = document.getElementById("touch")
	el.addEventListener("touchstart", handleTouchStart, false)
	el.addEventListener("touchend", handleTouchEnd, false)
	el.addEventListener("touchmove", handleTouchMove, false)
	el = document.getElementById("scroll")
	el.addEventListener("touchstart", handleScrollStart, false)
	el.addEventListener("touchend", handleScrollMove, false)
	el.addEventListener("touchmove", handleScrollMove, false)
	el = document.getElementById("textForm")
	el.onsubmit = formSubmit
}

function formSubmit(e) {
	e.preventDefault()
	var foo = {text: e.target[0].value}
	e.target.reset()

	// construct an HTTP request
	var xhr = new XMLHttpRequest();
	xhr.open("post", "/inputtext", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');

	// send the collected data as JSON
	xhr.send(JSON.stringify(foo));
}

function handleScrollStart(e) {
	lastScrollY = Math.floor(e.touches[0].screenY)
}

function handleScrollMove(e) {
	e.preventDefault()
	y = Math.floor(e.touches[0].screenY)
	dy = y - lastScrollY
	if(Math.abs(dy) > 5) {
		// construct an HTTP request
		var xhr = new XMLHttpRequest();
		xhr.open(
			"post",
			dy > 0 ? "/scrollup" : "/scrolldown"
			, true);
		//send the request
		xhr.send();

		lastScrollY = y
	}
}

function handleTouchStart(e) {
	lastX = Math.floor(e.touches[0].screenX)
	lastY = Math.floor(e.touches[0].screenY)
	lastT = Date.now()
	downT = lastT
}

function handleTouchEnd(e) {
	if(Date.now() - downT < 100) {
		// construct an HTTP request
		var xhr = new XMLHttpRequest();
		xhr.open("post", "/clickmouse", true);
		// send the request
		xhr.send();
	}
}

function handleTouchMove(e) {
	e.preventDefault()
	t = Date.now()
	dt = t - lastT

	x = e.touches[0].screenX
	y = e.touches[0].screenY
	dx = (x - lastX)// / dt
	dy = (y - lastY)// / dt
	dx = Math.ceil(dx * Math.abs(dx) / 4)
	dy = Math.ceil(dy * Math.abs(dy) / 4)
	var delta = {dx: dx, dy: dy}

	// construct an HTTP request
	var xhr = new XMLHttpRequest();
	xhr.open("post", "/movemouse", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');

	// send the collected data as JSON
	xhr.send(JSON.stringify(delta));

	lastX = x
	lastY = y
	lastT = t
}
