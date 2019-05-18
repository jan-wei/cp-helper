(function() {
	const loadUrl = "https://codeforces.com/problemset/problem/+";

	if(window.location.href.match(loadUrl)) {
		// get the problem number
		let urlSplit = window.location.href.split('/');
		let number = "cf" + urlSplit[urlSplit.length - 2] + urlSplit[urlSplit.length - 1];
		
		// get the proble title
		let titleElement = document.querySelector(".title");
		// do not take in the number, just the title
		let title = titleElement.innerText.substring(3);
		
		// get the test cases
		let inputElements = document.getElementsByClassName("input");
		let outputElements = document.getElementsByClassName("output");
		
		if(inputElements.length != outputElements.length) {
			console.error("Input elements' length does not match output elements' length");
			return;
		}
		
		let numTestCases = inputElements.length;
		
		let testCases = [];
		for(let i = 0; i < numTestCases; i++) {
			testCases.push({
				input: inputElements[i].getElementsByTagName("pre")[0].innerText,
				output: outputElements[i].getElementsByTagName("pre")[0].innerText
			});
		}
		
		// zip into json object
		let problem = JSON.stringify({
			number: number,
			title: title,
			link: window.location.href,
			numTestCases: numTestCases,
			testCases: testCases
		});
		
		// post object to local server
		const url = "http://localhost:8080/problem/new"
		
		let xhr = new XMLHttpRequest();
		xhr.open("POST", url, true);
		
		// request body is serialized as json
		xhr.setRequestHeader("Content-Type", "application/json");
		
		xhr.onreadystatechange = function() {
			if(xhr.readyState == 4 && xhr.status != 200) {
				// error
				console.error(xhr.responseText);
			}
		};
		
		xhr.send(problem);
		return;
	}
})();
	
