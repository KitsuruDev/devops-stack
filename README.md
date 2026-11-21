# DevOps-stack

Демонстрационный проект: Go-бэкенд, SPA-фронтенд и nginx-реверс-прокси, развёртываемые через набор docker-compose файлов с blue/green и canary схемами, CI/CD на GitLab.

## Структура

- `backend/` - Go-сервис, multi-stage Dockerfile, non-root, HEALTHCHECK
- `frontend/` - Node build → nginx runtime, статика SPA
- `nginx/` - Dockerfile с BuildKit secret для TLS, конфиг реверс-прокси
- `deploy/` - compose-файлы и скрипты (smoke-color, canary-check, switch)
- `.env.example` - шаблон переменных окружения
- `.gitlab-ci.yml` - пайплайн

## Как устроен деплой

Единого compose-файла нет - их несколько, и они решают разные задачи:
- **Базовый** (`db`, `backend`, `frontend`, `nginx`): постгрес с volume и healthcheck, backend ждёт `service_healthy` у db, nginx слушает `8443`.
- **Green**: тот же набор в изолированной сети `green-net` (external), nginx слушает `9445:8443`.
- **Canary**: только `backend-canary` на `8444:8000` во внешней сети `devops-stack_default`.

Переключение трафика между blue и green выполняет `deploy/switch.sh`: он пишет `/etc/nginx/conf.d/gateway.conf` на хосте и делает `nginx -s reload`. Gateway-nginx находится вне compose, поэтому blue/green - два стека + хостовый шлюз, который на них указывает.

## Особенности docker-compose
1. Несколько файлов вместо одного: базовый, green, canary. Сценарий выбирается тем, какой файл запускается.
2. Сети `green-net` и `devops-stack_default` объявлены `external: true`: их нужно создать на хосте заранее, стек не самодостаточен.
3. `pull_policy` задаётся осознанно: у `db` - `if_not_present`, у `nginx` - `never`. Образ `notes-nginx` собирается локально на `web-server-1` и в реестр не публикуется.
4. Разные порты у blue/green/canary (`8443` / `9445` / `8444`) позволяют держать их одновременно и переключать шлюз без простоя.
5. Healthcheck на уровне compose - у `db` и `backend`. Nginx опирается на `HEALTHCHECK` из Dockerfile.
6. Canary-файл намеренно минимален: только backend, подключается к уже поднятой инфраструктуре.
7. Реальные секреты не коммитятся: `.env.example` - шаблон, TLS-ключи прокидываются как BuildKit secret на этапе сборки nginx.

## Особенности Dockerfile
- Все образы multi-stage, runtime slim (alpine).
- Backend: кэш go-модулей до копирования исходников, `CGO_ENABLED=0`, `-ldflags="-s -w"`, пользователь `app`, вызов `/app/healthcheck.sh`.
- Frontend: `node:22-alpine` для сборки, `nginx:1.27-alpine` в runtime, копируется только `dist`.
- Nginx: сертификаты монтируются через `--mount=type=secret` в builder-стадии, в runtime попадают только ключ и сертификат с правами `600`/`644`, `HEALTHCHECK` ходит на `https://127.0.0.1:8443/healthz`.

## Скрипты деплоя
- `smoke-color.sh <port>` - smoke-проверка конкретного цвета: `/health`, SPA, `/api/notes`.
- `canary-check.sh <host>` - 100 запросов на `/health` canary-backend, падает при успехе < 99.
- `switch.sh <blue|green>` - переключение gateway-nginx на нужный порт.
