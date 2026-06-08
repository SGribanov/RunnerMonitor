# RunnerMonitor

[English version](README.md)

RunnerMonitor -- легковесное TUI- и CLI-приложение для мониторинга и управления
CI runner'ами. Оно сделано для рабочей станции, где одновременно
есть несколько GitHub Actions runner'ов на Windows и WSL/Linux, с перспективой
переезда runner'ов на отдельную машину в локальной сети.

Приложение объединяет локальное состояние runner'а, статус GitHub, очередь
workflow, принадлежность к проекту, безопасные команды управления и read-only
видимость GitHub-hosted jobs и удаленных self-hosted runner'ов.

## Назначение

- Автоматически находит GitHub Actions runner'ы в настроенных Windows, WSL и
  Linux каталогах.
- Показывает проект/репозиторий, к которому относится runner.
- Объединяет локальное состояние службы/процесса со статусом GitHub через
  `gh api`.
- Показывает queued и in-progress GitHub-hosted workflow jobs как read-only
  строки.
- Показывает self-hosted runner'ы, зарегистрированные в отслеживаемых GitHub
  репозиториях, даже если их runner-папок нет на этой машине.
- Показывает busy-статус, количество queued workflow и stale queued workflow.
- Автоматически обновляет состояние TUI каждые 5 секунд, не скрывая старую
  таблицу на время refresh.
- Один раз при запуске TUI проверяет GitHub Releases и показывает короткое
  уведомление, если доступна новая версия.
- Ограничивает повторные GitHub status-запросы во время автообновления TUI,
  чтобы снизить нагрузку на процессы и API.
- Позволяет запускать, останавливать, перезапускать, чистить, удалять и
  перенастраивать runner'ы.
- Защищает опасные операции dry-run режимом, проверкой busy-состояния и
  явным подтверждением.
- Дает project-scoped команды, которые может запускать Codex или оператор перед
  push/ожиданием CI.
- Поддерживает сохраненные SSH-профили удаленной runner-машины.
- Оставляет OneDev как будущий provider-oriented этап.

## Стек

- Go 1.26+
- Charmbracelet Bubble Tea для TUI
- Charmbracelet Bubbles для таблицы и поля ввода
- Charmbracelet Lip Gloss для оформления терминала
- GitHub CLI (`gh`) для данных GitHub Actions
- Windows PowerShell для Windows Services
- WSL и Linux systemd для Linux runner'ов
- SSH для доступа к удаленной runner-машине

## Требования

- Go 1.26+
- Авторизованный GitHub CLI с доступом к нужным репозиториям:

  ```powershell
  gh auth status
  ```

- Windows PowerShell для локальных Windows runner'ов.
- WSL для Linux runner'ов на текущей машине.
- `git` для определения текущего GitHub-проекта.

## Быстрый старт

Скачайте готовый Windows-пакет из
[GitHub Releases](https://github.com/SGribanov/RunnerMonitor/releases/latest):

```text
RunnerMonitor-v0.5.1-windows-x64.zip
```

Распакуйте ZIP и запустите TUI:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\runner-monitor.ps1
```

Или соберите из исходников локально:

Собрать exe и создать локальный config рядом с exe:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1
.\runner-monitor.ps1 --show-config
```

Открыть TUI:

```powershell
.\runner-monitor.ps1
```

Один раз вывести inventory:

```powershell
.\runner-monitor.ps1 --once
```

Провести аудит runner'ов:

```powershell
.\runner-monitor.ps1 --audit
```

Запустить runner'ы для текущего GitHub-проекта:

```powershell
.\runner-monitor.ps1 --start-current
```

## Конфигурация

RunnerMonitor читает настройки из `runner-monitor.json`, который лежит рядом с
исполняемым файлом. Для локальной сборки по умолчанию это:

```text
C:\Repos\RunnerMonitor\bin\runner-monitor.json
```

Создать или посмотреть config:

```powershell
.\runner-monitor.ps1 --init-config
.\runner-monitor.ps1 --show-config
```

`--show-config` маскирует `wslSudoPassword` как `<set>` или `<empty>` и никогда
не печатает реальное значение.

Пример config:

```json
{
  "projectsRoot": "C:\\Repos",
  "windowsRunnerRoots": [
    "C:\\Runners"
  ],
  "wslRunnerRoots": [
    "/home/gsv777/Runners"
  ],
  "linuxRunnerRoots": [
    "/opt/Runners",
    "/srv/Runners"
  ],
  "githubHostedRepos": [],
  "tuiRefreshIntervalSeconds": 5,
  "wslSudoPassword": ""
}
```

Поля:

| Поле | Назначение |
|---|---|
| `projectsRoot` | Корень проектов для `--project`, например `C:\Repos`. |
| `windowsRunnerRoots` | Windows-каталоги runner'ов для поиска и безопасного удаления папок. |
| `wslRunnerRoots` | WSL-каталоги runner'ов для поиска. |
| `linuxRunnerRoots` | Linux-каталоги runner'ов для отдельной Linux runner-машины. |
| `githubHostedRepos` | Дополнительные `owner/repo` или GitHub URL, для которых нужно мониторить GitHub-hosted workflow jobs даже без локального runner'а. |
| `tuiRefreshIntervalSeconds` | Интервал автообновления TUI в секундах. По умолчанию `5`, если поле отсутствует или некорректно. |
| `wslSudoPassword` | Значение sudo-пароля для WSL fallback. Хранить только в app-local config. |

Для тестов или особых запусков можно переопределить путь:

```powershell
$env:RUNNER_MONITOR_CONFIG = "D:\Temp\runner-monitor.json"
```

Не коммитьте настоящий `runner-monitor.json` и sudo-пароли.

## Использование TUI

Команды внутри TUI:

Информация о runner'ах автоматически обновляется каждые 5 секунд по умолчанию.
Другой интервал можно задать через `tuiRefreshIntervalSeconds` в
`runner-monitor.json`. Команда `refresh` запускает немедленное обновление; если
данные уже есть, таблица остается видимой до прихода новой информации.
Автообновление кратко переиспользует свежий GitHub status, чтобы не запускать
`gh` для каждого репозитория каждые 5 секунд; ручной `refresh` все равно
запрашивает свежий GitHub status.

GitHub-hosted Actions jobs отображаются как строки с host `github` и local
state `hosted`. RunnerMonitor находит hosted jobs для репозиториев, которые уже
есть в локальных runner'ах, а также для репозиториев из `githubHostedRepos`.
Эти строки read-only: lifecycle, cleanup, remove, delete и reprovision команды
намеренно пропускаются, потому что GitHub-hosted runners являются временной
инфраструктурой GitHub, а не локальными машинами.

Self-hosted runner'ы, которые GitHub возвращает для отслеживаемых репозиториев,
тоже отображаются, если RunnerMonitor не нашел соответствующую локальную
runner-папку. Такие партнерские или иные удаленные строки используют host
`github`, local state `remote` и path `(not local)`. Они показывают GitHub
status, busy, labels, OS, version и очередь. Lifecycle, cleanup, logs и
reprovision команды пропускаются, потому что RunnerMonitor не может управлять
runner-папкой или сервисом на чужой машине. `remove [N] confirm` и
`delete [N] confirm` могут отвязать такие remote-only self-hosted runner
регистрации через GitHub API после явного подтверждения; локальная папка и
сервис при этом не трогаются.

При запуске TUI один раз делает best-effort проверку GitHub Releases через
`gh`. Если интернета нет или проверка не удалась, RunnerMonitor запускается
обычно и не показывает ошибку. Если есть более новый релиз, короткое
уведомление появляется над строкой ввода команд.

Нажмите `h` или `?`, либо введите `help`, чтобы открыть краткую справку прямо
в TUI. В ней есть короткая инструкция и описание команд. `Esc` сначала закрывает
справку; если справка закрыта, `Esc` выходит из TUI.

| Команда | Описание |
|---|---|
| `h`, `?`, `help` | Показать или скрыть встроенную справку. |
| `refresh` | Немедленно обновить локальное и GitHub-состояние runner'ов. |
| Стрелки | Выбрать строку runner'а. |
| `start [N]` | Запустить runner `N` или выбранную строку без номера. |
| `stop [N]` | Остановить runner `N` или выбранную строку без номера. |
| `restart [N]` | Перезапустить runner `N` или выбранную строку без номера. |
| `force-stop [N]` | Остановить даже если GitHub показывает busy. Использовать осторожно. |
| `force-restart [N]` | Перезапустить даже если GitHub показывает busy. Использовать осторожно. |
| `clear [N]` | Безопасно очистить idle runner `N` или выбранную строку. |
| `clear idle` | Очистить все idle runner'ы. Busy runner'ы пропускаются. |
| `auto-clear on` | Включить safe cleanup после refresh. |
| `auto-clear off` | Выключить auto cleanup. |
| `remove [N]` | Dry-run отвязки runner'а от GitHub. |
| `remove [N] confirm` | Выполнить отвязку runner'а после подтверждения. |
| `delete [N] confirm` | Отвязать runner и удалить безопасную папку runner'а; для remote-only GitHub-строк только отвязать регистрацию. |
| `logs [N]` | Открыть логи runner'а. |
| `connect remote NAME` | Открыть удаленный RunnerMonitor TUI по SSH. |
| `q`, `quit`, `exit`, `Esc`, `Ctrl+C` | Выйти из TUI. |

Таблица адаптируется к размеру терминала. На узком окне скрываются
низкоприоритетные колонки, а ключевые поля project/runner/status/busy/queue
остаются выровненными. В TUI-таблице колонка `Busy` использует обычный текст
`true`/`false`, чтобы статусные колонки оставались выровненными и читаемыми.

После lifecycle-команд `start`, `stop` и `restart` TUI сразу запускает
обновление статуса, чтобы таблица подтянула состояние service/GitHub без
ожидания следующего auto-refresh. Командная строка принимает `d` первой
буквой, поэтому `delete N confirm` можно вводить нормально.

## CLI-команды

Inventory и аудит:

```powershell
.\runner-monitor.ps1 --once
.\runner-monitor.ps1 --audit
```

Управление runner'ами по репозиторию:

```powershell
.\runner-monitor.ps1 --start-repo SGribanov/DeltaG
.\runner-monitor.ps1 --stop-repo SGribanov/DeltaG
.\runner-monitor.ps1 --restart-repo SGribanov/DeltaG
```

Управление одним runner'ом по имени:

```powershell
.\runner-monitor.ps1 --start-runner ideabox-runner --repo SGribanov/IdeaBox
.\runner-monitor.ps1 --stop-runner ideabox-runner --repo SGribanov/IdeaBox
.\runner-monitor.ps1 --restart-runner ideabox-runner --repo SGribanov/IdeaBox
```

Управление runner'ами для текущего Git-проекта:

```powershell
.\runner-monitor.ps1 --start-current
.\runner-monitor.ps1 --stop-current
.\runner-monitor.ps1 --restart-current
```

`start` для service-managed runner'ов теперь является проверкой готовности, а
не только запросом к service manager. Для systemd runner'ов, включая WSL,
RunnerMonitor включает unit, запускает его, ждет локальный `active`, а затем
ждет, пока GitHub покажет runner как `online`.

`stop` для service-managed runner'ов ждет, пока локальная служба остановится, а
затем ждет, пока GitHub покажет runner как offline. Если Windows service runner
останавливается из обычного, не elevated TUI, RunnerMonitor открывает elevated
PowerShell helper через UAC; после подтверждения helper выполняет stop, а TUI
продолжает refresh, пока строка не обновится.

Безопасная очистка:

```powershell
.\runner-monitor.ps1 --clear-repo SGribanov/DeltaG
.\runner-monitor.ps1 --clear-current
.\runner-monitor.ps1 --clear-idle
.\runner-monitor.ps1 --clear-runner ideabox-runner
```

Отключение автозапуска:

```powershell
.\runner-monitor.ps1 --disable-autostart
```

Config:

```powershell
.\runner-monitor.ps1 --init-config
.\runner-monitor.ps1 --init-config --overwrite-config
.\runner-monitor.ps1 --show-config
```

Удаленная машина:

```powershell
.\runner-monitor.ps1 --configure-remote runnerbox
.\runner-monitor.ps1 --connect-remote runnerbox
```

Удаление runner'а по умолчанию работает как dry-run:

```powershell
.\runner-monitor.ps1 --remove-runner ideabox-runner --repo SGribanov/IdeaBox
.\runner-monitor.ps1 --remove-runner ideabox-runner --repo SGribanov/IdeaBox --confirm
.\runner-monitor.ps1 --remove-runner ideabox-runner --repo SGribanov/IdeaBox --confirm --delete-folder
```

Добавление runner'а настраивает существующую подготовленную папку runner
distribution. `--project` -- это имя папки проекта внутри `projectsRoot`.

```powershell
.\runner-monitor.ps1 --add-runner runner-monitor-win --project RunnerMonitor --runner-folder C:\Runners\SGribanov-RunnerMonitor\runner-monitor-win --labels "self-hosted,Windows,X64"
.\runner-monitor.ps1 --add-runner runner-monitor-win --project RunnerMonitor --runner-folder C:\Runners\SGribanov-RunnerMonitor\runner-monitor-win --labels "self-hosted,Windows,X64" --confirm --replace
```

С установкой и запуском службы после настройки:

```powershell
.\runner-monitor.ps1 --add-runner runner-monitor-win --project RunnerMonitor --runner-folder C:\Runners\SGribanov-RunnerMonitor\runner-monitor-win --labels "self-hosted,Windows,X64" --confirm --replace --service
```

## Удаленная runner-машина

После переезда runner'ов на отдельную машину RunnerMonitor запускается на этой
машине, а локально к нему можно подключаться по SSH.

Настроить или обновить профиль:

```powershell
.\runner-monitor.ps1 --configure-remote runnerbox
```

Промпт спросит:

- имя remote-профиля;
- SSH host/alias;
- OS хоста: `windows` или `linux`;
- путь к RunnerMonitor на удаленной машине;
- путь проекта по умолчанию.

Открыть сохраненный remote TUI:

```powershell
.\runner-monitor.ps1 --connect-remote runnerbox
```

Эквивалентная Windows SSH-команда:

```powershell
ssh -t runnerbox "powershell -NoProfile -ExecutionPolicy Bypass -File C:/Repos/RunnerMonitor/runner-monitor.ps1"
```

Запустить runner'ы на удаленной машине перед push/ожиданием CI:

```powershell
ssh runnerbox "cd C:/Repos/DeltaG; powershell -NoProfile -ExecutionPolicy Bypass -File C:/Repos/RunnerMonitor/runner-monitor.ps1 --start-current"
```

Linux host пример:

```powershell
ssh -t runnerbox "cd /opt/RunnerMonitor && ./runner-monitor"
ssh runnerbox "cd /srv/DeltaG && /opt/RunnerMonitor/runner-monitor --start-current"
```

Remote-профили хранятся отдельно от app-local runner config в пользовательском
каталоге config как `RunnerMonitor\remote-hosts.json`.

## Каталоги runner'ов

Предпочтительная текущая структура:

```text
C:\Runners\<owner>-<repo>\<runner-name>
/home/gsv777/Runners/<owner>-<repo>/<runner-name>
```

Для будущего выделенного Linux host:

```text
/opt/Runners/<owner>-<repo>/<runner-name>
/srv/Runners/<owner>-<repo>/<runner-name>
```

Миграция папок runner'ов ведется отдельно и только по одному runner'у за раз.
Не перемещайте и не удаляйте busy runner без явного подтверждения.

## Модель безопасности

- Busy runner'ы защищены по умолчанию.
- Удаление и reprovisioning по умолчанию dry-run.
- Удаление папки требует `--delete-folder` и разрешено только внутри
  настроенных safe runner roots.
- Удаление папки отвергает сами configured root-каталоги и нормализованный path
  traversal.
- Cleanup удаляет безопасные generated файлы вроде содержимого `_work` и
  installer archives, но сохраняет регистрацию и бинарники runner'а.
- `--show-config` никогда не печатает реальный WSL sudo password.
- Локальные config-файлы и секреты нельзя коммитить.
- GitHub registration/remove token редактируется из ошибок команд; передача
  token убрана из parent process arguments там, где это позволяют runner config
  scripts.

## Сборка и тесты

Тесты:

```powershell
go test ./...
```

Сборка:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1
```

Перегенерировать иконки:

```powershell
uv run --with pillow python .\scripts\generate-icon-assets.py
go run github.com/akavel/rsrc@v0.10.2 -arch amd64 -ico .\assets\runner-monitor-hourglass.ico -o .\cmd\runner-monitor\runner-monitor_windows_amd64.syso
```

## Структура репозитория

```text
cmd/runner-monitor/     Точка входа приложения
internal/app/           Discovery, GitHub integration, lifecycle, TUI, cleanup
scripts/                Build, migration, cleanup и automation helpers
assets/                 Иконки песочных часов
reports/                Runner audit reports
research/               Долгоживущие project insights
tasks/                  Планы и статусы по GitHub issues
```

## Участие

См. [CONTRIBUTING.md](CONTRIBUTING.md). Любые изменения должны сохранять
модель безопасности: никаких незащищенных destructive runner operations,
никаких секретов в git и никаких silent-изменений runner registration. Также
соблюдайте [Code of Conduct](CODE_OF_CONDUCT.md).

## Безопасность

См. [SECURITY.md](SECURITY.md). Не публикуйте в issues секреты, runner tokens,
sudo passwords и приватные hostnames.

## Лицензия

RunnerMonitor распространяется по [MIT License](LICENSE).
