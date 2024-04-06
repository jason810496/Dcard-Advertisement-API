import { randomIntBetween, randomItem } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import http from 'k6/http';
import exec from "k6/execution";
import { sleep } from 'k6';

export const options = {
	scenarios: {
		constant_request_rate: {
			executor: 'constant-arrival-rate',
			rate: 20000,
			timeUnit: '1s',
			duration: __ENV.DURATION,
			preAllocatedVUs: 100,
			maxVUs: 10000,
		},
	},
};

// export const options = {
//     vus: 500, // users
//     duration: "10s",
//     // rps: 10500,
// };

let urlString = `http://${__ENV.API_HOST}:${__ENV.API_PORT}/api/v1/ad`;

const limits = [10, 20, 30, 40, 50];
const offsets = [0, 10, 20, 30, 40, 50];
const genders = ["F", "M"];
const countries = ["HK", "JP", "US", "KR"]; // TW
const platforms = ["ios", "android"]; // "web"

export default function () {
	const limit = randomItem(limits);
	const offset = randomItem(offsets);
	const gender = randomItem(genders);

	const countryRand = randomIntBetween(1, 10);
	const platformRand = randomIntBetween(1, 10);

	let age = -1;
	let country;
	let platform;

	// simulate real-world traffic age distribution
	if (exec.vu.idInTest % 10 < 1) { // 10%
		age = randomIntBetween(1, parseInt(__ENV.AGE_END)-1);
	}
	else if (exec.vu.idInTest % 10 < 9) { // 80%
		age = randomIntBetween(parseInt(__ENV.AGE_START), parseInt(__ENV.AGE_END));
	}
	else{ // 10%
		age = randomIntBetween(parseInt(__ENV.AGE_START)+1, 100);
	}

	// simulate real-world traffic country distribution
	if (countryRand < 8) { // 80%
		country = "TW";
	}
	else { // 20%
		country = randomItem(countries);
	}

	// simulate real-world traffic platform distribution
	if (platformRand < 8) { // 80%
		platform = randomItem(platforms);
	}
	else { // 20%
		platform = "web";
	}

	http.get(
		`${urlString}?limit=${limit}&offset=${offset}&age=${age}&gender=${gender}&country=${country}&platform=${platform}`,
		{
			tags: {
				name: 'GetAd'
			}
		}
	);

	sleep(1);
}