/**
 * ProfilePage - Интеграционные тесты
 * 
 * Тестируют страницу профиля:
 * - Отображение информации о пользователе
 * - Аватар (буква или картинка)
 * - Переключение табов
 * - Заглушки на табах в разработке
 * - Вкладка настроек (задизейбленные поля)
 * - Навигация обратно в лобби
 * - Logout
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ProfilePage } from './ProfilePage';

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

// Мокаем получение пользователя
let mockUser: { id: string; username: string; avatarUrl?: string; createdAt: string } | null = {
  id: 'user-1',
  username: 'TestUser',
  createdAt: '2024-06-15T10:00:00Z',
};
let mockUserLoading = false;

vi.mock('@/entities/user', () => ({
  useCurrentUser: () => ({
    data: mockUser,
    isLoading: mockUserLoading,
  }),
}));

// Мокаем avatar utility
vi.mock('@/shared/lib/avatar', () => ({
  getAvatarUrl: (url: string | null | undefined) => {
    if (!url) return null;
    if (url.startsWith('http')) return url;
    return `http://localhost:9000/avatars/${url}`;
  },
}));

const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: false } },
});

const renderProfilePage = () => {
  return render(
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <ProfilePage />
      </BrowserRouter>
    </QueryClientProvider>
  );
};

describe('ProfilePage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockUser = {
      id: 'user-1',
      username: 'TestUser',
      createdAt: '2024-06-15T10:00:00Z',
    };
    mockUserLoading = false;
  });

  /**
   * ТЕСТ: Отображение спиннера во время загрузки
   */
  it('показывает спиннер при загрузке пользователя', () => {
    mockUserLoading = true;
    renderProfilePage();

    expect(document.querySelector('.spinner')).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Отображение имени пользователя
   */
  it('отображает имя пользователя', () => {
    renderProfilePage();

    expect(screen.getByText('TestUser')).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Отображение даты регистрации
   */
  it('отображает дату регистрации', () => {
    renderProfilePage();

    expect(screen.getByText(/В игре с/)).toBeInTheDocument();
    expect(screen.getByText(/июня 2024/)).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Аватар с буквой когда нет avatarUrl
   */
  it('показывает первую букву username когда нет аватарки', () => {
    renderProfilePage();

    expect(screen.getByText('T')).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Аватар с картинкой когда есть avatarUrl (полный URL)
   */
  it('показывает картинку когда есть полный avatarUrl', () => {
    mockUser = {
      ...mockUser!,
      avatarUrl: 'https://example.com/avatar.jpg',
    };
    renderProfilePage();

    const avatarImg = screen.getByAltText('TestUser');
    expect(avatarImg).toBeInTheDocument();
    expect(avatarImg).toHaveAttribute('src', 'https://example.com/avatar.jpg');
  });

  /**
   * ТЕСТ: Аватар с картинкой когда приходит avatar_id
   */
  it('строит URL аватарки когда приходит avatar_id', () => {
    mockUser = {
      ...mockUser!,
      avatarUrl: '550e8400-e29b-41d4-a716-446655440000',
    };
    renderProfilePage();

    const avatarImg = screen.getByAltText('TestUser');
    expect(avatarImg).toBeInTheDocument();
    expect(avatarImg).toHaveAttribute('src', 'http://localhost:9000/avatars/550e8400-e29b-41d4-a716-446655440000');
  });

  /**
   * ТЕСТ: Кнопка изменения аватара задизейблена
   */
  it('кнопка изменения аватара задизейблена', () => {
    renderProfilePage();

    const avatarButton = screen.getByTitle(/изменение аватара скоро/i);
    expect(avatarButton).toBeDisabled();
  });

  /**
   * ТЕСТ: Заголовок страницы
   */
  it('отображает заголовок "Профиль игрока"', () => {
    renderProfilePage();

    expect(screen.getByRole('heading', { name: 'Профиль игрока' })).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Кнопка "В лобби"
   */
  it('отображает кнопку возврата в лобби', () => {
    renderProfilePage();

    expect(screen.getByRole('button', { name: /в лобби/i })).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Навигация в лобби по кнопке
   */
  it('перенаправляет в лобби при клике на кнопку', async () => {
    const user = userEvent.setup();
    renderProfilePage();

    const backButton = screen.getByRole('button', { name: /в лобби/i });
    await user.click(backButton);

    expect(mockNavigate).toHaveBeenCalledWith('/lobby');
  });

  /**
   * ТЕСТ: Кнопка выхода
   */
  it('отображает кнопку выхода', () => {
    renderProfilePage();

    expect(screen.getByRole('button', { name: /выход/i })).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Logout работает корректно
   */
  it('выполняет logout при клике на кнопку выхода', async () => {
    const user = userEvent.setup();
    renderProfilePage();

    const logoutButton = screen.getByRole('button', { name: /выход/i });
    await user.click(logoutButton);

    expect(mockSetAuthenticated).toHaveBeenCalledWith(false);
    expect(mockNavigate).toHaveBeenCalledWith('/login');
  });

  /**
   * ТЕСТ: Отображение всех табов
   */
  it('отображает все 4 таба', () => {
    renderProfilePage();

    expect(screen.getByRole('button', { name: /статистика/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /достижения/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /история игр/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /настройки/i })).toBeInTheDocument();
  });

  /**
   * ТЕСТ: По умолчанию активен таб "Настройки"
   */
  it('по умолчанию активен таб настроек', () => {
    renderProfilePage();

    const settingsTab = screen.getByRole('button', { name: /настройки/i });
    expect(settingsTab).toHaveClass('profile-tabs__tab--active');
  });

  /**
   * ТЕСТ: Переключение на таб "Статистика"
   */
  it('переключается на таб статистики', async () => {
    const user = userEvent.setup();
    renderProfilePage();

    const statsTab = screen.getByRole('button', { name: /статистика/i });
    await user.click(statsTab);

    expect(statsTab).toHaveClass('profile-tabs__tab--active');
    expect(screen.getByText(/этот раздел находится в разработке/i)).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Переключение на таб "Достижения"
   */
  it('переключается на таб достижений', async () => {
    const user = userEvent.setup();
    renderProfilePage();

    const achievementsTab = screen.getByRole('button', { name: /достижения/i });
    await user.click(achievementsTab);

    expect(achievementsTab).toHaveClass('profile-tabs__tab--active');
    expect(screen.getByText(/этот раздел находится в разработке/i)).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Переключение на таб "История игр"
   */
  it('переключается на таб истории игр', async () => {
    const user = userEvent.setup();
    renderProfilePage();

    const historyTab = screen.getByRole('button', { name: /история игр/i });
    await user.click(historyTab);

    expect(historyTab).toHaveClass('profile-tabs__tab--active');
    expect(screen.getByText(/этот раздел находится в разработке/i)).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Заглушка "Статистика скоро появится"
   */
  it('показывает заглушку статистики в карточке профиля', () => {
    renderProfilePage();

    expect(screen.getByText(/статистика скоро появится/i)).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Вкладка настроек - поле имени пользователя задизейблено
   */
  it('поле имени пользователя в настройках задизейблено', () => {
    renderProfilePage();

    const usernameInput = screen.getByDisplayValue('TestUser');
    expect(usernameInput).toBeInTheDocument();
    expect(usernameInput).toBeDisabled();
  });

  /**
   * ТЕСТ: Вкладка настроек - секция безопасности с задизейбленными полями
   */
  it('поля смены пароля задизейблены', () => {
    renderProfilePage();

    const passwordInputs = screen.getAllByPlaceholderText(/пароль|••••/i);
    passwordInputs.forEach(input => {
      expect(input).toBeDisabled();
    });
  });

  /**
   * ТЕСТ: Вкладка настроек - кнопка изменения пароля задизейблена
   */
  it('кнопка изменения пароля задизейблена', () => {
    renderProfilePage();

    const changePasswordButton = screen.getByRole('button', { name: /изменить пароль/i });
    expect(changePasswordButton).toBeDisabled();
  });

  /**
   * ТЕСТ: Вкладка настроек - секция уведомлений активна
   */
  it('отображает секцию уведомлений', () => {
    renderProfilePage();

    expect(screen.getByText(/уведомления/i)).toBeInTheDocument();
    expect(screen.getByText(/звуковые уведомления в игре/i)).toBeInTheDocument();
    expect(screen.getByText(/показывать приглашения в комнаты/i)).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Вкладка настроек - тоггл уведомлений работает
   */
  it('переключает тоггл уведомлений', async () => {
    const user = userEvent.setup();
    renderProfilePage();

    const toggles = screen.getAllByRole('checkbox');
    expect(toggles.length).toBeGreaterThan(0);
    
    // Первый тоггл должен быть включен по умолчанию
    expect(toggles[0]).toBeChecked();
    
    // Кликаем по тогглу
    await user.click(toggles[0]);
    
    expect(toggles[0]).not.toBeChecked();
  });

  /**
   * ТЕСТ: Вкладка настроек - опасная зона
   */
  it('отображает секцию опасной зоны', () => {
    renderProfilePage();

    expect(screen.getByText(/опасная зона/i)).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Кнопка удаления аккаунта задизейблена
   */
  it('кнопка удаления аккаунта задизейблена', () => {
    renderProfilePage();

    const deleteButton = screen.getByRole('button', { name: /удалить аккаунт/i });
    expect(deleteButton).toBeDisabled();
  });

  /**
   * ТЕСТ: Бейджи "скоро" отображаются
   */
  it('отображает бейджи "скоро" на недоступных секциях', () => {
    renderProfilePage();

    const badges = screen.getAllByText('скоро');
    expect(badges.length).toBeGreaterThanOrEqual(3); // Профиль, Безопасность, Опасная зона
  });

  /**
   * ТЕСТ: Иконка в заглушке "В разработке"
   */
  it('показывает иконку стройки в заглушке', async () => {
    const user = userEvent.setup();
    renderProfilePage();

    const statsTab = screen.getByRole('button', { name: /статистика/i });
    await user.click(statsTab);

    expect(screen.getByText('🚧')).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Подсказка про изменение имени
   */
  it('показывает подсказку что изменение имени недоступно', () => {
    renderProfilePage();

    expect(screen.getByText(/изменение имени пользователя пока недоступно/i)).toBeInTheDocument();
  });
});
