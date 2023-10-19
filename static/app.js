import mini from './mini/framework.js'
import { mapStateUpdate, drawGameboard } from './mapupdates.js'
import { drawStartMenu } from './start.js'

import Player from './player.js'

// App
const app = mini.createApp('app')

// Views
const startMenu = drawStartMenu()
mini.render(app, startMenu)

// Game
export const mapState = mini.createState([])
export const playerState = mini.createState([])
export const blockUpdates = mini.createState([])
export const playerCount = mini.createState(0)
export const playerID = mini.createState()

let gameboard = undefined
let info = undefined
let chat = undefined
const playerElements = {}

export function gameloop(updateType) {
	console.log(playerID.value)
	const gameLoopTypes = {
		'status': () => console.log(playerCount.value),
		'new_game': newGame,
		'player_state_update': () => updatePlayerPosition(gameboard),
		'map_state_update': () => mapStateUpdate(gameboard, blockUpdates),
		'end_game': () => {
			return
		},
	}

	const action = gameLoopTypes[updateType]
	if (action) action()
}

function newGame() {
	if (!gameboard) gameboard = drawGameboard(mapState.value)
	const display = mini.div({ class: 'display' })
	const game = mini.div({ class: 'game' })
	info = mini.div({ class: 'info' })
	game.appendChild(info)
	game.appendChild(gameboard)
	display.appendChild(game)
	chat = mini.div({ class: 'chat' })
	display.appendChild(chat)
	updatePlayerPosition(gameboard)
	mini.render(app, display)
}

function updatePlayerPosition(gameboard) {
	playerState.value.forEach((player) => {
		if (playerID == player.id) {
			let hearts = ``
			for (let i = 1; i <= player.lives; i++) {
				hearts += `<div class="heart"></div>`
			}
			info.innerHTML = 'Lives: ' + hearts
		}

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
