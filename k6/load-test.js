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
			duration: '5m',
			preAllocatedVUs: 1000,
			maxVUs: 4000,
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
const genders = ["F", "M"];
const countries = ["TW", "HK", "JP", "US", "KR"];
const platforms = ["ios", "android", "web"];

export default function () {
	const limit = randomItem(limits);
	let age = -1;
	const gender = randomItem(genders);
	const country = randomItem(countries);
	const platform = randomItem(platforms);

	if (exec.vu.idInTest % 10 < 1) {
		// 10%
		age = randomIntBetween(1, parseInt(__ENV.AGE_END));
	}
	else if (exec.vu.idInTest % 10 < 8) {
		// 80%
		age = randomIntBetween(parseInt(__ENV.AGE_START), parseInt(__ENV.AGE_END));
	}
	else{
		// 10%
		age = randomIntBetween(parseInt(__ENV.AGE_START), 100);
	}

	for (let i = 0; i < 10; i++) {
		http.get(
			`${urlString}?limit=${limit}&offset=${i}&age=${age}&gender=${gender}&country=${country}&platform=${platform}`,
			{
				tags: {
					name: 'GetAd'
				}
			}
		);
	}

	sleep(1);
}