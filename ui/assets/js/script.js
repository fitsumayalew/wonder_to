const phone = document.getElementById('phone')
const password = document.getElementById('password')
const form = document.getElementById('form')
const email = document.getElementById('email')
const errorElement = document.getElementById('err')
form.addEventListener('submit', (e) => {
    let messages = []

    if (messages.length > 0) {
        e.preventDefault()
        errorElement.innerText = messages.join(',')

    }
    if (password.value.length > 0) {
        alert('Password must be longer than 8 characters ')
    }
})