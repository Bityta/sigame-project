/**
 * LobbyPage - Интеграционные тесты
 * 
 * Тестируют страницу лобби:
 * - Отображение информации о пользователе
 * - Создание комнаты
 * - Вход по коду комнаты
 * - Выход из активной комнаты
 * - Logout
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { LobbyPage } from './LobbyPage';

// Мокаем навигацию
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

// Мокаем auth store
const mockSetAuthenticated = vi.fn();
vi.mock('@/features/auth', () => ({
  useLogout: (options: any) => ({
    mutate: () => options?.onSuccess?.(),
    isPending: false,
  }),
  useAuthStore: (selector: any) => selector({ setAuthenticated: mockSetAuthenticated }),
}));

// Мокаем error store
const mockSetError = vi.fn();
vi.mock('@/shared/lib/error-store', () => ({
  useErrorStore: (selector: any) => selector({ setError: mockSetError }),
}));

// Мокаем получение пользователя
let mockUser = { id: 'user-1', username: 'TestUser' };
let mockUserLoading = false;

vi.mock('@/entities/user', () => ({
  useCurrentUser: () => ({
    data: mockUser,
    isLoading: mockUserLoading,
  }),
}));

// Мокаем хуки комнаты
let mockActiveRoom: any = null;
let mockActiveRoomLoading = false;
const mockLeaveRoomMutate = vi.fn();
const mockGetRoomByCode = vi.fn();

vi.mock('@/entities/room', () => ({
  roomApi: {
    getRoomByCode: (code: string) => mockGetRoomByCode(code),
  },
  useMyActiveRoom: () => ({
    data: mockActiveRoom,
    isLoading: mockActiveRoomLoading,
  }),
  useLeaveRoom: () => ({
    mutate: mockLeaveRoomMutate,
    isPending: false,
  }),
}));

// Мокаем RoomList компонент
vi.mock('@/features/room', () => ({
  RoomList: ({ activeRoom, onLeaveRoom, isLeavingRoom }: any) => (
    <div data-testid="room-list">
      Room List
      {activeRoom && (
        <div data-testid="active-room-indicator">Active: {activeRoom.name}</div>
      )}
    </div>
  ),
}));

const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: false } },
});

const renderLobbyPage = () => {
  return render(
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <LobbyPage />
      </BrowserRouter>
    </QueryClientProvider>
  );
};

describe('LobbyPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockUser = { id: 'user-1', username: 'TestUser' };
    mockUserLoading = false;
    mockActiveRoom = null;
    mockActiveRoomLoading = false;
  });

  /**
   * ТЕСТ: Отображение спиннера во время загрузки пользователя
   * 
   * Проверяет что пока данные пользователя загружаются,
   * отображается индикатор загрузки
   */
  it('показывает спиннер при загрузке пользователя', () => {
    mockUserLoading = true;
    renderLobbyPage();

    expect(document.querySelector('.spinner')).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Отображение имени пользователя в приветствии
   * 
   * Проверяет что имя авторизованного пользователя
   * отображается в заголовке лобби
   */
  it('отображает приветствие с именем пользователя', () => {
    renderLobbyPage();

    expect(screen.getByText(/привет.*TestUser/i)).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Отображение кнопки "Создать комнату"
   * 
   * Проверяет наличие кнопки для создания новой комнаты
   */
  it('отображает кнопку создания комнаты', () => {
    renderLobbyPage();

    expect(screen.getByRole('button', { name: /создать комнату/i })).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Переход на страницу создания комнаты
   * 
   * Проверяет что клик на "Создать комнату"
   * перенаправляет на страницу создания
   */
  it('перенаправляет на страницу создания комнаты', async () => {
    const user = userEvent.setup();
    renderLobbyPage();

    const createButton = screen.getByRole('button', { name: /создать комнату/i });
    await user.click(createButton);

    expect(mockNavigate).toHaveBeenCalledWith('/lobby/create');
  });

  /**
   * ТЕСТ: Наличие поля для ввода кода комнаты
   * 
   * Проверяет что есть поле для ввода кода
   * и кнопка для входа по коду
   */
  it('отображает поле ввода кода комнаты', () => {
    renderLobbyPage();

    expect(screen.getByPlaceholderText(/код комнаты/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /войти по коду/i })).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Вход в комнату по коду
   * 
   * Проверяет что при вводе кода и клике на кнопку
   * происходит поиск комнаты и редирект
   */
  it('ищет комнату по коду и перенаправляет', async () => {
    mockGetRoomByCode.mockResolvedValue({ id: 'room-found', name: 'Found Room' });
    const user = userEvent.setup();
    renderLobbyPage();

    const codeInput = screen.getByPlaceholderText(/код комнаты/i);
    const joinButton = screen.getByRole('button', { name: /войти по коду/i });

    await user.type(codeInput, 'ABC123');
    await user.click(joinButton);

    await waitFor(() => {
      expect(mockGetRoomByCode).toHaveBeenCalledWith('ABC123');
      expect(mockNavigate).toHaveBeenCalledWith('/room/room-found');
    });
  });

  /**
   * ТЕСТ: Автоматический uppercase для кода комнаты
   * 
   * Проверяет что введенный код автоматически
   * преобразуется в верхний регистр
   */
  it('преобразует код комнаты в верхний регистр', async () => {
    const user = userEvent.setup();
    renderLobbyPage();

    const codeInput = screen.getByPlaceholderText(/код комнаты/i);
    await user.type(codeInput, 'abc123');

    expect(codeInput).toHaveValue('ABC123');
  });

  /**
   * ТЕСТ: Отправка кода по Enter
   * 
   * Проверяет что форму можно отправить
   * нажатием Enter в поле ввода кода
   */
  it('отправляет форму по Enter', async () => {
    mockGetRoomByCode.mockResolvedValue({ id: 'room-xyz', name: 'XYZ Room' });
    const user = userEvent.setup();
    renderLobbyPage();

    const codeInput = screen.getByPlaceholderText(/код комнаты/i);
    await user.type(codeInput, 'XYZ789{enter}');

    await waitFor(() => {
      expect(mockGetRoomByCode).toHaveBeenCalledWith('XYZ789');
    });
  });

  /**
   * ТЕСТ: Ошибка при несуществующем коде
   * 
   * Проверяет что при вводе несуществующего кода
   * отображается ошибка
   */
  it('показывает ошибку при несуществующем коде', async () => {
    mockGetRoomByCode.mockRejectedValue({ response: { status: 404 } });
    const user = userEvent.setup();
    renderLobbyPage();

    const codeInput = screen.getByPlaceholderText(/код комнаты/i);
    const joinButton = screen.getByRole('button', { name: /войти по коду/i });

    await user.type(codeInput, 'WRONG1');
    await user.click(joinButton);

    await waitFor(() => {
      expect(mockSetError).toHaveBeenCalledWith(
        expect.stringContaining('не найдена'),
        '404'
      );
    });
  });

  /**
   * ТЕСТ: Кнопка "Войти по коду" неактивна без ввода
   * 
   * Проверяет что кнопка заблокирована
   * пока не введен код
   */
  it('кнопка "Войти по коду" неактивна без ввода', () => {
    renderLobbyPage();

    const joinButton = screen.getByRole('button', { name: /войти по коду/i });
    expect(joinButton).toBeDisabled();
  });

  /**
   * ТЕСТ: Кнопка "Войти по коду" активна после ввода
   * 
   * Проверяет что кнопка активируется
   * после ввода кода
   */
  it('кнопка "Войти по коду" активна после ввода кода', async () => {
    const user = userEvent.setup();
    renderLobbyPage();

    const codeInput = screen.getByPlaceholderText(/код комнаты/i);
    const joinButton = screen.getByRole('button', { name: /войти по коду/i });

    await user.type(codeInput, 'TEST');

    expect(joinButton).not.toBeDisabled();
  });

  /**
   * ТЕСТ: Кнопка выхода из системы
   * 
   * Проверяет наличие кнопки logout
   */
  it('отображает кнопку выхода', () => {
    renderLobbyPage();

    expect(screen.getByRole('button', { name: /выход/i })).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Logout работает корректно
   * 
   * Проверяет что при клике на "Выход"
   * сбрасывается аутентификация и происходит редирект
   */
  it('выполняет logout при клике на кнопку выхода', async () => {
    const user = userEvent.setup();
    renderLobbyPage();

    const logoutButton = screen.getByRole('button', { name: /выход/i });
    await user.click(logoutButton);

    expect(mockSetAuthenticated).toHaveBeenCalledWith(false);
    expect(mockNavigate).toHaveBeenCalledWith('/login');
  });

  /**
   * ТЕСТ: Блокировка создания комнаты при активной комнате
   * 
   * Проверяет что если пользователь уже в комнате,
   * кнопка "Создать комнату" заблокирована
   */
  it('блокирует создание комнаты при наличии активной', () => {
    mockActiveRoom = { id: 'active-room', name: 'Моя активная комната' };
    renderLobbyPage();

    const createButton = screen.getByRole('button', { name: /создать комнату/i });
    expect(createButton).toBeDisabled();
  });

  /**
   * ТЕСТ: Блокировка входа по коду при активной комнате
   * 
   * Проверяет что если пользователь уже в комнате,
   * вход по коду заблокирован
   */
  it('блокирует вход по коду при наличии активной комнаты', async () => {
    mockActiveRoom = { id: 'active-room', name: 'Моя активная комната' };
    const user = userEvent.setup();
    renderLobbyPage();

    const codeInput = screen.getByPlaceholderText(/код комнаты/i);
    const joinButton = screen.getByRole('button', { name: /войти по коду/i });

    expect(codeInput).toBeDisabled();
    expect(joinButton).toBeDisabled();
  });

  /**
   * ТЕСТ: Tooltip при блокировке создания
   * 
   * Проверяет что заблокированная кнопка имеет
   * подсказку о причине блокировки
   */
  it('показывает tooltip на заблокированной кнопке создания', () => {
    mockActiveRoom = { id: 'active-room', name: 'Active Room' };
    renderLobbyPage();

    const createButton = screen.getByRole('button', { name: /создать комнату/i });
    expect(createButton).toHaveAttribute('title', expect.stringContaining('покиньте'));
  });

  /**
   * ТЕСТ: Отображение RoomList компонента
   * 
   * Проверяет что список комнат отображается на странице
   */
  it('отображает список комнат', () => {
    renderLobbyPage();

    expect(screen.getByTestId('room-list')).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Заголовок приложения
   * 
   * Проверяет что отображается название приложения
   */
  it('отображает заголовок приложения', () => {
    renderLobbyPage();

    expect(screen.getByRole('heading', { level: 1 })).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Ограничение длины кода (6 символов)
   * 
   * Проверяет что поле ввода кода имеет
   * ограничение maxLength=6
   */
  it('ограничивает длину кода 6 символами', async () => {
    const user = userEvent.setup();
    renderLobbyPage();

    const codeInput = screen.getByPlaceholderText(/код комнаты/i);
    await user.type(codeInput, 'ABCDEFGHIJ');

    expect(codeInput).toHaveValue('ABCDEF');
  });

  /**
   * ТЕСТ: Секция быстрых действий
   * 
   * Проверяет наличие блока "Быстрые действия"
   */
  it('отображает секцию быстрых действий', () => {
    renderLobbyPage();

    expect(screen.getByText(/быстрые действия/i)).toBeInTheDocument();
  });
});

