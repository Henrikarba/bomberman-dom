import { playerState, blockUpdates, mapState, gameloop, playerID, playerCount, chatArea } from './app.js'
import mini from './mini/framework.js'

// update values here for gameloop

const createWebSocket = () => {
	const actionMap = {
		'playerID': updatePlayerID,
		'status': updatePlayerCount,
		'server_full': showServerFull,
		'new_game': initNewGame,
		'player_state_update': updatePlayerState,
		'map_state_update': updateMapState,
		'game_over': endGame,
		'message': chatAreaHandler,
	}

	const socket = new WebSocket('ws://localhost:5000/ws')

	socket.onmessage = (e) => {
		const data = JSON.parse(e.data)
		console.log(data)

		const action = actionMap[data.type]
		if (action) {
			requestAnimationFrame(() => {
				action(data)
				gameloop(data.type)
			})
		}
	}
	return socket
}

export const socket = createWebSocket()

function chatAreaHandler(data) {
	const name = mini.span({ style: 'color: orange;' }, `${data.name}`)
	const msg = mini.div({ style: 'word-wrap: break-word;' }, name, `: ${data.message}`)
	chatArea.appendChild(msg)
}

function updatePlayerID(data) {
	playerID.value = parseInt(data.message)
}

function updatePlayerCount(data) {
	playerCount.value = data.player_count
}

function showServerFull(data) {
	app.innerHTML = data.message
}

function endGame(data) {
	document.removeEventListener('keydown', keyDownHandler)
	document.removeEventListener('keyup', keyUpHandler)
}

function updatePlayerState(data) {
	playerState.value = data.players
}

function updateMapState(data) {
	blockUpdates.value = data.block_updates
}

function manageKeys(socket, sending, activeKeys) {
	if (activeKeys.size > 0) {
		socket.send(JSON.stringify({ type: 'keydown', keys: Array.from(activeKeys) }))
		sending = true
	} else {
		if (sending) {
			socket.send(JSON.stringify({ type: 'keyup', keys: [] }))
			sending = false
		}
	}
}

let activeKeys = new Set()
let sending = false
function initNewGame(data) {
	mapState.value = data.map
	playerState.value = data.players

	document.addEventListener('keydown', keyDownHandler)
	document.addEventListener('keyup', keyUpHandler)

	setInterval(() => {
		manageKeys(socket, sending, activeKeys)
	}, 50)
}

function keyDownHandler(e) {
	// console.log('KEY DOWN')
	const keyDown = e.key.toLowerCase()
	if (!isExcludedInput(e.target)) {
		if (!'wsad '.includes(keyDown)) return
		activeKeys.add(keyDown)
	}
}

function keyUpHandler(e) {
	const keyUp = e.key.toLowerCase()
	if (!isExcludedInput(e.target)) {
		if (!'wsad '.includes(keyUp)) return
		activeKeys.delete(keyUp)
	}
}

function isExcludedInput(target) {
	return target.tagName === 'INPUT'
}
