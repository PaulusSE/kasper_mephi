#!/usr/bin/env python3
# coding: utf-8

"""
parse_and_embed.py
------------------
1) Чтение JSON (articles_copy.json) или парсинг статей из какого-либо источника.
2) Предобработка (clean_text, lemmatize_ru, remove_stopwords).
3) Построение эмбеддингов (Sentence-BERT).
4) Сохранение результата (DataFrame и эмбеддинги) в pickle.
"""

import re
import json
import pickle
import pandas as pd
import nltk
import matplotlib.pyplot as plt

from nltk.corpus import stopwords
from pymystem3 import Mystem
from razdel import sentenize
from wordcloud import WordCloud
from sentence_transformers import SentenceTransformer
import torch

nltk.download('stopwords')
russian_stopwords = set(stopwords.words('russian'))
mystem = Mystem()

def clean_text(text: str) -> str:
    """Удаляет HTML-теги, неалфавитные символы, лишние пробелы."""
    if not isinstance(text, str):
        return ""
    text = re.sub(r"<.*?>", "", text)
    text = re.sub(r"[^а-яА-Яa-zA-Z0-9\s]", " ", text)
    text = re.sub(r"\s+", " ", text)
    return text.strip()

def lemmatize_ru(text: str) -> str:
    text = text.lower()
    lemmas = mystem.lemmatize(text)
    lemmas = [lemma.strip() for lemma in lemmas if lemma.strip()]
    return " ".join(lemmas)

def remove_stopwords(text: str) -> str:
    words = text.split()
    words = [w for w in words if w not in russian_stopwords]
    return " ".join(words)

def split_into_sentences(text: str):
    return [sent.text for sent in sentenize(text)]

def preprocess_text(text: str) -> str:
    """Объединённая функция очистки + лемматизации + удаления стоп-слов."""
    t = clean_text(text)
    t = lemmatize_ru(t)
    t = remove_stopwords(t)
    return t

def main():
    # 1. Читаем JSON со статьями (или парсим каким-то иным способом)
    json_file_path = "articles_copy.json"
    with open(json_file_path, "r", encoding="utf-8") as f:
        data = json.load(f)
    print(f"Загружено статей: {len(data)}")

    # 2. В DataFrame
    df = pd.DataFrame(data)
    print("Колонки в DataFrame:", df.columns.tolist())

    # 3. Предобработка
    df["abstract_clean"] = df["abstract"].apply(clean_text)
    df["abstract_lemmatized"] = df["abstract_clean"].apply(lemmatize_ru)
    df["abstract_final"] = df["abstract_lemmatized"].apply(remove_stopwords)

    # (Опционально) аналогично и для title:
    df["title_clean"] = df["title"].apply(clean_text)
    df["title_lemmatized"] = df["title_clean"].apply(lemmatize_ru)
    df["title_final"] = df["title_lemmatized"].apply(remove_stopwords)

    # Посмотрим статистику
    df["word_count"] = df["abstract_final"].apply(lambda x: len(x.split()))
    print("\n--- Статистика по текстам ---")
    print(f"Всего статей: {len(df)}")
    print(f"Среднее число слов в аннотации: {df['word_count'].mean():.2f}")

    # Пример: построим облако слов (wordcloud)
    all_words = []
    for text in df["abstract_final"]:
        tokens = text.split()
        all_words.extend(tokens)
    all_text = " ".join(all_words)
    if all_text.strip():
        wordcloud = WordCloud(width=800, height=400, background_color='white').generate(all_text)
        plt.figure(figsize=(10,5))
        plt.imshow(wordcloud, interpolation='bilinear')
        plt.axis("off")
        plt.title("Word Cloud по очищенным аннотациям", fontsize=14)
        # plt.show()  # Если у вас нет GUI, можете сохранить в файл
        plt.savefig("wordcloud.png")
        print("Сохранено облако слов в 'wordcloud.png'")
    else:
        print("Нет слов для WordCloud — возможно, данные пустые или сильно отфильтрованы.")

    # 4. Разделение на предложения (если нужно)
    df["abstract_sentences"] = df["abstract_final"].apply(split_into_sentences)

    # 5. Строим эмбеддинги (Sentence-BERT)
    model_name = "sentence-transformers/distiluse-base-multilingual-cased"
    print(f"\nЗагружаем модель {model_name}...")
    model = SentenceTransformer(model_name)
    abstracts_list = df["abstract_final"].tolist()
    print(f"Считаем эмбеддинги для {len(abstracts_list)} статей...")
    embeddings = model.encode(abstracts_list, convert_to_tensor=True)  # torch.Tensor

    # Превратим в numpy при желании
    embeddings_np = embeddings.detach().cpu().numpy()

    # 6. Сохраняем результат в pickle
    #    Чтобы потом recommend.py мог это использовать
    articles_info = []
    for i, row in df.iterrows():
        articles_info.append({
            "title": row["title"],
            "link": row.get("link", ""),
            # храним torch.Tensor (или можно embeddings_np[i], если хотим NumPy)
            "embedding": embeddings[i],
        })

    with open("parsed_articles.pkl", "wb") as f:
        pickle.dump({
            "df": df,
            "articles_info": articles_info,
            "embeddings_np": embeddings_np
        }, f)
    print("\n[OK] Данные и эмбеддинги сохранены в 'parsed_articles.pkl'.")

if __name__ == "__main__":
    main()
