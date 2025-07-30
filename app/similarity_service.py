from sentence_transformers import SentenceTransformer, util
import json
import sys
import logging
from typing import Dict, Any

logging.basicConfig(level = logging.INFO, format = '%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class SimilarityService:
    def __init__(self, model_name: str = 'sentence-transformers/all-MiniLM-L6-v2'):
        try:
            logger.info(f"Loading model: {model_name}")
            self.model = SentenceTransformer(model_name)
            logger.info("Model loaded successfully")
        except Exception as e:
            logger.error(f"Failed to load model: {e}")
            raise

    def calculate_similarity(self, sentence1: str, sentence2: str) -> float:
        try:
            embedding1 = self.model.encode(sentence1, convert_to_tensor=True)
            embedding2 = self.model.encode(sentence2, convert_to_tensor=True)
            cosine_score = util.cos_sim(embedding1, embedding2)
            similarity = float(cosine_score.item())
            similarity = max(0.0, min(1.0, similarity))
            logger.info(f"Calculated similarity: {similarity:.4f}")
            return similarity
            
        except Exception as e:
            logger.error(f"Error calculating similarity: {e}")
            raise

def process_request(service: SimilarityService, request_data: Dict[str, Any]) -> Dict[str, Any]:
    try:
        sentence1 = request_data.get('sentence1', '').strip()
        sentence2 = request_data.get('sentence2', '').strip()
        
        if not sentence1 or not sentence2:
            return {"error": "Both sentence1 and sentence2 must be provided and non-empty"}
        similarity = service.calculate_similarity(sentence1, sentence2)
        
        return {
            "similarity": round(similarity, 6)
        }
        
    except Exception as e:
        logger.error(f"Error processing request: {e}")
        return {"error": f"Processing failed: {str(e)}"}

def main():
    try:
        service = SimilarityService()
        input_data = sys.stdin.read().strip()
        if not input_data:
            response = {"error": "No input data received"}
        else:
            try:
                request_data = json.loads(input_data)
                response = process_request(service, request_data)
            except json.JSONDecodeError as e:
                response = {"error": f"Invalid JSON input: {str(e)}"}
        
        print(json.dumps(response))
        
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        error_response = {"error": f"Service error: {str(e)}"}
        print(json.dumps(error_response))
        sys.exit(1)

if __name__ == "__main__":
    main()