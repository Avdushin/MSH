## Установка

Для Windows:

Откройте Реестр (Registry Editor) нажав Win + R и введя regedit.

Перейдите в раздел HKEY_CLASSES_ROOT\*\shell.

Создайте новый ключ (папку) с именем, которое вы хотите видеть в контекстном меню (например, "Открыть с помощью Media Sherlock").

В созданном ключе создайте еще один ключ command.

В правой части окна укажите значение "C:\путь_к_вашему_скрипту\MediaSherlock.bat" "%1" для созданного ключа command. Убедитесь, что путь к вашему .bat файлу корректный.

Теперь при щелчке правой кнопкой мыши по файлу, у вас должен появиться пункт "Открыть с помощью Media Sherlock" в контекстном меню. Когда пользователь выберет этот пункт, ваша программа будет вызвана с передачей пути к файлу в качестве аргумента.