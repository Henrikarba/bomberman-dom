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

const deepCopyAttributes = (oldElem, newElem) => {
	if (!oldElem || !newElem) return

	if (oldElem.className) {
		newElem.className = oldElem.className
	}
	for (const attr of oldElem.attributes) {
		newElem.setAttribute(attr.name, attr.value)
	}
	if (oldElem.type === 'checkbox') {
		newElem.checked = oldElem.checked
	}
	const minChildrenLength = Math.min(oldElem.children.length, newElem.children.length)
	for (let i = 0; i < minChildrenLength; i++) {
		deepCopyAttributes(oldElem.children[i], newElem.children[i])
	}
}

const bindToDOM = (getter, state) => {
	let element = getter()
	if (!element) {
		element = document.createComment('')
	}
	state.subscribe(() => {
		const newElement = getter()
		console.log(element, newElement)
		deepCopyAttributes(element, newElement)
		element.replaceWith(newElement)
		element = newElement
	})
	return element
}
const createState = (initialValue) => {
	let value = initialValue
	const listeners = []

	return {
		get value() {
			return value
		},
		set value(newValue) {
			value = newValue
			listeners.forEach((listener) => listener(value))
		},
		subscribe(listener) {
			listeners.push(listener)
			return () => {
				const index = listeners.indexOf(listener)
				listeners.splice(index, 1)
			}
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
