const url = 'http://localhost:5000' // This would need to be set from configuration.
const cookieSessionID = 'session_id'

$(document).ready(() => {
	// Get session ID from server. We expect a session ID to be returned here.
	// Ideally this request would be handled by a (non-existent) 'client' microservice 
	// instead of turned to the server as this generates a direct dependency.
	//
	// Since we're doing an async request, there's events we could miss here because we're
	// doing an expensive POST request. Ideally, we could use something (maybe a queue)
	// to store any events but for time's sake, I won't be doing that here.
	$.ajax(url + '/new', {
		type: 'POST',
		data: JSON.stringify({ websiteURL: window.location.href }),
		contentType: 'application/json',
		success: (data, status, _) => {
			// Request successful, save session id cookie and start listeners
			Cookies.set(cookieSessionID, 'ay-hello-im-a-session-' + Math.random())
			listenForFirstResize()
			listenForFieldCopyPaste()
			listenForTimeToSubmit()
		},
		error: (_, status, err) => {
			// Similarly, we could also be a bit more clever here (retry strategies, etc)
			console.log('Failed request with status ' + status + ' and error: ' + err)
		},
	})
})

// Adds listener for resize event and removes itself after the first resize
function listenForFirstResize() {
	const resizeTimeoutMs = 1000
	const namespacedResizeEvent = 'resize.listenonce'
	const originalW = $(window).width()
	const originalH = $(window).height()

	let timer
	$(window).on(namespacedResizeEvent, (e) => {
		clearTimeout(timer)
		timer = setTimeout(() => {
			const w = $(window).width()
			const h = $(window).height()

			// Remove listener as we only care about the first resize.
			$(window).off(namespacedResizeEvent);

			// Create event and send it off to the backend.
			ev = {
				eventType: 'windowResize',
				websiteURL: window.location.href,
				sessionID: Cookies.get(cookieSessionID),
				resizeFrom: {
					width: originalW,
					height: originalH,
				},
				resizeTo: {
					width: w,
					height: h,
				},
			}
			postEvent(ev)
		}, resizeTimeoutMs)
	});
}

function listenForFieldCopyPaste() {
	var pasted = new Map(); // Map of input ID to boolean, good for lookups.
	$('input').bind('paste', (e) => {
		if (pasted.get(e.target.id) == null) {
			pasted.set(e.target.id, true)

			ev = {
				eventType: 'copyAndPaste',
				websiteURL: window.location.href,
				sessionID: Cookies.get(cookieSessionID),
				inputID: e.target.id,
			}
			postEvent(ev)
		}
	});
}

function listenForTimeToSubmit() {
	const namespacedKeyUpEvent = 'keyup.listenonce'

	$('input').on(namespacedKeyUpEvent, (_) => {
		const startTime = Date.now()

		// Ignore subsequent `keyup` events
		$('input').off(namespacedKeyUpEvent);

		// Start listening on form submit
		$('form').on('submit', (_) => {
			ev = {
				eventType: 'timeTaken',
				websiteURL: window.location.href,
				sessionID: Cookies.get(cookieSessionID),
				timeSeconds: (Date.now() - startTime) / 1000,
			}
			postEvent(ev)
		})
	})
}

function postEvent(ev) {
	$.ajax(url + '/new_event', {
		type: 'POST',
		data: JSON.stringify(ev),
		contentType: 'application/json',
		error: (_, status, err) => {
			// Similarly here, we could handle this better (if it errored we wouldn't remove it from the queue, for example)
			// For time's sake, I'll just dump the error and keep going
			console.log('Failed request with status ' + status + ' and error: ' + err)
		},
	})
}
