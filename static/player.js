import mini from './mini/framework.js'
import { socket } from './app.js'
let frameIndex = 0

const Player = (id) => {
	if (!id || id < 1 || id > 4) throw new Error(`Invalid player id: ${id}`)

	const startingPosition = {
		1: { x: 0, y: 0 },
		2: { x: 14, y: 12 },
		3: { x: 0, y: 12 },
		4: { x: 14, y: 0 },
	}

	const sprite = mini.div({
		id: `player${id}`,
		'data-player-id': id,
		style: `left: ${startingPosition[id].x * 64}px; top: ${startingPosition[id].y * 64}px`,
	})

	let lastUpdateTime = 0
	const animationDelay = 200 // 200ms delay between frames
	const updateSprite = (direction) => {
		const currentTime = performance.now()
		if (currentTime - lastUpdateTime < animationDelay) return
		lastUpdateTime = currentTime
		let rowIndex
		switch (direction) {
			case 'up':
				rowIndex = 4
				break
			case 'down':
				rowIndex = 5
				break
			case 'left':
				rowIndex = 0
				break
			case 'right':
				rowIndex = 1
				break
			default:
				rowIndex = 0
		}

		const frameCount = 4
		const frameWidth = 64
		const frameHeight = 64
		const xOffset = frameIndex * -frameWidth
		const yOffset = rowIndex * -frameHeight

		sprite.style.backgroundPosition = `${xOffset}px ${yOffset}px`
		frameIndex = (frameIndex + 1) % frameCount
	}
	const keyToDirection = {
		'w': 'up',
		's': 'down',
		'a': 'left',
		'd': 'right',
	}
	function keyDownHandler(e) {
		if (!'wsad'.includes(e.key)) return
		const direction = keyToDirection[e.key]
		const movement = { Direction: direction }
		socket.send(JSON.stringify({ type: 'keydown', movement }))
	}

	function keyUpHandler(e) {
		if (!'wsad'.includes(e.key)) return
		const direction = keyToDirection[e.key]
		const movement = { Direction: direction }

		socket.send(JSON.stringify({ type: 'keyup', movement }))
		console.log(movement)
	}
	const position = mini.createState(startingPosition[id])

	document.addEventListener('keydown', keyDownHandler)
	document.addEventListener('keyup', keyUpHandler)

	return {
		id,
		getSprite: () => sprite,
		updateSprite,
		position,
		detach: () => document.removeEventListener('keydown', keyDownHandler),
	}
}

export default Player
