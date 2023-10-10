// Import the framework
import mini from './mini/framework.js'

// Select your apps main container
// This is where your app will live
// Router will use this element to change content based on url
const container = document.getElementById('app')

// Create router
const router = mini.router(container)

// Create home view
const Home = () => {
	return [
		mini.section(
			{ class: 'main' },
			// Add attr data-use-router: true for router, other anchor tags will have default behavior
			mini.a({ href: '/about', 'data-use-router': true }, 'To about')
		),
		mini.footer({}, mini.p({ class: 'text-lg' }, 'paragraph tag with class text-lg inside footer tag')),
	]
}

const About = () => {
	return mini.section(
		{ class: 'main' },
		// Add attr data-use-router: true for router, other anchor tags will have default behavior
		mini.a({ href: '/', 'data-use-router': true }, 'To home'),

		mini.p({}, 'This is about page!'),
		mini.a({ href: 'http://google.com' }, 'To google')
	)
}

// Register handler
router.registerRoute('/', Home)
router.registerRoute('/about', About)

// Initializes the router, add it in the end
router.navigateTo()
