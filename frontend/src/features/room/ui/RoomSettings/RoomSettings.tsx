/**
 * Room Feature - RoomSettings
 * Настройки комнаты (для хоста)
 */

import { useState, useEffect } from 'react';
import { useUpdateRoomSettings } from '@/entities/room';
import { Button, Card, Toggle } from '@/shared/ui';
import { DEFAULT_ROOM_SETTINGS } from '@/shared/config';
import type { GameRoom, RoomSettings } from '@/shared/types';
import './RoomSettings.css';

interface RoomSettingsProps {
  room: GameRoom;
  isHost: boolean;
}

export const RoomSettingsComponent = ({ room, isHost }: RoomSettingsProps) => {
  const [settings, setSettings] = useState<RoomSettings>(
    room.settings || DEFAULT_ROOM_SETTINGS
  );

  const updateSettingsMutation = useUpdateRoomSettings();

  // Синхронизируем с room.settings
  useEffect(() => {
    if (room.settings) {
      setSettings(room.settings);
    }
  }, [room.settings]);

  const handleSave = () => {
    updateSettingsMutation.mutate({
      id: room.id,
      settings,
    });
  };

  const hasChanges = 
    settings.timeForAnswer !== room.settings?.timeForAnswer ||
    settings.timeForChoice !== room.settings?.timeForChoice ||
    settings.allowWrongAnswer !== room.settings?.allowWrongAnswer ||
    settings.showRightAnswer !== room.settings?.showRightAnswer;

  if (!isHost) {
    return (
      <Card className="room-settings" padding="medium">
        <h3 className="room-settings__title">Настройки игры</h3>
        <div className="room-settings__view">
          <div className="room-settings__row">
            <span>Время на ответ:</span>
            <span>{settings.timeForAnswer} сек</span>
          </div>
          <div className="room-settings__row">
            <span>Время на выбор:</span>
            <span>{settings.timeForChoice} сек</span>
          </div>
          <div className="room-settings__row">
            <span>Неправильные ответы:</span>
            <span>{settings.allowWrongAnswer ? 'Разрешены' : 'Запрещены'}</span>
          </div>
          <div className="room-settings__row">
            <span>Показывать ответ:</span>
            <span>{settings.showRightAnswer ? 'Да' : 'Нет'}</span>
          </div>
        </div>
      </Card>
    );
  }

  return (
    <Card className="room-settings" padding="medium">
      <h3 className="room-settings__title">Настройки игры</h3>
      
      <div className="room-settings__form">
        <div className="room-settings__field">
          <label>
            Время на ответ: {settings.timeForAnswer} сек
          </label>
          <input
            type="range"
            min={10}
            max={60}
            step={5}
            value={settings.timeForAnswer}
            onChange={(e) =>
              setSettings({ ...settings, timeForAnswer: Number(e.target.value) })
            }
            className="room-settings__slider"
          />
        </div>

        <div className="room-settings__field">
          <label>
            Время на выбор вопроса: {settings.timeForChoice} сек
          </label>
          <input
            type="range"
            min={5}
            max={30}
            step={5}
            value={settings.timeForChoice}
            onChange={(e) =>
              setSettings({ ...settings, timeForChoice: Number(e.target.value) })
            }
            className="room-settings__slider"
          />
        </div>

        <div className="room-settings__field">
          <Toggle
            label="Разрешить неправильные ответы"
            checked={settings.allowWrongAnswer}
            onChange={(checked) =>
              setSettings({ ...settings, allowWrongAnswer: checked })
            }
          />
        </div>

        <div className="room-settings__field">
          <Toggle
            label="Показывать правильный ответ"
            checked={settings.showRightAnswer}
            onChange={(checked) =>
              setSettings({ ...settings, showRightAnswer: checked })
            }
          />
        </div>

        <Button
          onClick={handleSave}
          variant="primary"
          fullWidth
          disabled={!hasChanges}
          isLoading={updateSettingsMutation.isPending}
        >
          Сохранить настройки
        </Button>
      </div>
    </Card>
  );
};

