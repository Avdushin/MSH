# Media Sherlock

<iframe width="560" height="315" src="https://www.youtube.com/embed/5sXCtyXH3B0?si=rtS_WDzFB-nBlnc9" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>

## Установка

### Для Windows:

### Обычная установка

1. Создайте папку распакуйте архив `ffmpegs.zip` по пути `C:\Program Files\`
2. Выполните файл `install.reg`
3. Назначте путь `C:\Program Files\ffmpegs\` в PATH
   1. Нажмите клавишу `Win`, наберите и откройте "Изменение системных переменных среды"
   2. В открывшемся окне выберите пункт "Переменные среды"
   3. В пункте "Системные переменные" найдите переменную "Path" и дважды кликните по ней
   4. Нажмите кнопку "Создать"
   5. Вставьте путь до папки ffmpegs: `C:\Program Files\ffmpegs`
4. Установка завершена
---

### Ручная настройка Regedit

Данный пункт, необходимо использовать, если вы не воспроизвели файл `install.reg`

Переместите папку `ffmpegs` по пути `C:\Program Files\`

Откройте Реестр (Registry Editor) нажав Win + R и введя `regedit`.

Перейдите в раздел `HKEY_CLASSES_ROOT\*\shell`

Создайте новый ключ (папку) с именем, которое вы хотите видеть в контекстном меню (например, "Open with Media Sherlock").

В созданном ключе создайте еще один ключ (папку) `command`.

В правой части окна укажите значение `"C:\путь_к_прграмме\MediaSherlock.exe" "%1"`` для созданного ключа command.

#### Актуальный путь к программе:

```shell
"C:\Program Files\ffmpegs\MSH.exe" "%1"
```

Теперь при щелчке правой кнопкой мыши по файлу, у вас должен появиться пункт "Open with Media Sherlock" в контекстном меню. Когда пользователь выберет этот пункт, ваша программа будет вызвана с передачей пути к файлу в качестве аргумента.