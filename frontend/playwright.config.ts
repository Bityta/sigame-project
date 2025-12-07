/**
 * Playwright Configuration
 * 
 * Конфигурация для E2E тестирования
 */

import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  // Директория с тестами
  testDir: './e2e',
  
  // Параллельный запуск тестов
  fullyParallel: true,
  
  // Не разрешать test.only в CI
  forbidOnly: !!process.env.CI,
  
  // Количество повторов при падении (в CI)
  retries: process.env.CI ? 2 : 0,
  
  // Количество воркеров
  workers: process.env.CI ? 1 : undefined,
  
  // Репортер
  reporter: 'html',
  
  // Общие настройки для всех тестов
  use: {
    // Базовый URL приложения
    baseURL: 'http://localhost:5173',
    
    // Записывать trace при первом повторе
    trace: 'on-first-retry',
    
    // Скриншоты при падении
    screenshot: 'only-on-failure',
  },

  // Проекты (браузеры)
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    // Можно добавить другие браузеры:
    // {
    //   name: 'firefox',
    //   use: { ...devices['Desktop Firefox'] },
    // },
    // {
    //   name: 'webkit',
    //   use: { ...devices['Desktop Safari'] },
    // },
  ],

  // Запуск dev сервера перед тестами
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
    timeout: 120 * 1000,
  },
});


