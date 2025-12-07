/**
 * ProfilePage
 * Страница профиля игрока
 */

import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCurrentUser } from '@/entities/user';
import { useLogout, useAuthStore } from '@/features/auth';
import { Button, Card, Spinner } from '@/shared/ui';
import { ROUTES } from '@/shared/config';
import { getAvatarUrl } from '@/shared/lib/avatar';
import './ProfilePage.css';

type TabType = 'stats' | 'achievements' | 'history' | 'settings';

/**
 * Компонент заглушки "В разработке"
 */
const ComingSoon = ({ title }: { title: string }) => (
  <div className="coming-soon">
    <div className="coming-soon__icon">🚧</div>
    <h3 className="coming-soon__title">{title}</h3>
    <p className="coming-soon__text">Этот раздел находится в разработке и скоро будет доступен</p>
  </div>
);

/**
 * Бейдж "Скоро" для недоступных функций
 */
const ComingSoonBadge = () => (
  <span className="coming-soon-badge">скоро</span>
);

export const ProfilePage = () => {
  const navigate = useNavigate();
  const { data: user, isLoading } = useCurrentUser();
  const setAuthenticated = useAuthStore((state) => state.setAuthenticated);
  const [activeTab, setActiveTab] = useState<TabType>('settings');

  const logoutMutation = useLogout({
    onSuccess: () => {
      setAuthenticated(false);
      navigate(ROUTES.LOGIN);
    },
  });

  if (isLoading) {
    return (
      <div className="profile-page">
        <Spinner center size="large" />
      </div>
    );
  }

  // Получаем URL аватарки (поддерживает как полный URL, так и avatar_id)
  const avatarUrl = getAvatarUrl(user?.avatarUrl);

  return (
    <div className="profile-page">
      <header className="profile-page__header">
        <div className="profile-page__header-content">
          <button 
            className="profile-page__back-btn"
            onClick={() => navigate(ROUTES.LOBBY)}
          >
            ← В лобби
          </button>
          <h1 className="profile-page__title">Профиль игрока</h1>
          <Button
            variant="ghost"
            size="small"
            onClick={() => logoutMutation.mutate()}
            isLoading={logoutMutation.isPending}
            className="profile-page__logout-btn"
          >
            Выход
          </Button>
        </div>
      </header>

      <div className="profile-page__content">
        {/* Карточка профиля */}
        <Card className="profile-card" padding="large">
          <div className="profile-card__avatar-section">
            <div className="profile-card__avatar">
              {avatarUrl ? (
                <img 
                  src={avatarUrl} 
                  alt={user?.username}
                  className="profile-card__avatar-image"
                />
              ) : (
                <span className="profile-card__avatar-letter">
                  {user?.username?.charAt(0).toUpperCase() || '?'}
                </span>
              )}
              <button 
                className="profile-card__avatar-edit profile-card__avatar-edit--disabled" 
                title="Изменение аватара скоро будет доступно"
                disabled
              >
                📷
              </button>
            </div>
            <div className="profile-card__info">
              <h2 className="profile-card__username">{user?.username}</h2>
              <p className="profile-card__member-since">
                В игре с {user?.createdAt ? new Date(user.createdAt).toLocaleDateString('ru-RU', { 
                  year: 'numeric', 
                  month: 'long' 
                }) : 'недавно'}
              </p>
            </div>
          </div>

          <div className="profile-card__stats-placeholder">
            <span className="profile-card__stats-placeholder-text">
              📊 Статистика скоро появится
            </span>
          </div>
        </Card>

        {/* Навигация по табам */}
        <div className="profile-tabs">
          <button 
            className={`profile-tabs__tab ${activeTab === 'stats' ? 'profile-tabs__tab--active' : ''}`}
            onClick={() => setActiveTab('stats')}
          >
            📊 Статистика
          </button>
          <button 
            className={`profile-tabs__tab ${activeTab === 'achievements' ? 'profile-tabs__tab--active' : ''}`}
            onClick={() => setActiveTab('achievements')}
          >
            🏆 Достижения
          </button>
          <button 
            className={`profile-tabs__tab ${activeTab === 'history' ? 'profile-tabs__tab--active' : ''}`}
            onClick={() => setActiveTab('history')}
          >
            📜 История игр
          </button>
          <button 
            className={`profile-tabs__tab ${activeTab === 'settings' ? 'profile-tabs__tab--active' : ''}`}
            onClick={() => setActiveTab('settings')}
          >
            ⚙️ Настройки
          </button>
        </div>

        {/* Контент табов */}
        <div className="profile-tab-content">
          {activeTab === 'stats' && (
            <Card className="tab-card" padding="large">
              <ComingSoon title="Статистика" />
            </Card>
          )}

          {activeTab === 'achievements' && (
            <Card className="tab-card" padding="large">
              <ComingSoon title="Достижения" />
            </Card>
          )}

          {activeTab === 'history' && (
            <Card className="tab-card" padding="large">
              <ComingSoon title="История игр" />
            </Card>
          )}

          {activeTab === 'settings' && (
            <div className="settings-section">
              <Card className="settings-card" padding="medium">
                <h3 className="settings-card__title">
                  👤 Профиль
                  <ComingSoonBadge />
                </h3>
                <div className="settings-card__content">
                  <div className="settings-field">
                    <label className="settings-field__label">Имя пользователя</label>
                    <input 
                      type="text" 
                      className="settings-field__input"
                      value={user?.username || ''} 
                      disabled
                      readOnly
                    />
                  </div>
                  <p className="settings-card__hint">
                    Изменение имени пользователя пока недоступно
                  </p>
                </div>
              </Card>

              <Card className="settings-card" padding="medium">
                <h3 className="settings-card__title">
                  🔐 Безопасность
                  <ComingSoonBadge />
                </h3>
                <div className="settings-card__content">
                  <div className="settings-field">
                    <label className="settings-field__label">Текущий пароль</label>
                    <input 
                      type="password" 
                      className="settings-field__input" 
                      placeholder="••••••••" 
                      disabled 
                    />
                  </div>
                  <div className="settings-field">
                    <label className="settings-field__label">Новый пароль</label>
                    <input 
                      type="password" 
                      className="settings-field__input" 
                      placeholder="Введите новый пароль" 
                      disabled 
                    />
                  </div>
                  <div className="settings-field">
                    <label className="settings-field__label">Подтвердите пароль</label>
                    <input 
                      type="password" 
                      className="settings-field__input" 
                      placeholder="Повторите новый пароль" 
                      disabled 
                    />
                  </div>
                  <div className="settings-card__actions">
                    <Button variant="primary" size="small" disabled>
                      Изменить пароль
                    </Button>
                  </div>
                </div>
              </Card>

              <Card className="settings-card" padding="medium">
                <h3 className="settings-card__title">🔔 Уведомления</h3>
                <div className="settings-card__content">
                  <label className="settings-toggle">
                    <input type="checkbox" defaultChecked />
                    <span className="settings-toggle__slider"></span>
                    <span className="settings-toggle__label">Звуковые уведомления в игре</span>
                  </label>
                  <label className="settings-toggle">
                    <input type="checkbox" defaultChecked />
                    <span className="settings-toggle__slider"></span>
                    <span className="settings-toggle__label">Показывать приглашения в комнаты</span>
                  </label>
                </div>
              </Card>

              <Card className="settings-card settings-card--danger" padding="medium">
                <h3 className="settings-card__title">
                  ⚠️ Опасная зона
                  <ComingSoonBadge />
                </h3>
                <div className="settings-card__content">
                  <p className="settings-card__warning">
                    Удаление аккаунта приведёт к потере всех данных. Это действие необратимо.
                  </p>
                  <Button variant="danger" size="small" disabled>
                    Удалить аккаунт
                  </Button>
                </div>
              </Card>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
