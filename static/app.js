import mini from './mini/framework.js'
import { drawStartMenu } from './start.js'
import Player from './player.js'

// App
const app = mini.createApp('app')
export const socket = new WebSocket('ws://localhost:5000/ws')

// Views
const startMenu = drawStartMenu()
mini.render(app, startMenu)

// Game
const mapState = mini.createState([])
const playerState = mini.createState([])
const blockUpdates = mini.createState([])
let gameboard = undefined
const playerElements = {}
let playerCount = 0

function gameloop(updateType) {
	switch (updateType) {
		case 'status':
			console.log(playerCount)
			break
		case 'new_game':
			if (!gameboard) gameboard = drawGameboard(mapState.value)
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
				}
				if (update.block == 'ex') {
					const RemoveBomb = document.querySelector(`.bomb[x="${update.x}"][y="${update.y}"]`)
					if (RemoveBomb) {
						RemoveBomb.remove()
					}
					const explosion = createExplosionElement(update.x, update.y)
					gameboard.appendChild(explosion)
				}
				if (update.block == 'f') {
					const flame = createFlameElement(update.x, update.y)
					gameboard.appendChild(flame)
				}
				if (update.block == 'e') {
					const changeCell = document.querySelector(`.destroyable[data-row="${update.y}"][data-cell="${update.x}"]`)
					if (changeCell) {
						changeCell.className = 'cell'
					}
					const RemovePower1 = document.querySelector(`.power1[x="${update.x}"][y="${update.y}"]`)
					if (RemovePower1) {
						RemovePower1.remove()
					}
					const RemovePower2 = document.querySelector(`.power2[x="${update.x}"][y="${update.y}"]`)
					if (RemovePower2) {
						RemovePower2.remove()
					}
					const RemovePower3 = document.querySelector(`.power3[x="${update.x}"][y="${update.y}"]`)
					if (RemovePower3) {
						RemovePower3.remove()
					}
					const RemoveFlame = document.querySelector(`.explosion[x="${update.x}"][y="${update.y}"]`)
					if (RemoveFlame) {
						RemoveFlame.remove()
					}
					const RemoveExplosion = document.querySelector(`.explosion[x="${update.x}"][y="${update.y}"]`)
					if (RemoveExplosion) {
						RemoveExplosion.remove()
					}
				}
				if (update.block == 'p1') {
					const changeCell = document.querySelector(`.destroyable[data-row="${update.y}"][data-cell="${update.x}"]`)
					if (changeCell) {
						changeCell.className = 'cell'
					}
					const power = createPowerElement(update.x, update.y, 1)
					gameboard.appendChild(power)
					console.log(power)
				}
				if (update.block == 'p2') {
					const changeCell = document.querySelector(`.destroyable[data-row="${update.y}"][data-cell="${update.x}"]`)
					if (changeCell) {
						changeCell.className = 'cell'
					}
					const power = createPowerElement(update.x, update.y, 2)
					gameboard.appendChild(power)
					console.log(power)
				}
				if (update.block == 'p3') {
					const changeCell = document.querySelector(`.destroyable[data-row="${update.y}"][data-cell="${update.x}"]`)
					if (changeCell) {
						changeCell.className = 'cell'
					}
					const power = createPowerElement(update.x, update.y, 3)
					gameboard.appendChild(power)
					console.log(power)
				}
			})
			break
	}

	function updatePlayerPosition(gameboard) {
		playerState.value.forEach((player) => {
			let playerElement = playerElements[player.id]
			if (!playerElement && player.lives > 0) {
				playerElement = Player(player)
				playerElements[player.id] = playerElement
				gameboard.appendChild(playerElement.getSprite())
			}

			if (player.lives <= 0) {
				let removePlayer = document.getElementById(`player${player.id}`)
				removePlayer.style.display = 'none'
			}

			if (playerElement) {
				const sprite = playerElement.getSprite()
				playerElement.updateSprite(player.direction)
				sprite.style.left = player.x * 64 + 'px'
				sprite.style.top = player.y * 64 + 'px'
				if (player.damaged) {
					sprite.classList.add('damaged')
					setTimeout(() => {
						sprite.classList.remove('damaged')
					}, 2000)
				}
			}
		})
	}

	function createPowerElement(x, y, power) {
		let powerElement
		if (power === 1) {
			powerElement = mini.div({
				class: 'power1',
				style: `left: ${x * 64}px; top: ${y * 64}px`,
				x: x,
				y: y,
			})
		} else if (power === 2) {
			powerElement = mini.div({
				class: 'power2',
				style: `left: ${x * 64}px; top: ${y * 64}px`,
				x: x,
				y: y,
			})
		} else if (power === 3) {
			powerElement = mini.div({
				class: 'power3',
				style: `left: ${x * 64}px; top: ${y * 64}px`,
				x: x,
				y: y,
			})
		}
		return powerElement
	}

	function createBombElement(x, y) {
		return mini.div({
			class: 'bomb',
			style: `left: ${x * 64}px; top: ${y * 64}px`,
			x: x,
			y: y,
		})
	}

	function createExplosionElement(x, y) {
		return mini.div({
			class: 'explosion',
			style: `left: ${x * 64}px; top: ${y * 64}px`,
			x: x,
			y: y,
		})
	}

	function createFlameElement(x, y) {
		return mini.div({
			class: 'explosion',
			style: `left: ${x * 64}px; top: ${y * 64}px`,
			x: x,
			y: y,
		})
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
		case 'status':
			playerCount = data.player_count
			break
		case 'server_full':
			app.innerHTML = data.message
			break
		case 'new_game':
			mapState.value = data.map
			playerState.value = data.players
			let sending = true

			document.addEventListener('keydown', keyDownHandler)
			document.addEventListener('keyup', keyUpHandler)

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
	if (!'wsad '.includes(keyDown)) return
	activeKeys.add(keyDown)
}

function keyUpHandler(e) {
	const keyUp = e.key.toLowerCase()
	if (!'wsad '.includes(keyUp)) return
	activeKeys.delete(keyUp)
}
