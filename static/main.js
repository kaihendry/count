function f (button) {
  fetch('/inc/')
  .then((response) => {
  return response.json()
  }).then((json) => {
  button.innerHTML = json.count
  console.log('parsed json', json)
  }).catch((ex) => {
  button.innerHTML = 'Total fail'
  console.log('parsing failed', ex)
})
}
