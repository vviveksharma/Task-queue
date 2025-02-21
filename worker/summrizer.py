import sys
import json
from transformers import pipeline
import torch

device = "mps" if torch.backends.mps.is_available() else "cpu"
summarizer = pipeline("summarization", model="facebook/bart-large-cnn")

def summarize_text(text):
    text_length = len(text.split()) 
    if text_length < 50:
        return "Text is too short for summarization."
    
    max_length = min(400, text_length // 2) 
    min_length = max(30, max_length // 3) 

    summary = summarizer(text, max_length=max_length, min_length=min_length, do_sample=False)
    
    return summary[0]["summary_text"]

if __name__ == "__main__":
    text = sys.stdin.read().strip()  # Read text from stdin
    summary = summarize_text(text)
    print(json.dumps({"summary": summary}))  # Output as JSON
