/**
 * @file auth.spec.js
 * E2E tests: Auth flow con Playwright
 *
 * Covers:
 * [E01] Login con credenziali valide → redirect al dashboard di ruolo
 * [E02] Login con credenziali errate → messaggio di errore visibile
 * [E03] Brute force: 5 tentativi falliti → UI lockout visibile
 * [E04] Token scaduto → redirect a /login?reason=session_expired
 * [E05] Logout → redirect a /login e token rimosso
 * [E06] URL diretto su route protetta senza sessione → redirect /login
 * [E07] /register apre il dialog segreteria
 */

import { test, expect } from '@playwright/test'

test.describe('Auth Flow', () => {

  test('E01 — valid admin login redirects to /admin/dashboard', async ({ page }) => {
    await page.goto('/login')
    await page.fill('input[type="email"]', 'admin@school.it')
    await page.fill('input[type="password"]', 'Admin1234!')
    await page.click('[data-testid="submit-login"]')
    await expect(page).toHaveURL(/\/admin\/dashboard/, { timeout: 10000 })
  })

  test('E02 — wrong credentials show error message', async ({ page }) => {
    await page.goto('/login')
    await page.fill('input[type="email"]', 'wrong@school.it')
    await page.fill('input[type="password"]', 'WrongPass')
    await page.click('[data-testid="submit-login"]')
    await expect(page.locator('[data-testid="error-message"], .q-banner')).toBeVisible({ timeout: 5000 })
  })

  test('E03 — 5 failed attempts shows lockout message', async ({ page }) => {
    await page.goto('/login')
    for (let i = 0; i < 5; i++) {
      await page.fill('input[type="email"]', 'hacker@school.it')
      await page.fill('input[type="password"]', `bad${i}`)
      await page.click('[data-testid="submit-login"]')
      await page.waitForTimeout(300)
    }
    const errorText = await page.locator('[data-testid="error-message"]').textContent()
    expect(errorText).toMatch(/Riprova tra|tooManyAttempts/i)
  })

  test('E04 — expired token in sessionStorage → redirect /login?reason=session_expired', async ({ page }) => {
    // Inject an expired token into sessionStorage before navigation
    await page.goto('/login')
    const expiredToken = `h.${btoa(JSON.stringify({ sub: 'u1', exp: Math.floor((Date.now() - 60_000) / 1000) }))}.s`
    await page.evaluate((tok) => {
      sessionStorage.setItem('auth_token', tok)
    }, expiredToken)
    await page.goto('/teacher/grades')
    await expect(page).toHaveURL(/\/login\?reason=session_expired/, { timeout: 5000 })
  })

  test('E05 — logout redirects to /login', async ({ page }) => {
    // Login first
    await page.goto('/login')
    await page.fill('input[type="email"]', 'teacher@school.it')
    await page.fill('input[type="password"]', 'Teacher1234!')
    await page.click('[data-testid="submit-login"]')
    await expect(page).toHaveURL(/\/teacher/, { timeout: 10000 })
    // Logout
    await page.click('[data-testid="logout-button"]')
    await expect(page).toHaveURL(/\/login/, { timeout: 5000 })
  })

  test('E06 — direct protected URL without session → /login', async ({ page }) => {
    await page.context().clearCookies()
    await page.evaluate(() => { sessionStorage.clear(); localStorage.clear() })
    await page.goto('/admin/dashboard')
    await expect(page).toHaveURL(/\/login/, { timeout: 5000 })
  })

  test('E07 — /register opens secretary dialog', async ({ page }) => {
    await page.goto('/register')
    await expect(page.locator('.q-dialog, [data-testid="secretary-dialog"]')).toBeVisible({ timeout: 5000 })
  })
})
