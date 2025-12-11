/**
 * RoomPage - Интеграционные тесты
 * 
 * Тестируют страницу комнаты:
 * - Загрузка и отображение комнаты
 * - Отображение списка игроков
 * - Копирование кода комнаты
 * - Функционал хоста (старт игры, кик игроков)
 * - Выход из комнаты
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { RoomPage } from './RoomPage';

// Мокаем навигацию
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

// Мокаем хуки комнаты
const mockLeaveRoomMutate = vi.fn();
const mockStartGameMutate = vi.fn();
const mockKickPlayerMutate = vi.fn();
const mockTransferHostMutate = vi.fn();
const mockJoinRoomMutate = vi.fn();
const mockRefetch = vi.fn();

let mockRoom: any = null;
let mockIsLoading = false;
let mockUser = { id: 'user-1', username: 'TestUser' };

vi.mock('@/entities/room', () => ({
  useRoom: () => ({
    data: mockRoom,
    isLoading: mockIsLoading,
    refetch: mockRefetch,
  }),
  useLeaveRoom: () => ({
    mutate: mockLeaveRoomMutate,
    isPending: false,
  }),
  useStartGame: (options: any) => ({
    mutate: (roomId: string) => {
      mockStartGameMutate(roomId);
      options?.onSuccess?.({ gameId: 'game-123' });
    },
    isPending: false,
  }),
  useKickPlayer: () => ({
    mutate: mockKickPlayerMutate,
    isPending: false,
  }),
  useTransferHost: () => ({
    mutate: mockTransferHostMutate,
    isPending: false,
  }),
  useJoinRoom: () => ({
    mutate: mockJoinRoomMutate,
    isPending: false,
  }),
  useRoomEvents: () => {},
}));

vi.mock('@/entities/user', () => ({
  useCurrentUser: () => ({
    data: mockUser,
  }),
}));

vi.mock('@/features/room', () => ({
  RoomSettingsComponent: ({ room, isHost }: any) => (
    <div data-testid="room-settings">
      Settings (isHost: {isHost ? 'true' : 'false'})
    </div>
  ),
}));


const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: false } },
});

const renderRoomPage = (roomId = 'room-123') => {
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[`/room/${roomId}`]}>
        <Routes>
          <Route path="/room/:roomId" element={<RoomPage />} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>
  );
};

// Вспомогательная функция для создания мок-комнаты
const createMockRoom = (overrides?: any) => ({
  id: 'room-123',
  roomCode: 'ABC123',
  hostId: 'user-1',
  packId: 'pack-1',
  name: 'Тестовая комната',
  status: 'waiting',
  maxPlayers: 6,
  currentPlayers: 2,
  isPublic: true,
  hasPassword: false,
  settings: {
    timeForAnswer: 30,
    timeForChoice: 15,
  },
  players: [
    { userId: 'user-1', username: 'TestUser', role: 'host', joinedAt: new Date().toISOString() },
    { userId: 'user-2', username: 'Player2', role: 'player', joinedAt: new Date().toISOString() },
  ],
  createdAt: new Date().toISOString(),
  ...overrides,
});

describe('RoomPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockRoom = createMockRoom();
    mockIsLoading = false;
    mockUser = { id: 'user-1', username: 'TestUser' };
  });

  /**
   * ТЕСТ: Отображение спиннера во время загрузки
   * 
   * Проверяет что пока данные комнаты загружаются,
   * отображается индикатор загрузки (Spinner)
   */
  it('показывает спиннер при загрузке', () => {
    mockIsLoading = true;
    mockRoom = null;
    renderRoomPage();

    // Проверяем наличие spinner (по классу или роли)
    expect(document.querySelector('.spinner')).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Отображение ошибки когда комната не найдена
   * 
   * Проверяет что если комната не существует,
   * показывается сообщение об ошибке и кнопка возврата в лобби
   */
  it('показывает ошибку когда комната не найдена', () => {
    mockRoom = null;
    mockIsLoading = false;
    renderRoomPage();

    expect(screen.getByText(/комната не найдена/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /вернуться в лобби/i })).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Корректное отображение названия комнаты
   * 
   * Проверяет что название комнаты отображается в заголовке
   */
  it('отображает название комнаты', () => {
    renderRoomPage();

    expect(screen.getByRole('heading', { name: 'Тестовая комната' })).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Отображение кода комнаты
   * 
   * Проверяет что код комнаты отображается для пользователей
   */
  it('отображает код комнаты', () => {
    renderRoomPage();

    expect(screen.getByText('ABC123')).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Клик на код комнаты вызывает копирование
   * 
   * Проверяет что при клике на блок с кодом комнаты
   * происходит попытка копирования (через визуальную обратную связь)
   */
  it('клик на код комнаты вызывает визуальную обратную связь копирования', async () => {
    const user = userEvent.setup();
    renderRoomPage();

    const codeBlock = screen.getByTitle(/нажмите.*скопировать/i);
    
    // Изначально нет класса copied
    expect(codeBlock).not.toHaveClass('room-page__code--copied');
    
    await user.click(codeBlock);

    // После клика появляется класс copied (обратная связь)
    await waitFor(() => {
      expect(codeBlock).toHaveClass('room-page__code--copied');
    });
  });

  /**
   * ТЕСТ: Визуальная обратная связь после копирования
   * 
   * Проверяет что после успешного копирования
   * блок с кодом получает CSS-класс "copied"
   */
  it('показывает визуальную обратную связь после копирования', async () => {
    const user = userEvent.setup();
    renderRoomPage();

    const codeBlock = screen.getByTitle(/нажмите.*скопировать/i);
    await user.click(codeBlock);

    await waitFor(() => {
      expect(codeBlock).toHaveClass('room-page__code--copied');
    });
  });

  /**
   * ТЕСТ: Отображение списка игроков
   * 
   * Проверяет что все игроки комнаты отображаются
   * с их именами и ролями
   */
  it('отображает список игроков', () => {
    renderRoomPage();

    expect(screen.getByText(/TestUser/)).toBeInTheDocument();
    expect(screen.getByText(/Player2/)).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Отображение короны у хоста
   * 
   * Проверяет что рядом с именем хоста
   * отображается иконка короны (👑)
   */
  it('показывает корону у хоста', () => {
    renderRoomPage();

    const hostPlayer = screen.getByText(/TestUser/);
    expect(hostPlayer.textContent).toContain('👑');
  });

  /**
   * ТЕСТ: Отображение счетчика игроков
   * 
   * Проверяет что отображается текущее и максимальное
   * количество игроков в комнате
   */
  it('отображает счетчик игроков', () => {
    renderRoomPage();

    // Проверяем наличие текста с количеством игроков
    expect(screen.getByText(/игроки.*2.*6/i)).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Хост видит кнопку "Начать игру"
   * 
   * Проверяет что хост комнаты видит кнопку
   * для запуска игры
   */
  it('хост видит кнопку "Начать игру"', () => {
    renderRoomPage();

    expect(screen.getByRole('button', { name: /начать игру/i })).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Кнопка "Начать игру" активна при 2+ игроках
   * 
   * Проверяет что кнопка старта активна
   * когда в комнате 2 или более игроков
   */
  it('кнопка "Начать игру" активна при 2+ игроках', () => {
    renderRoomPage();

    const startButton = screen.getByRole('button', { name: /начать игру/i });
    expect(startButton).not.toBeDisabled();
  });

  /**
   * ТЕСТ: Кнопка "Начать игру" неактивна при 1 игроке
   * 
   * Проверяет что кнопка старта заблокирована
   * если в комнате только 1 игрок
   */
  it('кнопка "Начать игру" неактивна при 1 игроке', () => {
    mockRoom = createMockRoom({
      currentPlayers: 1,
      players: [
        { userId: 'user-1', username: 'TestUser', role: 'host', joinedAt: new Date().toISOString() },
      ],
    });
    renderRoomPage();

    const startButton = screen.getByRole('button', { name: /начать игру/i });
    expect(startButton).toBeDisabled();
  });

  /**
   * ТЕСТ: Показ подсказки о минимальном количестве игроков
   * 
   * Проверяет что при 1 игроке отображается
   * подсказка о необходимости минимум 2 игроков
   */
  it('показывает подсказку когда нужно больше игроков', () => {
    mockRoom = createMockRoom({
      currentPlayers: 1,
      players: [
        { userId: 'user-1', username: 'TestUser', role: 'host', joinedAt: new Date().toISOString() },
      ],
    });
    renderRoomPage();

    expect(screen.getByText(/минимум 2 игрока/i)).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Обычный игрок видит "Ожидание старта"
   * 
   * Проверяет что не-хост видит сообщение об ожидании
   * вместо кнопки старта
   */
  it('обычный игрок видит сообщение ожидания', () => {
    mockUser = { id: 'user-2', username: 'Player2' }; // Не хост
    renderRoomPage();

    expect(screen.getByText(/ожидание начала игры/i)).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: /начать игру/i })).not.toBeInTheDocument();
  });

  /**
   * ТЕСТ: Клик на "Начать игру" запускает игру
   * 
   * Проверяет что при клике на кнопку "Начать игру"
   * вызывается мутация и происходит редирект на игру
   */
  it('запускает игру при клике на "Начать игру"', async () => {
    const user = userEvent.setup();
    renderRoomPage();

    const startButton = screen.getByRole('button', { name: /начать игру/i });
    await user.click(startButton);

    expect(mockStartGameMutate).toHaveBeenCalledWith('room-123');
    expect(mockNavigate).toHaveBeenCalledWith('/game/game-123');
  });

  /**
   * ТЕСТ: Выход из комнаты
   * 
   * Проверяет что при клике на "Покинуть комнату"
   * вызывается мутация выхода
   */
  it('позволяет покинуть комнату', async () => {
    const user = userEvent.setup();
    renderRoomPage();

    const leaveButton = screen.getByRole('button', { name: /покинуть/i });
    await user.click(leaveButton);

    expect(mockLeaveRoomMutate).toHaveBeenCalled();
  });

  /**
   * ТЕСТ: Хост видит кнопки управления игроками
   * 
   * Проверяет что хост видит кнопки кика и передачи хоста
   * для других игроков
   */
  it('хост видит кнопки управления игроками', () => {
    renderRoomPage();

    // Кнопки должны быть у игрока Player2, но не у самого хоста
    const kickButtons = screen.getAllByTitle(/выгнать/i);
    const transferButtons = screen.getAllByTitle(/передать хоста/i);

    expect(kickButtons.length).toBeGreaterThan(0);
    expect(transferButtons.length).toBeGreaterThan(0);
  });

  /**
   * ТЕСТ: Обычный игрок не видит кнопки управления
   * 
   * Проверяет что обычный игрок не может
   * кикать или передавать хоста
   */
  it('обычный игрок не видит кнопки управления', () => {
    mockUser = { id: 'user-2', username: 'Player2' };
    renderRoomPage();

    expect(screen.queryByTitle(/выгнать/i)).not.toBeInTheDocument();
    expect(screen.queryByTitle(/передать хоста/i)).not.toBeInTheDocument();
  });

  /**
   * ТЕСТ: Кнопка "Вернуться в лобби" при ошибке работает
   * 
   * Проверяет что при клике на кнопку возврата в лобби
   * (на странице ошибки) происходит навигация
   */
  it('кнопка возврата в лобби работает', async () => {
    const user = userEvent.setup();
    mockRoom = null;
    mockIsLoading = false;
    renderRoomPage();

    const returnButton = screen.getByRole('button', { name: /вернуться в лобби/i });
    await user.click(returnButton);

    expect(mockNavigate).toHaveBeenCalledWith('/lobby');
  });

  /**
   * ТЕСТ: Компонент настроек передает правильный isHost
   * 
   * Проверяет что RoomSettingsComponent получает
   * корректное значение isHost
   */
  it('передает isHost=true в RoomSettings для хоста', () => {
    renderRoomPage();

    const settings = screen.getByTestId('room-settings');
    expect(settings).toHaveTextContent('isHost: true');
  });

  /**
   * ТЕСТ: Компонент настроек получает isHost=false для не-хоста
   * 
   * Проверяет что обычный игрок передается с isHost=false
   */
  it('передает isHost=false в RoomSettings для обычного игрока', () => {
    mockUser = { id: 'user-2', username: 'Player2' };
    renderRoomPage();

    const settings = screen.getByTestId('room-settings');
    expect(settings).toHaveTextContent('isHost: false');
  });
});

