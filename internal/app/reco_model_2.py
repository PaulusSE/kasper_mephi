#!/usr/bin/env python3
# coding: utf-8

import json
import sys
import pickle
import torch
import numpy as np
import spacy
from collections import Counter
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity
import re
import random
from typing import List, Dict, Any
import os
import warnings
import logging

# 1. Отключить Xet Storage предупреждения и все Hf hub ворнинги
os.environ["HF_HUB_DISABLE_XET"] = "1"
os.environ["TRANSFORMERS_VERBOSITY"] = "error"
os.environ["TRANSFORMERS_NO_ADVISORY_WARNINGS"] = "1"
warnings.filterwarnings("ignore")
logging.basicConfig(stream=sys.stderr)
try:
    from transformers import logging as hf_logging
    hf_logging.set_verbosity_error()
except Exception:
    pass

# Загрузка модели spaCy для русского языка
nlp = spacy.load("ru_core_news_sm")

# ------------- NLP preprocessing -------------
def clean_text(text: str) -> str:
    return ' '.join(re.findall(r'\b\w+\b', text.lower()))

def lemmatize_ru(text: str) -> str:
    doc = nlp(text)
    return ' '.join([token.lemma_ for token in doc if not token.is_stop and not token.is_punct])

def preprocess(text: str) -> str:
    return lemmatize_ru(clean_text(text))

def extract_keywords(text: str, top_k: int = 10) -> List[str]:
    words = [w for w in preprocess(text).split() if len(w) > 4]
    return [w for w, _ in Counter(words).most_common(top_k)]

# ----------- Категоризация статьи ----------
def guess_category_from_text(title: str, abstract: str) -> str:
    text = (title or '') + " " + (abstract or '')
    text = text.lower()
    if any(w in text for w in ['обзор', 'survey', 'review']):
        return "review"
    if any(w in text for w in ['метод', 'method', 'approach', 'алгоритм']):
        return "method"
    if any(w in text for w in ['эксперимент', 'experiment', 'case', 'case study', 'benchmark']):
        return "experiment"
    if any(w in text for w in ['применение', 'application', 'use', 'applied']):
        return "application"
    return "other"

# ----------- BOOST-функции -----------
def progress_boost(article: Dict[str, Any], progressiveness: int) -> (float, List[str]):
    boost, reasons = 0.0, []
    year = int(article.get('year', 0))
    cat = article.get('category', guess_category_from_text(article.get('title', ''), article.get('abstract', '')))
    if progressiveness < 50 and cat == 'review':
        boost += 0.13
        reasons.append("Обзор для раннего этапа")
    elif progressiveness > 80 and (cat == 'experiment' or (year >= 2023)):
        boost += 0.12
        reasons.append("Свежая/экспериментальная статья для высокого прогресса")
    return boost, reasons

def quality_boost(article: Dict[str, Any]) -> (float, List[str]):
    boost, reasons = 0.0, []
    if article.get('scopus', False):   boost += 0.07; reasons.append("Scopus")
    if article.get('wos', False):      boost += 0.08; reasons.append("WoS")
    if float(article.get('impact', 0)) > 1.5:
        boost += 0.05; reasons.append("Высокий импакт-фактор")
    if int(article.get('citations', 0)) > 20:
        boost += 0.05; reasons.append("Высокая цитируемость")
    if int(article.get('year', 0)) < 2018:
        boost -= 0.10; reasons.append("Старая статья")
    return boost, reasons

def collaborative_boost(article_id: str, interactions_log: Dict[str, List[str]], similar_users_scores: Dict[str, float]) -> (float, List[str]):
    boost, reasons = 0.0, []
    weighted_count = sum(
        score for user, score in similar_users_scores.items()
        if article_id in interactions_log.get(user, [])
    )
    if weighted_count > 0:
        boost += 0.05 * weighted_count
        reasons.append(f"Коллаборативный выбор (вес: {weighted_count:.2f})")
    return boost, reasons

def novelty_boost(article: Dict[str, Any], interactions_log: Dict[str, List[str]], user_id: str) -> (float, List[str]):
    boost, reasons = 0.0, []
    if user_id not in interactions_log or get_article_id(article) not in interactions_log[user_id]:
        boost += 0.02
        reasons.append("Новинка для вас")
    return boost, reasons

# ----------- LOG взаимодействий и профили -----------
def log_interaction(user_id: str, article_id: str, interactions_log_path: str = "interactions.json"):
    try:
        with open(interactions_log_path, 'r', encoding='utf-8') as f:
            log = json.load(f)
    except Exception:
        log = {}
    if user_id not in log:
        log[user_id] = []
    if article_id not in log[user_id]:
        log[user_id].append(article_id)
    with open(interactions_log_path, 'w', encoding='utf-8') as f:
        json.dump(log, f, ensure_ascii=False, indent=2)

def load_interactions_log(interactions_log_path: str = "interactions.json") -> Dict[str, List[str]]:
    try:
        with open(interactions_log_path, 'r', encoding='utf-8') as f:
            return json.load(f)
    except Exception:
        return {}

# ----------- Рекомендации -----------
def flatten_text(val):
    if isinstance(val, list):
        return '\n'.join(flatten_text(v) for v in val)
    if isinstance(val, dict):
        return ' '.join(flatten_text(x) for x in val.values())
    return str(val) if val is not None else ''

def get_article_id(art):
    return art.get('id') or art.get('doi') or art.get('link')

def recommend(
    user_id: str,
    profile: Dict[str, str],
    dissertation: str,
    conferences: List[Dict[str, str]],
    patents: List[Dict[str, str]],
    comments: List[str],
    articles: List[Dict[str, Any]],
    model,
    top_n: int,
    progressiveness: int,
    interactions_log: Dict[str, List[str]],
    user_profiles: Dict[str, List[str]],
    device: str,
    exploration: bool = True,
    explain: bool = True
) -> Dict[str, Any]:
    texts = [
        flatten_text(dissertation),
        flatten_text(profile.get('publications', '')),
        flatten_text(profile.get('comments', '')),
        flatten_text(conferences),
        flatten_text(patents),
        flatten_text(comments),
    ]
    profile_text = '\n'.join(texts)
    profile_emb = model.encode(preprocess(profile_text), convert_to_tensor=True, device=device)

    article_embs = torch.stack([torch.tensor(art['embedding']) for art in articles]).to(device)
    base_scores = cosine_similarity(profile_emb.unsqueeze(0).cpu(), article_embs.cpu())[0]

    user_keywords = extract_keywords(profile_text, top_k=15)
    similar_users_scores = {
        uid: len(set(user_keywords) & set(keywords)) / len(set(user_keywords) | set(keywords))
        for uid, keywords in user_profiles.items() if uid != user_id
    }
    similar_users_scores = {uid: score for uid, score in similar_users_scores.items() if score >= 0.2}

    recommendations = []
    explored_articles = set()
    for idx, art in enumerate(articles):
        reasons = []
        art_keywords = extract_keywords((art.get('title', '') or '') + " " + (art.get('abstract', '') or ''))
        overlap = len(set(user_keywords) & set(art_keywords)) / (len(set(user_keywords) | set(art_keywords)) or 1)
        if overlap < 0.03:
            continue
        if explain:
            reasons.append(f"Пересечение тем: {overlap:.2f}")

        score = base_scores[idx] + overlap * 0.2

        p_boost, p_reasons = progress_boost(art, progressiveness)
        score += p_boost
        if explain: reasons += p_reasons

        q_boost, q_reasons = quality_boost(art)
        score += q_boost
        if explain: reasons += q_reasons

        c_boost, c_reasons = collaborative_boost(get_article_id(art), interactions_log, similar_users_scores)
        score += c_boost
        if explain: reasons += c_reasons

        n_boost, n_reasons = novelty_boost(art, interactions_log, user_id)
        score += n_boost
        if explain: reasons += n_reasons

        category = art.get('category', guess_category_from_text(art.get('title', ''), art.get('abstract', '')))
        recommendations.append({
            'title': art['title'],
            'url': art['link'],
            'score': round(score, 4),
            'category': category,
            'reasons': reasons,
            'overlap': round(overlap, 3),
        })
        explored_articles.add(get_article_id(art))

    if exploration:
        not_seen = [art for art in articles if get_article_id(art) not in explored_articles and (user_id not in interactions_log or get_article_id(art) not in interactions_log[user_id])]
        if not_seen:
            rec = random.choice(not_seen)
            recommendations.append({
                'title': rec['title'],
                'url': rec['link'],
                'score': 0.01,
                'category': rec.get('category', guess_category_from_text(rec.get('title', ''), rec.get('abstract', ''))),
                'reasons': ["Статья добавлена для расширения круга поиска (exploration)"],
                'overlap': 0.0,
            })

    categories = {}
    for rec in recommendations:
        categories.setdefault(rec['category'], []).append(rec)

    result = {
        "all": sorted(recommendations, key=lambda x: x['score'], reverse=True)[:top_n],
        "by_category": {cat: sorted(lst, key=lambda x: x['score'], reverse=True)[:max(1, top_n // 2)]
                        for cat, lst in categories.items()}
    }
    return result


# ------------- Example usage -------------
if __name__ == "__main__":
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument('-i', '--input', type=str, required=True)
    parser.add_argument('-n', '--top_n', type=int, default=5)
    args = parser.parse_args()

    profile_json = sys.stdin.read()
    profile = json.loads(profile_json)

    profile["patents"] = profile.get("patents") or []
    profile["conferences"] = profile.get("conferences") or []
    profile["comments"] = profile.get("comments") or []

    # --- Читаем текст диссертации, если файл реально существует ---
    dissertation_text = ""
    dissertation_path = profile.get("dissertation")
    if dissertation_path and os.path.exists(dissertation_path):
        try:
            with open(dissertation_path, "r", encoding="utf-8", errors="ignore") as f:
                dissertation_text = f.read()
        except Exception as e:
            dissertation_text = ""
            print(f"Warning: не удалось прочитать файл диссертации: {e}", file=sys.stderr)
    else:
        dissertation_text = ""

    # 1. Загрузка статей
    with open(args.input, 'rb') as f:
        data = pickle.load(f)

    meta = data["meta"]
    embeddings = data["embeddings"]

    assert len(meta) == embeddings.shape[0], f"meta: {len(meta)}, embeddings: {embeddings.shape}"

    articles = []
    for i, meta_article in enumerate(meta):
        article = dict(meta_article)  # копия меты
        article['embedding'] = embeddings[i].tolist()  # numpy -> python list
        articles.append(article)

    # 2. Если используешь sentence-transformers, загрузи модель (пример):
    model = SentenceTransformer('distiluse-base-multilingual-cased')
    device = 'cpu'  # или 'cuda' если поддерживается

    # 3. Опционально — лог взаимодействий, профили пользователей
    interactions_log = {}
    user_profiles = {}

    # 4. Вызов основной функции рекомендаций
    result = recommend(
        user_id=profile["student_id"],
        profile=profile,
        dissertation=dissertation_text,  # <-- вот тут теперь текст, а не путь!
        conferences=profile.get("conferences", []),
        patents=profile.get("patents", []),
        comments=profile.get("comments", []),
        articles=articles,  # теперь это список статей с embedding!
        model=model,
        top_n=args.top_n,
        progressiveness=profile.get("progressiveness", 0),
        interactions_log=interactions_log,
        user_profiles=user_profiles,
        device=device,
        exploration=True,
        explain=True
    )

    # Добавить объяснение, если нет диссертации
    if not dissertation_text:
        for rec in result["all"]:
            rec["reasons"].append("Рекомендация без учета диссертации (файл отсутствует)")

    while len(result["all"]) < args.top_n:
        not_seen = [art for art in articles if art.get('link') not in [r['url'] for r in result["all"]]]
        if not_seen:
            rec = random.choice(not_seen)
            result["all"].append({
                'title': rec.get('title', ''),
                'url': rec.get('link', ''),
                'score': 0.0,
                'category': rec.get('category', guess_category_from_text(rec.get('title', ''), rec.get('abstract', ''))),
                'reasons': ["Заполнение выдачи случайной статьёй"],
                'overlap': 0.0,
            })
        else:
            break

    # 5. Выводим только массив рекомендаций
    print(json.dumps(result["all"], ensure_ascii=False))

