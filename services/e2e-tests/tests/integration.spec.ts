import { test, expect } from '@playwright/test';
import { registerUser, generateUsername } from './helpers/auth';
import { createRoom, joinRoom, setReady } from './helpers/room';
import { waitForGameStart } from './helpers/game';

test.describe('Интеграционные тесты', () => {
  test('Auth Service - регистрация через API', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    const response = await page.request.post('/auth/register', {
      data: {
        username,
        password
      }
    });
    
    expect(response.status()).toBe(201);
    
    const body = await response.json();
    expect(body).toHaveProperty('user');
    expect(body).toHaveProperty('access_token');
    expect(body).toHaveProperty('refresh_token');
  });

  test('Auth Service - логин через API', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    await page.getByRole('button', { name: /выход/i }).click();
    
    const response = await page.request.post('/auth/login', {
      data: {
        username,
        password
      }
    });
    
    expect(response.status()).toBe(200);
    
    const body = await response.json();
    expect(body).toHaveProperty('access_token');
  });

  test('Lobby Service - создание комнаты через API', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    const token = await page.evaluate(() => localStorage.getItem('sigame_access_token'));
    
    const packsResponse = await page.request.get('/api/packs', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    const packs = await packsResponse.json();
    
    if (packs.packs && packs.packs.length > 0) {
      const packId = packs.packs[0].id;
      
      const response = await page.request.post('/api/lobby/rooms', {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        data: {
          name: 'Test Room',
          packId,
          maxPlayers: 4,
          isPublic: true
        }
      });
      
      expect(response.status()).toBe(201);
      
      const room = await response.json();
      expect(room).toHaveProperty('id');
      expect(room).toHaveProperty('roomCode');
    }
  });

  test('Lobby Service - получение списка комнат через API', async ({ page }) => {
    const response = await page.request.get('/api/lobby/rooms');
    
    expect(response.status()).toBe(200);
    
    const body = await response.json();
    expect(body).toHaveProperty('rooms');
    expect(Array.isArray(body.rooms)).toBe(true);
  });

  test('Game Service - создание игры через API', async ({ page, browser }) => {
    const hostUsername = generateUsername();
    const playerUsername = generateUsername();
    const password = 'testpass123';

    await registerUser(page, hostUsername, password);
    const roomId = await createRoom(page);
    
    const playerContext = await browser.newContext();
    const playerPage = await playerContext.newPage();
    await registerUser(playerPage, playerUsername, password);
    await joinRoom(playerPage, roomId);
    
    await setReady(page);
    await setReady(playerPage);
    
    await waitForGameStart(page);
    
    const gameUrl = page.url();
    const gameId = gameUrl.split('/game/')[1]?.split('?')[0];
    
    if (gameId) {
      const response = await page.request.get(`/api/game/${gameId}`);
      
      expect(response.status()).toBe(200);
      
      const game = await response.json();
      expect(game).toHaveProperty('game_id');
      expect(game).toHaveProperty('status');
    }
    
    await playerContext.close();
  });

  test('Pack Service - получение списка паков через API', async ({ page }) => {
    const username = generateUsername();
    const password = 'testpass123';

    await registerUser(page, username, password);
    
    const token = await page.evaluate(() => localStorage.getItem('access_token'));
    
    const response = await page.request.get('/api/packs', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    expect(response.status()).toBe(200);
    
    const body = await response.json();
    expect(body).toHaveProperty('packs');
    expect(Array.isArray(body.packs)).toBe(true);
  });
});

