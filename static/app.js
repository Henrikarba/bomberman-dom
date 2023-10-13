import mini from './mini/framework.js'
import Player from './player.js'

// App
const app = mini.createApp('app')
export const socket = new WebSocket('ws://localhost:5000/ws')

// Views
let currentView = 1
const drawStartMenu = null
const views = [drawStartMenu, drawGameboard]

// Game
const mapState = mini.createState([])

// Players
const player = Player(1)
const sprite = player.getSprite()

socket.onmessage = (e) => {
	const data = JSON.parse(e.data)
	console.log(data)
	if (Array.isArray(data)) {
		mapState.value = data
	} else if (data?.type && data.type == 'player_position_update') {
		player.position.value.x = data.players[0].x
		player.position.value.y = data.players[0].y
		player.updateSprite(data.players[0].direction)
		requestAnimationFrame(gameloop)
	}
}

mapState.subscribe(() => {
	const updatedBoard = drawGameboard()
	mini.render(app, updatedBoard)
})

requestAnimationFrame(gameloop)
function gameloop() {
	sprite.style.left = player.position.value.x + 'px'
	sprite.style.top = player.position.value.y + 'px'
}

function drawGameboard() {
	return mini.div(
		{ id: 'game' },
		...mapState.value.map((row, rowIndex) => {
			return mini.div(
				{ class: 'row flex' },
				...row.map((cell, cellIndex) => {
					let cellClass = 'cell'

					if (cell === '0') {
						cellClass = 'wall'
					} else if (cell === 'd') {
						cellClass = 'destroyable'
					}

					return mini.div({ class: cellClass, 'data-row': rowIndex, 'data-cell': cellIndex })
				})
			)
		}),
		sprite
	)
}
