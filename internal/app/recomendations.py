#!/usr/bin/env python3
# coding: utf-8
"""
recommendations.py
------------------
1) Загружает результаты из pickle-файла (DataFrame и эмбеддинги).
2) Принимает текст пользователя и параметры через CLI.
3) Предобрабатывает текст, строит эмбеддинг.
4) Считает косинусное сходство с эмбеддингами статей.
5) Выводит JSON-массив рекомендаций с полями title и url.
"""
import argparse, json, logging, pickle, re, sys
import pandas as pd
import torch
import numpy as np
import nltk
from nltk.corpus import stopwords
from pymystem3 import Mystem
from sentence_transformers import SentenceTransformer
from PyPDF2 import PdfReader

# Инициализация
nltk.download('stopwords', quiet=True)
russian_stopwords = set(stopwords.words('russian'))
mystem = Mystem()

def parse_args():
    parser = argparse.ArgumentParser(description='Get article recommendations')
    parser.add_argument('-i','--input', required=True, help='path to parsed_articles.pkl')
    parser.add_argument('-q','--query', help='text query (ignored if --file is set)')
    parser.add_argument('-f','--file', help='path to PDF file to extract query from')
    parser.add_argument('-n','--top-n', type=int, default=5, help='number of recommendations')
    parser.add_argument('-m','--model', default='sentence-transformers/distiluse-base-multilingual-cased')
    parser.add_argument('--device', choices=['cpu','cuda'], default=None)
    return parser.parse_args()

def clean_text(text: str) -> str:
    if not isinstance(text, str):
        return ""
    text = re.sub(r"<.*?>", "", text)
    text = re.sub(r"[^а-яА-Яa-zA-Z0-9\s]", " ", text)
    text = re.sub(r"\s+", " ", text)
    return text.strip()

def lemmatize_ru(text: str) -> str:
    try:
        text = text.lower()
        lemmas = mystem.lemmatize(text)
        return " ".join(lemma for lemma in lemmas if lemma.strip())
    except Exception:
        return text  # fallback

def remove_stopwords(text: str) -> str:
    try:
        return " ".join(word for word in text.split() if word not in russian_stopwords)
    except Exception:
        return text

def preprocess(text: str) -> str:
    t = clean_text(text)
    t = lemmatize_ru(t)
    t = remove_stopwords(t)
    return t

def recommend(query: str, articles: list, model: SentenceTransformer, top_n: int, device: str):
    # Предобработка и эмбеддинг запроса
    q = preprocess(query)
    q_emb = model.encode(q, convert_to_tensor=True, device=device)

    # Собираем эмбеддинги статей (torch.Tensor или np.ndarray)
    emb_list = []
    for art in articles:
        emb = art['embedding']
        if isinstance(emb, np.ndarray):
            emb = torch.tensor(emb, dtype=q_emb.dtype)
        emb_list.append(emb)
    article_embs = torch.stack(emb_list).to(q_emb.device)

    # Косинусное сходство всей пачкой
    scores = torch.nn.functional.cosine_similarity(q_emb.unsqueeze(0), article_embs, dim=1).cpu().numpy()

    # Сортировка и выбор топ-N
    top_idx = scores.argsort()[::-1][:top_n]
    recs = [{'title': articles[i]['title'], 'url': articles[i].get('link', '')} for i in top_idx]
    return recs

def main():
    args = parse_args()
    logging.basicConfig(format='%(asctime)s %(levelname)s: %(message)s', level=logging.INFO)

    # Загрузка данных
    try:
        with open(args.input, 'rb') as f:
            data = pickle.load(f)
    except Exception as e:
        print(json.dumps({'error': f'Failed to load pickle: {e}'}), file=sys.stdout)
        sys.exit(0)

    articles = data.get('articles_info') or data.get('articles')
    if not articles:
        print(json.dumps({'error': 'No articles_info found in pickle'}), file=sys.stdout)
        sys.exit(0)

    # Извлечение текста из PDF или параметра query
    if args.file:
        try:
            reader = PdfReader(args.file)
            pages = [page.extract_text() or '' for page in reader.pages]
            query = '\n'.join(pages)
        except Exception as e:
            print(json.dumps({'error': f"Failed to read PDF '{args.file}': {e}"}), file=sys.stdout)
            sys.exit(0)
    else:
        query = args.query or ''

    if not query.strip():
        print(json.dumps({'error': 'Empty query: provide --query or --file'}), file=sys.stdout)
        sys.exit(0)

    # Модель и устройство
    device = args.device or ('cuda' if torch.cuda.is_available() else 'cpu')
    model = SentenceTransformer(args.model)
    model.to(device)

    # Генерация рекомендаций
    try:
        recommendations = recommend(query, articles, model, args.top_n, device)
        print(json.dumps(recommendations, ensure_ascii=False))
    except Exception as e:
        print(json.dumps({'error': f'Failed to generate recommendations: {e}'}), file=sys.stdout)
        sys.exit(0)

if __name__ == '__main__':
    main()