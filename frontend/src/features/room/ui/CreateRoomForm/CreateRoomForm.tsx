/**
 * Room Feature - CreateRoomForm
 * Форма создания новой комнаты
 */

import { useState } from 'react';
import type { FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCreateRoom } from '@/entities/room';
import { usePacks } from '@/entities/pack';
import { Button, Input, Card } from '@/shared/ui';
import { ROUTES, LIMITS, DEFAULT_ROOM_SETTINGS } from '@/shared/config';
import type { CreateRoomRequest } from '@/shared/types';
import './CreateRoomForm.css';

export const CreateRoomForm = () => {
  const navigate = useNavigate();
  const { data: packs, isLoading: packsLoading } = usePacks();
  
  const [name, setName] = useState('');
  const [packId, setPackId] = useState('');
  const [maxPlayers, setMaxPlayers] = useState(4);
  const [isPublic, setIsPublic] = useState(true);
  const [password, setPassword] = useState('');
  const [localError, setLocalError] = useState('');

  const createRoomMutation = useCreateRoom({
    onSuccess: (room) => {
      navigate(ROUTES.ROOM(room.id));
    },
  });

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    setLocalError('');

    // Валидация
    if (name.length < LIMITS.ROOM_NAME.MIN || name.length > LIMITS.ROOM_NAME.MAX) {
      setLocalError(`Название должно быть от ${LIMITS.ROOM_NAME.MIN} до ${LIMITS.ROOM_NAME.MAX} символов`);
      return;
    }
    if (!packId) {
      setLocalError('Выберите пак вопросов');
      return;
    }
    if (!isPublic && !password) {
      setLocalError('Для приватной комнаты укажите пароль');
      return;
    }

    const data: CreateRoomRequest = {
      name,
      packId,
      maxPlayers,
      isPublic,
      password: isPublic ? undefined : password,
      settings: DEFAULT_ROOM_SETTINGS,
    };

    createRoomMutation.mutate(data);
  };

  return (
    <Card className="create-room-form">
      <h2 className="create-room-form__title">Создать комнату</h2>
      
      <form onSubmit={handleSubmit} className="create-room-form__form">
        <Input
          label="Название комнаты"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Введите название"
          fullWidth
          disabled={createRoomMutation.isPending}
        />

        <div className="create-room-form__field">
          <label className="create-room-form__label">Пак вопросов</label>
          <select
            className="create-room-form__select"
            value={packId}
            onChange={(e) => setPackId(e.target.value)}
            disabled={createRoomMutation.isPending || packsLoading}
          >
            <option value="">Выберите пак</option>
            {packs?.map((pack) => (
              <option key={pack.id} value={pack.id}>
                {pack.name} (автор: {pack.author})
              </option>
            ))}
          </select>
        </div>

        <div className="create-room-form__field">
          <label className="create-room-form__label">
            Максимум игроков: {maxPlayers}
          </label>
          <input
            type="range"
            min={LIMITS.MAX_PLAYERS.MIN}
            max={LIMITS.MAX_PLAYERS.MAX}
            value={maxPlayers}
            onChange={(e) => setMaxPlayers(Number(e.target.value))}
            className="create-room-form__slider"
            disabled={createRoomMutation.isPending}
          />
        </div>

        <div className="create-room-form__checkbox">
          <input
            type="checkbox"
            id="isPublic"
            checked={isPublic}
            onChange={(e) => setIsPublic(e.target.checked)}
            disabled={createRoomMutation.isPending}
          />
          <label htmlFor="isPublic">Публичная комната</label>
        </div>

        {!isPublic && (
          <Input
            label="Пароль"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Введите пароль"
            fullWidth
            disabled={createRoomMutation.isPending}
          />
        )}

        {localError && (
          <div className="create-room-form__error">{localError}</div>
        )}

        <div className="create-room-form__actions">
          <Button
            type="button"
            variant="ghost"
            onClick={() => navigate(ROUTES.LOBBY)}
            disabled={createRoomMutation.isPending}
          >
            Отмена
          </Button>
          <Button
            type="submit"
            variant="primary"
            isLoading={createRoomMutation.isPending}
          >
            Создать
          </Button>
        </div>
      </form>
    </Card>
  );
};

