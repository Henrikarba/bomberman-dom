import mini from './mini/framework.js'
let frameIndex = 0

const Player = (player) => {
	if (!player.id || player.id < 1 || player.id > 4) throw new Error(`Invalid player id: ${player.id}`)

	const sprite = mini.div({
		id: `player${player.id}`,
		'data-player-id': player.id,
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

	const position = mini.createState({})

	return {
		getSprite: () => sprite,
		updateSprite,
		position,
	}
}

export default Player
