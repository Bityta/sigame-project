/**
 * Room Feature - RoomSettings
 * Настройки комнаты (для хоста)
 */

import { useState, useEffect } from 'react';
import { useUpdateRoomSettings } from '@/entities/room';
import { Button, Card } from '@/shared/ui';
import { DEFAULT_ROOM_SETTINGS } from '@/shared/config';
import type { GameRoom, RoomSettings } from '@/shared/types';
import './RoomSettings.css';

interface RoomSettingsProps {
  room: GameRoom;
  isHost: boolean;
}

export const RoomSettingsComponent = ({ room, isHost }: RoomSettingsProps) => {
  const initialSettings = room.settings || DEFAULT_ROOM_SETTINGS;
  
  // Текущие настройки (редактируемые пользователем)
  const [settings, setSettings] = useState<RoomSettings>(initialSettings);
  
  // Сохранённые настройки (то, что сейчас на сервере)
  const [savedSettings, setSavedSettings] = useState<RoomSettings>(initialSettings);

  const updateSettingsMutation = useUpdateRoomSettings();

  // Синхронизируем с room.settings при изменении из вне
  useEffect(() => {
    if (room.settings) {
      setSettings(room.settings);
      setSavedSettings(room.settings);
    }
  }, [room.settings]);

  const handleSave = () => {
    updateSettingsMutation.mutate(
      {
        id: room.id,
        settings,
      },
      {
        onSuccess: (updatedRoom) => {
          // Синхронизируем оба стейта с ответом сервера
          if (updatedRoom.settings) {
            setSettings(updatedRoom.settings);
            setSavedSettings(updatedRoom.settings);
          }
        },
      }
    );
  };

  // Сравниваем с сохранёнными настройками, а не с пропсом
  const hasChanges = 
    settings.timeForAnswer !== savedSettings.timeForAnswer ||
    settings.timeForChoice !== savedSettings.timeForChoice;

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

