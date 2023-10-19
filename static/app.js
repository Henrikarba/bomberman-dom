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

//
const div = mini.div({})
const display = mini.div({ class: 'display' })

const info = mini.div({ class: 'info' })
const chat = drawchat()
const overlay = mini.div({ class: 'overlay' }, mini.h2({}, 'SPECTATING'))

let gameboard = undefined

const playerElements = {}

export function gameloop(updateType) {
	const gameLoopUpdates = {
		'status': () => console.log(playerCount.value),
		'playerID': newLobby,
		'new_game': newGame,
		'player_state_update': () => updatePlayerPosition(gameboard),
		'map_state_update': () => mapStateUpdate(gameboard, blockUpdates),
		'game_over': () => {
			gameboard.appendChild(overlay)
		},
	}

	const action = gameLoopUpdates[updateType]
	if (action) action()
}

function newLobby() {
	div.appendChild(info)
	div.appendChild(mini.div({ id: 'game' }))

	display.appendChild(div)
	display.appendChild(chat)
	mini.render(app, display)
}

function newGame() {
	if (!gameboard) gameboard = drawGameboard(mapState.value)
	if (overlay) overlay.remove()
	div.innerHTML = ''
	div.appendChild(info)
	div.appendChild(gameboard)

	updatePlayerPosition(gameboard)
}

function updatePlayerPosition(gameboard) {
	playerState.value.forEach((player) => {
		if (playerID.value == player.id) {
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
			info.innerHTML = 'DEAD'
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

function drawchat() {
	let inputLength = 0

	const handleInputChange = (e) => {
		inputLength = e.target.value.length
		const sendMessageButton = document.getElementById('send-message-button')
		if (inputLength > 0) {
			sendMessageButton.disabled = false
		} else {
			sendMessageButton.disabled = true
		}
	}
	const input = mini.input({
		id: 'message',
		oninput: handleInputChange,
	})

	return mini.div(
		{ class: 'chat' },
		mini.form(
			{ id: 'messageInput' },
			{
				style: 'display: flex; flex-direction: row;',
				onsubmit: (e) => {
					e.preventDefault()
					const sendMessage = {
						type: 'message',
						name: input.value,
					}
					socket.send(JSON.stringify(sendMessage))
					input.value = ''
				},
			},
			input,
			mini.button(
				{
					id: 'send-message-button',
					type: 'submit',
					disabled: true,
				},
				'SUBMIT'
			)
		)
	)
}
