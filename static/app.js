import mini from './mini/framework.js'
import Player from './player.js'

// App
const app = mini.createApp('app')
const socket = new WebSocket('ws://localhost:5000/ws')

// Game
const mapState = mini.createState([])
const playerState = mini.createState([])
const blockUpdates = mini.createState([])
let gameboard = undefined
const playerElements = {}

function gameloop(updateType) {
	if (!gameboard) gameboard = drawGameboard(mapState.value)

	switch (updateType) {
		case 'new_game':
			updatePlayerPosition(gameboard)
			mini.render(app, gameboard)
			break

		case 'player_state_update':
			updatePlayerPosition(gameboard)
			break
		case 'map_state_update':
			blockUpdates.value.forEach((update) => {
				if (update.block == 'B') {
					const bomb = createBombElement(update.x, update.y)
					gameboard.appendChild(bomb)
					console.log(bomb)
				}
			})
			break
	}

	function updatePlayerPosition(gameboard) {
		playerState.value.forEach((player) => {
			let playerElement = playerElements[player.id]
			if (!playerElement) {
				playerElement = Player(player)
				playerElements[player.id] = playerElement
				gameboard.appendChild(playerElement.getSprite())
			}

			if (playerElement) {
				const sprite = playerElement.getSprite()
				playerElement.updateSprite(player.direction)
				sprite.style.left = player.x * 64 + 'px'
				sprite.style.top = player.y * 64 + 'px'
			}
		})
	}

	function createBombElement(x, y) {
		return mini.div({ class: 'bomb', style: `left: ${x * 64}px; top: ${y * 64}px` })
	}
}

function drawGameboard(mapdata) {
	return mini.div(
		{ id: 'game' },
		...mapdata.map((row, rowIndex) => {
			return mini.div(
				{ class: 'row flex' },
				...row.map((cell, cellIndex) => {
					let cellClass = 'cell'

					if (cell === '0') {
						cellClass = 'wall'
					} else if (cell === 'd') {
						cellClass = 'destroyable'
					} else if (cell === 'b') {
						cellClass = 'bomb'
					}

					return mini.div({ class: cellClass, 'data-row': rowIndex, 'data-cell': cellIndex })
				})
			)
		})
	)
}

// Event Listeners
socket.onmessage = (e) => {
	const data = JSON.parse(e.data)
	console.log(data)
	switch (data.type) {
		case 'new_game':
			mapState.value = data.map
			playerState.value = data.players
			document.addEventListener('keydown', keyDownHandler)
			document.addEventListener('keyup', keyUpHandler)
			break
		case 'game_over':
			document.removeEventListener('keydown', keyDownHandler)
			document.removeEventListener('keyup', keyUpHandler)
			break
		case 'player_state_update':
			playerState.value = data.players
			break
		case 'map_state_update':
			blockUpdates.value = data.block_updates
			break
	}

	gameloop(data.type)
}

let activeKeys = new Set()

function keyDownHandler(e) {
	const keyDown = e.key.toLowerCase()
	if (!'wsadenter'.includes(keyDown)) return
	activeKeys.add(keyDown)
}

function keyUpHandler(e) {
	const keyUp = e.key.toLowerCase()
	if (!'wsadenter'.includes(keyUp)) return
	activeKeys.delete(keyUp)
}

let sending = true
setInterval(() => {
	if (activeKeys.size > 0) {
		socket.send(JSON.stringify({ type: 'keydown', keys: Array.from(activeKeys) }))
		sending = true
	} else {
		if (sending) {
			socket.send(JSON.stringify({ type: 'keyup', keys: [] }))
			sending = false
		}
	}
}, 50)
