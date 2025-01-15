# #!/usr/bin/env python3
# # coding: utf-8

# """
# recommend.py
# ------------
# 1) Загружает результаты из 'parsed_articles.pkl' (DataFrame и эмбеддинги).
# 2) Принимает текст пользователя (через аргументы командной строки или другой механизм).
# 3) Предобрабатывает этот текст, строит эмбеддинг (Sentence-BERT).
# 4) Считает косинусное сходство с эмбеддингами статей.
# 5) Выводит top-N рекомендаций (по умолчанию 5).
# """

# import sys
# import pickle
# import torch
# import pandas as pd
# from sentence_transformers import SentenceTransformer

# # Эти функции можно либо продублировать, либо вынести в общий utils-модуль:
# import re
# import nltk
# from nltk.corpus import stopwords
# from pymystem3 import Mystem

# nltk.download('stopwords')
# russian_stopwords = set(stopwords.words('russian'))
# mystem = Mystem()

# def clean_text(text: str) -> str:
#     if not isinstance(text, str):
#         return ""
#     text = re.sub(r"<.*?>", "", text)
#     text = re.sub(r"[^а-яА-Яa-zA-Z0-9\s]", " ", text)
#     text = re.sub(r"\s+", " ", text)
#     return text.strip()

# def lemmatize_ru(text: str) -> str:
#     text = text.lower()
#     lemmas = mystem.lemmatize(text)
#     lemmas = [lemma.strip() for lemma in lemmas if lemma.strip()]
#     return " ".join(lemmas)

# def remove_stopwords(text: str) -> str:
#     words = text.split()
#     words = [w for w in words if w not in russian_stopwords]
#     return " ".join(words)

# def preprocess_text(text: str) -> str:
#     t = clean_text(text)
#     t = lemmatize_ru(t)
#     t = remove_stopwords(t)
#     return t

# def recommend_articles(user_text, articles_info, model, top_n=5):
#     """
#     user_text: строка с текстом (запрос пользователя)
#     articles_info: список словарей [{title, link, embedding}, ...]
#     model: SentenceTransformer
#     top_n: число рекомендаций

#     Возвращает список словарей [{title, link, score}, ...].
#     """
#     # Предобработка текста пользователя
#     text_preproc = preprocess_text(user_text)

#     # Эмбеддинг для текста пользователя
#     user_emb = model.encode(text_preproc, convert_to_tensor=True)

#     # Считаем сходство для каждой статьи
#     scores = []
#     for article in articles_info:
#         article_emb = article["embedding"]  # torch.Tensor
#         score = torch.nn.functional.cosine_similarity(user_emb, article_emb, dim=0)
#         scores.append(score.item())

#     # Превращаем в DataFrame для сортировки
#     df = pd.DataFrame(articles_info)
#     df["score"] = scores

#     # Сортируем по убыванию
#     df_sorted = df.sort_values("score", ascending=False)

#     # Берём top-N
#     top_recs = df_sorted.head(top_n)
#     recommendations = []
#     for _, row in top_recs.iterrows():
#         recommendations.append({
#             "title": row["title"],
#             "link": row["link"],
#             "score": row["score"]
#         })
#     return recommendations

# def main():
#     # 1. Загружаем pickle, созданный parse_and_embed.py
#     with open("parsed_articles.pkl", "rb") as f:
#         data = pickle.load(f)
#     articles_info = data["articles_info"]
#     # df = data["df"]            # если вдруг понадобятся исходные поля DataFrame
#     # embeddings_np = data["embeddings_np"]  # тоже, если нужно

#     # 2. Загружаем ту же модель Sentence-BERT
#     #    (должна быть совпадающей с тем, что в parse_and_embed.py)
#     model_name = "sentence-transformers/distiluse-base-multilingual-cased"
#     model = SentenceTransformer(model_name)

#     # 3. Получаем текст пользователя из аргументов командной строки
#     #    Пример: python3 recommend.py "Текст для рекомендаций"
#     if len(sys.argv) < 2:
#         print("Использование: python3 recommend.py \"Текст для рекомендаций\"")
#         sys.exit(1)

#     user_text = sys.argv[1]

#     # 4. Получаем рекомендации
#     top_n = 5  # или взять из sys.argv[2] при желании
#     recs = recommend_articles(user_text, articles_info, model, top_n=top_n)

#     # 5. Выводим рекомендации
#     print(f"\n=== Рекомендации для текста: \"{user_text}\" ===")
#     for i, r in enumerate(recs, start=1):
#         print(f"{i}. {r['title']} (score={r['score']:.4f}) link={r['link']}")

# if __name__ == "__main__":
#     main()
import sys
import json
import random

def main():
    # 1. Читаем исходные статьи
    json_file_path = "articles.json"
    with open(json_file_path, "r", encoding="utf-8") as f:
        articles = json.load(f)
    
    # 2. Случайным образом выбираем 4 статьи
    # Если в articles < 4, нужно либо сделать проверку, либо взять len(articles)
    num_to_take = 4
    if len(articles) < num_to_take:
        # На всякий случай, если статей меньше 4, берём все
        num_to_take = len(articles)
    random_articles = random.sample(articles, num_to_take)

    # 3. Выводим результат в формате JSON (stdout)
    # Go-код сможет распарсить это как []Recommendation
    print(json.dumps(random_articles, ensure_ascii=False))

if __name__ == "__main__":
    main()