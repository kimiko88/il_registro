import http from 'k6/http';
import { check, sleep } from 'k6';

// ─── Scenario 2: Afternoon Parents & Students Rush (13:30 - 14:30) ──────────
// 1,500 concurrent users consulting grades, communications, and homework
export const options = {
  stages: [
    { duration: '30s', target: 200 },   // Initial wave
    { duration: '45s', target: 1500 },  // Peak traffic wave
    { duration: '2m', target: 1500 },   // Sustained afternoon load
    { duration: '30s', target: 0 },     // Ramp-down
  ],
  thresholds: {
    http_req_duration: ['p(95)<300', 'p(99)<800'],
    http_req_failed: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.API_BASE_URL || 'http://localhost:8080/api/v1';

export default function () {
  // 1. Health check
  const healthRes = http.get(`${BASE_URL}/health`);
  check(healthRes, {
    'health ok': (r) => r.status === 200,
  });

  // 2. Parent Login
  const loginPayload = JSON.stringify({
    email: __ENV.PARENT_EMAIL || 'genitore1@scuolaprova.it',
    password: __ENV.PARENT_PASSWORD || 'Password123!',
  });

  const loginRes = http.post(`${BASE_URL}/auth/login`, loginPayload, {
    headers: { 'Content-Type': 'application/json' },
  });

  check(loginRes, {
    'login response valid': (r) => r.status === 200 || r.status === 401 || r.status === 429,
  });

  if (loginRes.status === 200) {
    const token = loginRes.json('token');
    const authHeaders = {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    };

    // 3. Read Communications & Noticeboard (High Cache Hits)
    const commsRes = http.get(`${BASE_URL}/communications?limit=20`, { headers: authHeaders });
    check(commsRes, {
      'communications loaded': (r) => r.status === 200 || r.status === 304,
    });

    // 4. Read Student Grades & Average
    const gradesRes = http.get(`${BASE_URL}/grades`, { headers: authHeaders });
    check(gradesRes, {
      'grades loaded': (r) => r.status === 200 || r.status === 304,
    });

    // 5. Read Agenda & Homework for tomorrow
    const agendaRes = http.get(`${BASE_URL}/agenda`, { headers: authHeaders });
    check(agendaRes, {
      'agenda loaded': (r) => r.status === 200 || r.status === 304,
    });

    // 6. Check Absence Limit DPR 122/2009
    const attRes = http.get(`${BASE_URL}/attendance`, { headers: authHeaders });
    check(attRes, {
      'attendance record loaded': (r) => r.status === 200 || r.status === 304,
    });
  }

  // Realistic browsing delay for families (2 to 5 seconds)
  sleep(Math.random() * 3 + 2);
}
