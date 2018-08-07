var lastX, lastY, touchStartTime, lastScrollY

function startup() {
	el = document.getElementById("touch")
	el.addEventListener("touchstart", handleTouchStart, false)
	el.addEventListener("touchend", handleTouchEnd, false)
	el.addEventListener("touchmove", handleTouchMove, false)
	el = document.getElementById("scroll")
	el.addEventListener("touchstart", handleScrollStart, false)
	el.addEventListener("touchmove", handleScrollMove, false)
	el = document.getElementById("textForm")
	el.onsubmit = formSubmit
}

function formSubmit(e) {
	e.preventDefault()
	let foo = {text: e.target[0].value}
	e.target.reset()
	send("/inputtext", JSON.stringify(foo))
}

function handleScrollStart(e) {
	lastScrollY = Math.floor(e.touches[0].screenY)
}
function handleScrollMove(e) {
	e.preventDefault()
	let y = Math.floor(e.touches[0].screenY)
	let dy = y - lastScrollY
	if(Math.abs(dy) > 5) {
		send(dy > 0 ? "/scrollup" : "/scrolldown")
		lastScrollY = y
	}
}

function handleTouchStart(e) {
	lastX = Math.floor(e.touches[0].screenX)
	lastY = Math.floor(e.touches[0].screenY)
	touchStartTime = Date.now()
}
function handleTouchEnd(e) {
	if(Date.now() - touchStartTime < 100) {
		send("/clickmouse")
	}
}
function handleTouchMove(e) {
	e.preventDefault()
	let x = e.touches[0].screenX
	let y = e.touches[0].screenY

	let dx = (x - lastX)
	dx = dx * Math.abs(dx) / 4
	dx = dx > 0 ? Math.ceil(dx) : Math.floor(dx)

	let dy = (y - lastY)
	dy = dy * Math.abs(dy) / 4
	dy = dy > 0 ? Math.ceil(dy) : Math.floor(dy)

	let delta = {dx: dx, dy: dy}
	send("/movemouse", JSON.stringify(delta))

	lastX = x
	lastY = y
	lastT = t
}

function send(url, msg) {
	let xhr = new XMLHttpRequest()
	xhr.open("post", url, true)
	xhr.send(msg)
}
