import mini from './mini/framework.js'

export function drawGameboard(mapdata) {
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

export function mapStateUpdate(gameboard, blockUpdates) {
	const blockActionMap = {
		'B': createBombElement,
		'ex': createExplosionElement,
		'f': createFlameElement,
		'p1': (x, y) => createPowerElement(x, y, 1),
		'p2': (x, y) => createPowerElement(x, y, 2),
		'p3': (x, y) => createPowerElement(x, y, 3),
	}

	blockUpdates.value.forEach((update) => {
		const { block, x, y } = update

		if (blockActionMap[block]) {
			const elem = blockActionMap[block](x, y)
			gameboard.appendChild(elem)
		}

		if (block === 'ex') {
			removeElementsByClassAndCoordinates('bomb', x, y)
		}

		if (block === 'e') {
			const changeCell = document.querySelector(`.destroyable[data-row="${y}"][data-cell="${x}"]`)
			if (changeCell) {
				changeCell.className = 'cell'
			}
			const elems = ['power1', 'power2', 'power3', 'explosion']
			elems.forEach((className) => {
				removeElementsByClassAndCoordinates(className, x, y)
			})
		}
	})
}

function createPowerElement(x, y, power) {
	const changeCell = document.querySelector(`.destroyable[data-row="${y}"][data-cell="${x}"]`)
	if (changeCell) {
		changeCell.className = 'cell'
	}

	const powerMap = {
		1: 'power1',
		2: 'power2',
		3: 'power3',
	}

	const powerClass = powerMap[power]
	if (powerClass) {
		return mini.div({
			class: powerClass,
			style: `left: ${x * 64}px; top: ${y * 64}px`,
			x: x,
			y: y,
		})
	}
	return null
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

function removeElementsByClassAndCoordinates(className, x, y) {
	const elem = document.querySelector(`.${className}[x="${x}"][y="${y}"]`)
	if (elem) {
		elem.remove()
	}
}
