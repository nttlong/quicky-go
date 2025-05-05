// This is a sample load test script for k6.

import http from "k6/http";

export let options = {
 vus: 10, // Number of virtual users to simulate.
  duration: "1m", // How long to run the test.
  thresholds: {
    http_req_duration: ["p(95)<1000"], // 95% of requests must complete below 1000ms.
    http_reqs: ["rate<1"] // 50% of requests must be successful.
  }
};

// export let options = {
//   stages: [
//     { duration: '20s', target: 100 },   // Tăng dần lên 100 VUs trong 20s
//     { duration: '30s', target: 300 },   // Tăng tiếp lên 300 VUs trong 30s
//     { duration: '30s', target: 500 },   // Tăng đến 500 VUs trong 30s
//     { duration: '30s', target: 1000 },   // Giữ ổn định ở 500 VUs trong 30s
//     { duration: '20s', target: 0 },     // Giảm dần về 0 trong 20s
//   ]
// };

export default function() {
  let res = http.get("http://localhost:8080/api/test-004/auth/login");
  //let res = http.get("http://localhost:8080/health");
  if (res.status !== 200) {
    console.log(`Request failed with status ${res.status}`);
  }
}
//C:\Program Files\k6
//& "C:\Program Files\k6\k6" run e:\Docker\go\quicky-go\k6-loadtest\test.js

