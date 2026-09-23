#!/bin/sh
# Этап 2: проверка всех комбинаций параметров командной строки.
set -e

go build -o vfs-shell .

echo "=== 1. Без флагов вообще (программа должна открыться и ничего не выполнять) ==="
./vfs-shell &
sleep 2 && kill $! 2>/dev/null

echo "=== 2. Только --script, валидный базовый сценарий ==="
./vfs-shell --script=../tests/stage2/script_basic.txt &
sleep 2 && kill $! 2>/dev/null

echo "=== 3. Скрипт из одних комментариев ==="
./vfs-shell --script=../tests/stage2/script_comments.txt &
sleep 2 && kill $! 2>/dev/null

echo "=== 4. Скрипт с ошибками (не должен ронять программу) ==="
./vfs-shell --script=../tests/stage2/script_errors.txt &
sleep 2 && kill $! 2>/dev/null

echo "=== 5. Пустой файл скрипта ==="
./vfs-shell --script=../tests/stage2/script_empty.txt &
sleep 2 && kill $! 2>/dev/null

echo "=== 6. Несуществующий путь к скрипту (проверка обработки ошибки открытия файла) ==="
./vfs-shell --script=../tests/stage2/nonexistent.txt &
sleep 2 && kill $! 2>/dev/null

echo "=== 7. Только --vfs без --script (на этапе 2 VFS ещё не грузится, но флаг должен приниматься) ==="
./vfs-shell --vfs=../tests/stage3/vfs_minimal.xml &
sleep 2 && kill $! 2>/dev/null

echo "Все прогоны запущены. Проверь вывод в каждом открывшемся окне вручную."
