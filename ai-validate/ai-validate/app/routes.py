import time
from fastapi import APIRouter, HTTPException
from app.models import ApiKeyResponse, ResponseDetails, TextRequest


def build_router(model_pipeline):
    router = APIRouter()

    @router.get("/")
    async def root():
        return {"message": "API Key Detection Service is running"}

    @router.post("/detect-api-keys", response_model=ApiKeyResponse)
    async def detect_api_keys(request: TextRequest):
        if model_pipeline is None:
            raise HTTPException(status_code=500, detail="model pipeline not loaded")

        start_time = time.time()
        try:
            result = model_pipeline(request.text)
            processing_time = time.time() - start_time
            api_keys = [item["word"] for item in result if "word" in item]

            details = ResponseDetails(
                input_length=len(request.text),
                processing_time_seconds=round(processing_time, 4),
                valid_keys_count=len(api_keys),
            )

            return ApiKeyResponse(
                api_keys=api_keys,
                total_count=len(api_keys),
                details=details,
            )
        except Exception as e:
            raise HTTPException(status_code=500, detail=f"Error processing text: {str(e)}")

    @router.get("/health")
    async def health_check():
        return {"status": "healthy", "model_loaded": model_pipeline is not None}

    return router
