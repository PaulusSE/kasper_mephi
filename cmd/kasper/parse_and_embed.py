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
import argparse
import json
import pickle
import logging
import re
from pathlib import Path
import pandas as pd
import torch
from sentence_transformers import SentenceTransformer
import nltk
from nltk.corpus import stopwords
from pymystem3 import Mystem

try:
    russian_stopwords = set(stopwords.words('russian'))
except LookupError:
    nltk.download('stopwords')
    russian_stopwords = set(stopwords.words('russian'))

mystem = Mystem()


def clean_text(text: str) -> str:
    """Удаляет HTML-теги, неалфавитные символы и лишние пробелы."""
    if not isinstance(text, str):
        return ""
    text = re.sub(r"<.*?>", "", text)
    text = re.sub(r"[^а-яА-Яa-zA-Z0-9\s]", " ", text)
    text = re.sub(r"\s+", " ", text)
    return text.strip()


def lemmatize_ru(text: str) -> str:
    """Приводит текст к леммам с помощью pymystem3."""
    text = text.lower()
    lemmas = mystem.lemmatize(text)
    lemmas = [lemma.strip() for lemma in lemmas if lemma.strip()]
    return " ".join(lemmas)


def remove_stopwords(text: str) -> str:
    """Удаляет русские стоп-слова из текста."""
    words = text.split()
    words = [w for w in words if w not in russian_stopwords]
    return " ".join(words)


def preprocess_text(text: str) -> str:
    """Объединённая обработка: очистка, лемматизация, удаление стоп-слов."""
    t = clean_text(text)
    t = lemmatize_ru(t)
    t = remove_stopwords(t)
    return t


def parse_args():
    parser = argparse.ArgumentParser(description="Parse and embed articles")
    parser.add_argument(
        "-i", "--input",
        type=Path,
        required=True,
        help="Путь к JSON-файлу со статьями"
    )
    parser.add_argument(
        "-o", "--output",
        type=Path,
        required=True,
        help="Куда сохранить результирующий pickle-файл"
    )
    parser.add_argument(
        "-m", "--model",
        type=str,
        default="sentence-transformers/distiluse-base-multilingual-cased",
        help="Имя модели Sentence-BERT"
    )
    parser.add_argument(
        "-b", "--batch-size",
        type=int,
        default=32,
        help="Размер батча для кодирования эмбеддингов"
    )
    return parser.parse_args()


def main():
    args = parse_args()

    # Настройка логирования
    logging.basicConfig(
        format="%(asctime)s %(levelname)s: %(message)s",
        level=logging.INFO
    )

    logging.info(f"Чтение статей из {args.input}")
    with args.input.open("r", encoding="utf-8") as f:
        data = json.load(f)

    df = pd.DataFrame(data)
    logging.info(f"Загружено {len(df)} статей")

    # Предобработка текстов аннотаций
    logging.info("Начало предобработки текстов")
    df["abstract_final"] = df["abstract"].fillna("").apply(preprocess_text)

    # Подготовка списка текстов для эмбеддингов
    texts = df["abstract_final"].tolist()

    # Загрузка модели
    logging.info(f"Загрузка модели: {args.model}")
    model = SentenceTransformer(args.model)

    # Построение эмбеддингов батчами
    logging.info(f"Построение эмбеддингов батчами по {args.batch_size} шт.")
    embeddings = model.encode(
        texts,
        convert_to_tensor=True,
        batch_size=args.batch_size,
        show_progress_bar=True,
        device="cuda" if torch.cuda.is_available() else "cpu",
    )

    # Преобразование в NumPy (если нужно) и подготовка к сохранению
    embeddings_np = embeddings.detach().cpu().numpy()
    articles_info = []
    for row, emb in zip(data, embeddings):
        articles_info.append({
            "title": row.get("title", ""),
            "link": row.get("link", ""),
            "embedding": emb,
        })

    # Сохранение результатов
    logging.info(f"Сохранение результатов в {args.output}")
    with args.output.open("wb") as f:
        pickle.dump({
            "df": df,
            "articles_info": articles_info,
            "embeddings_np": embeddings_np,
        }, f)

    logging.info("Готово!")


if __name__ == "__main__":
    main()
