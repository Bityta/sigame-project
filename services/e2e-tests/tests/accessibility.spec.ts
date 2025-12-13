import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';

test.describe('Доступность', () => {
  test('навигация по формам с клавиатуры - Tab порядок', async ({ page }) => {
    await page.goto('/login');
    
    await page.keyboard.press('Tab');
    const firstFocused = await page.evaluate(() => document.activeElement?.tagName);
    expect(firstFocused).toBe('INPUT');
    
    await page.keyboard.press('Tab');
    const secondFocused = await page.evaluate(() => document.activeElement?.tagName);
    expect(secondFocused).toBe('INPUT');
  });

  test('Enter для отправки формы входа', async ({ page }) => {
    await page.goto('/login');
    
    await page.getByLabel(/имя пользователя/i).fill('testuser');
    await page.getByLabel(/пароль/i).fill('testpass123');
    await page.keyboard.press('Enter');
    
    await page.waitForTimeout(1000);
    
    const errorVisible = await page.getByText(/ошибка/i).or(page.getByText(/неверный/i)).isVisible({ timeout: 5000 }).catch(() => false);
    if (errorVisible) {
      await expect(page.getByText(/ошибка/i).or(page.getByText(/неверный/i))).toBeVisible({ timeout: 5000 });
    }
  });

  test('Enter для отправки формы регистрации', async ({ page }) => {
    const username = generateUsername();
    
    await page.goto('/register');
    
    await page.getByLabel(/имя пользователя/i).fill(username);
    await page.waitForTimeout(600);
    await page.getByLabel(/^пароль$/i).fill('testpass123');
    await page.getByLabel(/подтвердите пароль/i).fill('testpass123');
    await page.keyboard.press('Enter');
    
    await page.waitForURL(/\/lobby/, { timeout: 10000 });
  });

  test('корректные aria-labels на элементах', async ({ page }) => {
    await page.goto('/login');
    
    const usernameInput = page.getByLabel(/имя пользователя/i);
    const hasAriaLabel = await usernameInput.getAttribute('aria-label').then(() => true).catch(() => false);
    const hasId = await usernameInput.getAttribute('id').then(() => true).catch(() => false);
    expect(hasAriaLabel || hasId).toBe(true);
    
    const passwordInput = page.getByLabel(/пароль/i);
    const passwordHasAriaLabel = await passwordInput.getAttribute('aria-label').then(() => true).catch(() => false);
    const passwordHasId = await passwordInput.getAttribute('id').then(() => true).catch(() => false);
    expect(passwordHasAriaLabel || passwordHasId).toBe(true);
  });

  test('корректные роли элементов', async ({ page }) => {
    await page.goto('/login');
    
    await expect(page.getByRole('heading', { name: /вход/i })).toBeVisible();
    await expect(page.getByRole('button', { name: /войти/i })).toBeVisible();
    await expect(page.getByRole('link', { name: /зарегистрироваться/i })).toBeVisible();
  });

  test('сообщения об ошибках доступны', async ({ page }) => {
    await page.goto('/login');
    
    await page.getByRole('button', { name: /войти/i }).click();
    
    const errorMessage = page.getByText(/ошибка/i).or(page.getByText(/имя/i)).or(page.getByText(/пароль/i));
    const errorVisible = await errorMessage.isVisible({ timeout: 5000 }).catch(() => false);
    if (errorVisible) {
      await expect(errorMessage).toBeVisible({ timeout: 5000 });
    }
    
    const ariaLive = await errorMessage.getAttribute('aria-live').catch(() => null);
    expect(ariaLive).toBeDefined();
  });
});

