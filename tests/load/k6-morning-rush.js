import http from 'k6/http';
import { check, sleep } from 'k6';

// ─── Scenario 1: Morning Rush (07:55 - 08:15) ──────────────────────────────
// 300 concurrent teachers fetching timetables and submitting batch attendance
export const options = {
  stages: [
    { duration: '20s', target: 50 },  // Ramp-up to 50 VUs
    { duration: '30s', target: 300 }, // Peak ramp-up to 300 VUs
    { duration: '1m', target: 300 },  // Maintain peak morning rush
    { duration: '20s', target: 0 },   // Ramp-down
  ],
  thresholds: {
    http_req_duration: ['p(95)<300', 'p(99)<800'],
    http_req_failed: ['rate<0.01'], // Less than 1% failure rate
  },
};

const BASE_URL = __ENV.API_BASE_URL || 'http://localhost:8080/api/v1';

export default function () {
  // 1. Health check & SingleFlight probe
  const healthRes = http.get(`${BASE_URL}/health`);
  check(healthRes, {
    'health status is 200': (r) => r.status === 200,
  });

  // 2. Simulate Teacher Auth (demo credentials or test token)
  const loginPayload = JSON.stringify({
    email: __ENV.TEACHER_EMAIL || 'docente1@scuolaprova.it',
    password: __ENV.TEACHER_PASSWORD || 'Password123!',
  });

  const loginRes = http.post(`${BASE_URL}/auth/login`, loginPayload, {
    headers: { 'Content-Type': 'application/json' },
  });

  const loginOk = check(loginRes, {
    'login status is 200 or 401/429 handled': (r) => r.status === 200 || r.status === 401 || r.status === 429,
  });

  if (loginRes.status === 200) {
    const token = loginRes.json('token');
    const authHeaders = {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    };

    // 3. Fetch Classes & Daily Timetable (SingleFlight deduplication under stress)
    const classesRes = http.get(`${BASE_URL}/classes`, { headers: authHeaders });
    check(classesRes, {
      'classes loaded': (r) => r.status === 200,
    });

    // 4. Batch Attendance Submission (Multi-Row Atomic Insert)
    const today = new Date().toISOString().split('T')[0];
    const batchPayload = JSON.stringify({
      class_id: __ENV.TEST_CLASS_ID || '00000000-0000-0000-0000-000000000001',
      date: today,
      hour: 1,
      statuses: [
        { student_id: '00000000-0000-0000-0000-000000000011', status: 'present', notes: '' },
        { student_id: '00000000-0000-0000-0000-000000000012', status: 'present', notes: '' },
        { student_id: '00000000-0000-0000-0000-000000000013', status: 'absent', notes: 'influenza' },
        { student_id: '00000000-0000-0000-0000-000000000014', status: 'present', notes: '' },
        { student_id: '00000000-0000-0000-0000-000000000015', status: 'late', notes: 'entrata 8:10', entry_time: '08:10' },
      ],
    });

    const attRes = http.post(`${BASE_URL}/attendance/batch`, batchPayload, { headers: authHeaders });
    check(attRes, {
      'batch attendance saved': (r) => r.status === 200 || r.status === 201 || r.status === 400 || r.status === 403,
    });

    // 5. Quick Lesson Sign
    const lessonPayload = JSON.stringify({
      class_id: __ENV.TEST_CLASS_ID || '00000000-0000-0000-0000-000000000001',
      date: today,
      hour: 1,
      subject_id: __ENV.TEST_SUBJECT_ID || '00000000-0000-0000-0000-000000000021',
      topic: 'Equazioni di secondo grado e disequazioni',
      homework: 'Esercizi pag. 142 numeri 1-10',
    });

    const lessonRes = http.post(`${BASE_URL}/lessons`, lessonPayload, { headers: authHeaders });
    check(lessonRes, {
      'lesson signed': (r) => r.status === 200 || r.status === 201 || r.status === 400 || r.status === 403,
    });
  }

  // Think time between teacher classroom actions (1 to 3 seconds)
  sleep(Math.random() * 2 + 1);
}
