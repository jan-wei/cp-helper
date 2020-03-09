(function() {
	let parseCodeforces = function() {
		const problemsetUrl = "https://codeforces.com/problemset/problem/+";
		const contestUrl = "https://codeforces.com/contest/+";

		let urlMatchesProblemsetUrl = window.location.href.match(problemsetUrl);
		let urlMatchesContestUrl = window.location.href.match(contestUrl);

		if(urlMatchesProblemsetUrl || urlMatchesContestUrl) {
			// get the problem number
			let urlSplit = window.location.href.split('/');
			let number = urlSplit[urlSplit.length - (urlMatchesProblemsetUrl ? 2 : 3)] + urlSplit[urlSplit.length - 1];
			
			// get the problem title
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
				platform: "codeforces",
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
	}

	let parseAtcoder = function() {
		const contestUrl = "https://atcoder.jp/contests/+";

		if(window.location.href.match(contestUrl)) {
			// get the problem number
			let urlSplit = window.location.href.split('/');
			let number = urlSplit[urlSplit.length - 1];

			// get the problem title
			let titleElement = document.getElementsByClassName("h2")[0];
			// do not take in the number, just the title
			let title = titleElement.innerText.substring(4);

			// get the test cases
			let testCases = [];
			const testID = "pre-sample";
			// 20 is a random nummber as atcoder doesnt have a predictable naming of their test cases
			for(let number = 0; number < 20; number++) {
				let target = testID + number;
				let targetElement = document.getElementById(target);
				if(!targetElement) {
					continue;
				}
				// if element exists it must come in pairs (input / output)
				let testCaseInput = targetElement.innerText;
				number++;
				target = testID + number;
				targetElement = document.getElementById(target);
				if(!targetElement) {
					console.error("Failed to parse test cases");
					return;
				}
				let testCaseOutput = targetElement.innerText;
				testCases.push({
					input: testCaseInput,
					output: testCaseOutput
				});
			}

			// zip into json object
			let problem = JSON.stringify({
				platform: "atcoder",
				number: number,
				title: title,
				link: window.location.href,
				numTestCases: testCases.length,
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
	}

	const codeforcesHostname = "codeforces.com";
	const atcoderHost = "atcoder.jp";

	switch(window.location.hostname) {
		case codeforcesHostname:
			parseCodeforces();
			break;
		case atcoderHost:
			parseAtcoder();
			break;
	}
})();