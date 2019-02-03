const baseUrl = 'http://localhost:5000' // This would need to be set from configuration.
const cookieSessionID = 'session_id'

$(document).ready(() => {
	// Get session ID from server. We expect a session ID to be returned here.
	// Ideally this request would be handled by a (non-existent) 'client' microservice 
	// instead of turned to the server as this generates a direct dependency.
	//
	// Since we're doing an async request, there's events we could miss here because we're
	// doing an expensive POST request. Ideally, we could use something (maybe a queue)
	// to store any events but for time's sake, I won't be doing that here.
	$.ajax(baseUrl + '/new_session', {
		type: 'POST',
		data: JSON.stringify({ websiteURL: window.location.href }),
		contentType: 'application/json',
		success: (data) => {
			// Request successful, save session id cookie and start listeners
			Cookies.set(cookieSessionID, data.sessionID)

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
	const originalW = $(window).width().toString()
	const originalH = $(window).height().toString()

	let timer
	$(window).on(namespacedResizeEvent, (e) => {
		clearTimeout(timer)
		timer = setTimeout(() => {
			const w = $(window).width().toString()
			const h = $(window).height().toString()

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
			postEvent(ev, baseUrl + '/new_resize_event')
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
			postEvent(ev, baseUrl + '/new_cp_event')
		}
	});
}

function listenForTimeToSubmit() {
	const namespacedKeyUpEvent = 'keyup.listenonce'
	const namespacedSubmitEvent = 'submit.listenonce'

	$('input').on(namespacedKeyUpEvent, (_) => {
		const startTime = Date.now()

		// Ignore subsequent `keyup` events
		$('input').off(namespacedKeyUpEvent);

		// Start listening on form submit
		$('form').on(namespacedSubmitEvent, (e) => {
			e.preventDefault()

			ev = {
				eventType: 'timeTaken',
				websiteURL: window.location.href,
				sessionID: Cookies.get(cookieSessionID),
				timeSeconds: Math.round((Date.now() - startTime) / 1000),
			}
			postEvent(ev, baseUrl + '/new_time_taken_event', () => {
				// Request completed, submit form
				$('form').unbind(namespacedSubmitEvent).submit()
			})
		})
	})
}

function postEvent(ev, url, completeFn) {
	$.ajax(url, {
		type: 'POST',
		data: JSON.stringify(ev),
		contentType: 'application/json',
		complete: completeFn,
		error: (_, status, err) => {
			// Similarly here, we could handle this better (if it errored we wouldn't remove it from the queue, for example)
			// For time's sake, I'll just dump the error and keep going
			console.log('Failed request with status ' + status + ' and error: ' + err)
		},
	})
}
