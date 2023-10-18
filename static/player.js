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
		switch (direction) {
			case 'up':
				frameIndex = 2
				break
			case 'down':
				frameIndex = 1
				break
			case 'left':
				frameIndex = 0
				break
			case 'right':
				frameIndex = 3
				break
			default:
				frameIndex = 1
		}
		const frameWidth = 64
		const xOffset = frameIndex * -frameWidth

		sprite.style.backgroundPosition = `${xOffset}px 0`
	}

	const position = mini.createState({})

	return {
		getSprite: () => sprite,
		updateSprite,
		position,
	}
}

export default Player
