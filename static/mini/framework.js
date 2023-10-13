import { isValidHTMLElement } from './utils.js'
import * as el from './elements.js'
import Router from './router.js'

const router = Router

export function render(target, component) {
	if (!isValidHTMLElement(target)) throw new Error(`${target} is not a valid HTML element`)
	target.innerHTML = ''

	let elements
	if (typeof component === 'function') {
		elements = component()
	} else if (isValidHTMLElement(component)) {
		elements = component
	} else {
		throw new Error('Component must be either a function or a valid HTML element')
	}

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

	if (tag === 'a' && attrs['data-use-router']) {
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
	const keyMap = new Map()

	state.subscribe(() => {
		const newElement = getter()
		if (!newElement || !newElement.children) return

		const newChildren = Array.from(newElement.children)
		const newKeyMap = new Map()

		newChildren.forEach((child) => {
			const key = keyFn(child)
			const existingChild = keyMap.get(key)
			if (existingChild) {
				existingChild.checked = child.checked
				existingChild.classList = child.classList
			} else {
				newKeyMap.set(key, child)
			}
		})

		element.replaceWith(newElement)
		element = newElement
		keyMap.clear()

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

function createApp(element_id) {
	return document.getElementById(element_id)
}

const mini = {
	render,
	createApp,
	createElement,
	createState,
	router,
	bindToDOM,
	...el,
}
export default mini
