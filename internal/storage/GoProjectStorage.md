### Настройка моделей
- Добавлены Id
- Под каждую модель созданы кастомные ошибки
- Возможно для решения вопроса с id стоит добавить структуру без id и еще одну с прошлой структурой и полем id?

### Интерфейсы сущностей
Интерфейсы реализуют взаимодействие бд и наших моделей.
**Идея**: Максимально общие запросы к бд. Вся логика происходит на уровне выше. 

Имеют методы:
- Create
	Добавляет кортеж в бд на основе экземпляра go модели
- Get
	Дефолтная версия возвращает экземпляр go модели по id
- Update
	Обновляет данные кортежа с id на основе полей экземпляра go модели
- Delete
	Удаляет кортеж с id
- И кастомные под каждую сущность
#### User
Добавлен дополнительный метод GetByLogin для получения по логину. Так можно будет при входе получить по логину go модель.

В методах Create и Update. При передаче экземпляра их поле id является фиктивным. В методе Create еще и после создания изменяется id у экземпляра на присвоенный. 

#### Merch
Работает на id. Есть дополнительный метод, который возвращает лист всех мерчей.

#### Transactions
Сейчас реализован только метод Create. Он достаточно объемный. Стоит подумать над изменениями

## TO-DO
- Нужно настроить триггеры
- Создавать соединение в методах NewTStorage
