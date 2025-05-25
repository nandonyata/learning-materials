import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '3s', target: 2 }, // Ramp up to 5 virtual users over 10 seconds
    // { duration: '10s', target: 5 }, // Maintain 5 virtual users for 10 seconds
    // { duration: '10s', target: 0 }, // Ramp down to 0 virtual users over 10 seconds
  ],
};

// N: Number of requests per second per IP
const N = 6;

export default function () {
  const ips = [
    '192.168.1.1',
    '192.168.1.2',
    '10.0.0.1',
    '10.0.0.2',
    '172.16.0.1',
  ];

  // Select a random IP from the array of IPs
  const randomIP = ips[Math.floor(Math.random() * ips.length)];

  // Set the 'X-Forwarded-For' header to simulate a request from that IP
  const headers = {
    'X-Forwarded-For': randomIP,
  };

  // Send N requests per second
  for (let i = 0; i < N; i++) {
    const response = http.get('http://localhost:8080/ping', { headers });

    // Check if the response status is 200 OK
    check(response, {
      'is status 200': (r) => r.status === 200,
    });

    // Sleep to maintain N requests per second
    sleep(1 / N); // This will ensure that we send N requests in one second
  }
}
