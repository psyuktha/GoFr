from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from rag import db, chain, run

# Initialize FastAPI app
app = FastAPI()

# Pydantic model for API input
class QueryRequest(BaseModel):
    question: str

@app.post("/generate_post/")
async def generate_post(request: QueryRequest, template, prompt):
    try:
    # template = request.question
    # prompt = "Generate a social media post based on the following context."
        # # Template for generating posts
        template = """You are an AI content generator tasked with creating engaging, concise, and professional social media posts.

        Use only the information provided in the context.
        Do not add any extra details, assumptions, or speculative content.
        Maintain a tone suitable for [platform: e.g., LinkedIn, Twitter, Instagram, etc.].
        The post should be clear, concise, and adhere to any specified character limits or formatting guidelines.
        If the content is technical or professional, ensure the language is precise and jargon-free (if possible). For creative or casual posts, keep the tone friendly and approachable.

        Deliverables:
        Post text: A short, engaging post based on the context.
        Hashtags (if required): Relevant hashtags derived from the context.
        Call-to-action (optional): If appropriate, include a call-to-action to increase engagement.
        Note:
        Do not generate any content outside the context provided. If the context is insufficient, indicate that more information is required."""

        response=run(template,prompt)

        return {"status": "success", "response": response}

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
