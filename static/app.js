import mini from './mini/framework.js'
import { mapStateUpdate, drawGameboard } from './mapupdates.js'
import { drawStartMenu } from './start.js'
import { socket } from './socket.js'

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
export const chatArea = mini.div({ class: 'chat-area' })
const sendMessageButton = mini.button(
	{
		id: 'send-message-button',
		type: 'submit',
		disabled: true,
	},
	'SUBMIT'
)
const chat = drawchat()
export const overlay = mini.div({ class: 'overlay' }, mini.h2({}, 'SPECTATING'))

export let gameboard = undefined
let didGameStart = false

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
	didGameStart = false
	div.appendChild(info)
	div.appendChild(
		mini.div(
			{ id: 'lobby' },
			mini.div({}, mini.div({}, 'Players In Lobby'), mini.div({ id: 'lobby-player-count' })),
			mini.div({ id: 'lobby-timer' })
		)
	)

	display.appendChild(div)
	display.appendChild(chat)
	mini.render(app, display)
}

function newGame() {
	didGameStart = true
	if (!gameboard) gameboard = drawGameboard(mapState.value)
	if (overlay) overlay.remove()
	div.innerHTML = ''
	div.appendChild(info)
	div.appendChild(gameboard)

	updatePlayerPosition(gameboard)
}

function updatePlayerPosition(gameboard) {
	playerState.value.forEach((player) => {
		if (didGameStart) {
			if (playerID.value == player.id) {
				if (player.lives <= 0) {
					let removePlayer = document.getElementById(`player${player.id}`)
					removePlayer.style.display = 'none'
					info.innerHTML = 'DEAD'
				} else {
					let hearts = ``
					for (let i = 1; i <= player.lives; i++) {
						hearts += `<div class="heart"></div>`
					}
					info.innerHTML = 'Lives: ' + hearts
				}
			}

			let playerElement = playerElements[player.id]

			if (player.lives <= 0) {
				let removePlayer = document.getElementById(`player${player.id}`)
				removePlayer.style.display = 'none'
			}
			if (!playerElement && player.lives > 0) {
				playerElement = Player(player)
				playerElements[player.id] = playerElement
				gameboard.appendChild(playerElement.getSprite())
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
		}
	})
}

function drawchat() {
	let inputLength = 0

	const handleInputChange = (e) => {
		inputLength = e.target.value.length

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
		chatArea,
		mini.form(
			{
				id: 'messageInput',
				onsubmit: (e) => {
					e.preventDefault()
					const sendMessage = {
						type: 'message',
						message: input.value,
					}
					socket.send(JSON.stringify(sendMessage))
					input.value = ''
					sendMessageButton.disabled = true
				},
			},
			input,
			sendMessageButton
		)
	)
}
