import { isValidHTMLElement } from './utils.js'
import * as el from './elements.js'
import Router from './router.js'

const router = Router

export function render(target, component) {
	if (!isValidHTMLElement(target)) throw new Error(`${target} is not a valid HTML element`)
	target.innerHTML = ''

	const elements = component()
	if (Array.isArray(elements)) {
		elements.forEach((element) => {
			if (isValidHTMLElement(element)) {
				target.appendChild(element)
			}
		})
	} else if (isValidHTMLElement(elements)) {
		target.appendChild(elements)
	}
}

export const createElement = (tag, attrs = {}, ...children) => {
	const el = document.createElement(tag)
	for (const [key, value] of Object.entries(attrs)) {
		if (key.startsWith('on')) {
			const eventName = key.slice(2).toLowerCase()
			el.addEventListener(eventName, value)
		} else el.setAttribute(key, value)
	}

	if (tag === 'a' && attrs.href) {
		el.addEventListener('click', (e) => {
			e.preventDefault()
			const path = attrs.href
			history.pushState(null, '', path)
			const navigateEvent = new Event('navigate')
			window.dispatchEvent(navigateEvent)
		})
	}

	children.forEach((child) => {
		if (typeof child === 'string') {
			el.appendChild(document.createTextNode(child))
		} else if (child instanceof Node) {
			el.appendChild(child)
		}
	})
	return el
}

const bindToDOM = (getter, state, keyFn) => {
	let element = getter()
	if (!element) {
		element = document.createComment('')
	}

	// Store a map of keys to child elements
	const keyMap = new Map()

	state.subscribe(() => {
		const newElement = getter()
		const newChildren = Array.from(newElement.children)
		const newKeyMap = new Map()

		// Update existing children or add new ones
		newChildren.forEach((child) => {
			const key = keyFn(child)
			const existingChild = keyMap.get(key)
			if (existingChild) {
				// Update the existing DOM element's properties
				// (You can extend this to update other properties as needed)
				existingChild.checked = child.checked
				existingChild.classList = child.classList
			} else {
				newKeyMap.set(key, child)
			}
		})

		// Replace the old element with the new one
		element.replaceWith(newElement)
		element = newElement
		keyMap.clear()

		// Update the key map
		for (const [key, child] of newKeyMap) {
			keyMap.set(key, child)
		}
	})

	return element
}

const createState = (initialValue) => {
	let value = initialValue
	let previousValue = null
	const listeners = []

	return {
		get value() {
			return value
		},
		set value(newValue) {
			previousValue = value
			value = newValue
			listeners.forEach((listener) => listener(value, previousValue))
		},
		subscribe(listener) {
			listeners.push(listener)
			return () => {
				const index = listeners.indexOf(listener)
				listeners.splice(index, 1)
			}
		},
		get previousValue() {
			return previousValue
		},
	}
}

const mini = {
	render,
	createElement,
	createState,
	router,
	bindToDOM,
	...el,
}
export default mini
