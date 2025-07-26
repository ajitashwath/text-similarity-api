from sentence_transformers import SentenceTransformer, util

model = SentenceTransformer('sentence-transformers/all-MiniLM-L6-v2')

def calculate_similarity(sentence1: str, sentence2: str) -> float:
    embedding1 = model.encode(sentence1, convert_to_tensor = True)
    embedding2 = model.encode(sentence2, convert_to_tensor = True)

    cosine_score = util.cos_sim(embedding1, embedding2)
    return cosine_score.item()