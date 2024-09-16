#!/bin/bash

# Функция для генерации случайного имени длиной 4 буквы
generate_random_name() {
    cat /dev/urandom | tr -dc 'a-z' | fold -w 4 | head -n 1
}

# Цикл для создания 100 папок
for (( i=1; i<=100; i++ )); do
    folder_name=$(generate_random_name)
    mkdir -p "$folder_name"

    # В каждой папке создаем еще 10 случайных подпапок
    for (( j=1; j<=10; j++ )); do
        subfolder_name=$(generate_random_name)
        mkdir -p "$folder_name/$subfolder_name"
    done
done

echo "генерация папок завершена"