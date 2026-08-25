/**
 * @file rbac.spec.js
 * E2E tests: Role-Based Access Control con Playwright
 *
 * [R01] student non puo accedere a /admin/dashboard
 * [R02] parent non puo accedere a /teacher/grades
 * [R03] teacher non puo accedere a /secretary/users
 * [R04] admin accede a /admin/dashboard
 * [R05] secretary accede a /secretary
 * [R06] principal vede le quick actions corrette (no link /admin)
 * [R07] teacher vede azioni corrette (attendance, grades)
 */

import { test, expect } from '@playwright/test'

// Helper: login as a given role (assumes test accounts exist in the backend)
const CREDENTIALS = {
  admin: { email: 'admin@school.it', password: 'Admin1234!' },
  teacher: { email: 'teacher@school.it', password: 'Teacher1234!' },
  student: { email: 'student@school.it', password: 'Student1234!' },
  parent: { email: 'parent@school.it', password: 'Parent1234!' },
  secretary: { email: 'secretary@school.it', password: 'Secretary1234!' },
  principal: { email: 'principal@school.it', password: 'Principal1234!' },
}

async function loginAs(page, role) {
  const creds = CREDENTIALS[role]
  await page.goto('/login')
  await page.fill('input[type="email"]', creds.email)
  await page.fill('input[type="password"]', creds.password)
  await page.click('[data-testid="submit-login"]')
  await page.waitForURL(/\/(?!login)/, { timeout: 10000 })
}

test.describe('RBAC guard', () => {

  test('R01 — student redirected from /admin/dashboard', async ({ page }) => {
    await loginAs(page, 'student')
    await page.goto('/admin/dashboard')
    await expect(page).toHaveURL(/\/student/, { timeout: 5000 })
  })

  test('R02 — parent redirected from /teacher/grades', async ({ page }) => {
    await loginAs(page, 'parent')
    await page.goto('/teacher/grades')
    await expect(page).toHaveURL(/\/parent/, { timeout: 5000 })
  })

  test('R03 — teacher redirected from /secretary/users', async ({ page }) => {
    await loginAs(page, 'teacher')
    await page.goto('/secretary/users')
    await expect(page).toHaveURL(/\/teacher/, { timeout: 5000 })
  })

  test('R04 — admin can access /admin/dashboard', async ({ page }) => {
    await loginAs(page, 'admin')
    await page.goto('/admin/dashboard')
    await expect(page).toHaveURL(/\/admin\/dashboard/, { timeout: 5000 })
  })

  test('R05 — secretary can access /secretary', async ({ page }) => {
    await loginAs(page, 'secretary')
    await page.goto('/secretary')
    await expect(page).toHaveURL(/\/secretary/, { timeout: 5000 })
  })

  test('R06 — principal dashboard has no /admin links in actions', async ({ page }) => {
    await loginAs(page, 'principal')
    await page.goto('/')
    const actionLinks = await page.locator('[data-testid="quick-action"]').all()
    for (const link of actionLinks) {
      const href = await link.getAttribute('data-route') ?? ''
      expect(href).not.toMatch(/^\/admin/)
    }
  })

  test('R07 — teacher quick actions contain attendance and grades', async ({ page }) => {
    await loginAs(page, 'teacher')
    await page.goto('/')
    const actionLabels = await page.locator('[data-testid="quick-action"]').allTextContents()
    expect(actionLabels.some(l => /presenze|attendance/i.test(l))).toBe(true)
    expect(actionLabels.some(l => /voti|grades/i.test(l))).toBe(true)
  })
})
