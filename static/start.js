import mini from './mini/framework.js'
import { socket } from './socket.js'

export function drawStartMenu() {
	let inputLength = 0

	const handleInputChange = (e) => {
		inputLength = e.target.value.length
		const submitButton = document.getElementById('submit-button')
		if (inputLength >= 2 && inputLength <= 10) {
			submitButton.disabled = false
		} else {
			submitButton.disabled = true
		}
	}
	const input = mini.input({
		id: 'name',
		maxLength: '10',
		minLength: '2',
		oninput: handleInputChange,
	})

	return mini.div(
		{ id: 'menu' },
		'ENTER NAME',
		mini.form(
			{
				style: 'display: flex; flex-direction: column;',
				onsubmit: (e) => {
					e.preventDefault()
					const registerPlayer = {
						type: 'register',
						name: input.value,
					}
					socket.send(JSON.stringify(registerPlayer))
					input.value = ''
				},
			},
			input,
			mini.button(
				{
					id: 'submit-button',
					type: 'submit',
					disabled: true,
				},
				'SUBMIT'
			)
		)
	)
}
