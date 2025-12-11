"""Mock data for Pack Service - Hardcoded test packs"""
from typing import List, Dict, Any
from uuid import UUID


# Mock Pack IDs
PACK_ID_1 = "550e8400-e29b-41d4-a716-446655440001"
PACK_ID_2 = "550e8400-e29b-41d4-a716-446655440002"
PACK_ID_3 = "550e8400-e29b-41d4-a716-446655440003"
PACK_ID_4 = "550e8400-e29b-41d4-a716-446655440004"
PACK_ID_5 = "550e8400-e29b-41d4-a716-446655440005"


def get_pack_1() -> Dict[str, Any]:
    """Pack 1: Общие знания - полноразмерная игра"""
    return {
        "id": PACK_ID_1,
        "name": "Общие знания",
        "author": "SIGame Team",
        "description": "Классическая игра с вопросами по истории, географии, науке, спорту и культуре",
        "rounds_count": 1,
        "questions_count": 25,
        "created_at": "2024-01-01T00:00:00Z",
        "rounds": [
            {
                "id": "round-1-1",
                "round_number": 1,
                "name": "Первый раунд",
                "themes": [
                    {
                        "id": "theme-1-1",
                        "name": "История",
                        "questions": [
                            {"id": "q-1-1-1", "price": 100, "text": "В каком году началась Великая Отечественная война?", "answer": "1941", "media_type": "text"},
                            {"id": "q-1-1-2", "price": 200, "text": "Кто был первым президентом США?", "answer": "Джордж Вашингтон", "media_type": "text"},
                            {"id": "q-1-1-3", "price": 300, "text": "В каком году пала Берлинская стена?", "answer": "1989", "media_type": "text"},
                            {"id": "q-1-1-4", "price": 400, "text": "Кто открыл Америку в 1492 году?", "answer": "Христофор Колумб", "media_type": "text"},
                            {"id": "q-1-1-5", "price": 500, "text": "В каком году была Французская революция?", "answer": "1789", "media_type": "text"},
                        ]
                    },
                    {
                        "id": "theme-1-2",
                        "name": "География",
                        "questions": [
                            {"id": "q-1-2-1", "price": 100, "text": "Столица России", "answer": "Москва", "media_type": "text"},
                            {"id": "q-1-2-2", "price": 200, "text": "Самый большой океан на Земле", "answer": "Тихий океан", "media_type": "text"},
                            {"id": "q-1-2-3", "price": 300, "text": "На каком континенте находится Египет?", "answer": "Африка", "media_type": "text"},
                            {"id": "q-1-2-4", "price": 400, "text": "Самая высокая гора в мире", "answer": "Эверест (Джомолунгма)", "media_type": "text"},
                            {"id": "q-1-2-5", "price": 500, "text": "Самая длинная река в мире", "answer": "Нил или Амазонка", "media_type": "text"},
                        ]
                    },
                    {
                        "id": "theme-1-3",
                        "name": "Наука",
                        "questions": [
                            {"id": "q-1-3-1", "price": 100, "text": "Сколько планет в Солнечной системе?", "answer": "8", "media_type": "text"},
                            {"id": "q-1-3-2", "price": 200, "text": "Химический символ воды", "answer": "H2O", "media_type": "text"},
                            {"id": "q-1-3-3", "price": 300, "text": "Кто открыл закон всемирного тяготения?", "answer": "Исаак Ньютон", "media_type": "text"},
                            {"id": "q-1-3-4", "price": 400, "text": "Скорость света в вакууме (примерно)", "answer": "300 000 км/с", "media_type": "text"},
                            {"id": "q-1-3-5", "price": 500, "text": "Кто создал теорию относительности?", "answer": "Альберт Эйнштейн", "media_type": "text"},
                        ]
                    },
                    {
                        "id": "theme-1-4",
                        "name": "Спорт",
                        "questions": [
                            {"id": "q-1-4-1", "price": 100, "text": "Сколько игроков в футбольной команде на поле?", "answer": "11", "media_type": "text"},
                            {"id": "q-1-4-2", "price": 200, "text": "В каком городе проходили Олимпийские игры 2014 года?", "answer": "Сочи", "media_type": "text"},
                            {"id": "q-1-4-3", "price": 300, "text": "Сколько сетов нужно выиграть в теннисе на Большом шлеме?", "answer": "3 из 5", "media_type": "text"},
                            {"id": "q-1-4-4", "price": 400, "text": "Какая страна выиграла чемпионат мира по футболу в 2018 году?", "answer": "Франция", "media_type": "text"},
                            {"id": "q-1-4-5", "price": 500, "text": "Сколько золотых медалей завоевал Майкл Фелпс на Олимпиадах?", "answer": "23", "media_type": "text"},
                        ]
                    },
                    {
                        "id": "theme-1-5",
                        "name": "Культура",
                        "questions": [
                            {"id": "q-1-5-1", "price": 100, "text": "Кто написал роман 'Война и мир'?", "answer": "Лев Толстой", "media_type": "text"},
                            {"id": "q-1-5-2", "price": 200, "text": "Кто нарисовал Мону Лизу?", "answer": "Леонардо да Винчи", "media_type": "text"},
                            {"id": "q-1-5-3", "price": 300, "text": "Какой композитор написал 'Лунную сонату'?", "answer": "Людвиг ван Бетховен", "media_type": "text"},
                            {"id": "q-1-5-4", "price": 400, "text": "Автор трилогии 'Властелин колец'", "answer": "Джон Толкин", "media_type": "text"},
                            {"id": "q-1-5-5", "price": 500, "text": "В каком году вышел первый фильм о Гарри Поттере?", "answer": "2001", "media_type": "text"},
                        ]
                    }
                ]
            }
        ]
    }


def get_pack_2() -> Dict[str, Any]:
    """Pack 2: Быстрая игра - компактная версия для тестирования"""
    return {
        "id": PACK_ID_2,
        "name": "Быстрая игра",
        "author": "SIGame Team",
        "description": "Короткая игра из 9 вопросов для быстрого тестирования",
        "rounds_count": 1,
        "questions_count": 9,
        "created_at": "2024-01-01T00:00:00Z",
        "rounds": [
            {
                "id": "round-2-1",
                "round_number": 1,
                "name": "Единственный раунд",
                "themes": [
                    {
                        "id": "theme-2-1",
                        "name": "Простые вопросы",
                        "questions": [
                            {"id": "q-2-1-1", "price": 100, "text": "Сколько будет 2+2?", "answer": "4", "media_type": "text"},
                            {"id": "q-2-1-2", "price": 200, "text": "Столица Франции", "answer": "Париж", "media_type": "text"},
                            {"id": "q-2-1-3", "price": 300, "text": "Сколько дней в неделе?", "answer": "7", "media_type": "text"},
                        ]
                    },
                    {
                        "id": "theme-2-2",
                        "name": "Животные",
                        "questions": [
                            {"id": "q-2-2-1", "price": 100, "text": "Самое большое животное на Земле", "answer": "Синий кит", "media_type": "text"},
                            {"id": "q-2-2-2", "price": 200, "text": "Сколько ног у паука?", "answer": "8", "media_type": "text"},
                            {"id": "q-2-2-3", "price": 300, "text": "Какое животное является символом Австралии?", "answer": "Кенгуру", "media_type": "text"},
                        ]
                    },
                    {
                        "id": "theme-2-3",
                        "name": "Технологии",
                        "questions": [
                            {"id": "q-2-3-1", "price": 100, "text": "Кто основал компанию Apple?", "answer": "Стив Джобс", "media_type": "text"},
                            {"id": "q-2-3-2", "price": 200, "text": "Что означает WWW?", "answer": "World Wide Web", "media_type": "text"},
                            {"id": "q-2-3-3", "price": 300, "text": "В каком году был основан Google?", "answer": "1998", "media_type": "text"},
                        ]
                    }
                ]
            }
        ]
    }


def get_pack_3() -> Dict[str, Any]:
    """Pack 3: Расширенная игра - два раунда"""
    return {
        "id": PACK_ID_3,
        "name": "Расширенная игра",
        "author": "SIGame Team",
        "description": "Игра из двух раундов для тестирования переключения между раундами",
        "rounds_count": 2,
        "questions_count": 12,
        "created_at": "2024-01-01T00:00:00Z",
        "rounds": [
            {
                "id": "round-3-1",
                "round_number": 1,
                "name": "Первый раунд",
                "themes": [
                    {
                        "id": "theme-3-1",
                        "name": "Кино",
                        "questions": [
                            {"id": "q-3-1-1", "price": 100, "text": "Режиссер фильма 'Титаник'", "answer": "Джеймс Кэмерон", "media_type": "text"},
                            {"id": "q-3-1-2", "price": 200, "text": "Кто сыграл Железного человека в Marvel?", "answer": "Роберт Дауни-младший", "media_type": "text"},
                            {"id": "q-3-1-3", "price": 300, "text": "В каком году вышел первый фильм 'Звездные войны'?", "answer": "1977", "media_type": "text"},
                        ]
                    },
                    {
                        "id": "theme-3-2",
                        "name": "Музыка",
                        "questions": [
                            {"id": "q-3-2-1", "price": 100, "text": "Солист группы Queen", "answer": "Фредди Меркьюри", "media_type": "text"},
                            {"id": "q-3-2-2", "price": 200, "text": "Альбом Beatles, выпущенный в 1969 году", "answer": "Abbey Road", "media_type": "text"},
                            {"id": "q-3-2-3", "price": 300, "text": "Сколько струн у стандартной гитары?", "answer": "6", "media_type": "text"},
                        ]
                    }
                ]
            },
            {
                "id": "round-3-2",
                "round_number": 2,
                "name": "Второй раунд",
                "themes": [
                    {
                        "id": "theme-3-3",
                        "name": "Литература",
                        "questions": [
                            {"id": "q-3-3-1", "price": 200, "text": "Автор 'Мастера и Маргариты'", "answer": "Михаил Булгаков", "media_type": "text"},
                            {"id": "q-3-3-2", "price": 400, "text": "Главный герой романа 'Преступление и наказание'", "answer": "Родион Раскольников", "media_type": "text"},
                            {"id": "q-3-3-3", "price": 600, "text": "Сколько томов в 'Войне и мире'?", "answer": "4", "media_type": "text"},
                        ]
                    },
                    {
                        "id": "theme-3-4",
                        "name": "Искусство",
                        "questions": [
                            {"id": "q-3-4-1", "price": 200, "text": "Где находится Эрмитаж?", "answer": "Санкт-Петербург", "media_type": "text"},
                            {"id": "q-3-4-2", "price": 400, "text": "Кто написал 'Крик'?", "answer": "Эдвард Мунк", "media_type": "text"},
                            {"id": "q-3-4-3", "price": 600, "text": "Стиль живописи Пикассо", "answer": "Кубизм", "media_type": "text"},
                        ]
                    }
                ]
            }
        ]
    }


def get_pack_4() -> Dict[str, Any]:
    """Pack 4: Медиа-пак - вопросы с изображениями и видео для тестирования медиа"""
    return {
        "id": PACK_ID_4,
        "name": "Медиа-пак (тест)",
        "author": "SIGame Team",
        "description": "Пак с изображениями и видео для тестирования медиа-контента",
        "rounds_count": 1,
        "questions_count": 12,
        "created_at": "2024-01-01T00:00:00Z",
        "rounds": [
            {
                "id": "round-4-1",
                "round_number": 1,
                "name": "Медиа-раунд",
                "themes": [
                    {
                        "id": "theme-4-1",
                        "name": "Что на картинке?",
                        "questions": [
                            {
                                "id": "q-4-1-1",
                                "price": 100,
                                "text": "Что изображено на картинке?",
                                "answer": "Горы",
                                "media_type": "image",
                                "media_url": "https://picsum.photos/seed/mountains/800/600",
                                "media_duration_ms": 0
                            },
                            {
                                "id": "q-4-1-2",
                                "price": 200,
                                "text": "Какой город изображен на фото?",
                                "answer": "Город",
                                "media_type": "image",
                                "media_url": "https://picsum.photos/seed/city/800/600",
                                "media_duration_ms": 0
                            },
                            {
                                "id": "q-4-1-3",
                                "price": 300,
                                "text": "Что вы видите на этом изображении?",
                                "answer": "Природа",
                                "media_type": "image",
                                "media_url": "https://picsum.photos/seed/nature/800/600",
                                "media_duration_ms": 0
                            },
                        ]
                    },
                    {
                        "id": "theme-4-2",
                        "name": "Угадай место",
                        "questions": [
                            {
                                "id": "q-4-2-1",
                                "price": 100,
                                "text": "Где было сделано это фото?",
                                "answer": "Пляж",
                                "media_type": "image",
                                "media_url": "https://picsum.photos/seed/beach/800/600",
                                "media_duration_ms": 0
                            },
                            {
                                "id": "q-4-2-2",
                                "price": 200,
                                "text": "Какая архитектура показана?",
                                "answer": "Здание",
                                "media_type": "image",
                                "media_url": "https://picsum.photos/seed/architecture/800/600",
                                "media_duration_ms": 0
                            },
                            {
                                "id": "q-4-2-3",
                                "price": 300,
                                "text": "Что это за место?",
                                "answer": "Лес",
                                "media_type": "image",
                                "media_url": "https://picsum.photos/seed/forest/800/600",
                                "media_duration_ms": 0
                            },
                        ]
                    },
                    {
                        "id": "theme-4-3",
                        "name": "Видео вопросы",
                        "questions": [
                            {
                                "id": "q-4-3-1",
                                "price": 100,
                                "text": "Что показано в этом видео?",
                                "answer": "Кролики",
                                "media_type": "video",
                                "media_url": "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4",
                                "media_duration_ms": 10000
                            },
                            {
                                "id": "q-4-3-2",
                                "price": 200,
                                "text": "Какое событие происходит на видео?",
                                "answer": "Слёзы радости",
                                "media_type": "video",
                                "media_url": "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4",
                                "media_duration_ms": 10000
                            },
                            {
                                "id": "q-4-3-3",
                                "price": 300,
                                "text": "Что демонстрируется в ролике?",
                                "answer": "Стеклянные шары",
                                "media_type": "video",
                                "media_url": "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4",
                                "media_duration_ms": 15000
                            },
                        ]
                    },
                    {
                        "id": "theme-4-4",
                        "name": "Смешанные вопросы",
                        "questions": [
                            {
                                "id": "q-4-4-1",
                                "price": 100,
                                "text": "Обычный текстовый вопрос: Столица Италии?",
                                "answer": "Рим",
                                "media_type": "text",
                                "media_url": None,
                                "media_duration_ms": 0
                            },
                            {
                                "id": "q-4-4-2",
                                "price": 200,
                                "text": "Какое животное на картинке?",
                                "answer": "Животное",
                                "media_type": "image",
                                "media_url": "https://picsum.photos/seed/animal/800/600",
                                "media_duration_ms": 0
                            },
                            {
                                "id": "q-4-4-3",
                                "price": 300,
                                "text": "Ещё один текстовый: Сколько континентов на Земле?",
                                "answer": "7",
                                "media_type": "text",
                                "media_url": None,
                                "media_duration_ms": 0
                            },
                        ]
                    }
                ]
            }
        ]
    }


def get_pack_5() -> Dict[str, Any]:
    """Pack 5: Специальные вопросы - демонстрация типов secret, stake, forAll"""
    return {
        "id": PACK_ID_5,
        "name": "Специальные вопросы",
        "author": "SIGame Team",
        "description": "Пак со специальными типами вопросов: кот в мешке, ва-банк, вопрос для всех",
        "rounds_count": 1,
        "questions_count": 12,
        "created_at": "2024-01-15T00:00:00Z",
        "rounds": [
            {
                "id": "round-5-1",
                "round_number": 1,
                "name": "Специальный раунд",
                "themes": [
                    {
                        "id": "theme-5-1",
                        "name": "Обычные вопросы",
                        "questions": [
                            {"id": "q-5-1-1", "price": 100, "text": "Столица Японии", "answer": "Токио", "type": "normal", "media_type": "text"},
                            {"id": "q-5-1-2", "price": 200, "text": "Сколько цветов в радуге?", "answer": "7", "type": "normal", "media_type": "text"},
                            {"id": "q-5-1-3", "price": 300, "text": "Самый большой материк", "answer": "Евразия", "type": "normal", "media_type": "text"},
                        ]
                    },
                    {
                        "id": "theme-5-2",
                        "name": "Кот в мешке 🐱",
                        "questions": [
                            {"id": "q-5-2-1", "price": 100, "text": "Секретный вопрос: Какое животное считается королём зверей?", "answer": "Лев", "type": "secret", "media_type": "text"},
                            {"id": "q-5-2-2", "price": 200, "text": "Секретный вопрос: Сколько ног у паука?", "answer": "8", "type": "secret", "media_type": "text"},
                            {"id": "q-5-2-3", "price": 300, "text": "Секретный вопрос: Как называется детёныш собаки?", "answer": "Щенок", "type": "secret", "media_type": "text"},
                        ]
                    },
                    {
                        "id": "theme-5-3",
                        "name": "Ва-банк 💰",
                        "questions": [
                            {"id": "q-5-3-1", "price": 100, "text": "Ставка: В каком году основан Google?", "answer": "1998", "type": "stake", "media_type": "text"},
                            {"id": "q-5-3-2", "price": 200, "text": "Ставка: Сколько зубов у взрослого человека?", "answer": "32", "type": "stake", "media_type": "text"},
                            {"id": "q-5-3-3", "price": 300, "text": "Ставка: Назовите столицу Австралии", "answer": "Канберра", "type": "stake", "media_type": "text"},
                        ]
                    },
                    {
                        "id": "theme-5-4",
                        "name": "Вопрос для всех 👥",
                        "questions": [
                            {"id": "q-5-4-1", "price": 100, "text": "Для всех: Какой цвет получается при смешении синего и жёлтого?", "answer": "Зелёный", "type": "forAll", "media_type": "text"},
                            {"id": "q-5-4-2", "price": 200, "text": "Для всех: Назовите самую маленькую планету Солнечной системы", "answer": "Меркурий", "type": "forAll", "media_type": "text"},
                            {"id": "q-5-4-3", "price": 300, "text": "Для всех: В каком году закончилась Вторая мировая война?", "answer": "1945", "type": "forAll", "media_type": "text"},
                        ]
                    }
                ]
            }
        ]
    }


# All available packs
MOCK_PACKS = {
    PACK_ID_1: get_pack_1(),
    PACK_ID_2: get_pack_2(),
    PACK_ID_3: get_pack_3(),
    PACK_ID_4: get_pack_4(),
    PACK_ID_5: get_pack_5(),
}


def get_all_packs() -> List[Dict[str, Any]]:
    """Get list of all mock packs (without full content)"""
    return [
        {
            "id": pack["id"],
            "name": pack["name"],
            "author": pack["author"],
            "description": pack["description"],
            "rounds_count": pack["rounds_count"],
            "questions_count": pack["questions_count"],
            "created_at": pack["created_at"],
        }
        for pack in MOCK_PACKS.values()
    ]


def get_pack_by_id(pack_id: str) -> Dict[str, Any] | None:
    """Get pack by ID (without full content)"""
    pack = MOCK_PACKS.get(pack_id)
    if not pack:
        return None
    
    return {
        "id": pack["id"],
        "name": pack["name"],
        "author": pack["author"],
        "description": pack["description"],
        "rounds_count": pack["rounds_count"],
        "questions_count": pack["questions_count"],
        "created_at": pack["created_at"],
    }


def get_pack_content(pack_id: str) -> Dict[str, Any] | None:
    """Get full pack content with all questions"""
    return MOCK_PACKS.get(pack_id)


def pack_exists(pack_id: str) -> bool:
    """Check if pack exists"""
    return pack_id in MOCK_PACKS

