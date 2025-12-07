/**
 * LoginForm - Интеграционные тесты
 * 
 * Тестируют форму входа: валидацию, отправку данных, обработку ошибок
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { LoginForm } from './LoginForm';

// Мокаем навигацию
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

// Мокаем authStore
const mockSetAuthenticated = vi.fn();
vi.mock('../../model/authStore', () => ({
  useAuthStore: (selector: any) => selector({ setAuthenticated: mockSetAuthenticated }),
}));

// Мокаем мутацию логина
const mockLoginMutate = vi.fn();
let mockIsPending = false;

vi.mock('../../model/mutations', () => ({
  useLogin: (options: any) => ({
    mutate: (data: any) => {
      mockLoginMutate(data);
      if (data.username === 'testuser' && data.password === 'password123') {
        options?.onSuccess?.();
      } else {
        options?.onError?.();
      }
    },
    isPending: mockIsPending,
  }),
}));

const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: false } },
});

const renderLoginForm = () => {
  return render(
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <LoginForm />
      </BrowserRouter>
    </QueryClientProvider>
  );
};

describe('LoginForm', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockIsPending = false;
  });

  /**
   * ТЕСТ: Корректный рендеринг формы
   * 
   * Проверяет что форма отображает все необходимые элементы:
   * - Заголовок "Вход"
   * - Поле ввода имени пользователя
   * - Поле ввода пароля
   * - Кнопка "Войти"
   * - Ссылка на регистрацию
   */
  it('отображает все элементы формы', () => {
    renderLoginForm();

    expect(screen.getByRole('heading', { name: /вход/i })).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/введите имя пользователя/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/введите пароль/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /войти/i })).toBeInTheDocument();
    expect(screen.getByText(/нет аккаунта/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /зарегистрироваться/i })).toBeInTheDocument();
  });

  /**
   * ТЕСТ: Валидация короткого имени пользователя
   * 
   * Проверяет что при вводе слишком короткого имени (менее 3 символов)
   * показывается ошибка валидации и запрос на сервер не отправляется
   */
  it('показывает ошибку при коротком имени пользователя (менее 3 символов)', async () => {
    const user = userEvent.setup();
    renderLoginForm();

    const usernameInput = screen.getByPlaceholderText(/введите имя пользователя/i);
    const passwordInput = screen.getByPlaceholderText(/введите пароль/i);
    const submitButton = screen.getByRole('button', { name: /войти/i });

    await user.type(usernameInput, 'ab'); // Короткое имя
    await user.type(passwordInput, 'password123');
    await user.click(submitButton);

    expect(screen.getByText(/имя пользователя должно быть не менее 3 символов/i)).toBeInTheDocument();
    expect(mockLoginMutate).not.toHaveBeenCalled();
  });

  /**
   * ТЕСТ: Валидация короткого пароля
   * 
   * Проверяет что при вводе слишком короткого пароля (менее 6 символов)
   * показывается ошибка валидации и запрос на сервер не отправляется
   */
  it('показывает ошибку при коротком пароле (менее 6 символов)', async () => {
    const user = userEvent.setup();
    renderLoginForm();

    const usernameInput = screen.getByPlaceholderText(/введите имя пользователя/i);
    const passwordInput = screen.getByPlaceholderText(/введите пароль/i);
    const submitButton = screen.getByRole('button', { name: /войти/i });

    await user.type(usernameInput, 'testuser');
    await user.type(passwordInput, '12345'); // Короткий пароль
    await user.click(submitButton);

    expect(screen.getByText(/пароль должен быть не менее 6 символов/i)).toBeInTheDocument();
    expect(mockLoginMutate).not.toHaveBeenCalled();
  });

  /**
   * ТЕСТ: Успешная отправка формы
   * 
   * Проверяет что при валидных данных:
   * - Вызывается мутация логина с правильными данными
   * - После успеха вызывается setAuthenticated(true)
   * - Происходит редирект в лобби
   */
  it('отправляет форму и перенаправляет в лобби при успешном входе', async () => {
    const user = userEvent.setup();
    renderLoginForm();

    const usernameInput = screen.getByPlaceholderText(/введите имя пользователя/i);
    const passwordInput = screen.getByPlaceholderText(/введите пароль/i);
    const submitButton = screen.getByRole('button', { name: /войти/i });

    await user.type(usernameInput, 'testuser');
    await user.type(passwordInput, 'password123');
    await user.click(submitButton);

    expect(mockLoginMutate).toHaveBeenCalledWith({
      username: 'testuser',
      password: 'password123',
    });
    expect(mockSetAuthenticated).toHaveBeenCalledWith(true);
    expect(mockNavigate).toHaveBeenCalledWith('/lobby');
  });

  /**
   * ТЕСТ: Переход на страницу регистрации
   * 
   * Проверяет что клик на ссылку "Зарегистрироваться" 
   * перенаправляет на страницу регистрации
   */
  it('перенаправляет на страницу регистрации при клике на ссылку', async () => {
    const user = userEvent.setup();
    renderLoginForm();

    const registerLink = screen.getByRole('button', { name: /зарегистрироваться/i });
    await user.click(registerLink);

    expect(mockNavigate).toHaveBeenCalledWith('/register');
  });

  /**
   * ТЕСТ: Ввод данных в поля формы
   * 
   * Проверяет что пользователь может вводить текст в поля
   * и значения корректно обновляются
   */
  it('позволяет вводить данные в поля формы', async () => {
    const user = userEvent.setup();
    renderLoginForm();

    const usernameInput = screen.getByPlaceholderText(/введите имя пользователя/i);
    const passwordInput = screen.getByPlaceholderText(/введите пароль/i);

    await user.type(usernameInput, 'myusername');
    await user.type(passwordInput, 'mypassword');

    expect(usernameInput).toHaveValue('myusername');
    expect(passwordInput).toHaveValue('mypassword');
  });

  /**
   * ТЕСТ: Отправка формы по Enter
   * 
   * Проверяет что форму можно отправить нажатием Enter
   * после заполнения полей
   */
  it('отправляет форму при нажатии Enter', async () => {
    const user = userEvent.setup();
    renderLoginForm();

    const usernameInput = screen.getByPlaceholderText(/введите имя пользователя/i);
    const passwordInput = screen.getByPlaceholderText(/введите пароль/i);

    await user.type(usernameInput, 'testuser');
    await user.type(passwordInput, 'password123{enter}');

    expect(mockLoginMutate).toHaveBeenCalledWith({
      username: 'testuser',
      password: 'password123',
    });
  });
});

