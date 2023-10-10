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

2. **Creating elements**

The Mini framework provides a set of predefined elements that you can use right out of the box. These elements are essentially wrappers around the createElement function and are available as mini.<elementName>. For example, to create a section element, you can use mini.section().

```javascript
const mySection = mini.section({ class: 'my-class' }, 'This is a section')
```

Custom Elements

If the element you need is not available as a predefined function, you can use the createElement function to create your own. The createElement function takes three arguments:

    tag: The HTML tag name as a string (e.g., 'div', 'a', 'span').
    attrs: An object containing any attributes you want to set on the element. This can include event listeners, which should be prefixed with 'on' (e.g., 'onclick', 'oninput').
    children: Any child elements or text content. This can be a single value or an array.

````javascript
export const createElement = (tag, attrs = {}, ...children) => { /* ... */ }```
````
