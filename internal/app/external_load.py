import requests
from requests.auth import HTTPBasicAuth
import json

API_URL = "https://iasp.mephi.ru/ba_students_api/get_students_22_kaf"
LOGIN = ""
PASSWORD = ""

def fetch_students():
    try:
        response = requests.get(API_URL, auth=HTTPBasicAuth(LOGIN, PASSWORD))
        response.raise_for_status()

        # Преобразуем строку в JSON
        data = json.loads(response.text)
        # Получаем список студентов
        students = data.get("students", [])
        return students
    except requests.RequestException as e:
        print(f"Ошибка при запросе: {e}")
        return []
    except json.JSONDecodeError as e:
        print(f"Ошибка при декодировании JSON: {e}")
        return []

def save_to_json(data, filename="students.json"):
    try:
        with open(filename, "w", encoding="utf-8") as f:
            json.dump(data, f, ensure_ascii=False, indent=4)
        print(f"Данные успешно сохранены в {filename}")
    except Exception as e:
        print(f"Ошибка при сохранении файла: {e}")

def main():
    students = fetch_students()

    if not students:
        print("Студенты не найдены или произошла ошибка.")
        return

    # Сохраняем все данные в JSON
    save_to_json(students)

    # # Для наглядности можно вывести первые несколько студентов
    # for student in students[:5]:  # выводим только первые 5
    #     print("ФИО:", student.get("full_name"))
    #     print("Научный руководитель:", student.get("teacher"))
    #     print("Тема диссертации:", student.get("theme"))
    #     print("-" * 40)

if __name__ == "__main__":
    main()