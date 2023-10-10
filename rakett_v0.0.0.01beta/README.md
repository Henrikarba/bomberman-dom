# Getting started

## 0. **TodoMVC**

Open index.html to see [TodoMVC](https://todomvc.com/) built in mini-framework

## 1. **Basic starter example**

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

## 2. **Creating elements**

The Mini framework provides a set of predefined elements that you can use right out of the box. These elements are essentially wrappers around the createElement function and are available as mini.ELEMENT_NAME(). For example, to create a section element, you can use mini.section().

```javascript
const mySection = mini.section({ class: 'my-class' }, 'This is a section')
```

#### Custom Elements

If the element you need is not available as a predefined function, you can use the createElement function to create your own. The createElement function takes three arguments:

    tag: The HTML tag name as a string (e.g., 'div', 'a', 'span').
    attrs: An object containing any attributes you want to set on the element. This can include event listeners, which should be prefixed with 'on' (e.g., 'onclick', 'oninput').
    children: Any child elements or text content. This can be a single value or an array.

````javascript
export const createElement = (tag, attrs = {}, ...children) => { /* ... */ }```
````

## 3. **Special Attributes**

The createElement function also supports special attributes. For example, if you want to make an anchor tag work with the router, you can add a data-use-router attribute:

```javascript
const myLink = createElement('a', { href: '/home', 'data-use-router': true }, 'Go to Home')
// or
const myLink2 = mini.a({ href: '/home', 'data-use-router': true }, 'Go to home')
```

## 4. **State**

To create a state variable, use the mini.createState function and pass the initial value as an argument.

````javascript
const todos = mini.createState([
  { id: 0, name: 'Eat bananas', checked: false },
  { id: 1, name: 'Eat apples', checked: false },
  { id: 2, name: 'Finish this task', checked: false },
]);```
````

#### Accessing state

```javascript
console.log(todos.value) // Outputs the current value of todos
```

#### Updating state

```javascript
todos.value = [{ id: 3, name: 'New task', checked: false }]
```

#### Subscribing to state changes

You can subscribe to changes in a state variable using the .subscribe method. This method takes a function that will be called whenever the state changes.

```javascript
const completedCount = mini.createState(0)
todos.subscribe(updateCompletedCount)

function updateCompletedCount() {
	completedCount.value = todos.value.filter((todo) => !todo.checked).length
}
```

#### Using State in Components

```javascript
const todosView = () => {
	return mini.ul(
		{ class: 'todo-list' },
		...todos.value.map((todo, index) => {
			// Your component logic here
		})
	)
}
```

#### Binding to DOM

**Syntax**

```javascript
mini.bindToDOM(viewFunction, state, keyFunction)

//    viewFunction: A function that returns the HTML structure.
//    state: The state object that the view is dependent on.
//    keyFunction: A function to uniquely identify DOM elements.
```

**Examples**

```javascript
const main = () => {
	return mini.section(
		{ class: 'main' },
		mini.input({ id: 'toggle-all', class: 'toggle-all', type: 'checkbox' }),
		mini.label({ for: 'toggle-all', onclick: toggleAllTasks }, 'Mark all as complete'),
		mini.bindToDOM(todosView, todos, keyFn)
	)
}
```

```javascript
mini.bindToDOM(
	() => {
		clearBtn.style.display = todos.value.length !== 0 ? 'block' : 'none'
		return clearBtn
	},
	todos,
	() => 'clearButtonKey'
)
```
