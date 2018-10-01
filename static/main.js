function f (button) {
  fetch('/inc/')
    .then((response) => {
      return response.json()
    }).then((json) => {
      console.log('parsed json', json)
      button.innerHTML = json.toString()
    }).catch((ex) => {
      button.innerHTML = 'Total fail'
      console.log('parsing failed', ex)
    })
}
