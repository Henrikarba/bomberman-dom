import mini from './mini/framework.js'

const container = document.getElementById('app')

const todos = mini.createState([
	{ name: 'Eat bananas', checked: false },
	{ name: 'Eat apples', checked: false },
	{ name: 'Finish this task', checked: false },
])
let completedCount = mini.createState(todos.value.filter((item) => item.checked == false).length || 0)
let newTodo = ''

const header = () => {
	return mini.section(
		{ class: 'todoapp' },
		mini.header(
			{ class: 'header' },
			mini.h1({}, 'todos'),
			mini.input({
				class: 'new-todo',
				placeholder: 'What needs to be done?',
				autofocus: true,
				oninput: (e) => {
					newTodo = e.target.value
				},
				onkeydown: (e) => {
					if (e.key === 'Enter') {
						if (newTodo.trim() == '') return
						completedCount.value += 1
						const newTodos = [...todos.value, { name: newTodo.trim(), checked: false }]
						todos.value = newTodos
						newTodo = ''
						e.target.value = ''
					}
				},
			}),
			main(),
			footer()
		)
	)
}

const main = () => {
	return mini.section(
		{ class: 'main' },
		mini.input({ id: 'toggle-all', class: 'toggle-all', type: 'checkbox' }),
		mini.label({ for: 'toggle-all', onclick: toggleAllTasks }, 'Mark all as complete'),
		mini.bindToDOM(todosView, todos)
	)
}

const todosView = () => {
	return mini.ul(
		{ class: 'todo-list' },
		...todos.value.map((todo, index) =>
			mini.li(
				{},
				mini.div(
					{ class: 'view' },
					mini.input({ class: 'toggle', type: 'checkbox', onclick: (e) => toggleCompleted(e) }),
					mini.label({ ondblclick: (event) => handleEdit(event, index) }, todo.name),
					mini.button({ class: 'destroy', onclick: () => destroyTodo(index) })
				)
			)
		)
	)
}

const footer = () => {
	return mini.footer(
		{ class: 'footer' },
		mini.span({ class: 'todo-count' }, mini.bindToDOM(counter, completedCount), ' items left')
	)
}

const counter = () => {
	return mini.strong({}, String(completedCount.value))
}

const edit = (value) => {
	return mini.input({ id: 'edit', class: 'edit', value: value })
}

const ToDoApp = () => header()

function handleEdit(event, index) {
	let editorRemoved = false

	const liElement = event.target.closest('li')
	liElement.classList.add('editing')

	const editor = edit(todos.value[index].name)
	liElement.append(editor)
	editor.focus()

	const onBlur = () => {
		if (editorRemoved) return
		const newValue = editor.value.trim()
		if (newValue) {
			const newTodos = [...todos.value]
			newTodos[index].name = newValue
			todos.value = newTodos
		}
		liElement.classList.remove('editing')
		editorRemoved = true
		editor.removeEventListener('blur', onBlur)
		editor.removeEventListener('keydown', onKeyDown)
		editor.remove()
	}

	const onKeyDown = (e) => {
		if (editorRemoved) return
		if (e.key === 'Enter' || e.key === 'Escape') {
			editorRemoved = true
			const newValue = editor.value.trim()
			if (newValue) {
				const newTodos = [...todos.value]
				newTodos[index].name = newValue
				todos.value = newTodos
			}
			liElement.classList.remove('editing')
			editor.removeEventListener('blur', onBlur)
			editor.removeEventListener('keydown', onKeyDown)
			editor.remove()
		}
	}

	editor.addEventListener('blur', onBlur)
	editor.addEventListener('keydown', onKeyDown)
}

function toggleAllTasks() {
	const liElements = document.querySelectorAll('.todo-list li')
	const allChecked = Array.from(liElements).every((li) => li.querySelector('.toggle').checked)

	if (allChecked) {
		liElements.forEach((li) => {
			const checkbox = li.querySelector('.toggle')
			checkbox.checked = false
			li.classList.remove('completed')
		})
		completedCount.value = todos.value.length
	} else {
		liElements.forEach((li) => {
			const checkbox = li.querySelector('.toggle')
			checkbox.checked = true
			li.classList.add('completed')
		})
		completedCount.value = 0
	}
}

function destroyTodo(index) {
	const newTodos = [...todos.value]
	const isCompleted = newTodos[index].checked // Assuming each todo object has a 'checked' property

	newTodos.splice(index, 1)
	todos.value = newTodos

	if (isCompleted) {
		completedCount.value -= 1
	}
}

function toggleCompleted(e) {
	const liElement = e.target.closest('li')
	if (liElement) {
		if (e.target.checked) {
			liElement.classList.add('completed')
			completedCount.value -= 1
		} else {
			liElement.classList.remove('completed')
			completedCount.value += 1
		}
	}
}

const router = mini.router(container)
router.registerRoute('/', ToDoApp)
router.navigateTo()
