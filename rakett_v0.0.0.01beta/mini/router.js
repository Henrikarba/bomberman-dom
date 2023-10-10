import { isValidHTMLElement } from './utils.js'
import { render } from './framework.js'

const Router = (container) => {
	if (!isValidHTMLElement(container)) throw new Error(`${container} is not a valid HTML element`)

	const routes = {}
	let notFound = null

	const registerRoute = (path, handler) => {
		routes[path] = handler
	}

	const notFoundHandler = (handler) => {
		notFound = handler
	}

	const navigateTo = (e) => {
		const path = window.location.pathname
		const handler = routes[path] || notFound
		if (handler) render(container, handler)
		else console.error(`handler for ${path} not found. consider registering notFoundHandler`)
	}

	window.addEventListener('popstate', () => navigateTo())
	window.addEventListener('navigate', (e) => {
		console.log(e.target.attr)
		navigateTo(e)
	})

	return { registerRoute, notFoundHandler, navigateTo }
}

export default Router
