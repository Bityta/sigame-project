import { Page, expect } from '@playwright/test';

export async function registerUser(
  page: Page,
  username: string,
  password: string
): Promise<void> {
  await page.goto('/register');
  await page.waitForLoadState('networkidle');
  await page.waitForTimeout(1000);
  
  const usernameInput = page.getByPlaceholder(/от 3 до 20 символов/i);
  await usernameInput.waitFor({ state: 'visible', timeout: 30000 });
  
  if (username.length < 3 || username.length > 20) {
    throw new Error(`Invalid username length: ${username.length}. Username: ${username}`);
  }
  
  await usernameInput.fill(username, { delay: 50 });
  
  await page.waitForTimeout(800);
  
  const filledUsername = await usernameInput.inputValue();
  if (filledUsername !== username) {
    throw new Error(`Username not filled correctly. Expected: ${username}, Got: ${filledUsername}`);
  }
  
  if (filledUsername.length < 3 || filledUsername.length > 20) {
    throw new Error(`Filled username has invalid length: ${filledUsername.length}. Username: ${filledUsername}`);
  }
  
  await page.waitForTimeout(2000);
  
  const passwordInput = page.getByPlaceholder(/минимум 6 символов/i);
  await passwordInput.fill(password, { delay: 50 });
  await page.waitForTimeout(300);
  
  const confirmPasswordInput = page.getByPlaceholder(/повторите пароль/i);
  await confirmPasswordInput.fill(password, { delay: 50 });
  await page.waitForTimeout(300);
  
  const submitButton = page.getByRole('button', { name: /зарегистрироваться/i });
  await submitButton.waitFor({ state: 'visible' });
  
  const hintElement = page.locator('.input-hint').first();
  
  await page.waitForFunction(
    async () => {
      const hint = document.querySelector('.input-hint');
      const hintText = hint?.textContent || '';
      if (hintText.includes('занято') || hintText.includes('Имя занято')) {
        return false;
      }
      const button = document.querySelector('button[type="submit"]') as HTMLButtonElement;
      if (!button) return false;
      return !button.disabled;
    },
    { timeout: 25000 }
  ).catch(async () => {
    const hintText = await hintElement.textContent().catch(() => '');
    if (hintText.includes('занято') || hintText.includes('Имя занято')) {
      const newUsername = generateUsername();
      return registerUser(page, newUsername, password);
    }
    const isDisabled = await submitButton.isDisabled();
    const usernameValue = await usernameInput.inputValue();
    throw new Error(`Button still disabled. Username: ${usernameValue} (${usernameValue.length}), Hint: ${hintText}, Disabled: ${isDisabled}`);
  });
  
  const responsePromise = page.waitForResponse(
    response => response.url().includes('/auth/register'),
    { timeout: 30000 }
  ).catch(() => null);
  
  await page.waitForTimeout(500);
  await submitButton.click({ delay: 100 });
  
  await page.waitForTimeout(1500);
  
  const response = await responsePromise;
  
  if (response) {
    if (response.status() === 409) {
      const newUsername = generateUsername();
      return registerUser(page, newUsername, password);
    }
    
    if (response.status() !== 201) {
      const body = await response.text().catch(() => '');
      throw new Error(`Registration failed: ${response.status()} - ${body}`);
    }
  }
  
  await page.waitForURL(/\/lobby/, { timeout: 30000 }).catch(async () => {
    const currentUrl = page.url();
    const errorText = await page.locator('.input-error').first().textContent().catch(() => '');
    const hintText = await hintElement.textContent().catch(() => '');
    const buttonDisabled = await submitButton.isDisabled();
    throw new Error(`Did not navigate to lobby. URL: ${currentUrl}, Error: ${errorText}, Hint: ${hintText}, Button disabled: ${buttonDisabled}`);
  });
  
  await expect(page).toHaveURL(/\/lobby/);
  
  await page.waitForFunction(
    () => {
      const keys = Object.keys(localStorage);
      return keys.some(key => key.toLowerCase().includes('token'));
    },
    { timeout: 10000 }
  );
}

export async function loginUser(
  page: Page,
  username: string,
  password: string
): Promise<void> {
  await page.goto('/login');
  
  await page.getByPlaceholder(/введите имя пользователя/i).fill(username);
  await page.getByPlaceholder(/введите пароль/i).fill(password);
  
  await page.getByRole('button', { name: /войти/i }).click();
  
  await page.waitForURL(/\/lobby/, { timeout: 10000 });
  await expect(page).toHaveURL(/\/lobby/);
}

export async function logoutUser(page: Page): Promise<void> {
  await page.getByRole('button', { name: /выход/i }).click();
  await page.waitForURL(/\/login/);
  await expect(page).toHaveURL(/\/login/);
}

export function generateUsername(): string {
  const timestamp = Date.now().toString().slice(-8);
  const random = Math.random().toString(36).substring(2, 6);
  const username = `test_${timestamp}_${random}`;
  
  if (username.length > 20) {
    return username.substring(0, 20);
  }
  
  if (username.length < 3) {
    return `test_${timestamp}_${random}_x`.substring(0, 20);
  }
  
  return username;
}

