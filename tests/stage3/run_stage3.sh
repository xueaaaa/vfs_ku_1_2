#!/bin/sh
# Этап 3: загрузка VFS из XML и обработка ошибок загрузки.
set -e

go build -o vfs-shell .

echo "=== 1. Минимальная VFS ==="
./vfs-shell --vfs=../tests/stage3/vfs_minimal.xml --script=../tests/stage3/script_minimal.txt &
sleep 2 && kill $! 2>/dev/null

echo "=== 2. VFS с несколькими файлами в одной папке ==="
./vfs-shell --vfs=../tests/stage3/vfs_multiple.xml &
sleep 2 && kill $! 2>/dev/null

echo "=== 3. VFS с вложенностью >= 3 уровня ==="
./vfs-shell --vfs=../tests/stage3/vfs_nested.xml --script=../tests/stage3/script_nested.txt &
sleep 3 && kill $! 2>/dev/null

echo "=== 4. Ошибки навигации (cd в несуществующие папки) ==="
./vfs-shell --vfs=../tests/stage3/vfs_nested.xml --script=../tests/stage3/script_errors.txt &
sleep 2 && kill $! 2>/dev/null

echo "=== 5. Битый XML (невалидный формат) ==="
./vfs-shell --vfs=../tests/stage3/vfs_invalid.xml &
sleep 2 && kill $! 2>/dev/null

echo "=== 6. Несуществующий путь к VFS ==="
./vfs-shell --vfs=../tests/stage3/nonexistent.xml &
sleep 2 && kill $! 2>/dev/null

echo "Все прогоны запущены. Проверь вывод в каждом открывшемся окне вручную."
