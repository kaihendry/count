function f(button) {
	fetch("/inc/", { method: "GET" })
		.then(function(response) {
			return response.json()
		}).then(function(json) {
			button.innerHTML = json.count;
			console.log('parsed json', json)
		}).catch(function(ex) {
			button.innerHTML = "Total fail";
			console.log('parsing failed', ex)
		})
}
