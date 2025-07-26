from fastapi import FastAPI
from pydantic import BaseModel, Field
from . import model as model_handler

app = FastAPI(
    title = "Text Similarity API",
    description = "An API to compute the semantic similarity between two sentences.",
    version = "1.0.0"
)

class SentenceInput(BaseModel):
    sentence1: str = Field(..., description="The first sentence.", example = "AI is transforming the world.")
    sentence2: str = Field(..., description="The second sentence.", example = "Artificial intelligence is changing society.")

class SimilarityResponse(BaseModel):
    sentence1: str
    sentence2: str
    similarity: float = Field(..., description = "The similarity score between the two sentences, ranging from 0 to 1.")

@app.post("/similarity", response_model = SimilarityResponse, summary = "Calculate Sentence Similarity")
async def get_similarity(data: SentenceInput):
    similarity_score = model_handler.calculate_similarity(
        data.sentence1,
        data.sentence2
    )
    
    return {
        "sentence1": data.sentence1,
        "sentence2": data.sentence2,
        "similarity": round(similarity_score, 4)
    }

@app.get("/", summary = "Root Endpoint")
async def read_root():
    return {"message": "Welcome to the Text Similarity API. Go to /docs for API documentation."}