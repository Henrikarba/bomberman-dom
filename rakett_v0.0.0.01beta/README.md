## Getting started

1. **Basic starter example**

   ```javascript
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
   			mini.a({ href: '/', 'data-use-router': true }, 'To home')
   		),
   		mini.footer({}, mini.p({ class: 'text-lg' }, 'paragraph tag with class text-lg inside footer tag')),
   	]
   }

   // Register handler
   router.registerRoute('/', Home)

   // Initializes the router, add it in the end
   router.navigateTo()
   ```

   This will result in the following HTML:

   ```html
   <body id="app">
   	<section class="main">
   		<a href="/" data-use-router="true">My link</a>
   	</section>
   	<footer>
   		<p class="text-lg">paragraph tag with class text-lg inside footer tag</p>
   	</footer>
   </body>
   ```

2.
